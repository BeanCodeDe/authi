package api

import (
	"net/http"
	"time"

	"github.com/BeanCodeDe/SpaceLight-Auth/internal/core"
	"github.com/BeanCodeDe/SpaceLight-AuthMiddleware/authAdapter"
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
		Password string `json:"Password" validate:"required"`
	}

	userLoginDTO struct {
		ID       uuid.UUID `json:"ID" validate:"required"`
		Password string    `json:"Password" validate:"required"`
	}

	userResponseDTO struct {
		ID        uuid.UUID
		CreatedOn time.Time
		LastLogin time.Time
	}
)

func InitUserInterface(group *echo.Group) {
	group.GET("/login", login)
	group.PUT("", create, authAdapter.CheckToken, authAdapter.CheckRole(authAdapter.DataServiceRole))
}

func login(context echo.Context) error {
	log.Debugf("Login some user")
	userCore, err := bind(context, new(userLoginDTO))
	if err != nil {
		return err
	}

	token, err := userCore.Login()
	if err != nil {
		return err
	}

	log.Debugf("Logged in user %s", userCore.ID)
	return context.String(http.StatusOK, token)
}

func create(context echo.Context) error {
	log.Debugf("Create some user")
	userCore, err := bind(context, new(userCreateDTO))
	if err != nil {
		return err
	}
	createdUser, err := userCore.Create()
	if err != nil {
		return err
	}

	log.Debugf("Created user %s", createdUser.ID)
	userResponseDTO := mapToUserResponseDTO(createdUser)
	return context.JSON(http.StatusCreated, userResponseDTO)
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
