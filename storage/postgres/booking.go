package postgres

import (
	"fmt"

	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/jmoiron/sqlx"
)

type bookingRepo struct {
	db *sqlx.DB
}

func NewBookingStorage(db *sqlx.DB) repo.BookingStorageI {
	return &bookingRepo{
		db: db,
	}
}

func (bd *bookingRepo) Create(b *repo.Booking) (int64, error) {
	tr, err := bd.db.Begin()
	defer tr.Rollback()
	if err != nil {
		return 0, err
	}

	query := `
		INSERT INTO bookings (
			check_in,
			check_out,
			hotel_id,
			room_id,
			user_id
		) VALUES ($1, $2, $3, $4, $5)
		returning id
	`

	err = tr.QueryRow(
		query,
		b.CheckIn,
		b.CheckOut,
		b.HotelID,
		b.RoomID,
		b.UserID,
	).Scan(&b.ID)
	if err != nil {
		return 0, err
	}

	err = tr.Commit()
	if err != nil {
		return 0, err
	}

	return b.ID, nil
}

func (bd *bookingRepo) Get(booking_id int64) (*repo.Booking, error) {
	query := `
		SELECT
			id,
			check_in,
			check_out,
			hotel_id,
			room_id,
			user_id,
			booked_at
		FROM bookings WHERE id = $1
	`
	var res repo.Booking
	err := bd.db.QueryRow(
		query,
		booking_id,
	).Scan(
		&res.ID,
		&res.CheckIn,
		&res.CheckOut,
		&res.HotelID,
		&res.RoomID,
		&res.UserID,
		&res.BookedAt,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (bd *bookingRepo) Update(b *repo.Booking) error {
	query := `
		UPDATE bookings SET 
			check_in = $1,
			check_out = $2,
			room_id = $3
		WHERE id = $4
	`
	_, err := bd.db.Exec(
		query,
		b.CheckIn,
		b.CheckOut,
		b.RoomID,
		b.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (bd *bookingRepo) Delete(booking_id int64) error {
	query := " DELETE FROM bookings WHERE id = $1"
	_, err := bd.db.Exec(query, booking_id)
	if err != nil {
		return err
	}

	return nil
}

func (bd *bookingRepo) GetAll(params *repo.GetAllBookingsParams) (*repo.GetAllBookings, error) {
	var res repo.GetAllBookings
	res.Bookings = make([]*repo.Booking, 0)
	offset := (params.Page - 1) * params.Limit
	orderBy := " ORDER BY booked_at DESC"
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)
	if params.SortBy != "" {
		orderBy = " ORDER BY booked_at " + params.SortBy
	}
	query := `
		SELECT
			id,
			check_in,
			check_out,
			hotel_id,
			room_id,
			user_id,
			booked_at
		FROM bookings
	` + orderBy + limit
	rows, err := bd.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var b repo.Booking
		err := bd.db.QueryRow(
			query,
		).Scan(
			&b.ID,
			&b.CheckIn,
			&b.CheckOut,
			&b.HotelID,
			&b.RoomID,
			&b.UserID,
			&b.BookedAt,
		)
		if err != nil {
			return nil, err
		}
		res.Bookings = append(res.Bookings, &b)
	}
	queryCount := "SELECT count(1) FROM bookings"
	err = bd.db.QueryRow(queryCount).Scan(&res.Count)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

