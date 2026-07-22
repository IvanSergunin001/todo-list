package api

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

func (e *Env) tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	page := r.URL.Query().Get("page")

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if search != "" {
		dateCheck := SearchisValidDateFormat(search)
		if dateCheck == true {
			tasks, err := e.Store.GetByDate(ctx, search, 50)
			if err != nil {
				writeJson(w, http.StatusInternalServerError, err)
				return
			}
			writeJson(w, http.StatusOK, tasks)
		} else {
			tasks, err := e.Store.GetByTitle(ctx, search, 50)
			if err != nil {
				writeJson(w, http.StatusInternalServerError, err)
				return
			}
			writeJson(w, http.StatusOK, tasks)
		}
	} else {
		if page == "" {
			tasks, err := e.Store.GetTask(ctx, 50)
			if err != nil {
				writeJson(w, http.StatusInternalServerError, err)
				return
			}
			writeJson(w, http.StatusOK, tasks)
			return
		}
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, err)
			return
		}
		var offset int

		limit := 2

		offset = (pageNum - 1) * limit
		tasks, err := e.Store.GetByPagination(ctx, limit, offset)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, err)
			return
		}
		writeJson(w, http.StatusOK, tasks)
	}
}

func SearchisValidDateFormat(dateStr string) bool {
	layout := "02.01.2006"

	_, err := time.Parse(layout, dateStr)
	if err != nil {
		return false
	} else {
		return true
	}
}
