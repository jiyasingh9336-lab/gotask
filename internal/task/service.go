package task

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sudarshanmg/gotask/pkg/validation"
)

var validate = validator.New()

type TaskService interface {
	Create(req CreateTaskRequest, userID int64, role string) (*TaskResponse, error)
	GetAll(page, limit int, filter TaskFilter, userID int64, role string) ([]TaskResponse, int64, int, error)
	GetById(id int64, userID int64, role string) (*TaskResponse, error)
	Update(id int64, req UpdateTaskRequest, userID int64, role string) error
	Delete(id int64, userID int64, role string) error
}

type taskService struct {
	repo TaskRepository
}

func NewService(repo TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func mapTasktoResponse(task *Task) TaskResponse {
	res := TaskResponse{
		ID:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedBy:   task.CreatedBy,
		AssignedTo:  task.AssignedTo,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
	return res
}

func (s *taskService) Create(req CreateTaskRequest, userID int64, role string) (*TaskResponse, error) {
	if role == "employee" {
		return nil, errors.New("you are not allowed to create tasks")
	}

	if err := validate.Struct(req); err != nil {
		return nil, validation.FormatValidationError(err)
	}

	task := Task{
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedBy:   userID,
		AssignedTo:  req.AssignedTo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.repo.Create(&task)

	if err != nil {
		return nil, err
	}

	res := mapTasktoResponse(&task)
	return &res, nil
}

func (s *taskService) GetAll(page, limit int, filter TaskFilter, userID int64, role string) ([]TaskResponse, int64, int, error) {
	offset := (page - 1) * limit

	tasks, err := s.repo.FindAll(offset, limit, filter, userID, role)
	if err != nil {
		return nil, 0, 0, err
	}

	validateSortFields := map[string]bool{
		"id": true, "title": true, "created_at": true, "updated_at": true,
	}

	if !validateSortFields[filter.SortBy] {
		filter.SortBy = "id"
	}

	if filter.Order != "asc" && filter.Order != "desc" {
		filter.Order = "asc"
	}

	total, err := s.repo.CountAll()
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	responses := make([]TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		res := mapTasktoResponse(&t)
		responses = append(responses, res)
	}

	return responses, total, totalPages, nil
}

func (s *taskService) GetById(id int64, userID int64, role string) (*TaskResponse, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	task, err := s.repo.FindById(id)

	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrNotFound
	}

	if role == "employee" && task.AssignedTo != userID {
		return nil, ErrNotFound
	}

	res := mapTasktoResponse(task)
	return &res, nil
}

func (s *taskService) Update(id int64, req UpdateTaskRequest, userID int64, role string) error {
	if id <= 0 {
		return ErrInvalidID
	}

	if err := validate.Struct(req); err != nil {
		return validation.FormatValidationError(err)
	}
	task, err := s.repo.FindById(id)

	if err != nil {
		return err
	}

	if task == nil {
		return ErrNotFound
	}

	if role == "employee" {
		if task.AssignedTo != userID {
			return ErrNotFound
		}

		if req.Completed != nil {
			task.Completed = *req.Completed
		}

		return s.repo.Update(task)
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Completed != nil {
		task.Completed = *req.Completed
	}
	if req.AssignedTo != nil {
		task.AssignedTo = *req.AssignedTo
	}

	task.UpdatedAt = time.Now()
	return s.repo.Update(task)
}

func (s *taskService) Delete(id int64, userID int64, role string) error {
	if id <= 0 {
		return ErrInvalidID
	}
	task, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	if task == nil {
		return ErrNotFound
	}

	if role == "employee" {
		return errors.New("you are not allowed to delete tasks")
	}
	return s.repo.Delete(id)

}
