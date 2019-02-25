package settings

import (
	"context"
	"fmt"
)

// Load the values for a given batch of settings using
// the provided source.
func Load(ctx context.Context, s Source, settings []Setting) error {
	for _, setting := range settings {
		v, found := s.Get(ctx, setting.Name())
		if found {
			if err := setting.SetValue(v); err != nil {
				return fmt.Errorf("failed to load setting %s due to: %s", setting.Name(), err.Error())
			}
		}
	}
	return nil
}

type groupLoad struct {
	Path  []string
	Group Group
}

// LoadGroups works similarly to Load except that it will operate recursively
// on all settings and groups in the given group. Each group name will be
// added as a path segment leading to an individual setting.
func LoadGroups(ctx context.Context, s Source, groups []Group) error {
	stack := make([]groupLoad, 0, len(groups))
	for _, group := range groups {
		stack = append(stack, groupLoad{Path: []string{group.Name()}, Group: group})
	}
	for len(stack) > 0 {
		var current groupLoad
		current, stack = stack[len(stack)-1], stack[:len(stack)-1]
		for _, group := range current.Group.Groups() {
			newPath := make([]string, 0, len(current.Path)+1)
			newPath = append(newPath, current.Path...)
			newPath = append(newPath, group.Name())
			stack = append(stack, groupLoad{Group: group, Path: newPath})
		}
		err := Load(ctx, &PrefixSource{Source: s, Prefix: current.Path}, current.Group.Settings())
		if err != nil {
			return fmt.Errorf("failed to load group %s due to: %s", current.Group.Name(), err.Error())
		}
	}
	return nil
}
