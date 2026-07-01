package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TasksResp struct {
	Tasks []Task `json:"tasks"`
}

func AddTaskInDB(task Task, db *sql.DB) (string, error) {
	var id int64
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err != nil {
		return "", err
	}
	id, err = res.LastInsertId()
	idString := strconv.Itoa(int(id))
	return idString, nil
}

func SelectTask(db *sql.DB, limit int) (TasksResp, error) {
	var tasks TasksResp
	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
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

func GetTask(db *sql.DB, id int) (Task, error) { //исправить
	var task Task

	rows := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}
	return task, nil

}

func UpdateTaskInDB(id int, task Task, db *sql.DB) error { // исправить
	res, err := db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id", sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat), sql.Named("id", id))
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

func DeleteTaskInDB(id int, db *sql.DB) error { //кажется готово
	_, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}

func UpdateDate(id int, db *sql.DB, date string) error { //кажется готово
	res, err := db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
		sql.Named("date", date), sql.Named("id", id))

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

func SelectTaskForTitle(db *sql.DB, search string, limit int) (TasksResp, error) { // исправить
	tasks := TasksResp{Tasks: []Task{}}
	searchPattern := "%" + search + "%"
	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit", sql.Named("search", searchPattern), sql.Named("limit", limit))
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

func SelectTaskForDate(db *sql.DB, search string, limit int) (TasksResp, error) { // исправить
	tasks := TasksResp{Tasks: []Task{}}
	parsedTime, err := time.Parse("02.01.2006", search)
	if err != nil {
		return tasks, err
	}

	output := parsedTime.Format("20060102")

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date LIMIT :limit ", sql.Named("date", output), sql.Named("limit", limit))
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
