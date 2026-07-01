package api

import (
	database "Final_homework/pkg/db"
	"net/http"
	"strconv"
)

func (e *Env) getTaskHandler(res http.ResponseWriter, req *http.Request) {
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
	writeJson(res, http.StatusOK, task)
}
