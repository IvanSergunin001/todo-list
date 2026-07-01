package api

import (
	database "Final_homework/pkg/db"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func (e *Env) addTaskHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		err := errors.New("method not allowed: must be POST")
		writeJson(res, http.StatusInternalServerError, err)
		return
	}
	var newTask database.Task
	var dateCheck bool

	err := json.NewDecoder(req.Body).Decode(&newTask)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}
	defer req.Body.Close()

	if newTask.Title == "" {
		err = errors.New("wrong title")
		writeJson(res, http.StatusBadRequest, err)
		return
	}

	var targetDate time.Time
	targetDate = time.Now()                           //время сейчас
	targetDateString := targetDate.Format("20060102") //время сейчас в формате строки
	todayMidnight, _ := time.Parse("20060102", targetDateString)

	if newTask.Date == "" {
		newTask.Date = targetDateString
	} else {
		dateCheck = isValidDateFormat(newTask.Date)
		if dateCheck != true {
			err := errors.New("data not valid date format")
			writeJson(res, http.StatusBadRequest, err)
			return
		}
	}

	dateTimeFormat, err := time.Parse("20060102", newTask.Date) //дата в формате времени
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	if todayMidnight.After(dateTimeFormat) { //если настоящее время больше
		if newTask.Repeat == "" {
			newTask.Date = targetDateString
		} else if newTask.Repeat != "" {
			newTask.Date, err = NextDate(newTask.Date, targetDateString, newTask.Repeat) //
			if err != nil {
				writeJson(res, http.StatusInternalServerError, err)
				return
			}
		}
	}

	dateCheck = isValidDateFormat(newTask.Date)
	if dateCheck != true {
		err := errors.New("data not valid date format")
		writeJson(res, http.StatusBadRequest, err)
		return
	}

	newTask.ID, err = database.AddTaskInDB(newTask, e.DB) //
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}
	writeJson(res, http.StatusOK, newTask.ID)
}

func isValidDateFormat(dateStr string) bool {
	layout := "20060102"

	_, err := time.Parse(layout, dateStr)
	if err != nil {
		return false
	} else {
		return true
	}
}

func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")

	var serverAnswer any
	switch v := data.(type) {
	case error:
		serverAnswer = struct {
			Error string `json:"error"`
		}{
			Error: v.Error(),
		}
	case string:
		serverAnswer = struct {
			ID string `json:"id"`
		}{
			ID: v,
		}
	case database.TasksResp:
		if len(v.Tasks) == 0 {
			v.Tasks = []database.Task{}
		}
		serverAnswer = v
	case database.Task:
		serverAnswer = v
	default:
		w.WriteHeader(http.StatusInternalServerError)
		err := errors.New("Unsupported data type")
		serverAnswer = struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		}
	}

	jsonData, err := json.Marshal(serverAnswer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(jsonData)
}
