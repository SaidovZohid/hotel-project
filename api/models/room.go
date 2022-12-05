package models

type CreateOrUpdateRoom struct {
	RoomNumber    int64   `json:"room_number" binding:"required"`
	HotelID       int64   `json:"hotel_id" binding:"required"`
	Type          string  `json:"type" binding:"required,oneof=family single double" default:"single"`
	Description   string  `json:"description" binding:"required"`
	PricePerNight float64 `json:"price_per_night" binding:"required"`
}

type GetRoomInfo struct {
	ID            int64   `json:"id"`
	RoomNumber    int64   `json:"room_number"`
	HotelID       int64   `json:"hotel_id"`
	Type          string  `json:"type"`
	Description   string  `json:"description"`
	PricePerNight float64 `json:"price_per_night"`
	Status        bool    `json:"status"`
}

type GetAllParams struct {
	Limit  int64  `json:"limit" binding:"required" default:"10"`
	Page   int64  `json:"page" binding:"required" default:"1"`
	Search string `json:"search"`
	SortBy string `json:"sort_by" enums:"desc,asc" default:"desc"`
}

type GetAllRooms struct {
	Rooms []*GetRoomInfo `json:"rooms"`
	Count int64          `json:"count"`
}

type GetAllRoomsDates struct {
	CheckIn  string `json:"check_in" binding:"required"`
	CheckOut string `json:"check_out" binding:"required"`
}
