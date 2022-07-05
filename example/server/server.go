package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/axelzv9/errorx"
	"github.com/axelzv9/errorx/example/domain"
)

func UnauthorizedHandler(w http.ResponseWriter, _ *http.Request) {
	err := domain.NewErrUnauthorized()
	_ = writeError(w, err)
}

func InvalidParamsHandler(w http.ResponseWriter, _ *http.Request) {
	err := domain.NewErrInvalidParams(errors.New("some param is invalid"))
	_ = writeError(w, err)
}

// here you can manage output error texts and print to log full debug information
func writeError(w http.ResponseWriter, err error) error {
	if errX := errorx.AsError(err); errX != nil {
		str, _ := errX.JSONString()
		log.Printf("error: %s fields: %s", err, str)

		for _, field := range errX.Fields() {
			_ = field.Key
			_ = field.Value
		}

		switch errX.Type() {
		case domain.TypeApplication:
			w.WriteHeader(http.StatusOK)
		case domain.TypeInvalidParams:
			w.WriteHeader(http.StatusBadRequest)
		case domain.TypeUnauthorized:
			w.WriteHeader(http.StatusUnauthorized)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Add("Content-Type", "plain/text")
		_, e := w.Write([]byte(err.Error()))
		return e
	}

	log.Printf("unknown error type %s", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Content-Type", "plain/text")
	_, e := w.Write([]byte(err.Error()))
	return e
}
