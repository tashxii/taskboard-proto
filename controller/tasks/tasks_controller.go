package tasks

import (
	"net/http"
	"taskboard/controller/api"
	"taskboard/model"
	"taskboard/orm"
	"taskboard/service"

	"github.com/gin-gonic/gin"
)

// RegisterRoute registers API endpoints for tasks
func RegisterRoute(route *gin.RouterGroup) (err error) {
	route.GET("/tasks", list)
	route.POST("/tasks", create)
	route.GET("/tasks/:taskid", get)
	route.PUT("/tasks/:taskid", update)
	route.DELETE("/tasks/:taskid", delete)
	return
}

func list(c *gin.Context) {
	tx := orm.GetDB() // No transction
	srvc := service.NewTaskService(tx)
	tasks, serr := srvc.FindTasks([]string{"disp_order, created_date"})
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return
	}
	res := convertListTaskResponse(tasks)
	c.IndentedJSON(http.StatusOK, res)
}

func create(c *gin.Context) {
	task, serr := getTaskByCreateRequest(c)
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return
	}

	// create task
	tx := orm.GetDB().Begin()
	srvc := service.NewTaskService(tx)
	serr = srvc.CreateTask(task)
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

	res := convertTaskResponse(task)
	c.IndentedJSON(http.StatusOK, res)
}

func get(c *gin.Context) {
	tx := orm.GetDB() // No transaction
	srvc := service.NewTaskService(tx)
	find, err := findTaskByPathParameter(c, srvc)
	if err != nil {
		api.Rollback(tx)
		return
	}
	res := convertTaskResponse(find)
	c.IndentedJSON(http.StatusOK, res)
}

func findTaskByPathParameter(c *gin.Context, srvc *service.TaskService) (find *model.Task, serr error) {
	taskID, serr := api.GetPathParameter(c, "taskid")
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return nil, serr
	}
	find, serr = srvc.FindTask(&model.Task{ID: taskID})
	if serr != nil {
		api.SetErrorStatus(c, serr)
		return nil, serr
	}
	return
}

func update(c *gin.Context) {
	tx := orm.GetDB().Begin()
	srvc := service.NewTaskService(tx)
	find, err := findTaskByPathParameter(c, srvc)
	if err != nil {
		api.Rollback(tx)
		return
	}
	task, serr := getTaskByUpdateRequest(c, find)
	if serr != nil {
		api.Rollback(tx)
		api.SetErrorStatus(c, serr)
		return
	}

	// update task
	serr = srvc.UpdateTask(task)
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

	res := convertTaskResponse(task)
	c.IndentedJSON(http.StatusOK, res)
}

func delete(c *gin.Context) {
	tx := orm.GetDB().Begin()
	srvc := service.NewTaskService(tx)
	find, err := findTaskByPathParameter(c, srvc)
	if err != nil {
		api.Rollback(tx)
		return
	}
	// delete task
	serr := srvc.DeleteTask(find)
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
