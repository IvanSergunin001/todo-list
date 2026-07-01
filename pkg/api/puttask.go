package api

import (
	database "Final_homework/pkg/db"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

func (e *Env) putTaskHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		err := errors.New("method not allowed: must be PUT")
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	var newTask database.Task

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
		dateCheck := isValidDateFormat(newTask.Date)
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
			newTask.Date, err = NextDate(newTask.Date, targetDateString, newTask.Repeat)
			if err != nil {
				writeJson(res, http.StatusInternalServerError, err)
				return
			}
		}
	}

	dateCheck := isValidDateFormat(newTask.Date)
	if dateCheck != true {
		err := errors.New("data not valid date format")
		writeJson(res, http.StatusBadRequest, err)
		return
	}

	numId, err := strconv.Atoi(newTask.ID)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	err = database.UpdateTaskInDB(numId, newTask, e.DB)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	type Response struct{}
	data, err := json.Marshal(Response{})
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(data))

}
