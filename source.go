package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Source is the main entry point for fetching configuration
// values. The boolean value must be false if the source is
// unable to find an entry for the given path. If found, the
// raw value is returned and it is the responsibility of the
// consumer to verify the type or quality of the value.
type Source interface {
	Get(ctx context.Context, path ...string) (interface{}, bool)
}

// MapSource implements the Source interface for any map[string]interface{}.
// This implementation is intended to support the most common configuration
// sources which would be JSON, YAML, and ENV.
//
// Note: All keys should lower case (if applicable for the character set)
// 		 as lower case will also be applied to all lookup paths.
type MapSource struct {
	Map map[string]interface{}
}

func lowerCaseMap(m map[string]interface{}) map[string]interface{} {
	stack := []map[string]interface{}{m}
	for len(stack) > 0 {
		var current map[string]interface{}
		current, stack = stack[len(stack)-1], stack[:len(stack)-1]
		for k, v := range current {
			if mv, ok := v.(map[string]interface{}); ok {
				stack = append(stack, mv)
			}
			// The go language specification details the expected behavior when
			// modifying a map while iterating. This could be problematic in other
			// languages but the go specification allows it. New keys added *may*
			// be iterated over in the loop at some point and elements removed
			// are no longer iterated over.
			newKey := strings.ToLower(k)
			if newKey != k {
				// Only delete the element if we actually changing the key.
				// This is how we manage the fact that elements we previously
				// changed might actually show up again but in their already
				// correct form. Guarding with this switch prevents the key
				// from being removed or created which prevents replay.
				delete(current, k)
			}
			current[newKey] = v
		}
	}
	return m
}

// NewMapSource is the recommended way to create a MapSource instance.
// While they can be created with any map[string]interface{}, this constructor
// ensures that all keys of the map have a consistent case applied.
func NewMapSource(m map[string]interface{}) *MapSource {
	return &MapSource{Map: lowerCaseMap(m)}
}

// Get traverses a configuration map until it finds the requested element
// or reaches a dead end.
func (s *MapSource) Get(_ context.Context, path ...string) (interface{}, bool) {
	location := s.Map
	for x := 0; x < len(path)-1; x = x + 1 {
		pth := strings.ToLower(path[x])
		if _, found := location[pth]; !found {
			return nil, false
		}
		var ok bool
		if location, ok = location[pth].(map[string]interface{}); !ok {
			return nil, false
		}
	}
	v, ok := location[strings.ToLower(path[len(path)-1])]
	return v, ok
}

// NewJSONSource generates a config source from a JSON string.
func NewJSONSource(b []byte) (*MapSource, error) {
	v := make(map[string]interface{})
	err := json.Unmarshal(b, &v)
	return NewMapSource(v), err
}

// NewYAMLSource generates a config source from a YAML string.
func NewYAMLSource(b []byte) (*MapSource, error) {
	v := make(map[string]interface{})
	err := yaml.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}
	v, err = convertYamlMap(v)
	return NewMapSource(v), err
}

// convertYamlMap adapts the internal YAML types representing a map
// into the native go types so that the resulting map works with the
// type introspection used elsewhere.
func convertYamlMap(m map[string]interface{}) (map[string]interface{}, error) {
	stack := []map[string]interface{}{m}
	for len(stack) > 0 {
		var current map[string]interface{}
		current, stack = stack[len(stack)-1], stack[:len(stack)-1]
		for k, v := range current {
			switch vv := v.(type) {
			case map[interface{}]interface{}:
				tmp := make(map[string]interface{})
				for sk, sv := range vv {
					if _, ok := sk.(string); !ok {
						return nil, fmt.Errorf("map key %v is not a string", sk)
					}
					tmp[sk.(string)] = sv
				}
				delete(current, k)
				current[k] = tmp
				stack = append(stack, tmp)
			default:
			}
		}
	}
	return m, nil
}

// NewFileSource reads in the given file and parsing it with
// multiple encodings to find one that works.
func NewFileSource(path string) (*MapSource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// changed how we handled a deferred file closure due to go lint security check https://github.com/securego/gosec/issues/512#issuecomment-675286833
	defer func() {
		if cerr := f.Close(); cerr != nil {
			fmt.Printf("Error closing file: %s\n", cerr)
		}
	}()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	m, err := NewJSONSource(b)
	if err == nil {
		return m, nil
	}
	m, err = NewYAMLSource(b)
	if err == nil {
		return m, nil
	}
	return nil, fmt.Errorf("could not determine file format for %s", path)
}

// NewEnvSource uses the given environment to generate a configuration
// source. The "_" character is used as a delimeter and each one will
// result in a subtree.
func NewEnvSource(env []string) (*MapSource, error) {
	m := make(map[string]interface{})
	for _, envStr := range env {
		parts := strings.SplitN(envStr, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("failed to parse variable %s", envStr)
		}
		name := parts[0]
		value := parts[1]
		path := strings.Split(name, "_")
		location := m
		for x := 0; x < len(path)-1; x = x + 1 {
			part := path[x]
			if _, ok := location[part]; !ok {
				location[part] = make(map[string]interface{})
			}
			var nextLocation map[string]interface{}
			var ok bool
			if nextLocation, ok = location[part].(map[string]interface{}); !ok {
				// This condition indicates that there is a conflict while
				// interpreting the environment variables. Specifically, the
				// system is attempting to add an element to a subtree but one
				// of the parents in that tree is a value instead of a tree.
				// For example, A_B_C_D="value" will result in the string "value"
				// being stored in {"a": {"b": {"c": {"d": "value"}}}}. However, if
				// there was previously an entry like A_B="value2" then we would
				// have already created {"a": {"b": "value2"}}. The string "value2"
				// and the subtree {"c": {"d": "value"}} cannot exist under the
				// same parent of "b".
				//
				// This is not an issue for some of the other sources, like JSON
				// or YAML, because their structure prevents it. However, naive
				// key-value stores like ENV vars or Redis, potentially, do not
				// have an inherent hierarchy concept so this conflict can appear.
				// It's challenging to determine the best course of action here
				// because, while there are benefits to strict validation of the
				// data and returning an error for this condition, nearly every
				// environment will trigger this path. For now, we will choose to
				// ignore these conflicts. If the need arrises for a "strict mode"
				// then we can offer a separate component that returns an error.
				continue
			}
			location = nextLocation
		}
		location[path[len(path)-1]] = value
	}
	return NewMapSource(m), nil
}

// PrefixSource is a wrapper for other Source implementaions that adds
// a path element to the front of every lookup.
type PrefixSource struct {
	Source Source
	Prefix []string
}

// Get a value with a prefixed path.
func (s *PrefixSource) Get(ctx context.Context, path ...string) (interface{}, bool) {
	path = append(path, s.Prefix...)     // allocate correct size
	copy(path[len(s.Prefix):], path[0:]) // shift original elements
	copy(path, s.Prefix)                 // re-insert prefix values
	return s.Source.Get(ctx, path...)
}

// MultiSource is an ordered set of Sources from which to pull
// values. It will search until the first Source returns a found
// value or will return false for found.
type MultiSource []Source

// Get a value from the ordered set of Sources.
func (s MultiSource) Get(ctx context.Context, path ...string) (interface{}, bool) {
	for _, ss := range s {
		v, found := ss.Get(ctx, path...)
		if found {
			return v, found
		}
	}
	return nil, false
}
