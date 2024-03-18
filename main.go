package main

import Db "shitblog-server/db"
import "fmt"

func main() {
	db := Db.ConnectToDb()
	Db.RecreateTables(db)
	fmt.Println(Db.CreateUser("proggerx", db))
}
