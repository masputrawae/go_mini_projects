package repo

import (
	"errors"
	"login_portal/model"
	"sync"
)

// variabel global untuk error
// katanya best practice nya seperti ini (😙)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrNothingUpdate  = errors.New("nothing updated")
	ErrUsernameExists = errors.New("username already exists")
)

type (
	User struct {
		data []model.User
		mu   sync.RWMutex
	}

	// untuk update profile user
	// pakai pointer untuk pengecekan nilai nil
	// memungkinkan update hanya field tertentu saja
	UserUpdatePayload struct {
		Username  *string
		Password  *string
		FirstName *string
		LastName  *string
	}
)

// untuk inisialisasi repo user

func NewUser(data []model.User) *User {
	return &User{data: data}
}

// method untuk membuat user baru.
// parameter membutuhkan username, dan password.

func (u *User) Create(username, password string) (int, error) {
	// dapatkan id baru
	id := u.genID()

	u.mu.Lock()
	defer u.mu.Unlock()

	// cek username apakah sudah ada. username biasanya unik
	for i := range u.data {
		if u.data[i].Username == username {
			return -1, ErrUsernameExists
		}
	}

	// simpan user ke memori (karena tidak pakai database)
	u.data = append(u.data, model.User{
		ID:       id,
		Username: username,
		Password: password,

		// inisialisasi, biar nanti di isi saat update
		Profile: model.UserProfile{},
	})

	// kembalikan id, yang nanti berguna untuk
	// mengambil data user
	return id, nil
}

// method untuk mencari user berdasarkan id
// mengembalikan nilai (model.User, error)

func (u *User) FindUserByID(id int) (model.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	// loop sampai mendapatkan id yang sesuai
	for i := range u.data {
		if u.data[i].ID == id {

			// kembalikan data user ketika sudah ditemukan
			return u.data[i], nil
		}
	}

	// kalau sampai di sini, harusnya user dengan id yang
	// diminta tidak ditemukan 😃
	return model.User{}, ErrUserNotFound
}

// method untuk mencari user berdasarkan username
// mengembalikan nilai (model.User, error)

func (u *User) FindUserByUsername(username string) (model.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	// loop sampai mendapatkan username yang sesuai
	for i := range u.data {
		if u.data[i].Username == username {

			// kembalikan data user ketika sudah ditemukan
			return u.data[i], nil
		}
	}

	// kalau sampai di sini, harusnya user dengan username
	// yang diminta tidak ditemukan 😃
	return model.User{}, ErrUserNotFound
}

// method untuk update data user

func (u *User) Update(id int, p UserUpdatePayload) error {

	// cek semua untuk ketika semua kosong, langsung return
	// kesalahan tanpa perlu mengeksekusi langkah selanjutnya

	if p.Password == nil && p.Username == nil && p.FirstName == nil && p.LastName == nil {
		return ErrNothingUpdate
	}

	// kunci dan buka ketika selesai digunakan, pakai defer
	// karena katanya sudah pasti ditutup ketika selesai

	u.mu.Lock()
	defer u.mu.Unlock()

	for i := range u.data {
		if u.data[i].ID == id {

			if p.Password != nil {
				u.data[i].Password = *p.Password
			}

			if p.Username != nil {
				u.data[i].Username = *p.Username
			}

			if p.FirstName != nil {
				u.data[i].Profile.FirstName = *p.FirstName
			}

			if p.LastName != nil {
				u.data[i].Profile.LastName = *p.LastName
			}

			return nil
		}
	}

	return ErrUserNotFound
}

// method untuk menghapus user
func (u *User) Delete(id int) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	for i := range u.data {
		if u.data[i].ID == id {
			u.data = append(u.data[:i], u.data[i+1:]...)
			return nil
		}
	}

	return ErrUserNotFound
}

// method untuk generate id unik
// meniru auto increment (mungkin)

func (u *User) genID() int {

	u.mu.Lock()
	defer u.mu.Unlock()

	// cek jika data kosong, harusnya langsung return id 1 saja
	if len(u.data) == 0 {
		return 1
	}

	// inisialisasi id dan cari id dengan nilai terbesar
	maxID := 0
	for i := range u.data {
		if u.data[i].ID > maxID {
			maxID = u.data[i].ID
		}
	}

	// kembalikan id dan tambahkan 1, karena id masih dimiliki
	// user lain yang mempunyai nilai paling besar.
	return maxID + 1
}
