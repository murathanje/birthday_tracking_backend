# Birthday Tracking Backend

## Overview

Birthday Tracking Backend is a robust Go-based API service designed to help users manage and track birthdays with ease. Built using the Gin framework and PostgreSQL, this application provides secure user authentication, birthday management, and categorization features.

## Features

- 🔐 Secure User Authentication
  - JWT-based authentication for user operations
  - API Key authentication for admin endpoints
- 👥 User Management
  - User registration
  - Profile management
  - Account deletion
- 🎂 Birthday Management
  - Create, read, update, and delete birthday records
  - Categorize birthdays with custom string-based tags
  - Track birthdays for different groups (Family, Friends, Work, etc.)

## Technology Stack

- **Language**: Go (Golang)
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Documentation**: Swagger UI
- **Environment Management**: godotenv

## Prerequisites

- Go 1.16+
- PostgreSQL
- Git

## Installation

1. Clone the repository:
```bash
git clone https://github.com/murathanje/birthday_tracking_backend.git
cd birthday_tracking_backend
```

2. Create a `.env` file with the following configurations:
```env
# Database Configuration
DATABASE_HOST=localhost
DATABASE_USER=postgres
DATABASE_PASSWORD=your_password
DATABASE_NAME=birthday_db

# Server Configuration
SERVER_PORT=5050
GIN_MODE=debug

# Security
API_KEY=your_secret_api_key
JWT_SECRET=your_jwt_secret
```

3. Install dependencies:
```bash
go mod tidy
```

4. Generate Swagger documentation:
```bash
swag init -g cmd/server/main.go
```

## Running the Application

### Development Mode
```bash
go run cmd/server/main.go
```

### Production Mode
```bash
GIN_MODE=release go run cmd/server/main.go
```

## API Endpoints

### Authentication
- `POST /api/v1/register`: Create a new user account
- `POST /api/v1/login`: Authenticate and receive JWT token

### User Management
- `GET /api/v1/users/me`: Get current user profile
- `PUT /api/v1/users/me`: Update user profile
- `DELETE /api/v1/users/me`: Delete user account

### Birthday Management
- `POST /api/v1/birthdays`: Create a new birthday record
- `GET /api/v1/birthdays`: List all user's birthdays
- `GET /api/v1/birthdays/{id}`: Get a specific birthday
- `PUT /api/v1/birthdays/{id}`: Update a birthday record
- `DELETE /api/v1/birthdays/{id}`: Delete a birthday record

### Admin Endpoints
- `GET /api/v1/admin/users`: List all users (requires API Key)
- `GET /api/v1/admin/users/{id}`: Get user details (requires API Key)
- `PUT /api/v1/admin/users/{id}`: Update user (requires API Key)
- `DELETE /api/v1/admin/users/{id}`: Delete user (requires API Key)

## Swagger Documentation

Access the Swagger UI at: `https://managing-celle-trilema-d4ef42f0.koyeb.app/swagger/index.html`

## Environment Variables

| Variable           | Description                          | Default Value     |
|--------------------|--------------------------------------|-------------------|
| `DATABASE_HOST`    | PostgreSQL database host             | `localhost`       |
| `DATABASE_USER`    | PostgreSQL database username         | `postgres`        |
| `DATABASE_PASSWORD`| PostgreSQL database password         | `""`              |
| `DATABASE_NAME`    | PostgreSQL database name             | `birthday_db`     |
| `SERVER_PORT`      | Port for the API server              | `5050`            |
| `GIN_MODE`         | Gin framework mode (debug/release)   | `debug`           |
| `API_KEY`          | Secret key for admin operations      | `default-api-key` |
| `JWT_SECRET`       | Secret key for JWT token generation  | `default-jwt-secret` |

## Deployment

The application is currently deployed on Koyeb at:
`https://managing-celle-trilema-d4ef42f0.koyeb.app`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the MIT License.

## Contact

Project Link: [https://github.com/murathanje/birthday_tracking_backend](https://github.com/murathanje/birthday_tracking_backend)

---

Made with ❤️ by Murathan 