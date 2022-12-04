package postgres

import (
	"fmt"

	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/jmoiron/sqlx"
)

type roomRepo struct {
	db *sqlx.DB
}

func NewRoomStorage(db *sqlx.DB) repo.RoomStorageI {
	return &roomRepo{
		db: db,
	}
}

// This function for Creating room. It takes Room struct and returns the created room id and nil, if there is error it will return 0 and error
func (ud *roomRepo) Create(r *repo.Room) (int64, error) {
	tr, err := ud.db.Begin()
	defer tr.Rollback() // if i don't write tr.Commit() it will automatically rollback before canceling function
	if err != nil {
		return 0, err
	}
	query := `
		INSERT INTO rooms (
			room_number,
			hotel_id,
			type,
			description,
			price_per_night
		) VALUES($1, $2, $3, $4, $5)
		RETURNING id
	`
	err = tr.QueryRow(
		query,
		r.RoomNumber,
		r.HotelID,
		r.Type,
		r.Description,
		r.PricePerNight,
	).Scan(&r.ID)
	if err != nil {
		return 0, err
	}

	err = tr.Commit()
	if err != nil {
		return 0, err
	}
	return r.ID, nil
}

// This function for Updating Room. It takes Room struct and returns the nil, if there is error it will return error
func (ud *roomRepo) Update(r *repo.Room) error {
	tr, err := ud.db.Begin()
	defer tr.Rollback() // if i don't write tr.Commit() or while commiting takes error, it will automatically rollback before canceling function
	if err != nil {
		return err
	}
	query := `
		UPDATE rooms SET
			room_number = $1,
			type = $2,
			description = $3,
			price_per_night = $4
		WHERE id = $5
	`
	_, err = tr.Exec(
		query,
		r.RoomNumber,
		r.Type,
		r.Description,
		r.PricePerNight,
		r.ID,
	)
	if err != nil {
		return err
	}

	err = tr.Commit()
	if err != nil {
		return err
	}
	return nil
}

// This function for getting Room info. It takes room id and returns the Room struct and nil, if there is error it will return nil and error
func (ud *roomRepo) Get(room_id int64) (*repo.Room, error) {
	var r repo.Room
	query := `
		SELECT 
			id,
			room_number,
			hotel_id,
			type,
			description,
			price_per_night,
			status
		FROM rooms WHERE id = $1
	`
	err := ud.db.QueryRow(
		query,
		room_id,
	).Scan(
		&r.ID,
		&r.RoomNumber,
		&r.HotelID,
		&r.Type,
		&r.Description,
		&r.PricePerNight,
		&r.Status,
	)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// This function for deleting room. It takes room id and returns nil, if there is error it will error
func (ud *roomRepo) Delete(room_id int64) error {
	query := "DELETE FROM rooms WHERE id = $1"
	_, err := ud.db.Exec(
		query,
		room_id,
	)
	if err != nil {
		return err
	}

	return nil
}

// This function for getting all rooms info. It takes params and returns all rooms info and nil values, if it has an error it will return nil and error values
func (ud *roomRepo) GetAll(params *repo.GetAllRoomsParams) (*repo.GetAllRooms, error) {
	var res repo.GetAllRooms
	res.Rooms = make([]*repo.Room, 0)
	offset := (params.Page - 1) *params.Limit
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)
	orderBy := " ORDER BY price_per_night DESC " 
	filter := " WHERE status = true "
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(`
			AND type ILIKE '%s' 
			OR description ILIKE '%s' 
		`, str, str)
	}
	if params.SortBy != "" {
		orderBy = fmt.Sprintf(" ORDER BY price_per_night %s", params.SortBy)
	}
	query := `
		SELECT 
			id,
			room_number,
			hotel_id,
			type,
			description,
			price_per_night,
			status
		FROM rooms
	` + filter + orderBy + limit
	rows, err := ud.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r repo.Room
		err := rows.Scan(
			&r.ID,
			&r.RoomNumber,
			&r.HotelID,
			&r.Type,
			&r.Description,
			&r.PricePerNight,
			&r.Status,
		)
		if err != nil {
			return nil, err
		}
		res.Rooms = append(res.Rooms, &r)
	}
	queryCount := "SELECT count(1) FROM rooms " + filter
	err = ud.db.QueryRow(queryCount).Scan(&res.Count)
	if err != nil {
		return nil, err
	} 
	
	return &res, nil
}