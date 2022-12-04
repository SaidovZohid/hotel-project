package postgres_test

import (
	"testing"

	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createRoom(t *testing.T) int64 {
	hotel_id := createHotel(t)
	room_id, err := dbManager.Room().Create(&repo.Room{
		RoomNumber:    12,
		HotelID:       hotel_id,
		Type:          "single",
		Description:   faker.Sentence(),
		PricePerNight: 1000.34,
	})
	require.NoError(t, err)
	return room_id
}

func deleteRoom(t *testing.T, room_id int64) {
	err := dbManager.Room().Delete(room_id)
	require.NoError(t, err)
}

func TestCreateRoom(t *testing.T) {
	room_id := createRoom(t)
	deleteRoom(t, room_id)
}

func TestGetRoom(t *testing.T) {
	room_id := createRoom(t)
	room, err := dbManager.Room().Get(room_id)
	require.NoError(t, err)
	require.NotEmpty(t, room)
	deleteRoom(t, room_id)
}

func TestUpdateRoom(t *testing.T) {
	room_id := createRoom(t)
	err := dbManager.Room().Update(&repo.Room{
		ID:            room_id,
		RoomNumber:    10,
		Type:          "family",
		Description:   faker.Word(),
		PricePerNight: 100.90,
	})
	require.NoError(t, err)
	deleteRoom(t, room_id)
}

func TestDeleteRoom(t *testing.T) {
	room_id := createRoom(t)
	deleteRoom(t, room_id)
}

func TestGetAllRooms(t *testing.T) {
	room_id := createRoom(t)
	rooms, err := dbManager.Room().GetAll(&repo.GetAllRoomsParams{
		Limit: 10,
		Page:  1,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(rooms.Rooms), 1)
	deleteRoom(t, room_id)
}
