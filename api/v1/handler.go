package v1

import (
	"errors"
	"log"
	"time"

	"github.com/SaidovZohid/hotel-project/api/models"
	"github.com/SaidovZohid/hotel-project/config"
	emailPkg "github.com/SaidovZohid/hotel-project/pkg/email"
	"github.com/SaidovZohid/hotel-project/pkg/utils"
	"github.com/SaidovZohid/hotel-project/storage"
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
		cfg:  opt.Cfg,
		strg: opt.Strg,
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
