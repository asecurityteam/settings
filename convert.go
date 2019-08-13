package settings

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	timeName     = "time.Time"
	durationName = "time.Duration"
)

type namer interface {
	Name() string
}
type describer interface {
	Description() string
}

// fieledAndValue is used in the converter to bundle related information
// about a field from a struct.
type fieldAndValue struct {
	Field reflect.StructField
	Value reflect.Value
}

// Convert a struct into a Group. This function recurses over all nested
// structs which are gathered as sub-trees.
func Convert(v interface{}) (Group, error) {
	vv := reflect.Indirect(reflect.ValueOf(v))
	if v == nil {
		return nil, errors.New("nil value given to Convert")
	}
	if vv.Type().Kind() != reflect.Struct {
		return nil, fmt.Errorf("non-struct value %s given to Convert", vv.Type().String())
	}
	if !vv.CanAddr() {
		return nil, fmt.Errorf("unaddressable value %s given to Convert", vv.Type().String())
	}
	// Every struct may, optionally, provide the namer interface in order
	// to control the name of the group. If a struct does not expose the
	// name interface then we use the name of the struct with the package
	// name removed as an identfier. The same is true for description except
	// that we do not add a default description.
	nameTmp := strings.Split(vv.Type().Name(), ".")
	name := nameTmp[len(nameTmp)-1]
	if nr, ok := vv.Addr().Interface().(namer); ok {
		name = nr.Name()
	}
	desc := ""
	if nd, ok := vv.Addr().Interface().(describer); ok {
		desc = nd.Description()
	}
	g := &SettingGroup{
		NameValue:        name,
		DescriptionValue: desc,
	}
	// Now we process all of the fields in the struct. We're using a stack
	// here in order to handle cases of embedded structs which should not
	// result in sub-trees. Instead, we add all embedded struct fields to
	// the stack so they are populated as though they are top level fields.
	stack := make([]fieldAndValue, 0, vv.NumField())
	for x := 0; x < vv.NumField(); x = x + 1 {
		// The Value and Type versions of Field() return different types that
		// contain different information. The Type.Field() generates a TypeField
		// that contains the details needed to determin the field name and whether
		// it is embedded or not. The Value.Field() returns a Value that can be used
		// to manipulate the field.
		stack = append(stack, fieldAndValue{Value: vv.Field(x), Field: vv.Type().Field(x)})
	}
	for len(stack) > 0 {
		var current fieldAndValue
		current, stack = stack[len(stack)-1], stack[:len(stack)-1]
		currentF := current.Field
		currentV := current.Value
		currentVV := reflect.Indirect(currentV)
		desc := currentF.Tag.Get("description")

		if currentVV.Kind() == reflect.Struct && currentF.Anonymous {
			for x := 0; x < currentVV.NumField(); x = x + 1 {
				stack = append(stack, fieldAndValue{Value: currentVV.Field(x), Field: currentVV.Type().Field(x)})
			}
			continue
		}
		if !currentVV.CanAddr() {
			// Catch-all case for values that cannot be addressed. It's difficult to
			// enumerate all of the cases where this could be true. Generally speaking,
			// this is activated any time a pointer to the value cannot be created
			// which would prevent us from modifying the value later. Some examples
			// include most uses of Duration which are handled above, use of some
			// exported constants from a package, or when a non-pointer struct instances
			// is given. There's not much we can do except nudge the consumer to
			// change the value.
			return nil, fmt.Errorf("%s field %s.%s must be a pointer type", currentV.Type(), name, currentF.Name)
		}
		if currentVV.Kind() != reflect.Struct ||
			currentVV.Type().String() == timeName {
			set, err := settingFromValue(currentF.Name, desc, currentVV)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to convert %s.%s due to: %s",
					g.NameValue, currentF.Name, err.Error(),
				)
			}
			g.SettingValues = append(g.SettingValues, set)
			continue
		}

		sub, err := Convert(currentV.Interface())
		if err != nil {
			return nil, err
		}
		g.GroupValues = append(g.GroupValues, sub)
	}
	return g, nil
}

func settingFromValue(name string, description string, v reflect.Value) (Setting, error) {
	switch v.Type().String() {
	case timeName:
		s := &TimeSetting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("TimeValue").Set(v.Addr())
		return s, nil
	case durationName:
		s := &DurationSetting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("DurationValue").Set(v.Addr())
		return s, nil
	default:
	}
	switch v.Kind() {
	case reflect.Bool:
		s := &BoolSetting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("BoolValue").Set(v.Addr())
		return s, nil
	case reflect.Int8:
		s := &Int8Setting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("Int8Value").Set(v.Addr())
		return s, nil
	case reflect.Int16:
		s := &Int16Setting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("Int16Value").Set(v.Addr())
		return s, nil
	case reflect.Int32:
		s := &Int32Setting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("Int32Value").Set(v.Addr())
		return s, nil
	case reflect.Int64:
		s := &Int64Setting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("Int64Value").Set(v.Addr())
		return s, nil
	case reflect.Map:
		vTypeStored := v.Type()
		switch vTypeStored.String(){
		case "map[string][]string":
			s := &StringMapStringSliceSetting{
				BaseSetting: &BaseSetting{
					NameValue:        name,
					DescriptionValue: description,
				},
			}
			sv := reflect.Indirect(reflect.ValueOf(s))
			sv.FieldByName("StringMapStringSliceValue").Set(v.Addr())
			return s, nil
		default:
			return nil, fmt.Errorf("unknown map value type for setting %s", vTypeStored)
		}
	case reflect.Uint:
		s := &UintSetting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("UintValue").Set(v.Addr())
		return s, nil
	case reflect.Uint8:
		s := &Uint8Setting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("Uint8Value").Set(v.Addr())
		return s, nil
	case reflect.Uint16:
		s := &Uint16Setting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("Uint16Value").Set(v.Addr())
		return s, nil
	case reflect.Uint32:
		s := &Uint32Setting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("Uint32Value").Set(v.Addr())
		return s, nil
	case reflect.Uint64:
		s := &Uint64Setting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("Uint64Value").Set(v.Addr())
		return s, nil
	case reflect.Int:
		s := &IntSetting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("IntValue").Set(v.Addr())
		return s, nil
	case reflect.Float32:
		s := &Float32Setting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("Float32Value").Set(v.Addr())
		return s, nil
	case reflect.Float64:
		s := &Float64Setting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("Float64Value").Set(v.Addr())
		return s, nil
	case reflect.String:
		s := &StringSetting{
			BaseSetting: &BaseSetting{
				NameValue:        name,
				DescriptionValue: description,
			},
		}
		sv := reflect.Indirect(reflect.ValueOf(s))
		sv.FieldByName("StringValue").Set(v.Addr())
		return s, nil
	case reflect.Slice:
		if v.Type().Elem().String() == durationName {
			s := &DurationSliceSetting{
				BaseSetting: &BaseSetting{
					NameValue:        name,
					DescriptionValue: description,
				},
			}
			sv := reflect.Indirect(reflect.ValueOf(s))
			sv.FieldByName("DurationSliceValue").Set(v.Addr())
			return s, nil
		}
		switch v.Type().Elem().Kind() {
		case reflect.String:
			s := &StringSliceSetting{
				BaseSetting: &BaseSetting{
					NameValue:        name,
					DescriptionValue: description,
				},
			}
			sv := reflect.Indirect(reflect.ValueOf(s))
			sv.FieldByName("StringSliceValue").Set(v.Addr())
			return s, nil
		case reflect.Int:
			s := &IntSliceSetting{
				BaseSetting: &BaseSetting{
					NameValue:        name,
					DescriptionValue: description,
				},
			}
			sv := reflect.Indirect(reflect.ValueOf(s))
			sv.FieldByName("IntSliceValue").Set(v.Addr())
			return s, nil
		default:
			return nil, fmt.Errorf("unknown setting type []%s", v.Type().Elem().Kind())
		}
	default:
		return nil, fmt.Errorf("unknown setting type %s", v.Kind())
	}
}
