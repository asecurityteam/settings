package settings

import (
	"context"
	"fmt"
	"reflect"
)

// NewComponent is the entry point for the high-level api. This method manages much
// of the complexity of adding configuration to a system. The given context and
// source are used for all lookups of configuration data. The given value is an
// implementation of the component contract and the destination is a pointer created
// with new(T) where T is the output type (or equivalent convertible) to the element
// produced by the component contract implementation.
//
// The component contract is an interface that all input values must conform to and
// is roughly equivalent to the Factory concept. Each instance of the component
// contract must define two methods: Setting() C and New(context.Context, C) (T,
// error). Due to the lack of generics in go, there's no way to describe this
// contract as an actual go interface that would benefit from static typing support.
// As a result, this method uses reflect to enforce the contract in order to allow
// for C to be any type that is convertible to configuration via the Convert() method
// and for T to be any type that your use case requires.
//
// For example, the most minimal implementation of the contract would look like:
//
//  type Config struct {}
//  type Result struct {}
//  type Component struct {}
//  func (*Component) Settings() *Config { return &Config{} }
//  func (*Component) New(_ context.Context, c *Config) (*Result, error) {
//      return &Result{}, nil
//  }
//
// From here, any number of settings and sub-trees may be added to Config, any
// methods or attributes may be added to Result, and any complexity in the creation
// of Result maybe be added to the Component.New method. To then use this basic
// example as a component you would:
//
//  r := new(Result)
//  err := NewComponent(context.Background(), source, &Component{}, r)
//
// If the resulting error is nil then the destination value, r in this case, now
// points to the output of the Component.New method. The method returns an error any
// time the given component does not satisfy the contract, any time the configuration
// loading fails, or if the Component.New returns an error.
func NewComponent(ctx context.Context, s Source, v interface{}, destination interface{}) error {
	dv := reflect.ValueOf(destination)
	if dv.Kind() != reflect.Ptr {
		// The destination needs to be a pointer value in order for us to
		// set the value at that address with the resulting instance
		return fmt.Errorf("destination %s must be a pointer type", dv.Type())
	}
	if dv.IsNil() {
		return fmt.Errorf("destination %s cannot be nil. use new(T) to make a pointer", dv.Type())
	}

	if err := VerifyComponent(v); err != nil {
		return err
	}
	vv := reflect.ValueOf(v)
	sm := vv.MethodByName("Settings")
	nm := vv.MethodByName("New")
	smOut := sm.Call(nil)[0]

	g, err := Convert(smOut.Interface())
	if err != nil {
		return err
	}
	err = LoadGroups(ctx, s, []Group{g})
	if err != nil {
		return err
	}

	// Once the configuration struct is populated we can generate a new instance
	// of the component and set the destination pointer.
	nOuts := nm.Call([]reflect.Value{
		reflect.ValueOf(ctx).Convert(nm.Type().In(0)),
		smOut.Convert(nm.Type().In(1)),
	})
	nV, nErr := nOuts[0], nOuts[1]
	nV = reflect.Indirect(nV)

	if !nErr.IsNil() {
		return nErr.Interface().(error)
	}
	if !nV.Type().ConvertibleTo(dv.Elem().Type()) {
		return fmt.Errorf("cannot convert %s into %s", nV.Type(), dv.Elem().Type())
	}
	dv.Elem().Set(nV)
	return nil
}

// VerifyComponent checks if a given value implements the Component
// contract.
func VerifyComponent(v interface{}) error {
	// Check for all required method names.
	vv := reflect.ValueOf(v)
	var hasSettings bool
	var hasNew bool
	for x := 0; x < vv.Type().NumMethod(); x = x + 1 {
		fnName := vv.Type().Method(x).Name
		switch fnName {
		case "Settings":
			hasSettings = true
		case "New":
			hasNew = true
		default:
		}
	}
	if !hasSettings {
		return fmt.Errorf("type %s does not have a `Settings() T` method", vv.Type())
	}
	if !hasNew {
		return fmt.Errorf("type %s does not have a `New(ctx, T) (T2, error)` method", vv.Type())
	}
	sm := vv.MethodByName("Settings")
	nm := vv.MethodByName("New")

	// Check that Settings implements the correct signature.
	if sm.Type().NumIn() != 0 {
		return fmt.Errorf("method Settings for %s must not take arguments", vv.Type())
	}
	if sm.Type().NumOut() != 1 {
		return fmt.Errorf("method Settings for %s must return only one value", vv.Type())
	}

	// Grab the return type of Settings for use in validating the New method.
	smOut := sm.Type().Out(0)

	// Check that the New method implements the correct signature.
	if nm.Type().NumIn() != 2 {
		return fmt.Errorf("method New for %s must take only two arguments", vv.Type())
	}
	if nm.Type().NumOut() != 2 {
		return fmt.Errorf("method New for %s must return only two values", vv.Type())
	}
	if !reflect.TypeOf(context.Background()).ConvertibleTo(nm.Type().In(0)) {
		return fmt.Errorf("method New for %s must accept context as the first argument", vv.Type())
	}
	if !smOut.ConvertibleTo(nm.Type().In(1)) {
		return fmt.Errorf("method New for %s must accept an instance of Settings() return value", vv.Type())
	}
	if nm.Type().Out(1).Name() != "error" {
		return fmt.Errorf("method New for %s must return an error as the second value", vv.Type())
	}
	return nil
}

// GroupFromComponent works like Convert to change a struct into a Group but is able
// to do so with an implementation of the Component contract.
func GroupFromComponent(v interface{}) (Group, error) {
	if err := VerifyComponent(v); err != nil {
		return nil, err
	}
	vv := reflect.ValueOf(v)
	sm := vv.MethodByName("Settings")
	smOut := sm.Call(nil)[0]
	return Convert(smOut.Interface())
}
