package repo

import "time"

// This is users table struct
type User struct {
	ID          int64
	FirstName   string
	LastName    string
	Email       string
	Password    string
	PhoneNumber *string
	Type        string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}

type UserStorageI interface {
	Create(u *User) (int64, error)
	ChangeTypeUser(user_type string, user_id int64) error
	Get(user_id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(u *User) (*time.Time, error)
	UpdatePassword(u *UpdatePassword) error
	Delete(user_id int64) error
	GetAll(params *GetAllUsersParams) (*GetAllUsers, error)
}

type GetAllUsersParams struct {
	Limit  int64
	Page   int64
	Search string
	SortBy string
}

type GetAllUsers struct {
	Users []*User
	Count int64
}

type UpdatePassword struct {
	UserId   int64
	Password string
}
