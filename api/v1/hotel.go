package v1

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/SaidovZohid/hotel-project/api/models"
	"github.com/SaidovZohid/hotel-project/pkg/utils"
	"github.com/SaidovZohid/hotel-project/storage/postgres"
	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router /hotels [post]
// @Summary Create hotel
// @Description Create hotel
// @Tags hotel
// @Accept json
// @Produce json
// @Param data body models.CreatOrUpdateHotelReq true "Data"
// @Success 200 {object} models.GetIdAndToken
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) CreateHotel(ctx *gin.Context) {
	var (
		req models.CreatOrUpdateHotelReq
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
	res := repo.Hotel{
		HotelName:   req.HotelName,
		Description: req.Description,
		Address:     req.Address,
		ImageUrl:    req.ImageUrl,
		NumOfRooms:  req.NumOfRooms,
		ManagerID:   payload.UserID,
	}
	for _, image := range req.Images {
		i := repo.HotelImage{
			ImageUrl:       image.ImageUrl,
			SequenceNumber: image.SequenceNumber,
		}
		res.Images = append(res.Images, &i)
	}
	hotel_id, err := h.strg.Hotel().Create(&res)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = h.strg.User().ChangeTypeUser(postgres.ManagerType, payload.UserID)
	if err != nil {
		err := h.strg.Hotel().Delete(hotel_id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	token, _, err := utils.CreateToken(h.cfg, &utils.TokenParams{
		UserID:   payload.UserID,
		Email:    payload.Email,
		UserType: postgres.ManagerType,
		Duration: time.Hour * 24 * 30,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, models.GetIdAndToken{
		ID:          hotel_id,
		AccessToken: token,
	})
}

// @Router /hotels/{id} [get]
// @Summary Get hotel
// @Description Get hotel
// @Tags hotel
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.GetHotelInfo
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetHotel(ctx *gin.Context) {
	hotel_id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hotel, err := h.strg.Hotel().Get(hotel_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	res := models.GetHotelInfo{
		ID:          hotel.ID,
		HotelName:   hotel.HotelName,
		Description: hotel.Description,
		Address:     hotel.Address,
		ImageUrl:    hotel.ImageUrl,
		NumOfRooms:  hotel.NumOfRooms,
		ManagerID:   hotel.ManagerID,
	}
	for _, image := range hotel.Images {
		res.Images = append(res.Images, &models.HotelImage{
			ID:             image.ID,
			HotelID:        image.HotelID,
			ImageUrl:       image.ImageUrl,
			SequenceNumber: image.SequenceNumber,
		})
	}

	ctx.JSON(http.StatusOK, res)
}

// @Security ApiKeyAuth
// @Router /hotels/{id} [put]
// @Summary Update hotel
// @Description Update hotel
// @Tags hotel
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param data body models.CreatOrUpdateHotelReq true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) UpdateHotel(ctx *gin.Context) {
	hotel_id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hotel, err := h.strg.Hotel().Get(hotel_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errResponse(err))
		return
	}
	if payload.UserID != hotel.ManagerID {
		ctx.JSON(http.StatusForbidden, errResponse(err))
		return
	}

	var (
		req models.CreatOrUpdateHotelReq
	)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	var res repo.Hotel
	for _, image := range req.Images {
		i := repo.HotelImage{
			HotelID:        hotel_id,
			ImageUrl:       image.ImageUrl,
			SequenceNumber: image.SequenceNumber,
		}
		res.Images = append(res.Images, &i)
	}
	res = repo.Hotel{
		ID:          hotel_id,
		HotelName:   req.HotelName,
		Description: req.Description,
		Address:     req.Address,
		ImageUrl:    req.ImageUrl,
		NumOfRooms:  req.NumOfRooms,
		Images:      res.Images,
	}

	err = h.strg.Hotel().Update(&res)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseOK{
		Message: "Succesfuly updated!",
	})
}

// @Security ApiKeyAuth
// @Router /hotels/{id} [delete]
// @Summary Delete hotel
// @Description Delete hotel
// @Tags hotel
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.GetIdAndToken
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) DeleteHotel(ctx *gin.Context) {
	hotel_id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hotel, err := h.strg.Hotel().Get(hotel_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	if payload.UserID != hotel.ManagerID {
		ctx.JSON(http.StatusForbidden, errResponse(err))
		return
	}

	
	err = h.strg.Hotel().Delete(hotel_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = h.strg.User().ChangeTypeUser(postgres.UserType, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	token, _, err := utils.CreateToken(h.cfg, &utils.TokenParams{
		UserID:   payload.UserID,
		Email:    payload.Email,
		UserType: postgres.UserType,
		Duration: time.Hour * 24 * 30,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.GetIdAndToken{
		ID:          payload.UserID,
		AccessToken: token,
	})
}

// @Router /hotels [get]
// @Summary Get All hotel
// @Description Get All hotel
// @Tags hotel
// @Accept json
// @Produce json
// @Param param query models.GetAllHotelsParams false "Param"
// @Success 200 {object} models.GetAllHotels
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetAllHotels(ctx *gin.Context) {
	params, err := validateHotelParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	hotels, err := h.strg.Hotel().GetAll(&repo.GetAllHotelsParams{
		Limit:      params.Limit,
		Page:       params.Page,
		Search:     params.Search,
		NumOfRooms: params.NumOfRooms,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, getHotels(hotels))
}
