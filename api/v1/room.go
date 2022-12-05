package v1

import (
	"net/http"

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

	// I should add payload and database hotel check:

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
