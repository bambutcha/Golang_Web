// Assembly - XAMPP
// Server - Apache
// Database - MySQL

package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name string `json:"name"`
	Age uint16 `json:"age"`
}

func mysql() {
	// Подключение к базе данных
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3307)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Printf("Success connected to database\n\n")

	// Установка данных
	insert, err := db.Query("INSERT INTO `users` (`name`, `age`) VALUES ('Bob', 35)")
	if err != nil {
		panic(err)
	}
	defer insert.Close()
	fmt.Println("Data inserted")

	// Выборка данных
	res, err := db.Query("SELECT `name`, `age` FROM `users`")
	if err != nil {
		panic(err)
	}

	fmt.Println("Info about users")
	for res.Next() {
		var user User
		err := res.Scan(&user.Name, &user.Age)
		if err != nil {
			panic(err)
		}
		fmt.Printf("User: %s with age %d\n", user.Name, user.Age)
	}
}
