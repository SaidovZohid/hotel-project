package postgres_test

import (
	"testing"
	"time"

	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/stretchr/testify/require"
)

func createBooking(t *testing.T) int64 {
	user_id := createUser(t)
	hotel_id := createHotel(t)
	room_id := createRoom(t)
	booking_id, err := dbManager.Booking().Create(&repo.Booking{
		CheckIn:  time.Now(),
		CheckOut: time.Now().Add(time.Hour * 24 * 2),
		HotelID:  hotel_id,
		RoomID:   room_id,
		UserID:   user_id,
	})
	require.NoError(t, err)
	return booking_id
}

func deleteBooking(t *testing.T, booking_id int64) {
	err := dbManager.Booking().Delete(booking_id)
	require.NoError(t, err)
}

func TestCreateBooking(t *testing.T) {
	booking_id := createBooking(t)
	deleteBooking(t, booking_id)
}

func TestGetBooking(t *testing.T) {
	booking_id := createBooking(t)
	booking, err := dbManager.Booking().Get(booking_id)
	require.NoError(t, err)
	require.NotEmpty(t, booking)
	deleteBooking(t, booking_id)
}

func TestUpdateBooking(t *testing.T) {
	booking_id := createBooking(t)
	room_id := createRoom(t)
	err := dbManager.Booking().Update(&repo.Booking{
		ID:       booking_id,
		CheckIn:  time.Now().Add(time.Hour * 24),
		CheckOut: time.Now().Add(time.Hour * 24 * 3),
		RoomID:   room_id,
	})
	require.NoError(t, err)
	deleteBooking(t, booking_id)
	deleteRoom(t, room_id)
}

func TestDeleteBooking(t *testing.T) {
	booking_id := createBooking(t)
	deleteBooking(t, booking_id)
}

func TestGetAllBookings(t *testing.T) {
	booking_id := createBooking(t)
	bookings, err := dbManager.Booking().GetAll(&repo.GetAllBookingsParams{
		Limit: 10,
		Page:  1,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(bookings.Bookings), 1)
	deleteBooking(t, booking_id)
}
