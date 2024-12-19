package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Cтруктура для хранения данных статьи
type Article struct {
	Id       uint16
	Title    string
	Anons    string
	FullText string
}

var posts = []Article{}
var curPost = Article{}

/*
* Функция: connectDB
* Описание: Подключение к базе данных
* Возвращает:
*    @return: *sql.DB - объект для работы с базой данных
*    @return: error - объект для записи ошибки
 */
func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3307)/golang")
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to database")
	return db, nil
}

/*
* Функция: index
* Описание: Главная страница, отображение всех статей
* Параметры:
*    @param w http.ResponseWriter - объект для записи ответа на запрос
*    @param r *http.Request - объект запроса
* Возвращает:
*    @return: Нет возвращаемого значения
 */
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

/*
* Функция: create
* Описание: Создание страницы для создания статьи
* Параметры:
*    @param w http.ResponseWriter - объект для записи ответа на запрос
*    @param r *http.Request - объект запроса
* Возвращает:
*    @return: Нет возвращаемого значения
 */
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
* Описание: Сохранение статьи в базу данных
* Параметры:
*    @param w http.ResponseWriter - объект для записи ответа на запрос
*    @param r *http.Request - объект запроса
* Возвращает:
*    @return: Нет возвращаемого значения
 */
func saveArticle(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(r.FormValue("title"))
	anons := strings.TrimSpace(r.FormValue("anons"))
	fullText := strings.TrimSpace(r.FormValue("full_text"))

	if title == "" || anons == "" || fullText == "" ||
		isWhitespace(title) || isWhitespace(anons) || isWhitespace(fullText) {
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

func showPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	t, err := template.ParseFiles("pages/Show/show.html", "pages/templates/header.html", "pages/templates/footer.html")
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
	res, err := db.Query("SELECT * FROM `articles` WHERE `id` = ?", vars["id"])
	if err != nil {
		http.Error(w, "Failed to fetch articles", http.StatusInternalServerError)
		return
	}
	defer res.Close()

	curPost = Article{}
	for res.Next() {
		var post Article
		err := res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			http.Error(w, "Failed to scan article", http.StatusInternalServerError)
			return
		}

		curPost = post
	}

	fmt.Println()
	t.ExecuteTemplate(w, "show", curPost)

}

/**
* Функция: isWhitespace
* Описание: Проверка на пустоту строки
* Параметры:
*    @param s string - строка для проверки
* Возвращает:
*    @return: true если строка пуста, иначе false
 */
func isWhitespace(s string) bool {
	for _, r := range s {
		if r != ' ' {
			return false
		}
	}
	return true
}

/**
* Функция: handleFunc
* Описание: Обработка запросов (Маршрутизация)
* Возвращает:
*    @return: Нет возвращаемого значения
 */
func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", saveArticle).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", showPost).Methods("GET")

	http.Handle("/", rtr)
	http.Handle("/App/", http.StripPrefix("/App/", http.FileServer(http.Dir("App"))))
	// http.HandleFunc("/", index)
	// http.HandleFunc("/create", create)
	// http.HandleFunc("/save_article", saveArticle)

	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
