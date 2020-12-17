package biz

type User struct {
	ID int32
	Name string
	Age int32
}

type UserRepo interface {
	Save(*User) int32
}

func NewUserRegisterCase(repo UserRepo) *UserRegisterCase {
	return &UserRegisterCase{repo:repo}
}

type UserRegisterCase struct {
	repo UserRepo
}

func (s *UserRegisterCase) SaveUser(u *User) {
	id := s.repo.Save(u)
	u.ID = id
}