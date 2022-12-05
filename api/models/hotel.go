package models

type CreatOrUpdateHotelReq struct {
	HotelName   string                 `json:"hotel_name" binding:"required"`
	Description string                 `json:"description" binding:"required"`
	Address     string                 `json:"address" binding:"required"`
	ImageUrl    string                 `json:"image_url" binding:"required"`
	NumOfRooms  int64                  `json:"num_of_rooms" binding:"required"`
	Images      []*CreateHotelImageReq `json:"images" binding:"required"`
}

type CreateHotelImageReq struct {
	ImageUrl       string `json:"imagee_url" binding:"required"`
	SequenceNumber int64  `json:"sequence_number" binding:"required"`
}

type GetIdAndToken struct {
	ID          int64  `json:"id"`
	AccessToken string `json:"access_token"`
}

type GetHotelInfo struct {
	ID          int64         `json:"id"`
	HotelName   string        `json:"hotel_name"`
	Description string        `json:"description"`
	Address     string        `json:"address"`
	ImageUrl    string        `json:"image_url"`
	NumOfRooms  int64         `json:"num_of_rooms"`
	ManagerID   int64         `json:"manager_id"`
	Images      []*HotelImage `json:"images"`
}

type HotelImage struct {
	ID             int64  `json:"id"`
	HotelID        int64  `json:"hotel_id"`
	ImageUrl       string `json:"image_url"`
	SequenceNumber int64  `json:"sequence_number"`
}

type GetAllHotelsParams struct {
	Limit      int64  `json:"limit" binding:"required" default:"10"`
	Page       int64  `json:"page" binding:"required" default:"1"`
	Search     string `json:"search"`
	NumOfRooms int64  `json:"num_of_rooms"`
}

type GetAllHotels struct {
	Hotels []*GetHotelInfo `json:"hotels"`
	Count  int64           `json:"count"`
}
