package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

func (e *Env) doneTaskHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJson(res, http.StatusInternalServerError, errors.New("method not allowed: must be POST"))
		return
	}
	id := req.URL.Query().Get("id")

	var hollowResponse, _ = json.Marshal(struct{}{})

	numId, err := strconv.Atoi(id)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	task, err := e.Store.GetByID(ctx, numId)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	var targetDate time.Time
	targetDate = time.Now()
	now := targetDate.Format("20060102")

	if task.Repeat == "" {
		err = e.Store.DeleteTask(ctx, numId)
		if err != nil {
			writeJson(res, http.StatusInternalServerError, err)
			return
		} else {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(hollowResponse))
			return
		}
	}

	newDate, err := NextDate(task.Date, now, task.Repeat)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	err = e.Store.UpdateByDate(ctx, numId, newDate)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(hollowResponse))

}
