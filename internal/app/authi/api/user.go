package api

import (
	"errors"
	"net/http"

	"github.com/BeanCodeDe/authi/internal/app/authi/core"
	"github.com/BeanCodeDe/authi/internal/app/authi/errormessages"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

const (
	UserRootPath = "/user"
	userIdParam  = "userId"
)

type (
	user interface {
		mapToUserCore() core.User
	}
	userCreateDTO struct {
		Password string `json:"password" validate:"required"`
	}

	userLoginDTO struct {
		ID       uuid.UUID `json:"id" validate:"required"`
		Password string    `json:"password" validate:"required"`
	}
)

func InitUserInterface(group *echo.Group) {
	group.POST("", createUserId)
	group.PUT("/:"+userIdParam, create)
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
	userId, err := uuid.Parse(context.Param(userIdParam))
	if err != nil {
		log.Warnf("Error while binding userId: %v", err)
		return echo.ErrBadRequest
	}
	userCore.SetId(userId)
	if err := userCore.Create(); err != nil {
		if errors.Is(err, errormessages.UserAlreadyExists) {
			log.Warn("User with id %s already exists", userCore.GetId())
			return context.NoContent(http.StatusConflict)
		}
		return err
	}

	log.Debugf("Created user with id %s", userCore.GetId())
	return context.NoContent(http.StatusCreated)
}

func (user *userCreateDTO) mapToUserCore() core.User {
	return &core.UserCore{Password: user.Password}
}

func (user *userLoginDTO) mapToUserCore() core.User {
	return &core.UserCore{ID: user.ID, Password: user.Password}
}

func bind(context echo.Context, toBindUser user) (core.User, error) {
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
