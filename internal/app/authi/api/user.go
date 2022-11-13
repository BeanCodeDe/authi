package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/core"
	"github.com/BeanCodeDe/authi/internal/app/authi/errormessages"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

const UserRootPath = "/user"

type (
	user interface {
		mapToUserCore() *core.UserCore
	}
	userCreateDTO struct {
		ID       uuid.UUID `json:"id" validate:"required"`
		Password string    `json:"password" validate:"required"`
	}

	userLoginDTO struct {
		ID       uuid.UUID `json:"id" validate:"required"`
		Password string    `json:"password" validate:"required"`
	}

	userResponseDTO struct {
		ID        uuid.UUID `json:"id"`
		CreatedOn time.Time `json:"created_on"`
		LastLogin time.Time `json:"last_login"`
	}
)

func InitUserInterface(group *echo.Group) {
	group.POST("", createUserId)
	group.PUT("", create)
}

func createUserId(context echo.Context) error {
	log.Debug("Create User")
	return context.String(http.StatusOK, uuid.NewString())
}

func create(context echo.Context) error {
	log.Debugf("Create user")
	userCore, err := bind(context, new(userCreateDTO))
	if err != nil {
		log.Warnf("Error while binding user: %v", err)
		return echo.ErrBadRequest
	}
	if err := userCore.Create(); err != nil {
		if errors.Is(err, errormessages.UserAlreadyExists) {
			log.Warn("User with id %s already exists", userCore.ID)
			return context.NoContent(http.StatusConflict)
		}
		return err
	}

	log.Debugf("Created user with id %s", userCore.ID)
	return context.NoContent(http.StatusCreated)
}

func (user *userCreateDTO) mapToUserCore() *core.UserCore {
	return &core.UserCore{Password: user.Password}
}

func (user *userLoginDTO) mapToUserCore() *core.UserCore {
	return &core.UserCore{ID: user.ID, Password: user.Password}
}

func mapToUserResponseDTO(user *core.UserCore) *userResponseDTO {
	return &userResponseDTO{ID: user.ID, CreatedOn: user.CreatedOn, LastLogin: user.LastLogin}
}

func bind(context echo.Context, toBindUser user) (*core.UserCore, error) {
	log.Debugf("Bind context to user %v", context)
	if err := context.Bind(toBindUser); err != nil {
		log.Warnf("Could not bind user, %v", err)
		return nil, echo.ErrBadRequest
	}
	log.Debugf("User bind %v", toBindUser)
	if err := context.Validate(toBindUser); err != nil {
		log.Warnf("Could not validate user, %v", err)
		return nil, echo.ErrBadRequest
	}
	return toBindUser.mapToUserCore(), nil
}
