package tests

import (
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Task struct {
	ID      int64  `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}

func count(db *sqlx.DB) (int, error) {
	var count int
	return count, db.Get(&count, `SELECT count(id) FROM scheduler`)
}

func openDB(t *testing.T) *sqlx.DB {
	dsn := "host=localhost port=5432 user=postgres password=mypass dbname=scheduler sslmode=disable"

	envDSN := os.Getenv("TODO_DB_DSN")
	if envDSN != "" {
		dsn = envDSN
	}

	db, err := sqlx.Connect("pgx", dsn)
	require.NoError(t, err)

	return db
}

func TestDB(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	before, err := count(db)
	assert.NoError(t, err)

	today := time.Now().Format(`20060102`)

	var id int64

	err = db.QueryRow(`
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES ($1, 'Todo', 'Комментарий', '')
		RETURNING id
	`, today).Scan(&id)
	assert.NoError(t, err)

	var task Task
	err = db.Get(&task, `SELECT * FROM scheduler WHERE id = $1`, id)
	assert.NoError(t, err)
	assert.Equal(t, id, task.ID)
	assert.Equal(t, `Todo`, task.Title)
	assert.Equal(t, `Комментарий`, task.Comment)

	_, err = db.Exec(`DELETE FROM scheduler WHERE id = $1`, id)
	assert.NoError(t, err)

	after, err := count(db)
	assert.NoError(t, err)

	assert.Equal(t, before, after)
}
