package db

import (
	"errors"

	"github.com/axelzv9/errorx"
	"github.com/axelzv9/errorx/example/domain"
)

func SomeDataRequest(id int64) (result any, err error) {
	result, err = execSQL("SELECT * FROM table_name WHERE id = ?", id)
	if err != nil {
		return nil, domain.WrapDBError(err, errorx.WithInt("id", id))
	}
	return result, nil
}

func execSQL(sql string, args ...any) (result any, err error) {
	// execute SQL request
	_, _ = sql, args
	return nil, errors.New("any db error")
}
