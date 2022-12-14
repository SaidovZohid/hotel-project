package repo

type Hotel struct {
	ID          int64
	HotelName   string
	Description string
	Address     string
	ImageUrl    string
	NumOfRooms  int64
	ManagerID   int64
	Images      []*HotelImage
}

type HotelStorageI interface {
	Create(h *Hotel) (int64, error)
	Get(hotel_id int64) (*Hotel, error)
	GetByManagerID(manager_id int64) (*Hotel, error)
	Update(h *Hotel) error
	Delete(hotel_id int64) error
	GetAll(params *GetAllHotelsParams) (*GetAllHotels, error)
}

type GetAllHotelsParams struct {
	Limit      int64
	Page       int64
	Search     string
	NumOfRooms int64
}

type GetAllHotels struct {
	Hotels []*Hotel
	Count  int64
}

type HotelImage struct {
	ID             int64
	HotelID        int64
	ImageUrl       string
	SequenceNumber int64
}