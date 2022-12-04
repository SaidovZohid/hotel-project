package storage

import (
	"github.com/SaidovZohid/hotel-project/storage/postgres"
	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	User() repo.UserStorageI
	Hotel() repo.HotelStorageI
	Room() repo.RoomStorageI
	Booking() repo.BookingStorageI
}

type StoragePg struct {
	userRepo repo.UserStorageI
	hotelRepo repo.HotelStorageI
	roomRepo repo.RoomStorageI
	bookRepo repo.BookingStorageI
}

func NewStorage(db *sqlx.DB) StorageI {
	return &StoragePg{
		userRepo: postgres.NewUserStorage(db),
		hotelRepo: postgres.NewHotelStorage(db),
		roomRepo: postgres.NewRoomStorage(db),
		bookRepo: postgres.NewBookingStorage(db),
	}
}

func (s *StoragePg) User() repo.UserStorageI {
	return s.userRepo
}

func (s *StoragePg) Hotel() repo.HotelStorageI {
	return s.hotelRepo
}

func (s *StoragePg) Room() repo.RoomStorageI {
	return s.roomRepo
}

func (s *StoragePg) Booking() repo.BookingStorageI {
	return s.bookRepo
}