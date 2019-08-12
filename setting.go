package settings

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
)

// Setting is a generic container for all configuration values.
// Each implementation should maintain its own internal typing.
type Setting interface {
	Name() string
	Description() string
	Value() interface{}
	SetValue(v interface{}) error
}

// Group is a container for a collection of settings. The
// container can contain any number of nested sub-trees.
type Group interface {
	Name() string
	Description() string
	Groups() []Group
	Settings() []Setting
}

// SettingGroup an implementation of Group. This component is
// predominantly used by the struct converter to map native
// types into Setting and Group types.
type SettingGroup struct {
	NameValue        string
	DescriptionValue string
	GroupValues      []Group
	SettingValues    []Setting
}

// Name returns the group name as it appears in the configuration.
func (g *SettingGroup) Name() string {
	return g.NameValue
}

// Description returns the group description.
func (g *SettingGroup) Description() string {
	return g.DescriptionValue
}

// Groups returns any sub-trees in the current group.
func (g *SettingGroup) Groups() []Group {
	return g.GroupValues
}

// Settings returns any settings attached directly to this group.
func (g *SettingGroup) Settings() []Setting {
	return g.SettingValues
}

// BaseSetting implements the name and description aspects of
// any given setting.
type BaseSetting struct {
	NameValue        string
	DescriptionValue string
}

// Name returns the setting name as it appears in configuration.
func (s *BaseSetting) Name() string {
	return s.NameValue
}

// Description returns the human description of a setting for help text.
func (s *BaseSetting) Description() string {
	return s.DescriptionValue
}

// StringSetting manages an instance of string
type StringSetting struct {
	*BaseSetting
	StringValue *string
}

// NewStringSetting creates a StringSetting with the given default value.
func NewStringSetting(name string, description string, fallback string) *StringSetting {
	return &StringSetting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		StringValue: &fallback,
	}
}

// Value returns the underlying string.
func (s *StringSetting) Value() interface{} {
	return *s.StringValue
}

// SetValue changes the underlying string.
func (s *StringSetting) SetValue(v interface{}) error {
	var err error
	*s.StringValue, err = cast.ToStringE(v)
	return err
}

// BoolSetting manages an instance of bool
type BoolSetting struct {
	*BaseSetting
	BoolValue *bool
}

// NewBoolSetting creates a BoolSetting with the given default value.
func NewBoolSetting(name string, description string, fallback bool) *BoolSetting {
	return &BoolSetting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		BoolValue: &fallback,
	}
}

// Value returns the underlying bool.
func (s *BoolSetting) Value() interface{} {
	return *s.BoolValue
}

// SetValue changes the underlying bool.
func (s *BoolSetting) SetValue(v interface{}) error {
	var err error
	*s.BoolValue, err = cast.ToBoolE(v)
	return err
}

// IntSetting manages an instance of int
type IntSetting struct {
	*BaseSetting
	IntValue *int
}

// NewIntSetting creates a IntSetting with the given default value.
func NewIntSetting(name string, description string, fallback int) *IntSetting {
	return &IntSetting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		IntValue: &fallback,
	}
}

// Value returns the underlying int.
func (s *IntSetting) Value() interface{} {
	return *s.IntValue
}

// SetValue changes the underlying int.
func (s *IntSetting) SetValue(v interface{}) error {
	var err error
	*s.IntValue, err = cast.ToIntE(v)
	return err
}

// Int8Setting manages an instance of int8
type Int8Setting struct {
	*BaseSetting
	Int8Value *int8
}

// NewInt8Setting creates a Int8Setting with the given default value.
func NewInt8Setting(name string, description string, fallback int8) *Int8Setting {
	return &Int8Setting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		Int8Value: &fallback,
	}
}

// Value returns the underlying int8.
func (s *Int8Setting) Value() interface{} {
	return *s.Int8Value
}

// SetValue changes the underlying int8.
func (s *Int8Setting) SetValue(v interface{}) error {
	var err error
	*s.Int8Value, err = cast.ToInt8E(v)
	return err
}

// Int16Setting manages an instance of int16
type Int16Setting struct {
	*BaseSetting
	Int16Value *int16
}

// NewInt16Setting creates a Int16Setting with the given default value.
func NewInt16Setting(name string, description string, fallback int16) *Int16Setting {
	return &Int16Setting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		Int16Value: &fallback,
	}
}

// Value returns the underlying int16.
func (s *Int16Setting) Value() interface{} {
	return *s.Int16Value
}

// SetValue changes the underlying int16.
func (s *Int16Setting) SetValue(v interface{}) error {
	var err error
	*s.Int16Value, err = cast.ToInt16E(v)
	return err
}

// Int32Setting manages an instance of int32
type Int32Setting struct {
	*BaseSetting
	Int32Value *int32
}

// NewInt32Setting creates a Int32Setting with the given default value.
func NewInt32Setting(name string, description string, fallback int32) *Int32Setting {
	return &Int32Setting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		Int32Value: &fallback,
	}
}

// Value returns the underlying int32.
func (s *Int32Setting) Value() interface{} {
	return *s.Int32Value
}

// SetValue changes the underlying int32.
func (s *Int32Setting) SetValue(v interface{}) error {
	var err error
	*s.Int32Value, err = cast.ToInt32E(v)
	return err
}

// Int64Setting manages an instance of int64
type Int64Setting struct {
	*BaseSetting
	Int64Value *int64
}

// NewInt64Setting creates a Int64Setting with the given default value.
func NewInt64Setting(name string, description string, fallback int64) *Int64Setting {
	return &Int64Setting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		Int64Value: &fallback,
	}
}

// Value returns the underlying int64.
func (s *Int64Setting) Value() interface{} {
	return *s.Int64Value
}

// SetValue changes the underlying int64.
func (s *Int64Setting) SetValue(v interface{}) error {
	var err error
	*s.Int64Value, err = cast.ToInt64E(v)
	return err
}

// UintSetting manages an instance of uint
type UintSetting struct {
	*BaseSetting
	UintValue *uint
}

// NewUintSetting creates a UintSetting with the given default value.
func NewUintSetting(name string, description string, fallback uint) *UintSetting {
	return &UintSetting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		UintValue: &fallback,
	}
}

// Value returns the underlying uint.
func (s *UintSetting) Value() interface{} {
	return *s.UintValue
}

// SetValue changes the underlying uint.
func (s *UintSetting) SetValue(v interface{}) error {
	var err error
	*s.UintValue, err = cast.ToUintE(v)
	return err
}

// Uint8Setting manages an instance of uint8
type Uint8Setting struct {
	*BaseSetting
	Uint8Value *uint8
}

// NewUint8Setting creates a Uint8Setting with the given default value.
func NewUint8Setting(name string, description string, fallback uint8) *Uint8Setting {
	return &Uint8Setting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		Uint8Value: &fallback,
	}
}

// Value returns the underlying uint8.
func (s *Uint8Setting) Value() interface{} {
	return *s.Uint8Value
}

// SetValue changes the underlying uint8.
func (s *Uint8Setting) SetValue(v interface{}) error {
	var err error
	*s.Uint8Value, err = cast.ToUint8E(v)
	return err
}

// Uint16Setting manages an instance of uint16
type Uint16Setting struct {
	*BaseSetting
	Uint16Value *uint16
}

// NewUint16Setting creates a Uint16Setting with the given default value.
func NewUint16Setting(name string, description string, fallback uint16) *Uint16Setting {
	return &Uint16Setting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		Uint16Value: &fallback,
	}
}

// Value returns the underlying uint16.
func (s *Uint16Setting) Value() interface{} {
	return *s.Uint16Value
}

// SetValue changes the underlying uint16.
func (s *Uint16Setting) SetValue(v interface{}) error {
	var err error
	*s.Uint16Value, err = cast.ToUint16E(v)
	return err
}

// Uint32Setting manages an instance of uint32
type Uint32Setting struct {
	*BaseSetting
	Uint32Value *uint32
}

// NewUint32Setting creates a Uint32Setting with the given default value.
func NewUint32Setting(name string, description string, fallback uint32) *Uint32Setting {
	return &Uint32Setting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		Uint32Value: &fallback,
	}
}

// Value returns the underlying uint32.
func (s *Uint32Setting) Value() interface{} {
	return *s.Uint32Value
}

// SetValue changes the underlying uint32.
func (s *Uint32Setting) SetValue(v interface{}) error {
	var err error
	*s.Uint32Value, err = cast.ToUint32E(v)
	return err
}

// Uint64Setting manages an instance of uint64
type Uint64Setting struct {
	*BaseSetting
	Uint64Value *uint64
}

// NewUint64Setting creates a Uint64Setting with the given default value.
func NewUint64Setting(name string, description string, fallback uint64) *Uint64Setting {
	return &Uint64Setting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		Uint64Value: &fallback,
	}
}

// Value returns the underlying uint64.
func (s *Uint64Setting) Value() interface{} {
	return *s.Uint64Value
}

// SetValue changes the underlying uint64.
func (s *Uint64Setting) SetValue(v interface{}) error {
	var err error
	*s.Uint64Value, err = cast.ToUint64E(v)
	return err
}

// Float32Setting manages an instance of float32
type Float32Setting struct {
	*BaseSetting
	Float32Value *float32
}

// NewFloat32Setting creates a Float32Setting with the given default value.
func NewFloat32Setting(name string, description string, fallback float32) *Float32Setting {
	return &Float32Setting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		Float32Value: &fallback,
	}
}

// Value returns the underlying float32.
func (s *Float32Setting) Value() interface{} {
	return *s.Float32Value
}

// SetValue changes the underlying float32.
func (s *Float32Setting) SetValue(v interface{}) error {
	var err error
	*s.Float32Value, err = cast.ToFloat32E(v)
	return err
}

// Float64Setting manages an instance of float64
type Float64Setting struct {
	*BaseSetting
	Float64Value *float64
}

// NewFloat64Setting creates a Float64Setting with the given default value.
func NewFloat64Setting(name string, description string, fallback float64) *Float64Setting {
	return &Float64Setting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		Float64Value: &fallback,
	}
}

// Value returns the underlying float64.
func (s *Float64Setting) Value() interface{} {
	return *s.Float64Value
}

// SetValue changes the underlying float64.
func (s *Float64Setting) SetValue(v interface{}) error {
	var err error
	*s.Float64Value, err = cast.ToFloat64E(v)
	return err
}

// TimeSetting manages an instance of time.Time
type TimeSetting struct {
	*BaseSetting
	TimeValue *time.Time
}

// NewTimeSetting creates a TimeSetting with the given default value.
func NewTimeSetting(name string, description string, fallback time.Time) *TimeSetting {
	return &TimeSetting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		TimeValue: &fallback,
	}
}

// Value returns the underlying time.Time.
func (s *TimeSetting) Value() interface{} {
	return *s.TimeValue
}

// SetValue changes the underlying time.Time.
func (s *TimeSetting) SetValue(v interface{}) error {
	var err error
	*s.TimeValue, err = cast.ToTimeE(v)
	return err
}

// DurationSetting manages an instance of time.Duration
type DurationSetting struct {
	*BaseSetting
	DurationValue *time.Duration
}

// NewDurationSetting creates a DurationSetting with the given default value.
func NewDurationSetting(name string, description string, fallback time.Duration) *DurationSetting {
	return &DurationSetting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		DurationValue: &fallback,
	}
}

// Value returns the underlying time.Duration.
func (s *DurationSetting) Value() interface{} {
	return *s.DurationValue
}

// SetValue changes the underlying time.Duration.
func (s *DurationSetting) SetValue(v interface{}) error {
	var err error
	*s.DurationValue, err = cast.ToDurationE(v)
	return err
}

// BoolSliceSetting manages an instance of []bool
type BoolSliceSetting struct {
	*BaseSetting
	BoolSliceValue *[]bool
}

// NewBoolSliceSetting creates a BoolSliceSetting with the given default value.
func NewBoolSliceSetting(name string, description string, fallback []bool) *BoolSliceSetting {
	return &BoolSliceSetting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		BoolSliceValue: &fallback,
	}
}

// Value returns the underlying []bool.
func (s *BoolSliceSetting) Value() interface{} {
	return *s.BoolSliceValue
}

// SetValue changes the underlying []bool.
func (s *BoolSliceSetting) SetValue(v interface{}) error {
	var err error
	var tmp []string
	tmp, err = cast.ToStringSliceE(v)
	if err != nil {
		return fmt.Errorf("bool slice parsing failed at interim string step: %s", err.Error())
	}
	*s.BoolSliceValue, err = cast.ToBoolSliceE(tmp)
	return err
}

// DurationSliceSetting manages an instance of []time.Duration
type DurationSliceSetting struct {
	*BaseSetting
	DurationSliceValue *[]time.Duration
}

// NewDurationSliceSetting creates a DurationSliceSetting with the given default value.
func NewDurationSliceSetting(name string, description string, fallback []time.Duration) *DurationSliceSetting {
	return &DurationSliceSetting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		DurationSliceValue: &fallback,
	}
}

// Value returns the underlying []time.Duration.
func (s *DurationSliceSetting) Value() interface{} {
	return *s.DurationSliceValue
}

// SetValue changes the underlying []time.Duration.
func (s *DurationSliceSetting) SetValue(v interface{}) error {
	var err error
	var tmp []string
	tmp, err = cast.ToStringSliceE(v)
	if err != nil {
		return fmt.Errorf("duration slice parsing failed at interim string step: %s", err.Error())
	}
	*s.DurationSliceValue, err = cast.ToDurationSliceE(tmp)
	return err
}

// IntSliceSetting manages an instance of []int
type IntSliceSetting struct {
	*BaseSetting
	IntSliceValue *[]int
}

// NewIntSliceSetting creates a IntSliceSetting with the given default value.
func NewIntSliceSetting(name string, description string, fallback []int) *IntSliceSetting {
	return &IntSliceSetting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		IntSliceValue: &fallback,
	}
}

// Value returns the underlying []int.
func (s *IntSliceSetting) Value() interface{} {
	return *s.IntSliceValue
}

// SetValue changes the underlying []int.
func (s *IntSliceSetting) SetValue(v interface{}) error {
	var err error
	var tmp []string
	tmp, err = cast.ToStringSliceE(v)
	if err != nil {
		return fmt.Errorf("int slice parsing failed at interim string step: %s", err.Error())
	}
	*s.IntSliceValue, err = cast.ToIntSliceE(tmp)
	return err
}

// StringSliceSetting manages an instance of []string
type StringSliceSetting struct {
	*BaseSetting
	StringSliceValue *[]string
}

// NewStringSliceSetting creates a StringSliceSetting with the given default value.
func NewStringSliceSetting(name string, description string, fallback []string) *StringSliceSetting {
	return &StringSliceSetting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		StringSliceValue: &fallback,
	}
}

// Value returns the underlying []string.
func (s *StringSliceSetting) Value() interface{} {
	return *s.StringSliceValue
}

// SetValue changes the underlying []string.
func (s *StringSliceSetting) SetValue(v interface{}) error {
	var err error
	*s.StringSliceValue, err = cast.ToStringSliceE(v)
	return err
}

type StringMapStringSliceSetting struct {
	*BaseSetting
	StringMapStringSliceValue *map[string][]string
}

func (m *StringMapStringSliceSetting) Value() interface{} {
	return *m.StringMapStringSliceValue
}

func (m *StringMapStringSliceSetting) SetValue(v interface{}) error {
	var err error
	fmt.Println("####INCOMING VALUE######## : ", v)
	*m.StringMapStringSliceValue, err = cast.ToStringMapStringSliceE(v)
	return err
}
func NewStringMapStringSliceSetting(name string, description string, fallback map[string][]string) *StringMapStringSliceSetting {
	return &StringMapStringSliceSetting{
		BaseSetting: &BaseSetting{
			NameValue:        name,
			DescriptionValue: description,
		},
		StringMapStringSliceValue: &fallback,
	}
}
