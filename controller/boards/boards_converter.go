package boards

import (
	"taskboard/model"
	"taskboard/service"
	"time"

	"github.com/gin-gonic/gin"
)

// ID          string       `gorm:"primary_key;size:32"`
// Name        string       `gorm:"unique;size:255"`
// DispOrder   int          `gorm:"not null"`
// IsSystem    bool         `gorm:"not null"`
// IsClosed    bool         `gorm:"not null"`
// CreatedDate time.Time    `gorm:"not null"`
// Version     int          `gorm:"not null"` // Version for optimistic lock
// TaskEntries []*TaskEntry `gorm:"-"`        // Ignore, fetch by service layer when they needed

type boardResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DispOrder   int    `json:"dispOrder"`
	IsSystem    bool   `json:"isSystem"`
	IsClosed    bool   `json:"isClosed"`
	CreatedDate string `json:"createDate"`
	Version     int    `json:"version"`
}

type createRequest struct {
	Name     string `json:"name"`
	IsSystem bool   `json:"isSystem"`
	IsClosed bool   `json:"isClosed"`
}

type updateRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsSystem bool   `json:"isSystem"`
	IsClosed bool   `json:"isClosed"`
	Version  int    `json:"version"`
}

func convertBoardResponse(board *model.Board) *boardResponse {
	return &boardResponse{
		ID:          board.ID,
		Name:        board.Name,
		DispOrder:   board.DispOrder,
		IsSystem:    board.IsSystem,
		IsClosed:    board.IsClosed,
		CreatedDate: board.CreatedDate.Format(time.RFC3339),
		Version:     board.Version,
	}
}

func convertListBoardResponse(boards []model.Board) (res []*boardResponse) {
	res = make([]*boardResponse, 0, len(boards))
	for _, board := range boards {
		res = append(res, convertBoardResponse(&board))
	}
	return
}

func getBoardByCreateRequest(c *gin.Context) (*model.Board, error) {
	var req *createRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, service.NewBadRequestError(err)
	}
	board := model.NewBoard(
		req.Name,
		req.IsSystem,
		req.IsClosed,
		time.Now().UTC(),
	)
	return board, nil
}

func getBoardByUpdateRequest(c *gin.Context, find *model.Board) (*model.Board, error) {
	var req *updateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return nil, service.NewBadRequestError(err)
	}
	return &model.Board{
		ID:          find.ID,
		Name:        req.Name,
		IsSystem:    req.IsSystem,
		IsClosed:    req.IsClosed,
		Version:     req.Version,
		TaskEntries: []*model.TaskEntry{},
	}, nil
}
