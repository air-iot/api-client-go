package errors

import (
	"fmt"
	"testing"
)

func TestNewError(t *testing.T) {
	e := NewError(fmt.Errorf("err1"))
	t.Log(e.Error())
}

func TestNewMsg(t *testing.T) {
	e := NewMsg("a,%d", 2)
	t.Log(e.Error())
}

func TestNewErrorMsg(t *testing.T) {
	e := NewErrorMsg(fmt.Errorf("err1"), "a,%d", 2)
	t.Log(e.Error())

	e = NewErrorMsg(NewMsg("err2"), "2,%d", 3)
	t.Log(e.Error())
}
