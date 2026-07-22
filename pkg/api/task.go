package api

import (
	"net/http"
)

func (e *Env) taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		e.addTaskHandler(w, r)
	case http.MethodGet:
		e.getTaskHandler(w, r)
	case http.MethodPut:
		e.putTaskHandler(w, r)
	case http.MethodDelete:
		e.deleteTaskHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
