package domain

import "github.com/axelzv9/errorx"

const (
	TypeInvalidParams errorx.Type = iota + 1
	TypeUnauthorized
	TypeApplication
	// ... any error types in your application domain level
)

const (
	invalidParamsErrorText = "Invalid params"
	databaseErrorText      = "Unknown server error"
	unauthorizedErrorText  = "Unauthorized"
)

func WrapDBError(err error, opts ...errorx.Option) error {
	if err == nil {
		return nil
	}
	if len(opts) > 0 {
		opts = append(
			append(
				make([]errorx.Option, 0, 2+len(opts)),
				errorx.WithInternal(err), // wrapping real db error
				errorx.WithCaller(1),     // save the line of code where the error occurred
			),
			opts...,
		)
		return errorx.New(databaseErrorText, opts...)
	}
	return errorx.New(databaseErrorText, errorx.WithInternal(err), errorx.WithCaller(1))
}

func NewErrUnauthorized() error {
	return errorx.New(unauthorizedErrorText, errorx.WithType(TypeUnauthorized))
}

func NewErrInvalidParams(err error) error {
	return errorx.New(invalidParamsErrorText, errorx.WithType(TypeInvalidParams), errorx.WithInternal(err))
}
