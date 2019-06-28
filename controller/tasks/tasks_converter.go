package tasks

import (
	"taskboard/model"
	"taskboard/service"
	"time"

	"github.com/gin-gonic/gin"
)

// ID             string    `gorm:"primary_key;size:32"`
// Name           string    `gorm:"not null;size:255"`
// Description    string    `gorm:"size:8000"`
// AssigneeUserID string    `gorm:"size:32"`
// BoardID        string    `gorm:"size:32"`
// CreatedDate    time.Time `gorm:"not null"`
// IsClosed       bool      `gorm:"not null"`
// Version        int       `gorm:"not null"` // Version for optimistic lock
// EsitmateSize   int

type taskResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	AssigneeUserID string `json:"assigneeUserID"`
	BoardID        string `json:"boardID"`
	CreatedDate    string `json:"createDate"`
	IsClosed       bool   `json:"isClosed"`
	Version        int    `json:"version"`
	EsitmateSize   int    `json:"esitmateSize"`
}

type createRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	AssigneeUserID string `json:"assigneeUserID"`
	BoardID        string `json:"boardID"`
	CreatedDate    string `json:"createDate"`
	IsClosed       bool   `json:"isClosed"`
	EsitmateSize   int    `json:"esitmateSize"`
}

type updateRequest struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	AssigneeUserID string `json:"assigneeUserID"`
	BoardID        string `json:"boardID"`
	IsClosed       bool   `json:"isClosed"`
	Version        int    `json:"version"`
	EsitmateSize   int    `json:"esitmateSize"`
}

func convertTaskResponse(task *model.Task) *taskResponse {
	return &taskResponse{
		ID:             task.ID,
		Name:           task.Name,
		Description:    task.Description,
		AssigneeUserID: task.AssigneeUserID,
		BoardID:        task.BoardID,
		CreatedDate:    task.CreatedDate.Format(time.RFC3339),
		IsClosed:       task.IsClosed,
		Version:        task.Version,
		EsitmateSize:   task.EsitmateSize,
	}
}

func convertListTaskResponse(tasks []model.Task) (res []*taskResponse) {
	res = make([]*taskResponse, 0, len(tasks))
	for _, task := range tasks {
		res = append(res, convertTaskResponse(&task))
	}
	return
}

func getTaskByCreateRequest(c *gin.Context) (*model.Task, error) {
	var req *createRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, service.NewBadRequestError(err)
	}
	task := model.NewTask(
		req.Name,
		req.Description,
		req.IsClosed,
		time.Now().UTC(),
	)
	task.AssigneeUserID = req.AssigneeUserID
	task.BoardID = req.BoardID
	task.EsitmateSize = req.EsitmateSize
	return task, nil
}

func getTaskByUpdateRequest(c *gin.Context, find *model.Task) (*model.Task, error) {
	var req *updateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, service.NewBadRequestError(err)
	}
	return &model.Task{
		ID:             find.ID,
		Name:           req.Name,
		Description:    req.Description,
		AssigneeUserID: req.AssigneeUserID,
		BoardID:        req.BoardID,
		CreatedDate:    find.CreatedDate,
		IsClosed:       req.IsClosed,
		Version:        req.Version,
		EsitmateSize:   req.EsitmateSize,
	}, nil
}
