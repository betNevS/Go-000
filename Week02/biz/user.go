package biz

import (
	"github.com/betNevS/Go-000/Week02/model"
	"github.com/betNevS/Go-000/Week02/service"
)

func FindUser(name string) (*model.User, error) {
	return service.FindUserByName(name)
}
