# Thmanyah

A video management API built with Go and PostgreSQL.

## Prerequisites

- Go 1.24.3 or later
- PostgreSQL 12 or later
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI tool for database migrations

## Setup

### 1. Clone the repository

```bash
git clone https://github.com/JLL32/thmanyah.git
cd thmanyah
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Database setup

Create a PostgreSQL database:

```bash
createdb thmanyah
```

### 4. Environment configuration

Create a `.envrc` file in the project root with your database configuration:

```bash
export THMANYAH_DB_DSN="postgres://username:password@localhost:5432/thmanyah?sslmode=disable"
```

Replace `username` and `password` with your PostgreSQL credentials.

If you're using [direnv](https://direnv.net/), run:

```bash
direnv allow
```

Otherwise, source the file manually:

```bash
source .envrc
```

### 5. Run database migrations

```bash
make db/migrations/up
```

## Running the Application

### Development

Start the API server in development mode:

```bash
make run/api
```

The server will start on port 4000 by default. You can access the API at `http://localhost:4000`.

### Production Build

Build the application for production:

```bash
make build/api
```

This creates binaries in the `./bin/` directory for both local and Linux AMD64 architectures.

## Available Make Commands

- `make help` - Show all available commands
- `make run/api` - Run the API server
- `make db/psql` - Connect to the database using psql
- `make db/migrations/new name=<migration_name>` - Create a new database migration
- `make db/migrations/up` - Apply all pending migrations
- `make vendor` - Download and vendor dependencies
- `make audit` - Run code quality checks (format, vet, test)
- `make build/api` - Build the application

## Configuration

The application accepts the following command-line flags:

- `-port` - API server port (default: 4000)
- `-env` - Environment (development|staging|production) (default: development)
- `-db-dsn` - PostgreSQL connection string
- `-db-max-open-conns` - Maximum open database connections (default: 25)
- `-db-max-idle-conns` - Maximum idle database connections (default: 25)
- `-db-max-idle-time` - Maximum connection idle time (default: 15m)

## Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema management.

### Creating new migrations

```bash
make db/migrations/new name=add_users_table
```

### Applying migrations

```bash
make db/migrations/up
```

## Development Workflow

1. Make your changes
2. Run quality checks: `make audit`
3. Test your changes
4. Create migrations if needed: `make db/migrations/new name=<description>`
5. Apply migrations: `make db/migrations/up`

## Project Structure

```
thmanyah/
├── cmd/api/           # Application entry point
├── internal/          # Internal application code
├── migrations/        # Database migration files
├── Makefile          # Build and development commands
├── go.mod            # Go module file
└── README.md         # This file
```
