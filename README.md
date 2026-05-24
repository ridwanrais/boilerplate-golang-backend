# Go REST API Boilerplate

This project is a production-ready Golang boilerplate built using **Gin** (Routing), **Huma** (OpenAPI/Swagger Generation), **Ent** (Entity ORM), and **Atlas** (Versioned Migrations). It is organized using **Clean Architecture** principles and uses **PostgreSQL** via Docker Compose.

---

## 1. First Time Setup

To get this project running locally for the first time, follow these steps:

0. **Install Go:**
   Ensure you have Go installed on your machine (version 1.25+ recommended). You can download it from [golang.org](https://go.dev/doc/install).

1. **Start PostgreSQL Database:**
   Spin up the Postgres database in the background using Docker Compose.
   ```bash
   docker-compose up -d
   ```

2. **Configure Environment Variables:**
   Copy the example environment file.
   ```bash
   cp .env.example .env
   ```

3. **Install Dependencies:**
   ```bash
   make tidy
   ```

4. **Install Atlas CLI:**
   You need the Atlas CLI to apply database migrations locally.
   ```bash
   make atlas-install
   ```

5. **Apply Database Migrations:**
   Run the SQL migrations against your local Postgres database.
   ```bash
   make apply
   ```

6. **Start the server:**
   ```bash
   make run
   ```
   The server will connect to the Postgres DB and start on `http://localhost:8080`.

---

## 2. How to Change the Database Schema

This project uses a "schema-first" approach via [Ent](https://entgo.io/). 

1. **Edit the Schema:**
   Modify or create new models inside the `ent/schema/` directory (e.g., `ent/schema/user.go`). You define your database columns and relations using Go code.

2. **Generate the Go ORM Code:**
   After altering any schema file, you MUST run:
   ```bash
   make generate
   ```
   This regenerates the strongly-typed Go client (`ent/client.go`), builders, and structs so you can use the new fields in your repository layer.

3. **Generate the SQL Migration Files:**
   To physically reflect your Go schema changes in the database, you need to generate the `.sql` migration files. Run:
   ```bash
   make migrate name=describe_your_changes
   ```
   *Atlas will compare your Go schema against the current migrations and spin up a lightweight throwaway Docker Postgres container to generate a perfect Postgres-compatible `.sql` diff inside `ent/migrate/migrations/`.*

---

## 3. How to Apply Database Migrations

Whenever you generate new `.sql` migration files (or pull new ones from a teammate via Git), you must apply them to your database.

To apply migrations locally, run:
```bash
make apply
```

*Note: For production, you will configure your CI/CD pipeline or deployment scripts to run `atlas migrate apply` pointing to your production database URL.*

---

## 4. API Documentation (Auto-Managed!)

You do **not** need to manually manage OpenAPI or Swagger YAML files! 🎉

This project uses **Huma** tightly integrated with Gin. The documentation is automatically generated at runtime based on the Go structs you define for your HTTP Handlers and the `doc`, `example`, and validation tags attached to them.

**To view the API Docs:**
1. Start the server (`make run`).
2. Open your browser and navigate to: [http://localhost:8080/docs](http://localhost:8080/docs).

If you want to add a new route to the documentation, simply use `huma.Register` in `internal/delivery/http/user_handler.go` and define your request/response struct types. Huma handles the rest!
