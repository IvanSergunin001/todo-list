package api

import (
	database "Final_homework/pkg/db"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func (e *Env) putTaskHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		writeJson(res, http.StatusInternalServerError, errors.New("method not allowed: must be PUT"))
		return
	}

	var newTask database.Task
	var hollowResponse, _ = json.Marshal(struct{}{})

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

	dateCheck := isValidDateFormat(newTask.Date)
	if dateCheck != true {
		writeJson(res, http.StatusBadRequest, errors.New("data not valid date format"))
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()

	err = e.Store.UpdateTask(ctx, newTask)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(hollowResponse))

}
