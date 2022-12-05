package repo

import "time"

// This is rooms table struct
type Room struct {
	ID            int64
	RoomNumber    int64
	HotelID       int64
	Type          string
	Description   string
	PricePerNight float64
	Status        bool
}

type RoomStorageI interface {
	Create(u *Room) (int64, error)
	Get(room_id int64) (*Room, error)
	GetAllHotelRoomsAvailable(params *GetAllRoomsDates) (*GetAllRooms, error)
	Update(u *Room) error
	Delete(room_id int64) error
	GetAll(params *GetAllRoomsParams) (*GetAllRooms, error)
}

type GetAllRoomsParams struct {
	Limit  int64
	Page   int64
	Search string
	SortBy string
}

type GetAllRooms struct {
	Rooms []*Room
	Count int64
}

type GetAllRoomsDates struct {
	HotelId  int64     `json:"hotel_id"`
	CheckIn  time.Time `json:"check_in"`
	CheckOut time.Time `json:"check_out"`
}
