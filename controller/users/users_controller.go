package users

import (
	"net/http"
	"taskboard/controller/api"
	"taskboard/model"
	"taskboard/orm"
	"taskboard/service"

	"github.com/gin-gonic/gin"
)

type endPoint struct {
	login  string
	users  string
	userid string
}

// EndPoint presents boards endpoint
var EndPoint = endPoint{
	login:  "/login",
	users:  "/users",
	userid: "userid",
}

// RegisterRoute registers API endpoints for users
func (p *endPoint) RegisterRoute(route *gin.RouterGroup) (err error) {
	route.POST(p.login, login)
	route.GET(p.users, list)
	route.POST(p.users, create)
	route.GET(p.users+"/:"+p.userid, get)
	route.PUT(p.users+"/:"+p.userid, update)
	route.DELETE(p.users+"/:"+p.userid, delete)
	return
}

func login(c *gin.Context) {
	tx := orm.GetDB() // No transction
	req, serr := getLoginRequest(c)
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return
	}
	srvc := service.NewUserService(tx)
	user, serr := srvc.Login(req.Name, req.Password)
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return
	}
	res := convertUserResponse(user)
	c.IndentedJSON(http.StatusOK, res)
}

func list(c *gin.Context) {
	tx := orm.GetDB() // No transction
	srvc := service.NewUserService(tx)
	users, serr := srvc.FindUsers(&model.User{}, []string{"name"})
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return
	}
	res := convertListUserResponse(users)
	c.IndentedJSON(http.StatusOK, res)
}

func create(c *gin.Context) {
	user, serr := getUserByCreateRequest(c)
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return
	}

	// create user
	tx := orm.GetDB().Begin()
	srvc := service.NewUserService(tx)
	serr = srvc.CreateUser(user)
	if serr != nil {
		api.Rollback(tx)
		api.SetErrorStatus(c, serr)
		return
	}
	serr = api.Commit(tx)
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return
	}

	res := convertUserResponse(user)
	c.IndentedJSON(http.StatusOK, res)
}

func get(c *gin.Context) {
	tx := orm.GetDB() // No transaction
	srvc := service.NewUserService(tx)
	find, err := findUserByPathParameter(c, srvc)
	if err != nil {
		api.Rollback(tx)
		return
	}
	res := convertUserResponse(find)
	c.IndentedJSON(http.StatusOK, res)
}

func findUserByPathParameter(c *gin.Context, srvc *service.UserService) (find *model.User, serr error) {
	userID, serr := api.GetPathParameter(c, EndPoint.userid)
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return nil, serr
	}
	find, serr = srvc.FindUser(&model.User{ID: userID})
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return nil, serr
	}
	return
}

func update(c *gin.Context) {
	tx := orm.GetDB().Begin()
	srvc := service.NewUserService(tx)
	find, err := findUserByPathParameter(c, srvc)
	if err != nil {
		api.Rollback(tx)
		return
	}
	user, serr := getUserByUpdateRequest(c, find)
	if serr != nil {
		api.Rollback(tx)
		api.SetErrorStatus(c, serr)
		return
	}

	// update user
	serr = srvc.UpdateUser(user)
	if serr != nil {
		api.Rollback(tx)
		api.SetErrorStatus(c, serr)
		return
	}
	serr = api.Commit(tx)
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return
	}

	res := convertUserResponse(user)
	c.IndentedJSON(http.StatusOK, res)
}

func delete(c *gin.Context) {
	tx := orm.GetDB().Begin()
	srvc := service.NewUserService(tx)
	find, err := findUserByPathParameter(c, srvc)
	if err != nil {
		api.Rollback(tx)
		return
	}
	// delete user
	serr := srvc.DeleteUser(find)
	if serr != nil {
		api.Rollback(tx)
		api.SetErrorStatus(c, serr)
		return
	}
	serr = api.Commit(tx)
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return
	}
	c.Status(http.StatusOK)
}
