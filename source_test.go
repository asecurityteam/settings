package settings

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func Test_lowerCaseMap(t *testing.T) {
	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "empty",
			args: args{
				m: make(map[string]interface{}),
			},
			want: make(map[string]interface{}),
		},
		{
			name: "flat",
			args: args{
				m: map[string]interface{}{
					"ONE": nil,
					"two": nil,
				},
			},
			want: map[string]interface{}{
				"one": nil,
				"two": nil,
			},
		},
		{
			name: "nested",
			args: args{
				m: map[string]interface{}{
					"ONE": nil,
					"two": nil,
					"three": map[string]interface{}{
						"FOUR": nil,
						"Five": map[string]interface{}{
							"SIX": nil,
						},
					},
				},
			},
			want: map[string]interface{}{
				"one": nil,
				"two": nil,
				"three": map[string]interface{}{
					"four": nil,
					"five": map[string]interface{}{
						"six": nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lowerCaseMap(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lowerCaseMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapSource_Get(t *testing.T) {
	type fields struct {
		Map  map[string]interface{}
		JSON string
		ENV  []string
		YAML string
	}
	type args struct {
		path []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		want1  bool
	}{
		{
			name: "empty",
			fields: fields{
				Map:  map[string]interface{}{},
				JSON: `{}`,
				ENV:  []string{},
				YAML: `{}`,
			},
			args:  args{path: []string{"a", "b", "c"}},
			want:  nil,
			want1: false,
		},
		{
			name: "shallow",
			fields: fields{
				Map: map[string]interface{}{
					"a": "value",
					"b": map[string]interface{}{
						"bb": "value",
					},
					"c": map[string]interface{}{
						"cc": map[string]interface{}{
							"ccc": "value",
						},
					},
				},
				JSON: `{"a": "value", "b": {"bb": "value"}, "c": {"cc": {"ccc": "value"}}}`,
				ENV: []string{
					"A=value",
					"B_BB=value",
					"C_CC_CCC=value",
				},
				YAML: `
a: "value"
b:
  bb: "value"
c:
  cc:
    ccc: "value"
`,
			},
			args:  args{path: []string{"a"}},
			want:  "value",
			want1: true,
		},
		{
			name: "deep",
			fields: fields{
				Map: map[string]interface{}{
					"a": "value",
					"b": map[string]interface{}{
						"bb": "value",
					},
					"c": map[string]interface{}{
						"cc": map[string]interface{}{
							"ccc": "value",
						},
					},
				},
				JSON: `{"a": "value", "b": {"bb": "value"}, "c": {"cc": {"ccc": "value"}}}`,
				ENV: []string{
					"A=value",
					"B_BB=value",
					"C_CC_CCC=value",
				},
				YAML: `
a: "value"
b:
  bb: "value"
c:
  cc:
    ccc: "value"
`,
			},
			args:  args{path: []string{"c", "cc", "ccc"}},
			want:  "value",
			want1: true,
		},
		{
			name: "deep but missing",
			fields: fields{
				Map: map[string]interface{}{
					"a": "value",
					"b": map[string]interface{}{
						"bb": "value",
					},
					"c": map[string]interface{}{
						"cc": map[string]interface{}{
							"ccc": "value",
						},
					},
				},
				JSON: `{"a": "value", "b": {"bb": "value"}, "c": {"cc": {"ccc": "value"}}}`,
				ENV: []string{
					"A=value",
					"B_BB=value",
					"C_CC_CCC=value",
				},
				YAML: `
a: "value"
b:
  bb: "value"
c:
  cc:
    ccc: "value"
`,
			},
			args:  args{path: []string{"c", "ccd", "ccc"}},
			want:  nil,
			want1: false,
		},
		{
			name: "search of non-map",
			fields: fields{
				Map: map[string]interface{}{
					"a": "value",
					"b": map[string]interface{}{
						"bb": "value",
					},
					"c": map[string]interface{}{
						"cc": map[string]interface{}{
							"ccc": "value",
						},
					},
				},
				JSON: `{"a": "value", "b": {"bb": "value"}, "c": {"cc": {"ccc": "value"}}}`,
				ENV: []string{
					"A=value",
					"B_BB=value",
					"C_CC_CCC=value",
				},
				YAML: `
a: "value"
b:
  bb: "value"
c:
  cc:
    ccc: "value"
`,
			},
			args:  args{path: []string{"b", "bb", "bbb"}},
			want:  nil,
			want1: false,
		},
		{
			name: "caps",
			fields: fields{
				Map: map[string]interface{}{
					"a": "value",
					"b": map[string]interface{}{
						"bb": "value",
					},
					"c": map[string]interface{}{
						"cc": map[string]interface{}{
							"ccc": "value",
						},
					},
				},
				JSON: `{"a": "value", "b": {"bb": "value"}, "c": {"cc": {"ccc": "value"}}}`,
				ENV: []string{
					"A=value",
					"B_BB=value",
					"C_CC_CCC=value",
				},
				YAML: `
a: "value"
b:
  bb: "value"
c:
  cc:
    ccc: "value"
`,
			},
			args:  args{path: []string{"A"}},
			want:  "value",
			want1: true,
		},
		{
			name: "maps",
			fields: fields{
				Map: map[string]interface{}{
					"a": "value",
					"b": map[string]interface{}{
						"bb": "value",
					},
					"c": map[string]interface{}{
						"cc": map[string]interface{}{
							"ccc": "value",
						},
					},
				},
				JSON: `{"a": "value", "b": {"bb": "value"}, "c": {"cc": {"ccc": "value"}}}`,
				ENV: []string{
					"A=value",
					"B_BB=value",
					"C_CC_CCC=value",
				},
				YAML: `
a: "value"
b:
  bb: "value"
c:
  cc:
    ccc: "value"
`,
			},
			args: args{path: []string{"b"}},
			want: map[string]interface{}{
				"bb": "value",
			},
			want1: true,
		},
	}
	for _, tt := range tests {
		tFn := func(s *MapSource) func(t *testing.T) {
			return func(t *testing.T) {
				got, got1 := s.Get(context.Background(), tt.args.path...)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("MapSource.Get() got = %v, want %v", got, tt.want)
				}
				if got1 != tt.want1 {
					t.Errorf("MapSource.Get() got1 = %v, want %v", got1, tt.want1)
				}
			}
		}
		var s *MapSource
		var err error

		s = NewMapSource(tt.fields.Map)
		t.Run(fmt.Sprintf("%s: %s", "MAP", tt.name), tFn(s))

		s, err = NewJSONSource([]byte(tt.fields.JSON))
		if err != nil {
			t.Error(err.Error())
		}
		t.Run(fmt.Sprintf("%s: %s", "JSON", tt.name), tFn(s))
		s, err = NewEnvSource(tt.fields.ENV)
		if err != nil {
			t.Error(err.Error())
		}
		t.Run(fmt.Sprintf("%s: %s", "ENV", tt.name), tFn(s))
		s, err = NewYAMLSource([]byte(tt.fields.YAML))
		if err != nil {
			t.Error(err.Error())
		}
		t.Run(fmt.Sprintf("%s: %s", "YAML", tt.name), tFn(s))
	}
}

func TestPrefixSource(t *testing.T) {
	s := NewMapSource(map[string]interface{}{
		"a": map[string]interface{}{
			"aa": true,
		},
		"b": map[string]interface{}{
			"bb": map[string]interface{}{
				"bbb": true,
			},
		},
	})
	p := &PrefixSource{
		Source: s,
		Prefix: []string{"a"},
	}
	_, found := p.Get(context.Background(), "aa")
	if !found {
		t.Error("could not find a.aa with prefix source")
	}
	p = &PrefixSource{
		Source: &PrefixSource{
			Source: s,
			Prefix: []string{"b"},
		},
		Prefix: []string{"bb"},
	}
	_, found = p.Get(context.Background(), "bbb")
	if !found {
		t.Error("could not find b.bb.bbb with prefix source")
	}
	p = &PrefixSource{
		Source: s,
		Prefix: []string{"b", "bb"},
	}
	_, found = p.Get(context.Background(), "bbb")
	if !found {
		t.Error("could not find b.bb.bbb with prefix source")
	}
}

func TestMultiSource(t *testing.T) {
	s := MultiSource{
		NewMapSource(map[string]interface{}{}),
		NewMapSource(map[string]interface{}{}),
		NewMapSource(map[string]interface{}{
			"a": map[string]interface{}{
				"aa": true,
			},
		}),
	}
	_, found := s.Get(context.Background(), "a", "aa")
	if !found {
		t.Error("could not find a.aa with multi source")
	}
	s = MultiSource{}
	_, found = s.Get(context.Background(), "a", "aa")
	if found {
		t.Error("multi source claimed found when empty")
	}
}
