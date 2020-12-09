package service

import (
	"github.com/betNevS/Go-000/Week02/dao"
	"github.com/betNevS/Go-000/Week02/model"
)

func FindUserByName(name string) (*model.User, error) {
	user, err := dao.FindUserByName(name)
	if err != nil {
		return nil, err
	}
	user.Name += "是个好孩子！"
	return user, nil
}
