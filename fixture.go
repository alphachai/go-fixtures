package fixtures

import (
	"fmt"
	"reflect"
)

type BaseFixture struct{}

func (f *BaseFixture) Type() string {
	return fmt.Sprint(reflect.TypeOf(f).Elem())
}

type Fixture interface {
	Type() string
	SetUp() error
	TearDown() error
}
