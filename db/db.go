package db

import (
	"os"
	"fmt"
	"database/sql"
	"shitblog-server/utils"
	"strings"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
    host     = "db"
    port     = 5432
    user     = "blogger"
    dbname   = "blog"
)

func generateTokenFromUsername(username string) string {
	token, err := bcrypt.GenerateFromPassword([]byte(username), bcrypt.DefaultCost)
	utils.PanicIfError(err)
	return string(token)
}

func ConnectToDb() *sql.DB {
	password, _ := os.LookupEnv("POSTGRES_PASSWORD")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
        "password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        panic(err)
    }
	return db
}

func RecreateTables(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	utils.PanicIfError(err)
	_, err = db.Exec("DROP TABLE IF EXISTS posts")
	utils.PanicIfError(err)
	_, err = db.Exec("CREATE TABLE users (username TEXT, token TEXT)")
	utils.PanicIfError(err)
	_, err = db.Exec("CREATE TABLE posts (id BIGSERIAL PRIMARY KEY, title TEXT, text TEXT)")
	utils.PanicIfError(err)
}

func CreateUser(username string, db *sql.DB) int {
	if strings.Contains(username, "'") || strings.Contains(username, "\"") || strings.Contains(username, "=") || strings.Contains(username, ";") {
		return 2
	}
	res, err := db.Query("SELECT * FROM users WHERE username='" + username + "'")
	utils.PanicIfError(err)
	if res.Next() {
		return 1
	}
	_, err = db.Exec("INSERT INTO users (username, token) VALUES ('" + username + "', '" + generateTokenFromUsername(username) + "')")
	utils.PanicIfError(err)
	return 0
}
