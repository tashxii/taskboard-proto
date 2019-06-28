package service

import (
	"taskboard/model"
	"taskboard/orm"
	"taskboard/repository"

	"github.com/jinzhu/gorm"
)

// BoardService provides apis for board management.
type BoardService struct {
	tx            *gorm.DB
	boardRepo     *repository.BoardRepository
	taskRepo      *repository.TaskRepository
	taskEntryRepo *repository.TaskEntryRepository
}

// NewBoardService return new instance of BoardService.
func NewBoardService(tx *gorm.DB) *BoardService {
	return &BoardService{
		tx:            tx,
		boardRepo:     repository.NewBoardRepository(tx),
		taskRepo:      repository.NewTaskRepository(tx),
		taskEntryRepo: repository.NewTaskEntryRepository(tx),
	}
}

// FindBoard returns board matching specified condition
func (s *BoardService) FindBoard(condition interface{}) (*model.Board, error) {
	find, err := s.boardRepo.FindFirstBoard(condition, []string{"id"})
	if err != nil {
		if err == orm.ErrorRecordNotFound {
			return nil, NewSvcErrorf(ErrorCodeNotFound, err, "Board not found")
		}
		return nil, NewSvcError(ErrorCodeDB, err, "Failed to find board")
	}
	return &find, nil
}

// FindBoards finds all boards
func (s *BoardService) FindBoards(sortOrders []string) ([]model.Board, error) {
	boards, err := s.boardRepo.FindBoards(&model.Board{}, 0, orm.NoLimit, sortOrders)
	if err != nil {
		return nil, NewSvcError(ErrorCodeDB, err, "Failed to find boards")
	}
	return boards, nil
}

// CreateBoard creates new board
func (s *BoardService) CreateBoard(board *model.Board) error {
	err := s.boardRepo.CreateBoard(board)
	if err != nil {
		return NewSvcError(ErrorCodeDB, err, "Failed to create board")
	}
	return nil
}

// UpdateBoard updates specifed board
func (s *BoardService) UpdateBoard(board *model.Board) error {
	err := s.boardRepo.UpdateBoard(board)
	if err != nil {
		return NewSvcErrorf(ErrorCodeDB, err, "Failed to update board. ID:%s", board.ID)
	}
	return nil
}

// DeleteBoard deletes specifed board
func (s *BoardService) DeleteBoard(board *model.Board) error {
	err := s.boardRepo.DeleteBoard(board)
	if err != nil {
		return NewSvcErrorf(ErrorCodeDB, err, "Failed to delete board. ID:%s", board.ID)
	}
	err = s.taskEntryRepo.DeleteTaskEntriesByBoardID(board.ID)
	if err != nil {
		return NewSvcErrorf(ErrorCodeDB, err, "Failed to delete board's entries. TaskID:%s", board.ID)
	}
	err = s.taskRepo.ClearBoardID(board.ID)
	if err != nil {
		return NewSvcErrorf(ErrorCodeDB, err, "Failed to clear task's board ID. BoardID:%s", board.ID)
	}
	return nil
}
