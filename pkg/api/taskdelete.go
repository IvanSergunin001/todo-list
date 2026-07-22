package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

func (e *Env) deleteTaskHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		writeJson(res, http.StatusInternalServerError, errors.New("method not allowed: must be DELETE"))
		return
	}
	var hollowResponse, _ = json.Marshal(struct{}{})
	id := req.URL.Query().Get("id")

	numId, err := strconv.Atoi(id)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()

	err = e.Store.DeleteTask(ctx, numId)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(hollowResponse))

}
