package errorx

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	text := "test message"
	text2 := "test message 2"
	internal := errors.New("internal error message")
	err := New(text, WithInternal(internal), WithCaller(), WithStacktrace())

	err = New(text2, WithInternal(err))

	err = fmt.Errorf("wrapped err: %w", err)

	var errx *errorx
	if !errors.As(err, &errx) {
		t.Fatal("error is not assignable to extError")
	}
	if errx.Error() != text2 {
		t.Fatal("error text is incorrect")
	}
	if errx.Internal().Error() != internal.Error() {
		t.Fatal("internal error is incorrect")
	}
	if len(errx.Fields()) < 2 {
		t.Fatal("fields are incorrect")
	}
}
