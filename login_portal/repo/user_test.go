package repo

import (
	"login_portal/model"
	"testing"
)

func Test_UserRepo_NormalCreate(t *testing.T) {
	var users = []model.User{}
	r := NewUser(users)
	id, err := r.Create("jono", "secret")
	if err != nil {
		t.Error("harusnya tidak error")
	}
	if id == 0 {
		t.Error("harusnya mendapatkan id 1. karena data masih kosong")
	}
}

func Test_UserRepo_DuplicateUsername(t *testing.T) {
	var users = []model.User{
		{ID: 1, Username: "jono", Password: "secret"},
	}

	r := NewUser(users)
	_, err := r.Create("jono", "secret")
	if err == nil {
		t.Error("harusnya error:", ErrUsernameExists)
	}
}

func Test_UserRepo_NormalFindByID(t *testing.T) {
	var users = []model.User{
		{ID: 1, Username: "jono", Password: "secret"},
	}

	r := NewUser(users)
	user, err := r.FindUserByID(1)
	if err != nil {
		t.Error("harusnya tidak ada error")
	}

	if user.ID != 1 {
		t.Error("harusnya dapat user.ID = 1")
	}
}

func Test_UserRepo_NotFoundFindByID(t *testing.T) {
	var users = []model.User{
		{ID: 1, Username: "jono", Password: "secret"},
	}

	r := NewUser(users)
	_, err := r.FindUserByID(2)
	if err == nil {
		t.Error("harusnya error: ", ErrUserNotFound)
	}
}

func Test_UserRepo_NormalUpdate(t *testing.T) {
	var users = []model.User{
		{ID: 1, Username: "jono", Password: "secret", Profile: model.UserProfile{}},
	}

	r := NewUser(users)
	err := r.Update(1, UserUpdatePayload{Username: new("dudung")})
	if err != nil {
		t.Error("harusnya tidak error")
	}

	user, _ := r.FindUserByID(1)
	if user.Username != "dudung" {
		t.Error("harusnya username berubah")
	}

	if user.Profile.FirstName != "" {
		t.Error("harusnya firstname kosong, karena belum di update sama sekali")
	}

	if user.Profile.LastName != "" {
		t.Error("harusnya lastname kosong, karena belum di update sama sekali")
	}

	err = r.Update(1, UserUpdatePayload{
		FirstName: new("Dudung"),
		LastName:  new("Sujono"),
	})
	if err != nil {
		t.Error("harusnya tidak error")
	}

	nUser, _ := r.FindUserByID(1)
	if nUser.Profile.FirstName != "Dudung" {
		t.Error("harusnya profile firstname berubah")
	}

	if nUser.Profile.LastName != "Sujono" {
		t.Error("harusnya profile lastname berubah")
	}
}

func Test_UserRepo_NormalUpdateAll(t *testing.T) {
	var users = []model.User{
		{ID: 1, Username: "jono", Password: "secret", Profile: model.UserProfile{}},
	}

	r := NewUser(users)
	err := r.Update(1, UserUpdatePayload{
		Username:  new("dudung"),
		Password:  new("dudung-pw"),
		FirstName: new("Dudung"),
		LastName:  new("Sujono")},
	)
	if err != nil {
		t.Error("harusnya tidak error")
	}

	user, _ := r.FindUserByID(1)
	if user.Username != "dudung" {
		t.Error("harusnya username berubah")
	}

	if user.Password != "dudung-pw" {
		t.Error("harusnya password berubah")
	}

	if user.Profile.FirstName != "Dudung" {
		t.Error("harusnya profile firstname berubah")
	}

	if user.Profile.LastName != "Sujono" {
		t.Error("harusnya profile lastname berubah")
	}
}

func Test_UserRepo_ValueNilUpdate(t *testing.T) {
	var users = []model.User{
		{ID: 1, Username: "jono", Password: "secret", Profile: model.UserProfile{}},
	}

	r := NewUser(users)
	err := r.Update(1, UserUpdatePayload{})
	if err == nil {
		t.Error("harusnya error", ErrNothingUpdate)
	}
}

func Test_UserRepo_NormalDelete(t *testing.T) {
	var users = []model.User{
		{ID: 1, Username: "jono", Password: "secret", Profile: model.UserProfile{}},
	}

	r := NewUser(users)
	err := r.Delete(1)
	if err != nil {
		t.Error("harusnya tidak ada error")
	}
}

func Test_UserRepo_NotFoundDelete(t *testing.T) {
	var users = []model.User{
		{ID: 1, Username: "jono", Password: "secret", Profile: model.UserProfile{}},
	}

	r := NewUser(users)
	err := r.Delete(4)
	if err == nil {
		t.Error("harusnya ada error: ", ErrUserNotFound)
	}
}

func Test_UserRepo_NormalGenerateID(t *testing.T) {
	var users = []model.User{
		{ID: 1, Username: "jono", Password: "secret", Profile: model.UserProfile{}},
		{ID: 5, Username: "jimi", Password: "secret", Profile: model.UserProfile{}},
		{ID: 8, Username: "ani", Password: "secret", Profile: model.UserProfile{}},
		{ID: 3, Username: "sasa", Password: "secret", Profile: model.UserProfile{}},
	}

	r := NewUser(users)
	id := r.genID()
	if id != 9 {
		t.Error("harusnya dapet id 9")
	}
}
