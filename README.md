# Instashop API

A robust RESTful API for Instashop built with Go, featuring user authentication, product management, and order processing capabilities.

## 🚀 Features

### Authentication & Authorization
- User registration and login
- JWT-based authentication with access and refresh tokens
- Role-based access control (Admin/Customer)
- Session management

### Product Management
- CRUD operations for products
- Stock management
- Product categorization
- Price tracking

### Order Management
- Order creation and processing
- Order status tracking
- Order cancellation
- Multiple items per order

### Additional Features
- Structured error handling
- Comprehensive logging
- API documentation with Swagger
- Database transaction support 
- Input validation

## 🏁 Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 16
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/zde37/instashop-task.git
cd instashop-task
```

2. Setup Postgres manually on your pc or use docker(make sure docker is running)


`Create postgres docker container`

```bash
make postgres
```


`Create database`
```bash
make createdb
```

3. Create a `.env` file and enter these environments variables. You can use whatever value of your choice, just make sure the KEYS are the same
```bash
DSN=postgresql://root:4713a4cd628778cd1c37a95518f3eaf3@localhost:5432/Instashop_DB?sslmode=disable
PORT=4070 
JWY_SECRET_KEY=4713a4cd628778cd1c37a95518f3eaf3 
ENVIRONMENT=dev
```

4. Run the server
```bash
make run
```