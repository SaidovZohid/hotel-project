package v1

import (
	"net/http"
	"strconv"

	"github.com/SaidovZohid/hotel-project/api/models"
	"github.com/SaidovZohid/hotel-project/pkg/utils"
	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Security ApiKeyAuth
// @Router /users [post]
// @Summary Create user
// @Description Create user
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.User true "User"
// @Success 200 {object} models.ResponseId
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) CreateUser(ctx *gin.Context) {
	var (
		req models.User
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

	if payload.UserType != "superuser" {
		ctx.JSON(http.StatusForbidden, errResponse(ErrForbidden))
		return
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	user_id, err := h.strg.User().Create(&repo.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    hashedPassword,
		PhoneNumber: req.PhoneNumber,
		Type:        req.Type,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, models.ResponseId{
		Message: user_id,
	})
}

// @Router /users/{id} [get]
// @Summary Get user
// @Description Get user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.GetUser
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := h.strg.User().Get(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, models.GetUser{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Type:        user.Type,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	})
}

// @Security ApiKeyAuth
// @Router /users/profile [get]
// @Summary Get user profile
// @Description Get user profile
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.GetUser
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetProfileUser(ctx *gin.Context) {
	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	user, err := h.strg.User().Get(payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, models.GetUser{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Type:        user.Type,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	})
}

// @Security ApiKeyAuth
// @Router /users/{id} [put]
// @Summary Update user
// @Description Update user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param user body models.UpdateUser true "User"
// @Success 200 {object} models.ResponseUpdatedAt
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	var (
		req models.UpdateUser
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
	if payload.UserType != "manager" {
		ctx.JSON(http.StatusForbidden, errResponse(ErrForbidden))
		return
	}

	updated_at, err := h.strg.User().Update(&repo.User{
		ID:          id,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Type:        req.Type,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, models.ResponseUpdatedAt{
		UpdatedAt: updated_at,
	})
}

// @Security ApiKeyAuth
// @Router /users/{id} [delete]
// @Summary Delete user
// @Description Delete user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	if payload.UserType != "manager" {
		ctx.JSON(http.StatusForbidden, errResponse(ErrForbidden))
		return
	}

	err = h.strg.User().Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, models.ResponseOK{
		Message: "Succesfully deleted user!",
	})
}

// @Router /users [get]
// @Summary Get all user
// @Description Get all user
// @Tags user
// @Accept json
// @Produce json
// @Param params query models.GetAllParams false "Param"
// @Success 200 {object} models.GetAllUsers
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) GetAllUser(ctx *gin.Context) {
	params, err := validateRoomParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	users, err := h.strg.User().GetAll(&repo.GetAllUsersParams{
		Limit:  params.Limit,
		Page:   params.Page,
		Search: params.Search,
		SortBy: params.SortBy,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, getAllUsers(users))
}
