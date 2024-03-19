package db

import (
	"os"
	"fmt"
	"database/sql"
	"shitblog-server/utils"
	"strings"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"crypto/md5"
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

func RecreateTables() {
	db := ConnectToDb()
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	utils.PanicIfError(err)
	_, err = db.Exec("DROP TABLE IF EXISTS posts")
	utils.PanicIfError(err)
	_, err = db.Exec("CREATE TABLE users (username TEXT, token TEXT)")
	utils.PanicIfError(err)
	_, err = db.Exec("CREATE TABLE posts (id BIGSERIAL PRIMARY KEY, author TEXT, title TEXT, text TEXT)")
	utils.PanicIfError(err)
}

func CreateUser(username string) string {
	db := ConnectToDb()
	if strings.Contains(username, "'") || strings.Contains(username, "\"") || strings.Contains(username, "=") || strings.Contains(username, ";") {
		return "BAN"
	}
	res, err := db.Query("SELECT * FROM users WHERE username='" + username + "'")
	utils.PanicIfError(err)
	if res.Next() {
		return "username taken"
	}
	token := generateTokenFromUsername(username)
	hash := md5.Sum([]byte(token))
	_, err = db.Exec("INSERT INTO users (username, token) VALUES ('" + username + "', '" + fmt.Sprintf("%x", hash) + "')")
	utils.PanicIfError(err)
	return token
}

func DeleteUser(username string, token string) int {
	db := ConnectToDb()
	if strings.Contains(username, "'") || strings.Contains(username, "\"") || strings.Contains(username, "=") || strings.Contains(username, ";") {
		return 2 // BAN
	}
	res, err := db.Query("SELECT * FROM users WHERE username='" + username + "'")
	utils.PanicIfError(err)
	if !res.Next() {
		return 1 // No such user
	}
	var correct_hash string
	var usernam string
	res.Scan(&usernam, &correct_hash)
	hash := fmt.Sprintf("%x", md5.Sum([]byte(token)))
	if hash != correct_hash {
		return 2 // Incorrect token
	}
	_, err = db.Exec("DELETE FROM users WHERE username='" + username + "'")
	utils.PanicIfError(err)
	return 0 // Success
}

func GetUsers() []string {
	db := ConnectToDb()
	res, err := db.Query("SELECT * FROM users")
	var ret []string
	utils.PanicIfError(err)
	for res.Next() {
		var username, token string
		err = res.Scan(&username, &token)
		utils.PanicIfError(err)
		ret = append(ret, username)
	}
	return ret
}

func CreatePost(token string, title string, text string) int {
	db := ConnectToDb()
	hash := fmt.Sprintf("%x", md5.Sum([]byte(token)))
	res, err := db.Query("SELECT * FROM users WHERE token='" + hash + "'")
	utils.PanicIfError(err)
	if !res.Next() {
		return 1
	}
	var author, tok string
	err = res.Scan(&author, &tok)
	utils.PanicIfError(err)
	_, err = db.Exec("INSERT INTO posts (author, title, text) VALUES ('" + author + "', '" + title  + "', '" + text + "')")
	utils.PanicIfError(err)
	return 0
}

func GetPosts(count int, author string) []map[string]string {
	var ret []map[string]string
	db := ConnectToDb()
	var res *sql.Rows
	var err error
	if author == "" {
		res, err = db.Query("SELECT * FROM posts ORDER BY ID DESC LIMIT " + fmt.Sprint(count))
	} else {
		res, err = db.Query("SELECT * FROM posts WHERE author = '" + author + "' ORDER BY ID DESC LIMIT " + fmt.Sprint(count))
	}
	utils.PanicIfError(err)
	for res.Next() {
		cur := make(map[string]string)
		var id, author, title, text string
		err := res.Scan(&id, &author, &title, &text)
		utils.PanicIfError(err)
		cur["id"] = id
		cur["author"] = author
		cur["title"] = title
		cur["text"] = text
		ret = append(ret, cur)
	}
	return ret
}

func GetPost(id uint64) map[string]string {
	db := ConnectToDb()
	res, err := db.Query("SELECT * FROM posts WHERE id = '" + fmt.Sprint(id) + "'")
	utils.PanicIfError(err)
	if !res.Next() {
		return nil
	}
	var author, title, text string;
	err = res.Scan(&id, &author, &title, &text)
	ret := make(map[string]string)
	ret["id"] = fmt.Sprint(id)
	ret["author"] = author
	ret["title"] = title
	ret["text"] = text
	return ret
}

func DeletePost(token string, id uint64) int {
	db := ConnectToDb()
	var username, correct_hash string
	hash := fmt.Sprintf("%x", md5.Sum([]byte(token)))
	res, err := db.Query("SELECT * FROM users WHERE token='" + hash + "'")
	utils.PanicIfError(err)
	if !res.Next() {
		return 2
	}
	err = res.Scan(&username, &correct_hash)
	utils.PanicIfError(err)
	if correct_hash != hash {
		return 2
	}
	_, err = db.Exec("DELETE FROM posts WHERE id = '" + fmt.Sprint(id) + "' AND author = '" + username + "'")
	utils.PanicIfError(err)
	return 0
}
