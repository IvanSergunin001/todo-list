package main //переделать

import (
	"Final_homework/pkg/api"
	database "Final_homework/pkg/db"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := os.Getenv("TODO_DBFILE")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=mypass dbname=scheduler sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println("Не удалось подключиться к PostgreSQL:", err)
		return
	}

	if err := database.Init(db); err != nil {
		fmt.Println("Не удалось создать таблицы:", err)
		return
	}

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)
	store := database.NewCustomerRepository(db)
	api.Init(store)

	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = "7540" // порт по умолчанию
	}

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
