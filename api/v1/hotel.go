package v1

import (
	"net/http"

	"github.com/SaidovZohid/hotel-project/api/models"
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
// @Param data body models.CreatHotelReq true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) CreateHotel(ctx *gin.Context) {
	var (
		req models.CreatHotelReq
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

	ctx.JSON(http.StatusCreated, models.ResponseOK{
		Message: hotel_id,
	})
}