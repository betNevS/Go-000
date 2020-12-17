package data

import (
	"log"
	"math/rand"

	"github.com/betNevS/Go-000/Week04/internal/biz"
)

var _ biz.UserRepo = new(userRepo)

func NewUserRepo() biz.UserRepo {
	return &userRepo{}
}

type userRepo struct {
}

func (r *userRepo) Save(u *biz.User) int32 {
	id := rand.Intn(100)
	log.Printf("Register user, name: %s, age: %d, id: %d\n", u.Name, u.Age, id)
	return int32(id)
}
