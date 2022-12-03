package postgres

import (
	"fmt"
	"time"

	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/jmoiron/sqlx"
)

const (
	ManagerType   = "manager"
	SuperUserType = "superuser"
	UserType      = "user"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) repo.UserStorageI {
	return &userRepo{
		db: db,
	}
}

// This function for Creating user. It takes User struct and returns the created user id and nil, if there is error it will return 0 and error
func (ud *userRepo) Create(u *repo.User) (int64, error) {
	tr, err := ud.db.Begin()
	defer tr.Rollback() // if i don't write tr.Commit() it will automatically rollback before canceling function
	if err != nil {
		return 0, err
	}
	query := `
		INSERT INTO users (
			first_name,
			last_name,
			email,
			password,
			phone_number,
			type
		) VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err = tr.QueryRow(
		query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Password,
		u.PhoneNumber,
		u.Type,
	).Scan(&u.ID)
	if err != nil {
		return 0, err
	}

	err = tr.Commit()
	if err != nil {
		return 0, err
	}
	return u.ID, nil
}

// This function for Updating user. It takes User struct and returns the updated time and nil, if there is error it will return nil and error
func (ud *userRepo) Update(u *repo.User) (*time.Time, error) {
	tr, err := ud.db.Begin()
	defer tr.Rollback() // if i don't write tr.Commit() or while commiting takes error, it will automatically rollback before canceling function
	if err != nil {
		return nil, err
	}
	query := `
		UPDATE users SET
			first_name = $1,
			last_name = $2,
			email = $3,
			phone_number = $4,
			type = $5,
			updated_at = $6
		WHERE id = $7 and deleted_at IS NULL
		RETURNING updated_at
	`
	err = tr.QueryRow(
		query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.PhoneNumber,
		u.Type,
		time.Now(),
		u.ID,
	).Scan(&u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	err = tr.Commit()
	if err != nil {
		return nil, err
	}
	return u.UpdatedAt, nil
}

// This function for getting user info. It takes user id and returns the User struct and nil, if there is error it will return nil and error
func (ud *userRepo) Get(user_id int64) (*repo.User, error) {
	var u repo.User
	query := `
		SELECT 
			id,
			first_name,
			last_name,
			email,
			phone_number,
			type,
			created_at,
			updated_at
		FROM users WHERE id = $1
	`
	err := ud.db.QueryRow(
		query,
		user_id,
	).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.PhoneNumber,
		&u.Type,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// This function for deleting user. It takes user id and returns nil, if there is error it will error
func (ud *userRepo) Delete(user_id int64) error {
	query := `
		UPDATE users SET 
			deleted_at = $1
		WHERE id = $2
	`
	_, err := ud.db.Exec(
		query,
		time.Now(),
		user_id,
	)
	if err != nil {
		return err
	}

	return nil
}

// This function for getting all users info. It takes params and returns all users info and nil values, if it has an error it will return nil and error values
func (ud *userRepo) GetAll(params *repo.GetAllUsersParams) (*repo.GetAllUsers, error) {
	var res repo.GetAllUsers
	res.Users = make([]*repo.User, 0)
	offset := (params.Page - 1) *params.Limit
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)
	orderBy := " ORDER BY created_at DESC "
	filter := " WHERE deleted_at IS NULL"
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(`
			AND first_name ILIKE '%s' 
			OR last_name ILIKE '%s' 
			OR email ILIKE '%s'
			OR phone_number ILIKE '%s'
			OR type ILIKE '%s'
		`, str, str, str, str, str)
	}
	if params.SortBy != "" {
		orderBy = fmt.Sprintf(" ORDER BY created_at %s", params.SortBy)
	}
	query := `
		SELECT 
			id,
			first_name,
			last_name,
			email,
			phone_number,
			type,
			created_at,
			updated_at
		FROM users
	` + filter + orderBy + limit
	rows, err := ud.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var u repo.User
		err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.PhoneNumber,
			&u.Type,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		res.Users = append(res.Users, &u)
	}
	return &res, nil
}