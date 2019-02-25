package settings

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type testConf struct {
	Value string
}
type testItem struct {
	Value string
}

func (*testItem) Do() error { return nil }

type testInterface interface {
	Do() error
}

type missingSettings struct{}

func (*missingSettings) New(context.Context, *testConf) (*testItem, error) { return &testItem{}, nil }

type missingNew struct{}

func (*missingNew) Settings() *testConf { return &testConf{} }

type settingsTooManyArgs struct{}

func (*settingsTooManyArgs) Settings(context.Context) *testConf { return &testConf{} }
func (*settingsTooManyArgs) New(context.Context, *testConf) (*testItem, error) {
	return &testItem{}, nil
}

type settingsTooManyReturns struct{}

func (*settingsTooManyReturns) Settings() (*testConf, error) { return &testConf{}, nil }
func (*settingsTooManyReturns) New(context.Context, *testConf) (*testItem, error) {
	return &testItem{}, nil
}

type newNotEnoughArgs struct{}

func (*newNotEnoughArgs) New(*testConf) (*testItem, error) { return &testItem{}, nil }
func (*newNotEnoughArgs) Settings() *testConf              { return &testConf{} }

type newTooManyArgs struct{}

func (*newTooManyArgs) New(context.Context, int, *testConf) (*testItem, error) {
	return &testItem{}, nil
}
func (*newTooManyArgs) Settings() *testConf { return &testConf{} }

type newNotEnoughReturns struct{}

func (*newNotEnoughReturns) New(context.Context, *testConf) error {
	return nil
}
func (*newNotEnoughReturns) Settings() *testConf { return &testConf{} }

type newTooManyReturns struct{}

func (*newTooManyReturns) New(context.Context, *testConf) (*testItem, int, error) {
	return nil, 0, nil
}
func (*newTooManyReturns) Settings() *testConf { return &testConf{} }

type newWrongContextType struct{}

func (*newWrongContextType) New(int, *testConf) (*testItem, error) { return &testItem{}, nil }
func (*newWrongContextType) Settings() *testConf                   { return &testConf{} }

type newWrongErrorType struct{}

func (*newWrongErrorType) New(context.Context, *testConf) (*testItem, int) { return &testItem{}, 0 }
func (*newWrongErrorType) Settings() *testConf                             { return &testConf{} }

type settingsNewMismatch struct{}

func (*settingsNewMismatch) New(context.Context, int) (*testItem, error) { return &testItem{}, nil }
func (*settingsNewMismatch) Settings() *testConf                         { return &testConf{} }

type testComponent struct{}

func (*testComponent) New(_ context.Context, c *testConf) (*testItem, error) {
	return &testItem{Value: c.Value}, nil
}
func (*testComponent) Settings() *testConf { return &testConf{} }

type testComponentInterface struct{}

func (*testComponentInterface) New(_ context.Context, c *testConf) (testInterface, error) {
	return &testItem{Value: c.Value}, nil
}
func (*testComponentInterface) Settings() *testConf { return &testConf{} }

type testComponentErr struct{}

func (*testComponentErr) New(_ context.Context, c *testConf) (*testItem, error) {
	return nil, errors.New("error")
}
func (*testComponentErr) Settings() *testConf { return &testConf{} }

var (
	nilDst *testItem
)

func TestNewComponent(t *testing.T) {
	ifacePtr := new(testInterface)
	type args struct {
		ctx context.Context
		s   Source
		v   interface{}
		dst interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    interface{}
	}{
		{
			name: "missingSettings",
			args: args{
				ctx: context.Background(),
				s:   &MapSource{Map: make(map[string]interface{})},
				v:   &missingSettings{},
				dst: new(testItem),
			},
			wantErr: true,
		},
		{
			name: "missingNew",
			args: args{
				ctx: context.Background(),
				s:   &MapSource{Map: make(map[string]interface{})},
				v:   &missingNew{},
				dst: new(testItem),
			},
			wantErr: true,
		},
		{
			name: "settingsTooManyArgs",
			args: args{
				ctx: context.Background(),
				s:   &MapSource{Map: make(map[string]interface{})},
				v:   &settingsTooManyArgs{},
				dst: new(testItem),
			},
			wantErr: true,
		},
		{
			name: "settingsTooManyReturns",
			args: args{
				ctx: context.Background(),
				s:   &MapSource{Map: make(map[string]interface{})},
				v:   &settingsTooManyReturns{},
				dst: new(testItem),
			},
			wantErr: true,
		},
		{
			name: "newNotEnoughArgs",
			args: args{
				ctx: context.Background(),
				s:   &MapSource{Map: make(map[string]interface{})},
				v:   &newNotEnoughArgs{},
				dst: new(testItem),
			},
			wantErr: true,
		},
		{
			name: "newTooManyArgs",
			args: args{
				ctx: context.Background(),
				s:   &MapSource{Map: make(map[string]interface{})},
				v:   &newTooManyArgs{},
				dst: new(testItem),
			},
			wantErr: true,
		},
		{
			name: "newNotEnoughReturns",
			args: args{
				ctx: context.Background(),
				s:   &MapSource{Map: make(map[string]interface{})},
				v:   &newNotEnoughReturns{},
				dst: new(testItem),
			},
			wantErr: true,
		},
		{
			name: "newTooManyReturns",
			args: args{
				ctx: context.Background(),
				s:   &MapSource{Map: make(map[string]interface{})},
				v:   &newTooManyReturns{},
				dst: new(testItem),
			},
			wantErr: true,
		},
		{
			name: "newWrongContextType",
			args: args{
				ctx: context.Background(),
				s:   &MapSource{Map: make(map[string]interface{})},
				v:   &newWrongContextType{},
				dst: new(testItem),
			},
			wantErr: true,
		},
		{
			name: "newWrongErrorType",
			args: args{
				ctx: context.Background(),
				s:   &MapSource{Map: make(map[string]interface{})},
				v:   &newWrongErrorType{},
				dst: new(testItem),
			},
			wantErr: true,
		},
		{
			name: "settingsNewMismatch",
			args: args{
				ctx: context.Background(),
				s:   &MapSource{Map: make(map[string]interface{})},
				v:   &settingsNewMismatch{},
				dst: new(testItem),
			},
			wantErr: true,
		},
		{
			name: "nil destination",
			args: args{
				ctx: context.Background(),
				s:   &MapSource{Map: make(map[string]interface{})},
				v:   &settingsNewMismatch{},
				dst: nilDst,
			},
			wantErr: true,
		},
		{
			name: "non-pointer destination",
			args: args{
				ctx: context.Background(),
				s:   &MapSource{Map: make(map[string]interface{})},
				v:   &settingsNewMismatch{},
				dst: testItem{},
			},
			wantErr: true,
		},
		{
			name: "sets pointer on success",
			args: args{
				ctx: context.Background(),
				s: NewMapSource(map[string]interface{}{
					"testconf": map[string]interface{}{
						"value": "test",
					},
				}),
				v:   &testComponent{},
				dst: new(testItem),
			},
			wantErr: false,
			want:    &testItem{Value: "test"},
		},
		{
			name: "sets interface type pointer on success",
			args: args{
				ctx: context.Background(),
				s: NewMapSource(map[string]interface{}{
					"testconf": map[string]interface{}{
						"value": "test",
					},
				}),
				v:   &testComponentInterface{},
				dst: ifacePtr,
			},
			wantErr: false,
			want:    ifacePtr,
		},
		{
			name: "returns error if constructor fails",
			args: args{
				ctx: context.Background(),
				s:   &MapSource{Map: make(map[string]interface{})},
				v:   &testComponentErr{},
				dst: new(testItem),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = NewComponent(tt.args.ctx, tt.args.s, tt.args.v, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("NewComponent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.dst, tt.want) {
				t.Errorf("NewComponent() dst = %v, want %v", tt.args.dst, tt.want)
			}
		})
	}
}
