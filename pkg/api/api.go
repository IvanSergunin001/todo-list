package api

import (
	"database/sql"
	"net/http"
	"os"
)

type Env struct {
	DB *sql.DB
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwt string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}
			var valid bool
			valid = isValidToken(jwt, Secret)
			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized) //переделать
				return
			}
		}
		next(w, r) //
	})
}

func Init(db *sql.DB) {

	env := &Env{DB: db}

	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/signin", signInHandler)
	http.HandleFunc("/api/task", auth(env.taskHandler))
	http.HandleFunc("/api/tasks", auth(env.tasksHandler))
	http.HandleFunc("/api/task/done", auth(env.doneTaskHandler))

}
