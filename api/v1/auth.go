package v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SaidovZohid/hotel-project/api/models"
	"github.com/SaidovZohid/hotel-project/pkg/utils"
	"github.com/SaidovZohid/hotel-project/storage/postgres"
	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/gin-gonic/gin"
)

// @Router /auth/register [post]
// @Summary Register user
// @Description Register user
// @Tags register
// @Accept json
// @Produce json
// @Param data body models.RegisterRequest true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) Register(ctx *gin.Context) {
	var (
		req models.RegisterRequest
	)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	_, err := h.strg.User().GetByEmail(req.Email)
	if !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusBadRequest, errResponse(ErrEmailExists))
		return
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	user := repo.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Type:      postgres.UserType,
		Password:  hashedPassword,
	}
	userData, err := json.Marshal(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = h.inMemory.Set(RegisterInfo+req.Email, string(userData), 10*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	go func() {
		err := h.sendVerificationCode(RegisterCodeKey, req.Email)
		if err != nil {
			log.Printf("failed to send verification code: %v\n", err)
		}
	}()

	ctx.JSON(http.StatusCreated, models.ResponseOK{
		Message: "Verification code has been sent!",
	})
}

// @Router /auth/verify [post]
// @Summary Veriy your account with code which we sent
// @Description Veriy your account with code which we sent
// @Tags register
// @Accept json
// @Produce json
// @Param data body models.VerifyRequest true "Data"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) Verify(ctx *gin.Context) {
	var (
		req models.VerifyRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	userData, err := h.inMemory.Get(RegisterInfo + req.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errResponse(err))
		return
	}

	var user repo.User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errResponse(err))
		return
	}

	code, err := h.inMemory.Get(RegisterCodeKey + req.Email)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errResponse(ErrCodeExpired))
		return
	}

	if req.Code != code {
		ctx.JSON(http.StatusForbidden, errResponse(ErrIncorrectCode))
		return
	}

	result, err := h.strg.User().Create(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	token, _, err := utils.CreateToken(h.cfg, &utils.TokenParams{
		UserID:   result,
		Email:    req.Email,
		UserType: postgres.UserType,
		Duration: time.Hour * 2,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.AuthResponse{
		Id:          result,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Type:        user.Type,
		AccessToken: token,
	})
}

// @Router /auth/login [post]
// @Summary Login User
// @Description Login User
// @Tags register
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) Login(ctx *gin.Context) {
	var (
		req *models.LoginRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := h.strg.User().GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusForbidden, errResponse(ErrWrongEmailOrPassword))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errResponse(ErrWrongEmailOrPassword))
		return
	}

	token, _, err := utils.CreateToken(h.cfg, &utils.TokenParams{
		UserID:   user.ID,
		Email:    req.Email,
		UserType: postgres.UserType,
		Duration: time.Hour * 24 * 30,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.AuthResponse{
		Id:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Type:        user.Type,
		AccessToken: token,
	})
}

// @Router /auth/forgot-password [post]
// @Summary GetAllParams  password
// @Description Forgot  password
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.ForgotPasswordRequest true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) ForgotPassword(ctx *gin.Context) {
	var (
		req models.ForgotPasswordRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	_, err := h.strg.User().GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	go func() {
		err := h.sendVerificationCode(RegisterCodeKey, req.Email)
		if err != nil {
			fmt.Printf("failed to send verification code: %v", err)
		}
	}()

	ctx.JSON(http.StatusCreated, models.ResponseOK{
		Message: "Validation code has been sent",
	})
}

// @Router /auth/verify-forgot-password [post]
// @Summary Verify forgot password
// @Description Verify forgot password
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.VerifyRequest true "Data"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) VerifyForgotPassword(ctx *gin.Context) {
	var (
		req *models.VerifyRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	code, err := h.inMemory.Get(RegisterCodeKey + req.Email)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errResponse(ErrCodeExpired))
		return
	}

	if req.Code != code {
		ctx.JSON(http.StatusForbidden, errResponse(ErrIncorrectCode))
		return
	}

	result, err := h.strg.User().GetByEmail(req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	token, _, err := utils.CreateToken(h.cfg, &utils.TokenParams{
		UserID:   result.ID,
		Email:    result.Email,
		Duration: time.Minute * 30,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, models.AuthResponse{
		Id:          result.ID,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		Type:        result.Type,
		AccessToken: token,
	})
}

// @Security ApiKeyAuth
// @Router /auth/update-password [post]
// @Summary Update password
// @Description Update password
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.UpdatePasswordRequest true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ResponseError
func (h *handlerV1) UpdatePassword(ctx *gin.Context) {
	var (
		req models.UpdatePasswordRequest
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

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = h.strg.User().UpdatePassword(&repo.UpdatePassword{
		UserId:   payload.UserID,
		Password: hashedPassword,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, models.ResponseOK{
		Message: "Password has been updated!",
	})
}
