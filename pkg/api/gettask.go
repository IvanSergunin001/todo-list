package api

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

func (e *Env) getTaskHandler(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	numId, err := strconv.Atoi(id)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()

	task, err := e.Store.GetByID(ctx, numId)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}
	writeJson(res, http.StatusOK, task)
}
