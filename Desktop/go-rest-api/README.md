# Expense Tracker REST API

A simple REST API for managing personal expenses, built with Go and PostgreSQL.

The project provides basic CRUD operations for expenses, along with filtering and summary functionality. It uses Go's standard library for HTTP handling and PostgreSQL for data storage.

## Features

* Create, view, update, and delete expenses
* Filter expenses by category and date
* Validate request data and expense amounts
* Calculate basic expense summaries
* PostgreSQL database integration
* RESTful API using Go's `net/http`
* No external web framework

## Tech Stack

* **Go 1.21+**
* **net/http** — HTTP server and routing
* **PostgreSQL** — database
* **lib/pq** — PostgreSQL driver
* **database/sql** — database operations

## How to Run

### Prerequisites

Make sure you have the following installed:

* Go 1.21 or later
* PostgreSQL
* Git

### 1. Clone the repository

```bash
git clone <repository-url>
cd expense-tracker
```

### 2. Create the database

Open PostgreSQL and create a database named `expensetracker`:

```sql
CREATE DATABASE expensetracker;
```

Then connect to the database and create the `expenses` table:

```sql
CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    amount DECIMAL(10, 2) NOT NULL,
    category VARCHAR(100) NOT NULL,
    note TEXT,
    spent_on DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 3. Install dependencies

Run:

```bash
go mod tidy
```

The project uses `github.com/lib/pq` to connect Go with PostgreSQL.

### 4. Configure the database

Update the database connection settings in:

```text
config/db.go
```

with your local PostgreSQL credentials.

### 5. Start the API

```bash
go run main.go
```

Once the server starts, the API will be available on the configured port.

## API Overview

The API supports the following operations:

| Method   | Endpoint            | Description          |
| -------- | ------------------- | -------------------- |
| `POST`   | `/expenses`         | Create a new expense |
| `GET`    | `/expenses`         | Get all expenses     |
| `GET`    | `/expenses/{id}`    | Get an expense by ID |
| `PUT`    | `/expenses/{id}`    | Update an expense    |
| `DELETE` | `/expenses/{id}`    | Delete an expense    |
| `GET`    | `/expenses/summary` | Get expense summary  |

The exact request and response formats can be found in the handler implementation.

## Design Decisions

### Why Go's Standard Library?

I decided to use Go's standard `net/http` package instead of a third-party framework. For a small API like this, it keeps the project simple and makes the request handling and routing easier to understand.

It also gave me a chance to work directly with Go's built-in HTTP functionality instead of relying on a framework to handle everything.

### Why PostgreSQL?

PostgreSQL is a good fit for this project because expenses have a clear relational structure. It also makes filtering and aggregation queries straightforward, especially for generating expense summaries.

## Validation

The API validates incoming request data before saving it to the database.

For example:

* Amount must be greater than zero
* Required fields cannot be empty
* Dates must use the `YYYY-MM-DD` format
* Invalid request data returns `400 Bad Request`
* Requests for non-existing expenses return an appropriate `404 Not Found`

## Assumptions

For this version of the project, I kept the scope intentionally small:

* The API is designed for a single user.
* There is no authentication or authorization.
* Expense dates use the standard `YYYY-MM-DD` format.
* PostgreSQL is expected to be running locally.
* Database credentials are currently configured locally rather than loaded from environment variables.

## What I Would Improve

If I continued developing the project, I would make a few improvements:

* Move database credentials to environment variables
* Add unit and integration tests
* Add Docker and Docker Compose for easier setup
* Add pagination to the expenses list
* Improve API error responses with consistent JSON structures
* Add authentication if the API is extended to support multiple users
* Add more detailed expense statistics and reporting

## Project Structure

A simplified view of the project structure:

```text
expense-tracker/
│
├── config/
│   └── db.go
│
├── handlers/
│   └── ...
│
├── models/
│   └── ...
│
├── main.go
├── go.mod
└── go.sum
```

The project is intentionally kept small so that the main parts of a REST API — routing, request handling, validation, database operations, and responses — are easy to follow.
