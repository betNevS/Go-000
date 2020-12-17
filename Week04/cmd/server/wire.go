// +build wireinject

package main

import (
	"github.com/betNevS/Go-000/Week04/internal/biz"
	"github.com/betNevS/Go-000/Week04/internal/data"
	"github.com/google/wire"
)

func InitUserRegisterCase() *biz.UserRegisterCase {
	wire.Build(biz.NewUserRegisterCase, data.NewUserRepo)
	return &biz.UserRegisterCase{}
}

