package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/SaidovZohid/hotel-project/api/models"
	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router /bookings/ [post]
// @Summary Create booking
// @Description Create booking
// @Tags booking
// @Accept json
// @Produce json
// @Param "data" body models.CreateOrUpdateBooking true "Data"
// @Success 200 {object} models.ResponseId
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) CreateBooking(ctx *gin.Context) {
	var (
		req models.CreateOrUpdateBooking
	)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	check_in, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	check_out, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	booking_id, err := h.strg.Booking().Create(&repo.Booking{
		CheckIn:  check_in,
		CheckOut: check_out,
		HotelID:  req.HotelID,
		RoomID:   req.RoomID,
		UserID:   payload.UserID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, models.ResponseId{
		Message: booking_id,
	})
}

// @Router /bookings/{id} [get]
// @Summary Get booking
// @Description Get booking
// @Tags booking
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.GetBooking
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetBooking(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	booking, err := h.strg.Booking().Get(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, models.GetBooking{
		ID:       booking.ID,
		CheckIn:  booking.CheckIn,
		CheckOut: booking.CheckOut,
		HotelID:  booking.HotelID,
		RoomID:   booking.RoomID,
		UserID:   booking.UserID,
		BookedAt: booking.BookedAt,
	})
}

// @Security ApiKeyAuth
// @Router /bookings/{id} [put]
// @Summary Update booking
// @Description Update booking
// @Tags booking
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param data body models.CreateOrUpdateBooking true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) UpdateBooking(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	var (
		req models.CreateOrUpdateBooking
	)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	check_in, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	check_out, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = h.strg.Booking().Update(&repo.Booking{
		ID:       id,
		CheckIn:  check_in,
		CheckOut: check_out,
		HotelID:  req.HotelID,
		RoomID:   req.RoomID,
		UserID:   payload.UserID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, models.ResponseOK{
		Message: "Succesfully updated!",
	})
}

// @Security ApiKeyAuth
// @Router /bookings/{id} [delete]
// @Summary Delete booking
// @Description Delete booking
// @Tags booking
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) DeleteBooking(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	err = h.strg.Booking().Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, models.ResponseOK{
		Message: "Succesfully deleted!",
	})
}

// @Security ApiKeyAuth
// @Router /bookings [get]
// @Summary Get All bookings for superadmin
// @Description Get All bookings  for superadmin
// @Tags booking
// @Accept json
// @Produce json
// @Param param query models.GetAllParams false "Params"
// @Success 200 {object} models.GetAllBookings
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetAllBooking(ctx *gin.Context) {
	params, err := validateRoomParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	if payload.UserType != "superuser" {
		ctx.JSON(http.StatusForbidden, errResponse(ErrForbidden))
		return
	}

	bookings, err := h.strg.Booking().GetAll(&repo.GetAllBookingsParams{
		Limit:  params.Limit,
		Page:   params.Page,
		Search: params.Search,
		SortBy: params.SortBy,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, getBookings(bookings))
}

// @Security ApiKeyAuth
// @Router /bookings/hotel/{id} [get]
// @Summary Get All bookings for superadmin and manager
// @Description Get All bookings  for superadmin and manager
// @Tags booking
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param param query models.GetAllParams false "Params"
// @Success 200 {object} models.GetAllBookings
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetAllHotelsBooking(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	params, err := validateRoomParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	bookings, err := h.strg.Booking().GetAllHotelsBookings(&repo.GetAllHotelsBookingsParams{
		HotelID: id,
		Limit:  params.Limit,
		Page:   params.Page,
		Search: params.Search,
		SortBy: params.SortBy,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, getBookings(bookings))
}