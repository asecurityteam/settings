package settings

import (
	"reflect"
	"testing"
	"time"
)

type named struct{}

func (*named) Name() string {
	return "TEST"
}

type described struct {
	*named
}

func (*described) Description() string {
	return "TEST"
}

type inner struct {
	V string
}
type embedded struct {
	*inner
}
type nested struct {
	V *inner
}

func TestConvert(t *testing.T) {
	dur := time.Minute
	tests := []struct {
		name    string
		v       interface{}
		want    Group
		wantErr bool
	}{
		{
			name:    "non-struct/nil",
			v:       nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "non-struct/value",
			v:       "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "struct/non-addressable",
			v:       &(struct{ V struct{} }{V: struct{}{}}),
			want:    nil,
			wantErr: true,
		},

		{
			name: "struct/bool",
			v:    &(struct{ V bool }{V: true}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewBoolSetting("V", "", true),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/int",
			v:    &(struct{ V int }{V: 1}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewIntSetting("V", "", 1),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/int8",
			v:    &(struct{ V int8 }{V: 1}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewInt8Setting("V", "", 1),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/int16",
			v:    &(struct{ V int16 }{V: 1}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewInt16Setting("V", "", 1),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/int32",
			v:    &(struct{ V int32 }{V: 1}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewInt32Setting("V", "", 1),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/int64",
			v:    &(struct{ V int64 }{V: 1}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewInt64Setting("V", "", 1),
				},
			},
			wantErr: false,
		},

		{
			name: "struct/uint",
			v:    &(struct{ V uint }{V: 1}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewUintSetting("V", "", 1),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/uint8",
			v:    &(struct{ V uint8 }{V: 1}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewUint8Setting("V", "", 1),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/uint16",
			v:    &(struct{ V uint16 }{V: 1}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewUint16Setting("V", "", 1),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/uint32",
			v:    &(struct{ V uint32 }{V: 1}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewUint32Setting("V", "", 1),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/uint64",
			v:    &(struct{ V uint64 }{V: 1}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewUint64Setting("V", "", 1),
				},
			},
			wantErr: false,
		},

		{
			name: "struct/float32",
			v:    &(struct{ V float32 }{V: 1.0}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewFloat32Setting("V", "", 1.0),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/float64",
			v:    &(struct{ V float64 }{V: 1.0}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewFloat64Setting("V", "", 1.0),
				},
			},
			wantErr: false,
		},

		{
			name: "struct/string",
			v:    &(struct{ V string }{V: "a"}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewStringSetting("V", "", "a"),
				},
			},
			wantErr: false,
		},

		{
			name: "struct/duration",
			v:    &(struct{ V *time.Duration }{V: &dur}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewDurationSetting("V", "", time.Minute),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/time",
			v:    &(struct{ V time.Time }{V: time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewTimeSetting("V", "", time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
			wantErr: false,
		},

		{
			name: "struct/[]duration",
			v:    &(struct{ V []time.Duration }{V: []time.Duration{time.Minute}}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewDurationSliceSetting("V", "", []time.Duration{time.Minute}),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/[]string",
			v:    &(struct{ V []string }{V: []string{"a"}}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewStringSliceSetting("V", "", []string{"a"}),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/[]int",
			v:    &(struct{ V []int }{V: []int{1}}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewIntSliceSetting("V", "", []int{1}),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/map[string][]string",
			v: &(struct{ V map[string][]string }{
				V: map[string][]string{"letters": {"a", "b"}, "characters": {"!", "@"}}}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewStringMapStringSliceSetting("V", "",
						map[string][]string{"letters": {"a", "b"}, "characters": {"!", "@"}}),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/map[string]string",
			v: &(struct{ V map[string]string }{
				V: map[string]string{"letter": "a", "character": "!"}}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewStringMapStringSetting("V", "",
						map[string]string{"letter": "a", "character": "!"}),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/named",
			v:    &named{},
			want: &SettingGroup{
				NameValue: "TEST",
			},
			wantErr: false,
		},
		{
			name: "struct/described",
			v:    &described{&named{}},
			want: &SettingGroup{
				NameValue:        "TEST",
				DescriptionValue: "TEST",
			},
			wantErr: false,
		},

		{
			name: "struct/annotated",
			v: &(struct {
				V string `description:"a string"`
			}{V: "a"}),
			want: &SettingGroup{
				SettingValues: []Setting{
					NewStringSetting("V", "a string", "a"),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/embedded",
			v:    &embedded{&inner{V: "a"}},
			want: &SettingGroup{
				NameValue: "embedded",
				SettingValues: []Setting{
					NewStringSetting("V", "", "a"),
				},
			},
			wantErr: false,
		},
		{
			name: "struct/nested",
			v: &nested{
				V: &inner{
					V: "a",
				},
			},
			want: &SettingGroup{
				NameValue: "nested",
				GroupValues: []Group{
					&SettingGroup{
						NameValue: "inner",
						SettingValues: []Setting{
							NewStringSetting("V", "", "a"),
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}
