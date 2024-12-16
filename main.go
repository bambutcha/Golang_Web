package main

import (
	"fmt"
	"net/http"
	"html/template"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

)

type Article struct {
	Id       uint16
	Title    string
	Anons    string
	FullText string
}

var posts = []Article{} 

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3307)/golang")
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to database")
	return db, nil
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pages/Home/index.html", "pages/templates/header.html", "pages/templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	db, err := connectDB()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Выборка данных
	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		http.Error(w, "Failed to fetch articles", http.StatusInternalServerError)
		return
	}
	defer res.Close()

	posts = []Article{}
	for res.Next() {
		var post Article
		err := res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			http.Error(w, "Failed to scan article", http.StatusInternalServerError)
			return
		}

		posts = append(posts, post)
	}

	fmt.Println()
	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pages/Create/create.html", "pages/templates/header.html", "pages/templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	t.ExecuteTemplate(w, "create", nil)
}

/*
* Функция: saveArticle
* Описание: 
*/
func saveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fullText := r.FormValue("full_text")

	if title == "" || anons == "" || fullText == "" {
		fmt.Fprintf(w, `<script>alert("Все поля должны быть заполнены!"); window.history.back();</script>`)
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	db, err := connectDB()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Использование подготовленного запроса
	_, err = db.Exec("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES (?, ?, ?)", title, anons, fullText)
	if err != nil {
		http.Error(w, "Failed to insert data", http.StatusInternalServerError)
		return
	}
	fmt.Println("Data inserted")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Функция для обработки запроса (Маршрутизация)
func handleFunc() {
	http.Handle("/App/", http.StripPrefix("/App/", http.FileServer(http.Dir("App"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/save_article", saveArticle)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
