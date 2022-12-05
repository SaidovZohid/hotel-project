package models

import "time"

type CreateOrUpdateBooking struct {
	CheckIn  string `json:"check_in" binding:"required"`
	CheckOut string `json:"check_out" binding:"required"`
	HotelID  int64  `json:"hotel_id" binding:"required"`
	RoomID   int64  `json:"room_id" binding:"required"`
}

type GetBooking struct {
	ID       int64     `json:"id"`
	CheckIn  time.Time `json:"check_in"`
	CheckOut time.Time `json:"check_out"`
	HotelID  int64     `json:"hotel_id"`
	RoomID   int64     `json:"room_id"`
	UserID   int64     `json:"user_id"`
	BookedAt time.Time `json:"booked_at"`
}

type GetAllBookings struct {
	Bookings []*GetBooking `json:"bookings"`
	Count    int64         `json:"count"`
}
