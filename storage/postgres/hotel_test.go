package postgres_test

import (
	"testing"

	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createHotel(t *testing.T) int64 {
	manager_id := createUser(t)
	hotel_id, err := dbManager.Hotel().Create(&repo.Hotel{
		HotelName:   faker.Name(),
		Description: faker.Sentence(),
		Address:     faker.MacAddress(),
		ImageUrl:    faker.Word(),
		NumOfRooms:  12,
		ManagerID:   manager_id,
		Images: []*repo.HotelImage{
			{
				ImageUrl:       faker.Word(),
				SequenceNumber: 1,
			},
			{
				ImageUrl:       faker.Word(),
				SequenceNumber: 2,
			},
		},
	})
	require.NoError(t, err)
	return hotel_id
}

func deleteHotel(t *testing.T, hotel_id int64) {
	err := dbManager.Hotel().Delete(hotel_id)
	require.NoError(t, err)
}

func TestCreateHotel(t *testing.T) {
	hotel_id := createHotel(t)
	deleteHotel(t, hotel_id)
}

func TestGetHotel(t *testing.T) {
	hotel_id := createHotel(t)
	hotel, err := dbManager.Hotel().Get(hotel_id)
	require.NoError(t, err)
	require.NotEmpty(t, hotel)
	deleteHotel(t, hotel_id)
}

func TestUpdateHotel(t *testing.T) {
	hotel_id := createHotel(t)
	err := dbManager.Hotel().Update(&repo.Hotel{
		ID:          hotel_id,
		HotelName:   faker.Name(),
		Description: faker.Sentence(),
		Address:     faker.MacAddress(),
		ImageUrl:    faker.Word(),
		NumOfRooms:  12,
		Images: []*repo.HotelImage{
			{
				HotelID:        hotel_id,
				ImageUrl:       faker.Word(),
				SequenceNumber: 1,
			},
			{
				HotelID:        hotel_id,
				ImageUrl:       faker.Word(),
				SequenceNumber: 2,
			},
			{
				HotelID:        hotel_id,
				ImageUrl:       faker.Word(),
				SequenceNumber: 3,
			},
		},
	})
	require.NoError(t, err)
	deleteHotel(t, hotel_id)
}

func TestDeleteHotel(t *testing.T) {
	hotel_id := createHotel(t)
	err := dbManager.Hotel().Delete(hotel_id)
	require.NoError(t, err)
}

func TestGetAllHotel(t *testing.T) {
	hotel_id := createHotel(t)
	hotels, err := dbManager.Hotel().GetAll(&repo.GetAllHotelsParams{
		Limit: 10,
		Page: 1,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(hotels.Hotels), 1)
	deleteHotel(t, hotel_id)
}
