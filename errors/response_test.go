package errors

import (
	"testing"
)

func TestNewError(t *testing.T) {
	e := NewError(NewMsg("err1"))
	t.Log(e.Error())
}

func TestNewMsg(t *testing.T) {
	e := NewMsg("a,%d", 2)
	t.Log(e.Error())
}

func TestNewErrorMsg(t *testing.T) {
	e := NewErrorMsg(NewMsg("err1"), "a,%d", 2)
	t.Log(e.Error())

	e = NewErrorMsg(NewMsg("err2"), "2,%d", 3)
	t.Log(e.Error())
}
