package api

import (
	database "Final_homework/pkg/db"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

func (e *Env) doneTaskHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		err := errors.New("method not allowed: must be POST")
		writeJson(res, http.StatusInternalServerError, err)
		return
	}
	type Response struct{}
	data, err := json.Marshal(Response{})
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}
	id := req.URL.Query().Get("id")

	numId, err := strconv.Atoi(id)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	task, err := database.GetTask(e.DB, numId)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	var targetDate time.Time
	targetDate = time.Now()
	now := targetDate.Format("20060102")

	if task.Repeat == "" {
		err = database.DeleteTaskInDB(numId, e.DB)
		if err != nil {
			writeJson(res, http.StatusInternalServerError, err)
			return
		} else {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(data))
			return
		}
	}

	newDate, err := NextDate(task.Date, now, task.Repeat)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	err = database.UpdateDate(numId, e.DB, newDate)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(data))

}
