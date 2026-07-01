package main //переделать

import (
	"Final_homework/pkg/api"
	database "Final_homework/pkg/db"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	_, err = os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	if install == true {
		err = database.Init(db, dbFile)
	}

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)
	api.Init(db)

	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = "7540" // порт по умолчанию
	}

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
