package service

import (
	"errors"
	"login_portal/model"
	"login_portal/repo"
	"testing"
)

func Test_UserService_NormalRegister(t *testing.T) {
	users := []model.User{}
	r := repo.NewUser(users)
	s := NewUser(r)

	id, err := s.Register("john", "secret-123")
	if err != nil {
		t.Error("harusnya tidak error")
	}

	if id != 1 {
		t.Error("harusnya mendapatkan id 1")
	}
}

func Test_UserService_ShortUsernameRegister(t *testing.T) {
	users := []model.User{}
	r := repo.NewUser(users)
	s := NewUser(r)

	_, err := s.Register("jo", "secret-123")
	if !errors.Is(err, ErrUsernameToShort) {
		t.Error("harusnya error", ErrUsernameToShort)
	}

	if len(users) != 0 {
		t.Error("harusnya user tidak disimpan")
	}
}

func Test_UserService_ShortPasswordRegister(t *testing.T) {
	users := []model.User{}
	r := repo.NewUser(users)
	s := NewUser(r)

	_, err := s.Register("john", "secret")
	if !errors.Is(err, ErrPasswordToShort) {
		t.Error("harusnya error", ErrPasswordToShort)
	}

	if len(users) != 0 {
		t.Error("harusnya user tidak disimpan")
	}
}

func Test_UserService_SameUsernameRegister(t *testing.T) {
	users := []model.User{
		{ID: 1, Username: "john", Password: "secret-123"},
	}

	r := repo.NewUser(users)
	s := NewUser(r)

	_, err := s.Register("john", "secret-other-123")
	if !errors.Is(err, repo.ErrUsernameExists) {
		t.Error("harusnya mengembalikan error", repo.ErrUsernameExists)
	}

	if len(users) != 1 {
		t.Error("harusnya user tidak disimpan")
	}
}

func Test_UserService_NormalLogin(t *testing.T) {
	hash, _ := hashPassword("secret-123")
	users := []model.User{
		{ID: 1, Username: "john", Password: hash},
	}
	r := repo.NewUser(users)
	s := NewUser(r)

	id, err := s.Login("john", "secret-123")
	if err != nil {
		t.Error("harusnya tidak ada error")
	}

	if id != 1 {
		t.Error("harusnya mengembalikan id 1")
	}
}

func Test_UserService_PasswordNotMatchLogin(t *testing.T) {
	hash, _ := hashPassword("secret-123")
	users := []model.User{
		{ID: 1, Username: "john", Password: hash},
	}
	r := repo.NewUser(users)
	s := NewUser(r)

	_, err := s.Login("john", "secret-not-match")
	if !errors.Is(err, ErrUsernameOrPasswordNotMatch) {
		t.Error("harusnya mengembalikan error", ErrUsernameOrPasswordNotMatch)
	}
}

func Test_UserService_UsernameNotMatchLogin(t *testing.T) {
	hash, _ := hashPassword("secret-123")
	users := []model.User{
		{ID: 1, Username: "john", Password: hash},
	}
	r := repo.NewUser(users)
	s := NewUser(r)

	_, err := s.Login("john-e", "secret-123")
	if !errors.Is(err, ErrUsernameOrPasswordNotMatch) {
		t.Error("harusnya mengembalikan error: ", ErrUsernameOrPasswordNotMatch)
	}
}

func Test_UserService_NormalGetUserByID(t *testing.T) {
	users := []model.User{
		{ID: 1, Username: "john", Password: "secret-123"},
	}

	r := repo.NewUser(users)
	s := NewUser(r)

	user, err := s.GetUserByID(1)
	if err != nil {
		t.Error("harusnya tidak ada error")
	}

	if user.ID != 1 {
		t.Error("harusnya mendapatkan id 1")
	}
}

func Test_UserService_NotFoundGetUserByID(t *testing.T) {
	users := []model.User{
		{ID: 1, Username: "john", Password: "secret-123"},
	}

	r := repo.NewUser(users)
	s := NewUser(r)

	user, err := s.GetUserByID(2)
	if !errors.Is(err, repo.ErrUserNotFound) {
		t.Error("harusnya mengembalikan error", repo.ErrUserNotFound)
	}

	if user.ID == 1 {
		t.Error("harusnya tidak dapat id 1")
	}
}
