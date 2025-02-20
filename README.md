# Go Blog Platform

A full-featured blog platform built with Go, featuring article creation, viewing, and management capabilities. This project demonstrates the implementation of a web application using Go, MySQL, and the Gorilla Mux router.

![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8?style=flat&logo=go)
![MySQL](https://img.shields.io/badge/MySQL-8.0+-00000F?style=flat&logo=mysql)

---

## 🌟 Features

- Create and publish articles
- View list of all articles
- Detailed article view
- Server-side data validation
- Secure database interactions
- Responsive design
- Template-based views

---

## 🔧 Prerequisites

Before running this project, make sure you have the following installed:

- Go 1.16 or higher
- MySQL 8.0 or higher
- Git

---

## 📦 Dependencies

- [Gorilla Mux](https://github.com/gorilla/mux) - For HTTP routing
- [Go-MySQL-Driver](https://github.com/go-sql-driver/mysql) - MySQL driver for Go

---

## 🚀 Installation

1. Clone the repository:
```bash
git clone <your-repository-url>
cd www
```

2. Install the required Go packages:
```bash
go mod download
```

3. Set up the MySQL database:
```sql
CREATE DATABASE golang;
USE golang;

CREATE TABLE articles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    anons TEXT NOT NULL,
    full_text TEXT NOT NULL
);
```

4. Configure the database connection in `main.go`:
```go
db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3307)/golang")
```
Replace the connection string with your MySQL credentials if needed.

-

## 🎯 Running the Application

1. Start the server:
```bash
go run main.go
```

2. Open your web browser and navigate to:

---

## 📁 Project Structure

```
blog-platform/
├── main.go          # Main application file
├── pages/           # HTML templates
│   ├── Home/        # Home page templates
│   ├── Create/      # Article creation templates
│   ├── Show/        # Article view templates
│   └── templates/   # Shared templates
└── App/             # Static assets

---

## 🔍 API Endpoints

- GET / - Home page, displays all articles
- GET /create - Show article creation form
- POST /save_article - Save new article
- GET /post/{id} - View specific article

---

## 🛡️ Security Features
- SQL injection prevention using prepared statements
- Input validation
- Error handling
- Safe template rendering

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
