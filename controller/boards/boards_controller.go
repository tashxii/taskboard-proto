package boards

import (
	"net/http"
	"taskboard/controller/api"
	"taskboard/model"
	"taskboard/orm"
	"taskboard/service"

	"github.com/gin-gonic/gin"
)

// RegisterRoute registers API endpoints for boards
func RegisterRoute(route *gin.RouterGroup) (err error) {
	route.GET("/boards", list)
	route.POST("/boards", create)
	route.GET("/boards/:boardid", get)
	route.PUT("/boards/:boardid", update)
	route.DELETE("/boards/:boardid", delete)
	// route.GET("/boards/tasks/:boardid", getBoardTasks)
	// route.PUT("/orders/boards", updateBoardOrders)
	// route.PUT("/orders//tasks", updateBoardTaskOrders)
	return
}

func list(c *gin.Context) {
	tx := orm.GetDB() // No transction
	srvc := service.NewBoardService(tx)
	boards, serr := srvc.FindBoards([]string{"disp_order, created_date"})
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return
	}
	res := convertListBoardResponse(boards)
	c.IndentedJSON(http.StatusOK, res)
}

func create(c *gin.Context) {
	board, serr := getBoardByCreateRequest(c)
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return
	}

	// create board
	tx := orm.GetDB().Begin()
	srvc := service.NewBoardService(tx)
	serr = srvc.CreateBoard(board)
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

	res := convertBoardResponse(board)
	c.IndentedJSON(http.StatusOK, res)
}

func get(c *gin.Context) {
	tx := orm.GetDB() // No transaction
	srvc := service.NewBoardService(tx)
	find, err := findBoardByPathParameter(c, srvc)
	if err != nil {
		api.Rollback(tx)
		return
	}
	res := convertBoardResponse(find)
	c.IndentedJSON(http.StatusOK, res)
}

func findBoardByPathParameter(c *gin.Context, srvc *service.BoardService) (find *model.Board, serr error) {
	boardID, serr := api.GetPathParameter(c, "boardid")
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return nil, serr
	}
	find, serr = srvc.FindBoard(&model.Board{ID: boardID})
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return nil, serr
	}
	return
}

func update(c *gin.Context) {
	tx := orm.GetDB().Begin()
	srvc := service.NewBoardService(tx)
	find, err := findBoardByPathParameter(c, srvc)
	if err != nil {
		api.Rollback(tx)
		return
	}
	board, serr := getBoardByUpdateRequest(c, find)
	if serr != nil {
		api.Rollback(tx)
		api.SetErrorStatus(c, serr)
		return
	}

	// update board
	serr = srvc.UpdateBoard(board)
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

	res := convertBoardResponse(board)
	c.IndentedJSON(http.StatusOK, res)
}

func delete(c *gin.Context) {
	tx := orm.GetDB().Begin()
	srvc := service.NewBoardService(tx)
	find, err := findBoardByPathParameter(c, srvc)
	if err != nil {
		api.Rollback(tx)
		return
	}
	// delete board
	serr := srvc.DeleteBoard(find)
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
