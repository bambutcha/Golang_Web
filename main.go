package main

import (
	"fmt"
	"net/http"
	"html/template"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pages/Home/index.html", "pages/templates/header.html", "pages/templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	t.ExecuteTemplate(w, "index", nil)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pages/Create/create.html", "pages/templates/header.html", "pages/templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, `<script>alert("Все поля должны быть заполнены!"); window.history.back();</script>`)
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3307)/golang")
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	fmt.Println("Successfully connected to database")

	// Использование подготовленного запроса
	_, err = db.Exec("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES (?, ?, ?)", title, anons, full_text)
	if err != nil {
		http.Error(w, "Failed to insert data", http.StatusInternalServerError)
		return
	}
	fmt.Println("Data inserted")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleFunc() {
	http.Handle("/App/", http.StripPrefix("/App/", http.FileServer(http.Dir("App"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/save_article", save_article)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
