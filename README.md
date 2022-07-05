# Golang native error extension 

Helper for customization and wrapping actual errors, hiding from users

## Usage example

```go
// wrap inner error
err := errors.New("new inner error")
err = errorx.New("user_friendly_text", errorx.WithInner(err))
...
// add caller line
err = errorx.New("user_friendly_text", errorx.WithInner(err), errorx.WithCaller())
...
// add custom fields
err = errorx.New("user_friendly_text", 
	errorx.WithInner(err), 
	errorx.WithString("key", "value"),
	errorx.WithInt("key2", 123),
)
...
// typed error
const CustomType errorx.Type = 1

err = errorx.New("user_friendly_text", errorx.WithInner(err), errorx.WithType(CustomType))
...
// extract custom fields to json
if errX := errorx.AsError(err); errX != nil {
    str, _ := errX.JSONString()
    log.Printf("error: %s fields: %s", err, str)
}
...
// extract custom fields
if errX := errorx.AsError(err); errX != nil {
    for _, field := range errX.Fields() {
		var key, value string
		key = field.Key
        value = field.Value
    }
}
...
// check custom type
switch errX.Type() {
    case CustomType:
		// handle custom type err
    default:
        // do something else
}
```

See full integration example [here](https://github.com/axelzv9/errorx/tree/master/example).
