package api

import (
	database "Final_homework/pkg/db"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func (e *Env) deleteTaskHandler(res http.ResponseWriter, req *http.Request) { //кажется готово
	if req.Method != http.MethodDelete {
		err := errors.New("method not allowed: must be DELETE")
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	id := req.URL.Query().Get("id")

	numId, err := strconv.Atoi(id)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	err = database.DeleteTaskInDB(numId, e.DB)
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
