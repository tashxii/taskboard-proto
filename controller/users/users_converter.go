package users

import (
	"taskboard/model"
	"taskboard/service"

	"github.com/gin-gonic/gin"
)

// ID           string `gorm:"primary_key;size:32"`
// Name         string `gorm:"size:255;not null;unique"`
// PasswordHash string `gorm:"size:255;not null;"`
// Avator       string `gorm:"size:255"`
// Version      int    `gorm:"not null"` // Version for optimistic lock

type loginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type userResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Avator  string `json:"avator"`
	Version int    `json:"version"`
}

type createRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Avator   string `json:"avator"`
}

type updateRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Avator   string `json:"avator"`
	Version  int    `json:"version"`
}

func getLoginRequest(c *gin.Context) (*loginRequest, error) {
	var req loginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, service.NewBadRequestError(err)
	}
	return &req, nil
}

func convertUserResponse(user *model.User) *userResponse {
	return &userResponse{
		ID:      user.ID,
		Name:    user.Name,
		Avator:  user.Avator,
		Version: user.Version,
	}
}

func convertListUserResponse(users []model.User) (res []*userResponse) {
	res = make([]*userResponse, 0, len(users))
	for _, user := range users {
		res = append(res, convertUserResponse(&user))
	}
	return
}

func getUserByCreateRequest(c *gin.Context) (*model.User, error) {
	var req createRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, service.NewBadRequestError(err)
	}
	return model.NewUser(req.Name, req.Password, req.Avator), nil
}

func getUserByUpdateRequest(c *gin.Context, find *model.User) (*model.User, error) {
	var req updateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, service.NewBadRequestError(err)
	}
	user := &model.User{
		ID:      find.ID,
		Name:    req.Name,
		Avator:  req.Avator,
		Version: req.Version,
	}
	user.SetPassword(req.Password)
	return user, nil
}
