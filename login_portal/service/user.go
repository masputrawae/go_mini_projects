package service

import (
	"errors"
	"login_portal/model"
	"login_portal/repo"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameToShort            = errors.New("username is too short")
	ErrPasswordToShort            = errors.New("password is too short")
	ErrUsernameOrPasswordNotMatch = errors.New("username or password does not match")
)

type User struct {
	userRepo *repo.User
}

func NewUser(userRepo *repo.User) *User {
	return &User{userRepo: userRepo}
}

// method untuk buat akun baru.
func (u *User) Register(username, password string) (int, error) {
	if len(username) < 3 {
		return -1, ErrUsernameToShort
	}

	if len(password) < 8 {
		return -1, ErrPasswordToShort
	}

	hash, err := hashPassword(password)
	if err != nil {
		return -1, err
	}

	id, err := u.userRepo.Create(username, hash)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// method untuk login
func (u *User) Login(username, password string) (int, error) {
	user, err := u.userRepo.FindUserByUsername(username)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return -1, ErrUsernameOrPasswordNotMatch
		}

		return -1, err
	}

	if !checkPassword(user.Password, password) {
		return -1, ErrUsernameOrPasswordNotMatch
	}

	return user.ID, nil
}

func (u *User) GetUserByID(id int) (model.User, error) {
	user, err := u.userRepo.FindUserByID(id)
	if err != nil {
		return model.User{}, err
	}

	// hapus password agar tidak bocor
	user.Password = ""

	return user, nil
}

// method untuk update profile
func (u *User) Edit(id int, p repo.UserUpdatePayload) error {
	payload := repo.UserUpdatePayload{}

	if p.Password != nil && len(*p.Password) > 8 {
		hash, err := hashPassword(*p.Password)
		if err != nil {
			return err
		}
		payload.Password = new(hash)
	}

	if p.Username != nil && len(*p.Username) > 3 {
		payload.Username = p.Username
	}

	if p.FirstName != nil {
		payload.FirstName = p.FirstName
	}

	if p.LastName != nil {
		payload.LastName = p.LastName
	}

	return u.userRepo.Update(id, payload)
}

// fungsi kecil untuk hash dan check password
func hashPassword(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func checkPassword(h, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p)) == nil
}
