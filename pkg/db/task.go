package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskStorage interface {
	AddTaskInDB(ctx context.Context, task Task) (string, error)
	GetTask(ctx context.Context, limit int) (TasksResp, error)
	GetTaskForID(ctx context.Context, id int) (Task, error)
	UpdateTaskInDB(ctx context.Context, id int, task Task) error
	DeleteTaskInDB(ctx context.Context, id int) error
	UpdateDate(ctx context.Context, id int, date string) error
	GetTaskForTitle(ctx context.Context, search string, limit int) (TasksResp, error)
	GetTaskForDate(ctx context.Context, search string, limit int) (TasksResp, error)
	GetTaskByPagination(ctx context.Context, limit int, offset int) (TasksResp, error)
}

type TaskStoragePostgres struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) TaskStoragePostgres {
	return TaskStoragePostgres{db: db}
}

type TasksResp struct {
	Tasks []Task `json:"tasks"`
}

func (s *TaskStoragePostgres) AddTask(ctx context.Context, task Task) (string, error) {
	var id int64

	err := s.db.QueryRowContext(ctx, "INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)RETURNING id", task.Date, task.Title, task.Comment, task.Repeat).Scan(&id)
	if err != nil {
		return "", err
	}
	idString := strconv.Itoa(int(id))
	return idString, nil
}

func (s *TaskStoragePostgres) GetTask(ctx context.Context, limit int) (TasksResp, error) {
	var tasks TasksResp
	rows, err := s.db.QueryContext(ctx, "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT $1", limit)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		product := Task{}
		err := rows.Scan(&product.ID, &product.Date, &product.Title, &product.Comment, &product.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks.Tasks = append(tasks.Tasks, product)
	}
	if err := rows.Err(); err != nil {
		return tasks, err
	}
	return tasks, nil
}

func (s *TaskStoragePostgres) GetByID(ctx context.Context, id int) (Task, error) {
	var task Task

	rows := s.db.QueryRowContext(ctx, "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = $1", id)
	err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}
	return task, nil

}

func (s *TaskStoragePostgres) UpdateTask(ctx context.Context, task Task) error {
	numId, err := strconv.Atoi(task.ID)
	if err != nil {
		return err
	}
	res, err := s.db.ExecContext(ctx, "UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 WHERE id = $5", task.Date, task.Title, task.Comment, task.Repeat, numId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil

}

func (s *TaskStoragePostgres) DeleteTask(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM scheduler WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskStoragePostgres) UpdateByDate(ctx context.Context, id int, date string) error {
	res, err := s.db.ExecContext(ctx, "UPDATE scheduler SET date = $1 WHERE id = $2", date, id)

	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

func (s *TaskStoragePostgres) GetByTitle(ctx context.Context, search string, limit int) (TasksResp, error) {
	tasks := TasksResp{Tasks: nil}
	searchPattern := "%" + search + "%"
	rows, err := s.db.QueryContext(ctx, "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE $1 OR comment LIKE $2 ORDER BY date LIMIT $3", searchPattern, searchPattern, limit)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		product := Task{}
		err := rows.Scan(&product.ID, &product.Date, &product.Title, &product.Comment, &product.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks.Tasks = append(tasks.Tasks, product)

	}
	return tasks, nil
}

func (s *TaskStoragePostgres) GetByDate(ctx context.Context, search string, limit int) (TasksResp, error) {
	tasks := TasksResp{Tasks: nil}
	parsedTime, err := time.Parse("02.01.2006", search)
	if err != nil {
		return tasks, err
	}

	output := parsedTime.Format("20060102")

	rows, err := s.db.QueryContext(ctx, "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = $1 LIMIT $2", output, limit)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		product := Task{}
		err := rows.Scan(&product.ID, &product.Date, &product.Title, &product.Comment, &product.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks.Tasks = append(tasks.Tasks, product)

	}
	return tasks, nil
}

func (s *TaskStoragePostgres) GetByPagination(ctx context.Context, limit int, offset int) (TasksResp, error) {
	var tasks TasksResp
	rows, err := s.db.QueryContext(ctx, "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		product := Task{}
		err := rows.Scan(&product.ID, &product.Date, &product.Title, &product.Comment, &product.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks.Tasks = append(tasks.Tasks, product)
	}
	if err := rows.Err(); err != nil {
		return tasks, err
	}
	return tasks, nil
}
