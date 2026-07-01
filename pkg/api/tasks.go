package api

import (
	database "Final_homework/pkg/db"
	"net/http"
	"time"
)

func (e *Env) tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	if search != "" {
		dateCheck := SearchisValidDateFormat(search)
		if dateCheck == true {
			tasks, err := database.SelectTaskForDate(e.DB, search, 50)
			if err != nil {
				writeJson(w, http.StatusInternalServerError, err)
				return
			}
			writeJson(w, http.StatusOK, tasks)
		} else {
			tasks, err := database.SelectTaskForTitle(e.DB, search, 50)
			if err != nil {
				writeJson(w, http.StatusInternalServerError, err)
				return
			}
			writeJson(w, http.StatusOK, tasks)
		}
	} else {
		tasks, err := database.SelectTask(e.DB, 50)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, err)
			return
		}
		writeJson(w, http.StatusOK, tasks)
	}
}

func SearchisValidDateFormat(dateStr string) bool {
	layout := "02.01.2006" //переделать

	_, err := time.Parse(layout, dateStr)
	if err != nil {
		return false
	} else {
		return true
	}
}
