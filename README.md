# 🚀 GoTask — Secure Task Management API

> A secure and modular **Task Management REST API** built with **Go (Golang)** and **PostgreSQL**. 🔐
> The project implements **JWT Authentication, Refresh Tokens, Logout, Role-Based Access Control (RBAC), Task CRUD, Pagination, Filtering, and Sorting**.

---

## ✨ Features

### 🔐 Authentication & Security

* 👤 User Registration
* 🔑 Secure Login with JWT
* 🎫 Access Token Authentication
* 🔄 Refresh Token Support
* 🚪 Logout functionality
* 🚨 Logout from all devices
* 🔒 Protected API routes

### 👥 Role-Based Access Control (RBAC)

The system supports three user roles:

* 👑 **Admin**
* 🧑‍💼 **Manager**
* 👨‍💻 **Employee**

Users have different permissions based on their roles.

### 📋 Task Management

* ➕ Create Tasks
* 📖 Get All Tasks
* 🔎 Get Task by ID
* ✏️ Update Tasks
* 🗑️ Delete Tasks
* 👤 Assign tasks to employees
* 🔐 Role-based task access

### ⚙️ Additional Features

* 🏗️ Clean Architecture
  `Handler → Service → Repository`
* 🗃️ PostgreSQL Database
* 📄 Pagination
* 🔍 Filtering
* ↕️ Sorting
* ✅ Input Validation using `go-playground/validator`
* 🚨 Structured Error Handling
* 🌱 Environment-based Configuration

---

## 👥 RBAC Permissions

| Feature                | 👑 Admin | 🧑‍💼 Manager |    👨‍💻 Employee   |
| ---------------------- | :------: | :-----------: | :-----------------: |
| ➕ Create Task          |     ✅    |       ✅       |          ❌          |
| 📋 View All Tasks      |     ✅    |       ✅       | Assigned Tasks Only |
| 🔎 View Task by ID     |     ✅    |       ✅       | Assigned Tasks Only |
| ✏️ Update Task         |     ✅    |       ✅       | Assigned Tasks Only |
| 📝 Update Task Details |     ✅    |       ✅       |          ❌          |
| ☑️ Update Task Status  |     ✅    |       ✅       |          ✅          |
| 🗑️ Delete Task        |     ✅    |       ✅       |          ❌          |

---

## 🏗️ Project Architecture

```text
gotask/
│
├── 📂 cmd/
│   └── 📂 server/
│       └── main.go              # 🚀 Application entry point
│
├── 📂 internal/
│   │
│   ├── 📂 auth/                 # 🔐 Authentication & RBAC
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── middleware.go
│   │   ├── routes.go
│   │   └── context.go
│   │
│   └── 📂 task/                 # 📋 Task management logic
│       ├── handler.go
│       ├── service.go
│       ├── repository.go
│       ├── model.go
│       └── routes.go
│
├── 📂 pkg/
│   ├── 📂 config/               # ⚙️ Environment configuration
│   ├── 📂 db/                   # 🗃️ Database connection
│   ├── 📂 response/             # 📦 JSON response helpers
│   └── 📂 validation/           # ✅ Input validation
│
├── 🔒 .env                      # Environment variables
├── 📄 go.mod
└── 📖 README.md
```

---

## 🛠️ Tech Stack

| Technology                    | Purpose             |
| ----------------------------- | ------------------- |
| 🐹 **Go (Golang)**            | Backend Development |
| 🐘 **PostgreSQL**             | Database            |
| 🔐 **JWT**                    | Authentication      |
| 🧩 **Chi Router**             | HTTP Routing        |
| 🔒 **bcrypt**                 | Password Hashing    |
| ✅ **go-playground/validator** | Input Validation    |
| 🌱 **Environment Variables**  | Configuration       |

---

## 📋 Requirements

Before running the project, make sure you have:

* 🐹 Go **1.21+**
* 🐘 PostgreSQL
* 🧪 Postman or Thunder Client *(optional, for API testing)*

---

# 🚀 Getting Started

## 1️⃣ Clone the Repository

```bash
git clone https://github.com/jiyasingh9336-lab/gotask.git
cd gotask
```

## 2️⃣ Create Environment Variables

Create a `.env` file in the root directory:

```env
PORT=8080

URL=postgres://YOUR_USERNAME:YOUR_PASSWORD@localhost:5432/gotaskdb?sslmode=disable

JWT_SECRET=yourSuperSecretKey
```

> ⚠️ Never commit your `.env` file to GitHub.

---

## 3️⃣ Create the Database 🐘

```sql
CREATE DATABASE gotaskdb;
```

---

## 4️⃣ Create Required Tables 🗃️

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'employee',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_by BIGINT NOT NULL,
    assigned_to BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    revoked BOOLEAN NOT NULL DEFAULT FALSE
);
```

---

## 5️⃣ Run the Server 🚀

```bash
go run cmd/server/main.go
```

The server will start at:

```text
http://localhost:8080
```

🎉 **Your API is now running!**

---

# 🔌 API Endpoints

## 🔐 Authentication

| Method | Endpoint           | Description                        |
| ------ | ------------------ | ---------------------------------- |
| `POST` | `/auth/register`   | 👤 Register a new user             |
| `POST` | `/auth/login`      | 🔑 Login and receive tokens        |
| `POST` | `/auth/refresh`    | 🔄 Generate a new access token     |
| `POST` | `/auth/logout`     | 🚪 Logout and revoke refresh token |
| `POST` | `/auth/logout-all` | 🚨 Logout from all devices         |

---

## 📋 Task APIs

All task routes require an **Access Token** 🔐

```text
Authorization: Bearer <access_token>
```

| Method   | Endpoint      | Description       |
| -------- | ------------- | ----------------- |
| `POST`   | `/tasks`      | ➕ Create a task   |
| `GET`    | `/tasks`      | 📋 Get tasks      |
| `GET`    | `/tasks/{id}` | 🔎 Get task by ID |
| `PUT`    | `/tasks/{id}` | ✏️ Update a task  |
| `DELETE` | `/tasks/{id}` | 🗑️ Delete a task |

---

## 🔍 Pagination, Filtering & Sorting

Example request:

```text
GET /tasks?page=1&limit=10&completed=false&sort=created_at&order=desc
```

### Supported Query Parameters

| Parameter   | Description                           |
| ----------- | ------------------------------------- |
| `page`      | 📄 Page number                        |
| `limit`     | 🔢 Number of tasks per page           |
| `completed` | ☑️ Filter completed/uncompleted tasks |
| `sort`      | ↕️ Sort tasks by a field              |
| `order`     | ⬆️ `asc` or ⬇️ `desc`                 |

---

# 🔐 Authentication Flow

```text
👤 Register
     │
     ▼
🔑 Login
     │
     ▼
🎫 Access Token + 🔄 Refresh Token
     │
     ▼
🔒 Access Protected Routes
     │
     ▼
⏳ Access Token Expires
     │
     ▼
🔄 Use Refresh Token
     │
     ▼
✨ Receive New Access Token
```

---

# 🧠 What I Learned

While building **GoTask**, I worked with:

* 🐹 Building REST APIs using Go
* 🏗️ Clean Architecture principles
* 🗃️ PostgreSQL database integration
* 🔐 JWT Authentication
* 🔄 Access & Refresh Token flow
* 🔒 Password hashing with bcrypt
* 👥 Role-Based Access Control (RBAC)
* 🛡️ Authentication middleware
* 📄 Pagination
* 🔍 Filtering and sorting
* 🧩 Repository pattern
* 🚨 Error handling and validation
* 🌿 Environment-based configuration
* 🌳 Git and GitHub workflow

---

# 🗺️ Roadmap

### Completed 🎉

* [x] 🗃️ PostgreSQL Integration
* [x] 📋 Task CRUD APIs
* [x] 🔐 JWT Authentication
* [x] 🔄 Refresh Tokens
* [x] 🚪 Logout
* [x] 🚨 Logout from All Devices
* [x] 👥 Role-Based Access Control
* [x] 📄 Pagination
* [x] 🔍 Filtering
* [x] ↕️ Sorting

### Coming Next 🚀

* [ ] 📚 Swagger API Documentation
* [ ] 🐳 Docker & Docker Compose
* [ ] 🧪 Unit Tests
* [ ] 🔗 Integration Tests
* [ ] ⚙️ GitHub Actions CI/CD

---

## 👩‍💻 Author

Made with 💻, ☕ and lots of debugging 🐛 by **Jeeya Kumari** ✨

⭐ If you found this project interesting, consider giving it a star!
