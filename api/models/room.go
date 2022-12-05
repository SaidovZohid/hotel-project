package models

type CreateOrUpdateRoom struct {
	RoomNumber    int64   `json:"room_number" binding:"required"`
	HotelID       int64   `json:"hotel_id" binding:"required"`
	Type          string  `json:"type" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	PricePerNight float64 `json:"price_per_night" binding:"required"`
}
