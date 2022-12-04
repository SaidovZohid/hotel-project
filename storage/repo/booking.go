package repo

import "time"

type Booking struct {
	ID       int64
	CheckIn  time.Time
	CheckOut time.Time
	HotelID  int64
	RoomID   int64
	UserID   int64
	BookedAt time.Time
}

type BookingStorageI interface {
	Create(b *Booking) (int64, error)
	Get(booking_id int64) (*Booking, error)
	Update(b *Booking) error
	Delete(booking_id int64) error
	GetAll(params *GetAllBookingsParams) (*GetAllBookings, error)
}

type GetAllBookings struct {
	Bookings []*Booking
	Count    int64
}

type GetAllBookingsParams struct {
	Limit  int64
	Page   int64
	Search string
	SortBy string
}
