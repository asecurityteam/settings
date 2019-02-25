package settings

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func typeHint(v interface{}) string {
	t := reflect.TypeOf(v)
	tn := t.String()
	if t.Kind() == reflect.Slice {
		tn = fmt.Sprintf(`[]%s`, t.Elem().String())
	}
	return tn
}

func yamlTypeDisplay(v interface{}) string {
	t := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)
	display := fmt.Sprintf("%v", v)
	if t.Kind() == reflect.Slice {
		b := bytes.NewBufferString("\n")
		for x := 0; x < vv.Len(); x = x + 1 {
			b.WriteString(fmt.Sprintf("  - %s\n", yamlTypeDisplay(vv.Index(x).Interface())))
		}
		return b.String()
	}
	if t.Kind() == reflect.String {
		return `"` + display + `"`
	}
	if t.String() == "time.Duration" || t.String() == "time.Time" {
		return fmt.Sprintf("\"%s\"", v)
	}
	return display
}

func removeExtraLines(s string) string {
	var b bytes.Buffer
	scn := bufio.NewScanner(strings.NewReader(s))
	for scn.Scan() {
		txt := scn.Text()
		trm := strings.TrimSpace(txt)
		if len(trm) != strings.Count(trm, "\n") {
			_, _ = b.WriteString(txt)
			_, _ = b.WriteString("\n")
		}
	}
	return b.String()
}

// ExampleYamlGroups renders a Group to YAML.
func ExampleYamlGroups(gs []Group) string {
	var b bytes.Buffer
	for _, g := range gs {
		if len(g.Settings()) > 0 || len(g.Groups()) > 0 {
			// Skip empty sections entirely.
			_, _ = b.WriteString(fmt.Sprintf("%s:\n", g.Name()))
		}
		if len(g.Settings()) > 0 {
			sets := ExampleYamlSettings(g.Settings())
			sc := bufio.NewScanner(strings.NewReader(sets))
			for sc.Scan() {
				_, _ = b.WriteString("  " + sc.Text() + "\n")
			}
		}
		if len(g.Groups()) > 0 {
			grps := ExampleYamlGroups(g.Groups())
			sc := bufio.NewScanner(strings.NewReader(grps))
			for sc.Scan() {
				_, _ = b.WriteString("  " + sc.Text() + "\n")
			}
		}
	}
	return removeExtraLines(b.String())
}

// ExampleYamlSettings renders a collection of settings as YAML text.
func ExampleYamlSettings(settings []Setting) string {
	var b bytes.Buffer
	for _, setting := range settings {
		hint := typeHint(setting.Value())
		display := yamlTypeDisplay(setting.Value())
		_, _ = b.WriteString(fmt.Sprintf("# (%s) %s\n", hint, setting.Description()))
		displayName := strings.ToLower(setting.Name())
		if display[0] == '\n' {
			// Special case for things that appear on the next line so we can
			// trim the extra spaces after the name.
			_, _ = b.WriteString(fmt.Sprintf("%s:%v\n", displayName, display))
			continue
		}
		_, _ = b.WriteString(fmt.Sprintf("%s: %v\n", displayName, display))
	}
	return removeExtraLines(b.String())
}

func envTypeDisplay(v interface{}) string {
	t := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)
	display := fmt.Sprintf(`"%v"`, v)
	if t.Kind() == reflect.Slice {
		b := bytes.NewBufferString(`"`)
		for x := 0; x < vv.Len()-1; x = x + 1 {
			d := envTypeDisplay(vv.Index(x).Interface())
			_, _ = b.WriteString(fmt.Sprintf("%s ", strings.Trim(d, `"`)))
		}
		if vv.Len() > 0 {
			d := envTypeDisplay(vv.Index(vv.Len() - 1).Interface())
			_, _ = b.WriteString(strings.Trim(d, `"`))
		}
		_, _ = b.WriteString(`"`)
		return b.String()
	}
	if t.String() == "time.Duration" {
		return fmt.Sprintf(`"%s"`, v)
	}
	if t.String() == "time.Time" {
		return `"` + vv.Interface().(time.Time).Format(time.RFC3339Nano) + `"`
	}
	return display
}

// ExampleEnvGroups renders a Group to ENV vars.
func ExampleEnvGroups(groups []Group) string {
	var b bytes.Buffer
	stack := make([]Group, len(groups))
	copy(stack, groups)
	for len(stack) > 0 {
		var current Group
		current, stack = stack[len(stack)-1], stack[:len(stack)-1]
		for _, g := range current.Groups() {
			// Make a copy of the group with a prefixed name and drop
			// it in the stack. This will progressively build out the
			// env var name prefix for nexted groups with each group
			// still rendering individual settings with the right prefix.
			cpy := &SettingGroup{
				NameValue: strings.ToUpper(current.Name() + "_" + g.Name()),
			}
			if len(g.Settings()) > 0 {
				cpy.SettingValues = make([]Setting, len(g.Settings()))
				copy(cpy.SettingValues, g.Settings())
			}
			if len(g.Groups()) > 0 {
				cpy.GroupValues = make([]Group, len(g.Groups()))
				copy(cpy.GroupValues, g.Groups())
			}
			stack = append(stack, cpy)
		}
		sets := ExampleEnvSettings(current.Settings())
		sc := bufio.NewScanner(strings.NewReader(sets))
		for sc.Scan() {
			if strings.HasPrefix(sc.Text(), "#") {
				_, _ = b.WriteString(sc.Text() + "\n")
				continue
			}
			_, _ = b.WriteString(strings.ToUpper(current.Name()) + "_" + sc.Text() + "\n")
		}
	}
	return removeExtraLines(b.String())
}

// ExampleEnvSettings renders a collection of settings as ENV vars.
func ExampleEnvSettings(settings []Setting) string {
	var b bytes.Buffer
	for _, setting := range settings {
		hint := typeHint(setting.Value())
		display := envTypeDisplay(setting.Value())
		_, _ = b.WriteString(fmt.Sprintf("# (%s) %s\n", hint, setting.Description()))
		_, _ = b.WriteString(fmt.Sprintf("%s=%s\n", strings.ToUpper(setting.Name()), display))
	}
	return removeExtraLines(b.String())
}
