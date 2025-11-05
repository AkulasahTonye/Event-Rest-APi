# 🚀 Go REST API

# A simple Restful APi Built with GoLang for demonstration and learning purposes.

## 🧩 Features
- Create Event
- Get All Events
- Get Single Event
- Update Event
- Delete Event
- Register for Event
- Cancel Registration
- Login User
- Create User


## 🛠️ Tech Stack
- **Language:** Go (Golang)
- **Framework:** [Gin](https://github.com/gin-gonic/gin)
- **Database:**  SQLite *(depending on setup)*
- **Dependency Management:** Go Modules Go Sum
- **Makefile:** For streamlined development

## ⚙️ Prerequisites
Before you begin, ensure you have the following installed:

- [Go 1.21+](https://go.dev/dl/)

## 📦 Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/<your-username>/<your-repo-name>.git
   cd <your-repo-name>
   
2. **Install dependencies**
   ```bash
   go mod  tidy
   

   
3. **Run the application**
   ```bash
   go run cmd/main.go

4. **Run the application using Makefile**
   ```bash
   make run
   
5. **Access the API** 
      ```bash
   http://localhost:8080