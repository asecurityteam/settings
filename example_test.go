package settings

import (
	"testing"
	"time"

	"github.com/andreyvit/diff"
)

func Test_typeHint(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		want string
	}{
		{
			name: "bool",
			v:    true,
			want: "bool",
		},
		{
			name: "int",
			v:    int(1),
			want: "int",
		},
		{
			name: "float64",
			v:    float64(1.0),
			want: "float64",
		},
		{
			name: "string",
			v:    "value",
			want: "string",
		},
		{
			name: "duration",
			v:    time.Second,
			want: durationName,
		},
		{
			name: "time",
			v:    time.Now(),
			want: timeName,
		},
		{
			name: "string slice",
			v:    []string{"value"},
			want: "[]string",
		},
		{
			name: "duration slice",
			v:    []time.Duration{time.Second},
			want: "[]time.Duration",
		},
		{
			name: "time slice",
			v:    []time.Time{time.Now()},
			want: "[]time.Time",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := typeHint(tt.v); got != tt.want {
				t.Errorf("typeHint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_yamlTypeDisplay(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		want string
	}{
		{
			name: "bool",
			v:    true,
			want: "true",
		},
		{
			name: "int",
			v:    int(1),
			want: "1",
		},
		{
			name: "float64",
			v:    float64(1.5),
			want: "1.5",
		},
		{
			name: "string",
			v:    "value",
			want: `"value"`,
		},
		{
			name: "duration",
			v:    time.Second,
			want: `"1s"`,
		},
		{
			name: "time",
			v:    time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC),
			want: `"1999-01-01 00:00:00 +0000 UTC"`,
		},
		{
			name: "string slice",
			v:    []string{"value1", "value2"},
			want: "\n  - \"value1\"\n  - \"value2\"\n",
		},
		{
			name: "duration slice",
			v:    []time.Duration{time.Second, time.Minute},
			want: "\n  - \"1s\"\n  - \"1m0s\"\n",
		},
		{
			name: "time slice",
			v: []time.Time{
				time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			want: "\n  - \"1999-01-01 00:00:00 +0000 UTC\"\n  - \"2000-01-01 00:00:00 +0000 UTC\"\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := yamlTypeDisplay(tt.v); got != tt.want {
				t.Errorf("typeDiplay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExampleYamlSettings(t *testing.T) {
	tests := []struct {
		name     string
		settings []Setting
		want     string
	}{
		{
			name:     "empty",
			settings: nil,
			want:     "",
		},
		{
			name: "collection",
			settings: []Setting{
				NewBoolSetting("enabled", "is it on?", false),
				NewTimeSetting("when", "when does it happen?", time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)),
				NewStringSliceSetting("what", "do something with these", []string{"one", "two"}),
			},
			want: `# (bool) is it on?
enabled: false
# (time.Time) when does it happen?
when: "1999-01-01 00:00:00 +0000 UTC"
# ([]string) do something with these
what:
  - "one"
  - "two"
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExampleYamlSettings(tt.settings); got != tt.want {
				t.Errorf("ExampleYamlSettings() = %v, want %v\n%s", got, tt.want, diff.CharacterDiff(got, tt.want))
			}
		})
	}
}

func TestExampleYamlSettingGroups(t *testing.T) {
	tests := []struct {
		name   string
		groups []Group
		want   string
	}{
		{
			name: "single",
			groups: []Group{
				&SettingGroup{
					NameValue: "outer",
					SettingValues: []Setting{
						NewBoolSetting("enabled", "is it on?", false),
						NewTimeSetting("when", "when does it happen?", time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)),
						NewStringSliceSetting("what", "do something with these", []string{"one", "two"}),
					},
				},
			},
			want: `outer:
  # (bool) is it on?
  enabled: false
  # (time.Time) when does it happen?
  when: "1999-01-01 00:00:00 +0000 UTC"
  # ([]string) do something with these
  what:
    - "one"
    - "two"
`,
		},
		{
			name: "multiple",
			groups: []Group{
				&SettingGroup{
					NameValue: "outer",
					SettingValues: []Setting{
						NewBoolSetting("enabled", "is it on?", false),
						NewTimeSetting("when", "when does it happen?", time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)),
						NewStringSliceSetting("what", "do something with these", []string{"one", "two"}),
					},
					GroupValues: []Group{
						&SettingGroup{
							NameValue: "inner1",
							SettingValues: []Setting{
								NewBoolSetting("enabled", "is it on?", false),
								NewTimeSetting("when", "when does it happen?", time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)),
								NewStringSliceSetting("what", "do something with these", []string{"one", "two"}),
							},
						},
						&SettingGroup{
							NameValue: "inner2",
							SettingValues: []Setting{
								NewBoolSetting("enabled", "is it on?", false),
								NewTimeSetting("when", "when does it happen?", time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)),
								NewStringSliceSetting("what", "do something with these", []string{"one", "two"}),
							},
							GroupValues: []Group{
								&SettingGroup{
									NameValue: "inner2inner",
									SettingValues: []Setting{
										NewBoolSetting("enabled", "is it on?", false),
										NewTimeSetting("when", "when does it happen?", time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)),
										NewStringSliceSetting("what", "do something with these", []string{"one", "two"}),
									},
								},
							},
						},
					},
				},
			},
			want: `outer:
  # (bool) is it on?
  enabled: false
  # (time.Time) when does it happen?
  when: "1999-01-01 00:00:00 +0000 UTC"
  # ([]string) do something with these
  what:
    - "one"
    - "two"
  inner1:
    # (bool) is it on?
    enabled: false
    # (time.Time) when does it happen?
    when: "1999-01-01 00:00:00 +0000 UTC"
    # ([]string) do something with these
    what:
      - "one"
      - "two"
  inner2:
    # (bool) is it on?
    enabled: false
    # (time.Time) when does it happen?
    when: "1999-01-01 00:00:00 +0000 UTC"
    # ([]string) do something with these
    what:
      - "one"
      - "two"
    inner2inner:
      # (bool) is it on?
      enabled: false
      # (time.Time) when does it happen?
      when: "1999-01-01 00:00:00 +0000 UTC"
      # ([]string) do something with these
      what:
        - "one"
        - "two"
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExampleYamlGroups(tt.groups); got != tt.want {
				t.Errorf("ExampleYamlGroups() = %v, want %v\n%s", got, tt.want, diff.LineDiff(got, tt.want))
			}
		})
	}
}

func Test_envTypeDisplay(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		want string
	}{
		{
			name: "bool",
			v:    true,
			want: `"true"`,
		},
		{
			name: "int",
			v:    int(1),
			want: `"1"`,
		},
		{
			name: "float64",
			v:    float64(1.5),
			want: `"1.5"`,
		},
		{
			name: "string",
			v:    "value",
			want: `"value"`,
		},
		{
			name: "duration",
			v:    time.Second,
			want: `"1s"`,
		},
		{
			name: "time",
			v:    time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC),
			want: `"1999-01-01T00:00:00Z"`,
		},
		{
			name: "string slice",
			v:    []string{"value1", "value2"},
			want: `"value1 value2"`,
		},
		{
			name: "duration slice",
			v:    []time.Duration{time.Second, time.Minute},
			want: `"1s 1m0s"`,
		},
		{
			name: "time slice",
			v: []time.Time{
				time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			want: `"1999-01-01T00:00:00Z 2000-01-01T00:00:00Z"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := envTypeDisplay(tt.v); got != tt.want {
				t.Errorf("envTypeDisplay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExampleEnvSettings(t *testing.T) {
	tests := []struct {
		name     string
		settings []Setting
		want     string
	}{
		{
			name:     "empty",
			settings: nil,
			want:     "",
		},
		{
			name: "collection",
			settings: []Setting{
				NewBoolSetting("enabled", "is it on?", false),
				NewTimeSetting("when", "when does it happen?", time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)),
				NewStringSliceSetting("what", "do something with these", []string{"one", "two"}),
			},
			want: `# (bool) is it on?
ENABLED="false"
# (time.Time) when does it happen?
WHEN="1999-01-01T00:00:00Z"
# ([]string) do something with these
WHAT="one two"
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExampleEnvSettings(tt.settings); got != tt.want {
				t.Errorf("ExampleEnvSettings() = %v, want %v\n%s", got, tt.want, diff.CharacterDiff(got, tt.want))
			}
		})
	}
}

func TestExampleEnvSettingGroup(t *testing.T) {
	tests := []struct {
		name   string
		groups []Group
		want   string
	}{
		{
			name: "single",
			groups: []Group{
				&SettingGroup{
					NameValue: "outer",
					SettingValues: []Setting{
						NewBoolSetting("enabled", "is it on?", false),
						NewTimeSetting("when", "when does it happen?", time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)),
						NewStringSliceSetting("what", "do something with these", []string{"one", "two"}),
					},
				},
			},
			want: `# (bool) is it on?
OUTER_ENABLED="false"
# (time.Time) when does it happen?
OUTER_WHEN="1999-01-01T00:00:00Z"
# ([]string) do something with these
OUTER_WHAT="one two"
`,
		},
		{
			name: "multiple",
			groups: []Group{
				&SettingGroup{
					NameValue: "outer",
					SettingValues: []Setting{
						NewBoolSetting("enabled", "is it on?", false),
						NewTimeSetting("when", "when does it happen?", time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)),
						NewStringSliceSetting("what", "do something with these", []string{"one", "two"}),
					},
					GroupValues: []Group{
						&SettingGroup{
							NameValue: "inner1",
							SettingValues: []Setting{
								NewBoolSetting("enabled", "is it on?", false),
								NewTimeSetting("when", "when does it happen?", time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)),
								NewStringSliceSetting("what", "do something with these", []string{"one", "two"}),
							},
						},
						&SettingGroup{
							NameValue: "inner2",
							SettingValues: []Setting{
								NewBoolSetting("enabled", "is it on?", false),
								NewTimeSetting("when", "when does it happen?", time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)),
								NewStringSliceSetting("what", "do something with these", []string{"one", "two"}),
							},
							GroupValues: []Group{
								&SettingGroup{
									NameValue: "inner2inner",
									SettingValues: []Setting{
										NewBoolSetting("enabled", "is it on?", false),
										NewTimeSetting("when", "when does it happen?", time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)),
										NewStringSliceSetting("what", "do something with these", []string{"one", "two"}),
									},
								},
							},
						},
					},
				},
			},
			want: `# (bool) is it on?
OUTER_ENABLED="false"
# (time.Time) when does it happen?
OUTER_WHEN="1999-01-01T00:00:00Z"
# ([]string) do something with these
OUTER_WHAT="one two"
# (bool) is it on?
OUTER_INNER2_ENABLED="false"
# (time.Time) when does it happen?
OUTER_INNER2_WHEN="1999-01-01T00:00:00Z"
# ([]string) do something with these
OUTER_INNER2_WHAT="one two"
# (bool) is it on?
OUTER_INNER2_INNER2INNER_ENABLED="false"
# (time.Time) when does it happen?
OUTER_INNER2_INNER2INNER_WHEN="1999-01-01T00:00:00Z"
# ([]string) do something with these
OUTER_INNER2_INNER2INNER_WHAT="one two"
# (bool) is it on?
OUTER_INNER1_ENABLED="false"
# (time.Time) when does it happen?
OUTER_INNER1_WHEN="1999-01-01T00:00:00Z"
# ([]string) do something with these
OUTER_INNER1_WHAT="one two"
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExampleEnvGroups(tt.groups); got != tt.want {
				t.Errorf("ExampleEnvGroups() = %v, want %v\n%s", got, tt.want, diff.LineDiff(got, tt.want))
			}
		})
	}
}
