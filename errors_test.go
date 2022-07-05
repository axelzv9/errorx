package errorx

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	text := "test message"
	text2 := "test message 2"
	inner := errors.New("inner error message")
	err := New(text, WithInner(inner), WithCaller(), WithStacktrace())

	err = New(text2, WithInner(err))

	err = fmt.Errorf("wrapped err: %w", err)

	var ext *errorx
	if !errors.As(err, &ext) {
		t.Fatal("error is not assignable to extError")
	}
	if ext.Error() != text2 {
		t.Fatal("error text is incorrect")
	}
	if ext.Inner().Error() != inner.Error() {
		t.Fatal("inner error is incorrect")
	}
	if len(ext.Fields()) < 2 {
		t.Fatal("fields are incorrect")
	}
}
