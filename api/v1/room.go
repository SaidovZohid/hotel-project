package v1

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/SaidovZohid/hotel-project/api/models"
	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router /rooms [post]
// @Summary Create room
// @Description Create room
// @Tags room
// @Accept json
// @Produce json
// @Param data body models.CreateOrUpdateRoom true "Data"
// @Success 200 {object} models.ResponseId
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) CreateRoom(ctx *gin.Context) {
	var (
		req models.CreateOrUpdateRoom
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		// i should check id:
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	hotel, err := h.strg.Hotel().GetByManagerID(payload.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusForbidden, errResponse(err))
		return
	}
	log.Println(hotel.ManagerID, payload.UserID)
	switch payload.UserType {
	case "manager":
		if payload.UserID != hotel.ManagerID {
			ctx.JSON(http.StatusForbidden, errResponse(ErrForbidden))
			return
		}
	case "user":
		ctx.JSON(http.StatusForbidden, errResponse(ErrForbidden))
		return
	}

	if req.HotelID != hotel.ID {
		ctx.JSON(http.StatusForbidden, errResponse(errors.New("hotel id is not true")))
		return
	}

	room := repo.Room{
		RoomNumber:    req.RoomNumber,
		HotelID:       req.HotelID,
		Type:          req.Type,
		Description:   req.Description,
		PricePerNight: req.PricePerNight,
	}

	room_id, err := h.strg.Room().Create(&room)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseId{
		Message: room_id,
	})
}

// @Router /rooms/{id} [get]
// @Summary Get room
// @Description Get room
// @Tags room
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.GetRoomInfo
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetRoom(ctx *gin.Context) {
	room_id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	room, err := h.strg.Room().Get(room_id)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusForbidden, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.GetRoomInfo{
		ID:            room.ID,
		RoomNumber:    room.RoomNumber,
		HotelID:       room.HotelID,
		Type:          room.Type,
		Description:   room.Description,
		PricePerNight: room.PricePerNight,
		Status:        room.Status,
	})
}

// @Security ApiKeyAuth
// @Router /rooms/{id} [put]
// @Summary Update room
// @Description Update room
// @Tags room
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param data body models.CreateOrUpdateRoom true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) UpdateRoom(ctx *gin.Context) {
	room_id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	var (
		req models.CreateOrUpdateRoom
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
	hotel, err := h.strg.Hotel().GetByManagerID(payload.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusForbidden, errResponse(err))
		return
	}

	switch payload.UserType {
	case "manager":
		if payload.UserID != hotel.ManagerID {
			ctx.JSON(http.StatusForbidden, errResponse(ErrForbidden))
			return
		}
	case "user":
		ctx.JSON(http.StatusForbidden, errResponse(ErrForbidden))
		return
	}

	if req.HotelID != hotel.ID {
		ctx.JSON(http.StatusForbidden, errResponse(errors.New("hotel id is not true")))
		return
	}

	room := repo.Room{
		ID:            room_id,
		RoomNumber:    req.RoomNumber,
		Type:          req.Type,
		Description:   req.Description,
		PricePerNight: req.PricePerNight,
	}

	err = h.strg.Room().Update(&room)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseOK{
		Message: "Succesfully Updated!",
	})
}

// @Security ApiKeyAuth
// @Router /rooms/{id} [delete]
// @Summary Delete room
// @Description Delete room
// @Tags room
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) DeleteRoom(ctx *gin.Context) {
	room_id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	hotel, err := h.strg.Hotel().GetByManagerID(payload.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusForbidden, errResponse(err))
		return
	}

	switch payload.UserType {
	case "manager":
		if payload.UserID != hotel.ManagerID {
			ctx.JSON(http.StatusForbidden, errResponse(ErrForbidden))
			return
		}
	case "user":
		ctx.JSON(http.StatusForbidden, errResponse(ErrForbidden))
		return
	}

	err = h.strg.Room().Delete(room_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseOK{
		Message: "Succesfully Deleted!",
	})
}

// @Router /rooms/available/{id} [get]
// @Summary Get All rooms
// @Description Get All rooms
// @Tags room
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param times query models.GetAllRoomsDates true "Times"
// @Success 200 {object} models.GetAllRooms
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetAllHotelRoomsAvailable(ctx *gin.Context) {
	hotel_id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	check_in, err := time.Parse("2006-01-02", ctx.Query("check_in"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	check_out, err := time.Parse("2006-01-02", ctx.Query("check_out"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	rooms, err := h.strg.Room().GetAllHotelRoomsAvailable(&repo.GetAllRoomsDates{
		CheckIn:  check_in,
		CheckOut: check_out,
		HotelId:  hotel_id,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, getRooms(rooms))
}

// @Router /rooms [get]
// @Summary Get All rooms
// @Description Get All rooms
// @Tags room
// @Accept json
// @Produce json
// @Param params query models.GetAllParams false "Param"
// @Success 200 {object} models.GetAllRooms
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetAllRooms(ctx *gin.Context) {
	params, err := validateRoomParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	rooms, err := h.strg.Room().GetAll(&repo.GetAllParams{
		Limit:  params.Limit,
		Page:   params.Page,
		Search: params.Search,
		SortBy: params.SortBy,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, getRooms(rooms))
}
