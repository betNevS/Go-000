package dao

import (
	"database/sql"

	"github.com/betNevS/Go-000/Week02/model"
	"github.com/pkg/errors"
)

var dbName = "user"

func FindUserByName(name string) (*model.User, error) {
	db, err := NewDB(dbName)
	if err != nil {
		return nil, err
	}
	userSQL := "SELECT * FROM userinfo WHERE `name` = ?"
	user := &model.User{}
	err = db.QueryRow(userSQL, name).Scan(&user.Id, &user.Name, &user.Age)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrapf(ErrNotFound, "sql [%s]", userSQL)
		}
		return nil, errors.Wrapf(err, "get user name[%s] error", name)
	}
	return user, nil
}
