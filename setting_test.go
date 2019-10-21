package settings

import (
	"reflect"
	"testing"
	"time"
)

func getParsedYamlTestFixtureForMapStringToSlice(t *testing.T) map[string]interface{} {
	yamlTestString := `
fruits:
  - apple
  - orange
  - mango
vegetables:
  - corn
  - squash
`
	testS, err := NewYAMLSource([]byte(yamlTestString))
	if err != nil {
		t.Errorf("failed to parse YAML test fixture due to %s", err)
	}
	return testS.Map
}

func getParsedYamlTestFixtureForMapStringToString(t *testing.T) map[string]interface{} {
	yamlTestString := `
fruit: apple
vegetable: corn
`
	testS, err := NewYAMLSource([]byte(yamlTestString))
	if err != nil {
		t.Errorf("failed to parse YAML test fixture due to %s", err)
	}
	return testS.Map
}

func TestSetting(t *testing.T) {
	tests := []struct {
		name     string
		setting  Setting
		good     interface{}
		expected interface{}
		bad      interface{}
	}{
		{
			name:     "String",
			setting:  NewStringSetting("String", "", ""),
			good:     "anything",
			expected: "anything",
			bad:      "",
		},
		{
			name:     "Bool",
			setting:  NewBoolSetting("Bool", "", false),
			good:     "true",
			expected: true,
			bad:      "3",
		},
		{
			name:     "Int",
			setting:  NewIntSetting("Int", "", 0),
			good:     "1",
			expected: 1,
			bad:      "false",
		},
		{
			name:     "Int8",
			setting:  NewInt8Setting("Int8", "", 0),
			good:     "1",
			expected: int8(1),
			bad:      "false",
		},
		{
			name:     "Int16",
			setting:  NewInt16Setting("Int16", "", 0),
			good:     "1",
			expected: int16(1),
			bad:      "false",
		},
		{
			name:     "Int32",
			setting:  NewInt32Setting("Int32", "", 0),
			good:     "1",
			expected: int32(1),
			bad:      "false",
		},
		{
			name:     "Int64",
			setting:  NewInt64Setting("Int64", "", 0),
			good:     "1",
			expected: int64(1),
			bad:      "false",
		},
		{
			name:     "Uint",
			setting:  NewUintSetting("Uint", "", 0),
			good:     "1",
			expected: uint(1),
			bad:      "false",
		},
		{
			name:     "Uint8",
			setting:  NewUint8Setting("Uint8", "", 0),
			good:     "1",
			expected: uint8(1),
			bad:      "false",
		},
		{
			name:     "Uint16",
			setting:  NewUint16Setting("Uint16", "", 0),
			good:     "1",
			expected: uint16(1),
			bad:      "false",
		},
		{
			name:     "Uint32",
			setting:  NewUint32Setting("Uint32", "", 0),
			good:     "1",
			expected: uint32(1),
			bad:      "false",
		},
		{
			name:     "Uint64",
			setting:  NewUint64Setting("Uint64", "", 0),
			good:     "1",
			expected: uint64(1),
			bad:      "false",
		},
		{
			name:     "Float32",
			setting:  NewFloat32Setting("Float32", "", 0),
			good:     "1.0",
			expected: float32(1.0),
			bad:      "false",
		},
		{
			name:     "Float64",
			setting:  NewFloat64Setting("Float64", "", 0),
			good:     "1.0",
			expected: float64(1.0),
			bad:      "false",
		},
		{
			name:     "Time",
			setting:  NewTimeSetting("Time", "", time.Now()),
			good:     time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339Nano),
			expected: time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC),
			bad:      "false",
		},
		{
			name:     "Duration",
			setting:  NewDurationSetting("Duration", "", 0),
			good:     "1s",
			expected: time.Second,
			bad:      "false",
		},
		{
			name:     "BoolSlice",
			setting:  NewBoolSliceSetting("BoolSlice", "", nil),
			good:     "true true false",
			expected: []bool{true, true, false},
			bad:      "3",
		},
		{
			name:     "DurationSlice",
			setting:  NewDurationSliceSetting("DurationSlice", "", nil),
			good:     "1s 1m 1h",
			expected: []time.Duration{time.Second, time.Minute, time.Hour},
			bad:      "false",
		},
		{
			name:     "IntSlice",
			setting:  NewIntSliceSetting("IntSlice", "", nil),
			good:     "1 2 3",
			expected: []int{1, 2, 3},
			bad:      "false",
		},
		{
			name:     "StringSlice",
			setting:  NewStringSliceSetting("StringSlice", "", nil),
			good:     "one two three",
			expected: []string{"one", "two", "three"},
			bad:      make(map[string]interface{}),
		},
		{
			name:     "StringMapStringSlice from JSON",
			setting:  NewStringMapStringSliceSetting("StringMapStringSlice", "", nil),
			good:     `{"dogs": ["german shepard", "golden retriever"], "birds": ["eagle", "pigeon"]}`,
			expected: map[string][]string{"dogs": {"german shepard", "golden retriever"}, "birds": {"eagle", "pigeon"}},
			bad:      `{"animal": "dog"}`,
		},
		{
			name:     "StringMapStringSlice from YAML",
			setting:  NewStringMapStringSliceSetting("StringMapStringSlice", "", nil),
			good:     getParsedYamlTestFixtureForMapStringToSlice(t),
			expected: map[string][]string{"fruits": {"apple", "orange", "mango"}, "vegetables": {"corn", "squash"}},
			bad:      `- animal - dog`,
		},
		{
			name:     "StringMapString from JSON",
			setting:  NewStringMapStringSetting("StringMapString", "", nil),
			good:     `{"dog": "german shepard", "bird": "eagle"}`,
			expected: map[string]string{"dog": "german shepard", "bird": "eagle"},
			bad:      `{"animal": ["dog"]}`,
		},
		{
			name:     "StringMapString from YAML",
			setting:  NewStringMapStringSetting("StringMapString", "", nil),
			good:     getParsedYamlTestFixtureForMapStringToString(t),
			expected: map[string]string{"fruit": "apple", "vegetable": "corn"},
			bad:      `- animal - dog`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setting.SetValue(tt.good)
			if err != nil {
				t.Errorf("%s: %s", tt.name, err.Error())
			}
			if !reflect.DeepEqual(tt.setting.Value(), tt.expected) {
				t.Errorf("%s.Value() got = %v, want %v", tt.name, tt.setting.Value(), tt.expected)
			}
			if tt.bad != "" {
				err = tt.setting.SetValue(tt.bad)
				if err == nil {
					t.Errorf("%s: accepted bad input as good", tt.name)
				}
			}
		})
	}
}
