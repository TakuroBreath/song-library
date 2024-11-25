# Song Library API

## Overview

Song Library is a comprehensive web application designed to manage and retrieve song information. The application provides a robust API for adding, updating, deleting, and retrieving song details with support for pagination and filtering.

## Features

- 🎵 Add songs with detailed information
- 🔍 Filter and search songs
- 📄 Pagination support
- 🌐 External API integration for song details
- 📊 Swagger documentation

## Technology Stack

- **Language**: Go (Golang)
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: Standard library `database/sql`
- **Logging**: `log/slog`
- **API Documentation**: Swagger

## Prerequisites

- Go 1.23+
- PostgreSQL

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/TakuroBreath/song-library.git
cd song-library
```

### 2. Set Up Environment Variables

Create a `.env` file in the root directory with the following variables:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=song_library
ENV=local
API_URL=https://external-song-api.com
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Database Migration

The application uses an automatic migration system. Ensure PostgreSQL is running and the database is created.

### 5. Run the Application

```bash
go run cmd/song-library/main.go
```

## API Endpoints

### Songs

- `GET /api/songs`: Retrieve songs with filtering and pagination
- `GET /api/songs/verses`: Get song verses with pagination
- `POST /api/songs`: Add a new song
- `PUT /api/songs`: Update existing song details
- `DELETE /api/songs`: Remove a song

## Swagger Documentation

Access Swagger UI at: `http://localhost:8080/swagger/index.html`

## Environment Configurations

The application supports three environments:
- `local`: Debug logging, text handler
- `dev`: Debug logging, JSON handler
- `production`: Info logging, JSON handler

## Error Handling

The application provides detailed error responses and logs for:
- Database connection issues
- API integration errors
- Validation failures
- Resource not found scenarios

## Logging

Comprehensive logging is implemented using `slog` with different configurations for each environment:
- Detailed debug logs in local/dev environments
- Minimal info logs in production

## Security

- Input validation for all endpoints
- Parameterized database queries to prevent SQL injection
- Environment-based configuration management

## Performance Considerations

- Pagination support for large datasets
- Efficient database queries
- Caching potential for frequently accessed resources

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Swagger](https://swagger.io/)
- [PostgreSQL](https://www.postgresql.org/)