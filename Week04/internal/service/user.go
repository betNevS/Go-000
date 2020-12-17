package service

import (
	"context"
	v1 "github.com/betNevS/Go-000/Week04/api/user/v1"
	"github.com/betNevS/Go-000/Week04/internal/biz"
)

type UserService struct {
	v1.UnimplementedUserServer
	u *biz.UserRegisterCase // biz
}

func NewUserService(u *biz.UserRegisterCase) v1.UserServer {
	return &UserService{
		u:u,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, r *v1.UserRequest) (*v1.UserReply, error)  {
	u := &biz.User{Name:r.Name, Age:r.Age}
	s.u.SaveUser(u)
	return &v1.UserReply{Id:u.ID}, nil
}
