package models

type CreatHotelReq struct {
	HotelName   string                 `json:"hotel_name"`
	Description string                 `json:"description"`
	Address     string                 `json:"address"`
	ImageUrl    string                 `json:"image_url"`
	NumOfRooms  int64                  `json:"num_of_rooms"`
	Images      []*CreateHotelImageReq `json:"images"`
}

type CreateHotelImageReq struct {
	ImageUrl       string `json:"imagee_url"`
	SequenceNumber int64  `json:"sequence_number"`
}
