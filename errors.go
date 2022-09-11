package errorx

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

func New(text string, opts ...Option) error {
	return Wrap(errors.New(text), opts...)
}

func Wrap(err error, opts ...Option) error {
	errx := &errorx{inner: err}
	for _, opt := range opts {
		opt(errx)
	}
	return errx
}

func AsError(err error) *errorx {
	var e *errorx
	if errors.As(err, &e) {
		return e
	}
	return nil
}

type errorx struct {
	inner    error
	internal error
	fields   []Field
	typ      Type
}

func (e errorx) Error() string   { return e.inner.Error() }
func (e errorx) Internal() error { return e.internal }
func (e errorx) Fields() []Field { return e.fields }
func (e errorx) Type() Type      { return e.typ }

// Deprecated: use Internal instead
func (e errorx) Inner() error { return e.internal }

func (e errorx) JSONString() (string, error) {
	fields := map[string]string{}
	if internal := e.internal; internal != nil {
		fields["internal"] = internal.Error()
	}
	for _, field := range e.fields {
		fields[field.Key] = field.Value
	}
	b, err := json.Marshal(fields)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type Field struct {
	Key, Value string
}

type Type int

type Option func(*errorx)

func WithInternal(err error) Option {
	return func(e *errorx) {
		var errx *errorx
		if errors.As(err, &errx) {
			e.fields = append(e.fields, errx.fields...)
			if internal := errx.Internal(); internal != nil {
				e.internal = internal
				return
			}
			e.internal = errx.internal
			return
		}
		e.internal = err
	}
}

// Deprecated: use WithInternal instead
func WithInner(err error) Option {
	return WithInternal(err)
}

func WithAny(key string, value any) Option {
	return func(e *errorx) {
		v := fmt.Sprint(value)
		e.fields = append(e.fields, Field{
			Key:   key,
			Value: v,
		})
	}
}

func WithString(key, value string) Option {
	return func(e *errorx) {
		e.fields = append(e.fields, Field{
			Key:   key,
			Value: value,
		})
	}
}

func WithInt(key string, value int64) Option {
	return func(e *errorx) {
		e.fields = append(e.fields, Field{
			Key:   key,
			Value: strconv.FormatInt(value, 10),
		})
	}
}

func WithType(typ Type) Option {
	return func(e *errorx) {
		e.typ = typ
	}
}

func WithCaller(skip ...int) Option {
	return func(e *errorx) {
		s := 2
		if len(skip) > 0 {
			s += skip[0]
		}
		frame, ok := getCallerFrame(s)
		if !ok {
			return
		}
		e.fields = append(e.fields, Field{
			Key:   "caller",
			Value: fmt.Sprintf("%s:%d", trimFilePath(frame.File), frame.Line),
		})
	}
}

func WithStacktrace() Option {
	return func(e *errorx) {
		e.fields = append(e.fields, Field{
			Key:   "stacktrace",
			Value: string(debug.Stack()),
		})
	}
}

func getCallerFrame(skip int) (frame runtime.Frame, ok bool) {
	const skipOffset = 2 // skip getCallerFrame and Callers

	pc := make([]uintptr, 1)
	numFrames := runtime.Callers(skip+skipOffset, pc)
	if numFrames < 1 {
		return
	}

	frame, _ = runtime.CallersFrames(pc).Next()
	return frame, frame.PC != 0
}

func trimFilePath(filePath string) string {
	idx := strings.LastIndexByte(filePath, '/')
	if idx == -1 {
		return filePath
	}
	idx = strings.LastIndexByte(filePath[:idx], '/')
	if idx == -1 {
		return filePath
	}
	return filePath[idx+1:]
}
