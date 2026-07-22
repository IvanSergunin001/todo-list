package api

import (
	database "Final_homework/pkg/db"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func (e *Env) addTaskHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJson(res, http.StatusInternalServerError, errors.New("method not allowed: must be POST"))
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
		writeJson(res, http.StatusBadRequest, errors.New("wrong title"))
		return
	}

	var targetDate time.Time
	targetDate = time.Now()                           //время сейчас
	targetDateString := targetDate.Format("20060102") //время сейчас в формате строки
	todayMidnight, _ := time.Parse("20060102", targetDateString)

	finalDate, err := defaultDate(newTask.Date, targetDateString)
	if err != nil {
		writeJson(res, http.StatusBadRequest, err)
		return
	}
	newTask.Date = finalDate

	dateTimeFormat, err := time.Parse("20060102", newTask.Date) //дата в формате времени
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	if todayMidnight.After(dateTimeFormat) { //если настоящее время больше
		finalDate, err = resolveDateByRepeat(newTask.Repeat, newTask.Date, targetDateString)
		if err != nil {
			writeJson(res, http.StatusInternalServerError, err)
			return
		}
		newTask.Date = finalDate
	}

	dateCheck = isValidDateFormat(newTask.Date)
	if dateCheck != true {
		writeJson(res, http.StatusBadRequest, errors.New("data not valid date format"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	newTask.ID, err = e.Store.AddTask(ctx, newTask)
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
		json.NewEncoder(w).Encode(map[string]string{"error": "Unsupported data type"})
		return
	}

	jsonData, err := json.Marshal(serverAnswer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(jsonData)
}

func defaultDate(date string, defaultVal string) (string, error) {
	if date == "" {
		date = defaultVal
	} else {
		dateCheck := isValidDateFormat(date)
		if dateCheck != true {
			return "", errors.New("data not valid date format")
		}
	}

	return date, nil
}

func resolveDateByRepeat(repeat, date, defaultVal string) (string, error) {
	if repeat == "" {
		date = defaultVal
	} else if repeat != "" {
		dateCheck, err := NextDate(date, defaultVal, repeat)
		if err != nil {
			return "", err
		}
		date = dateCheck
	}

	return date, nil
}
