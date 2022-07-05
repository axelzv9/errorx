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
	err := &errorx{
		text: text,
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}

func AsError(err error) *errorx {
	var e *errorx
	if errors.As(err, &e) {
		return e
	}
	return nil
}

type errorx struct {
	text   string
	inner  error
	fields []Field
	typ    Type
}

func (e errorx) Error() string   { return e.text }
func (e errorx) Inner() error    { return e.inner }
func (e errorx) Fields() []Field { return e.fields }
func (e errorx) Type() Type      { return e.typ }

func (e errorx) JSONString() (string, error) {
	fields := map[string]string{}
	if inner := e.inner; inner != nil {
		fields["inner"] = inner.Error()
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

func WithInner(err error) Option {
	return func(e *errorx) {
		var ext *errorx
		if errors.As(err, &ext) {
			e.fields = append(e.fields, ext.fields...)
			if inner := ext.Inner(); inner != nil {
				e.inner = inner
				return
			}
			e.inner = errors.New(ext.text)
			return
		}
		e.inner = err
	}
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
