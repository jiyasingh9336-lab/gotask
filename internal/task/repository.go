package task

import (
	"database/sql"
	"errors"
	"time"
)

type TaskRepository interface {
	Create(task *Task) (int64, error)
	FindAll(offset, limit int, filter TaskFilter, userID int64, role string) ([]Task, error)
	FindById(id int64) (*Task, error)
	Update(task *Task) error
	Delete(id int64) error
	CountAll() (int64, error)
}

type PostgresTaskRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) TaskRepository {
	return &PostgresTaskRepository{DB: db}
}

func (r *PostgresTaskRepository) Create(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO tasks (title, description, completed, created_at, updated_at,created_by,assigned_to)
            VALUES ($1, $2, $3, $4, $5,$6,$7)
            RETURNING id;
          `

	task.Completed = false
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	err := r.DB.QueryRow(query, task.Title, task.Description, task.Completed, task.CreatedAt, task.UpdatedAt, task.CreatedBy, task.AssignedTo).Scan(&id)

	if err != nil {
		return 0, err
	}

	task.Id = id
	return id, nil
}

func (r *PostgresTaskRepository) FindAll(offset, limit int, filter TaskFilter, userID int64, role string) ([]Task, error) {
	if filter.SortBy == "" {
		filter.SortBy = "id"
	}
	if filter.Order != "asc" && filter.Order != "desc" {
		filter.Order = "asc"
	}

	validSortFields := map[string]bool{
		"id": true, "title": true, "created_at": true, "updated_at": true, "created_by": true, "assigned_to": true, "completed": true,
	}
	if !validSortFields[filter.SortBy] {
		filter.SortBy = "id"
	}

	var query string
	var rows *sql.Rows
	var err error

	if role == "employee" {
		query = `
		SELECT id, title, description, completed, created_at, updated_at, created_by, assigned_to
		FROM tasks
		WHERE ($1::bool IS NULL OR completed = $1)
		AND assigned_to = $2
		ORDER BY ` + filter.SortBy + ` ` + filter.Order + `
		LIMIT $3 OFFSET $4;`

		rows, err = r.DB.Query(query, filter.Completed, userID, limit, offset)

	} else {
		query = `
		SELECT id, title, description, completed, created_at, updated_at, created_by, assigned_to
		FROM tasks
		WHERE ($1::bool IS NULL OR completed = $1)
		ORDER BY ` + filter.SortBy + ` ` + filter.Order + `
		LIMIT $2 OFFSET $3;`

		rows, err = r.DB.Query(query, filter.Completed, limit, offset)
	}

	if err != nil {
		return nil, err
	}

	tasks := []Task{}

	for rows.Next() {
		task := Task{}
		if err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt, &task.CreatedBy, &task.AssignedTo); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return tasks, nil
}

func (r *PostgresTaskRepository) FindById(id int64) (*Task, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at,created_by,assigned_to FROM tasks WHERE id = $1;`

	task := Task{}
	err := r.DB.QueryRow(query, id).Scan(&task.Id, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt, &task.CreatedBy, &task.AssignedTo)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *PostgresTaskRepository) Update(task *Task) error {
	query := `UPDATE tasks
            SET title = $1, description = $2, completed = $3, updated_at = $4, assigned_to = $6
            WHERE id = $5;
          `

	task.UpdatedAt = time.Now()
	res, err := r.DB.Exec(query, task.Title, task.Description, task.Completed, task.UpdatedAt, task.AssignedTo, task.Id)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}

	return nil
}

func (r *PostgresTaskRepository) Delete(id int64) error {
	query := `DELETE FROM tasks WHERE id = $1;`

	res, err := r.DB.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no task deleted")
	}

	return nil
}

func (r *PostgresTaskRepository) CountAll() (int64, error) {
	var count int64
	err := r.DB.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
