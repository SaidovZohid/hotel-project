package v1

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/SaidovZohid/hotel-project/api/models"
	"github.com/SaidovZohid/hotel-project/config"
	emailPkg "github.com/SaidovZohid/hotel-project/pkg/email"
	"github.com/SaidovZohid/hotel-project/pkg/utils"
	"github.com/SaidovZohid/hotel-project/storage"
	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/gin-gonic/gin"
)

const (
	RegisterCodeKey string = "register_code_"
	RegisterInfo    string = "register_info_"
)

var (
	ErrWrongEmailOrPassword = errors.New("wrong email or password")
	ErrUserNotVerifid       = errors.New("user not verified")
	ErrEmailExists          = errors.New("email is already exists")
	ErrIncorrectCode        = errors.New("incorrect verification code")
	ErrCodeExpired          = errors.New("verification code is expired")
	ErrForbidden            = errors.New("forbidden")
)

type handlerV1 struct {
	cfg      *config.Config
	strg     storage.StorageI
	inMemory storage.InMemoryStorageI
}

type HandlerV1 struct {
	Cfg      *config.Config
	Strg     storage.StorageI
	InMemory storage.InMemoryStorageI
}

func New(opt *HandlerV1) *handlerV1 {
	return &handlerV1{
		cfg:      opt.Cfg,
		strg:     opt.Strg,
		inMemory: opt.InMemory,
	}
}

func errResponse(err error) models.ResponseError {
	return models.ResponseError{
		Error: err.Error(),
	}
}

func (h *handlerV1) sendVerificationCode(key, email string) error {
	code, err := utils.GenerateRandomCode(6)
	if err != nil {
		return err
	}

	err = h.inMemory.Set(key+email, code, time.Minute*2)
	if err != nil {
		return err
	}
	err = emailPkg.SendEmail(h.cfg, &emailPkg.SendEmailRequest{
		To:      []string{email},
		Subject: "Verification Code",
		Body: map[string]string{
			"code": code,
		},
		Type: emailPkg.VerificationEmail,
	})

	if err != nil {
		log.Print("Failed to sent code to email")
	}
	return nil
}

func validateHotelParams(ctx *gin.Context) (*models.GetAllHotelsParams, error) {
	var (
		limit        int64 = 10
		page         int64 = 1
		num_of_rooms int64
		err          error
	)
	if ctx.Query("limit") != "" {
		limit, err = strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if ctx.Query("page") != "" {
		page, err = strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if ctx.Query("num_of_rooms") != "" {
		num_of_rooms, err = strconv.ParseInt(ctx.Query("num_of_rooms"), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return &models.GetAllHotelsParams{
		Limit:      limit,
		Page:       page,
		Search:     ctx.Query("search"),
		NumOfRooms: num_of_rooms,
	}, nil
}

func validateRoomParams(ctx *gin.Context) (*models.GetAllParams, error) {
	var (
		limit int64 = 10
		page  int64 = 1
		err   error
	)
	if ctx.Query("limit") != "" {
		limit, err = strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if ctx.Query("page") != "" {
		page, err = strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return &models.GetAllParams{
		Limit:  limit,
		Page:   page,
		Search: ctx.Query("search"),
		SortBy: ctx.Query("sort_by"),
	}, nil
}

func getHotels(hotels *repo.GetAllHotels) models.GetAllHotels {
	var res models.GetAllHotels
	res.Hotels = make([]*models.GetHotelInfo, 0)
	res.Count = hotels.Count
	for _, hotel := range hotels.Hotels {
		h := models.GetHotelInfo{
			ID:          hotel.ID,
			HotelName:   hotel.HotelName,
			Description: hotel.Description,
			Address:     hotel.Address,
			ImageUrl:    hotel.ImageUrl,
			NumOfRooms:  hotel.NumOfRooms,
			ManagerID:   hotel.ManagerID,
		}
		for _, image := range hotel.Images {
			i := models.HotelImage{
				ID:             image.ID,
				HotelID:        image.HotelID,
				ImageUrl:       image.ImageUrl,
				SequenceNumber: image.SequenceNumber,
			}
			h.Images = append(h.Images, &i)
		}
		res.Hotels = append(res.Hotels, &h)
	}
	return res
}

func getRooms(rooms *repo.GetAllRooms) models.GetAllRooms {
	var res models.GetAllRooms
	res.Rooms = make([]*models.GetRoomInfo, 0)
	res.Count = rooms.Count
	for _, room := range rooms.Rooms {
		r := models.GetRoomInfo{
			ID:            room.ID,
			RoomNumber:    room.RoomNumber,
			HotelID:       room.HotelID,
			Type:          room.Type,
			Description:   room.Description,
			PricePerNight: room.PricePerNight,
			Status:        room.Status,
		}
		res.Rooms = append(res.Rooms, &r)
	}
	return res
}

func getAllUsers(users *repo.GetAllUsers) models.GetAllUsers {
	var res models.GetAllUsers
	res.Users = make([]*models.GetUser, 0)
	res.Count = users.Count
	for _, user := range users.Users {
		r := models.GetUser{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Type:        user.Type,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		}
		res.Users = append(res.Users, &r)
	}
	return res
}

func getBookings(bookings *repo.GetAllBookings) models.GetAllBookings {
	var res models.GetAllBookings
	res.Bookings = make([]*models.GetBooking, 0)
	res.Count = bookings.Count
	for _, booking := range bookings.Bookings {
		b := models.GetBooking{
			ID:       booking.ID,
			CheckIn:  booking.CheckIn,
			CheckOut: booking.CheckOut,
			HotelID:  booking.HotelID,
			RoomID:   booking.RoomID,
			UserID:   booking.UserID,
			BookedAt: booking.BookedAt,
		}
		res.Bookings = append(res.Bookings, &b)
	}
	return res
}
