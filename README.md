# Chirpy
📖 Overview
Chirpy is a lightweight microblogging platform inspired by Twitter, built with Go and PostgreSQL. This project serves as a learning exercise to understand backend development concepts including REST API design, JWT authentication, database operations, and webhook handling.
✨ Features

🔐 User Authentication - JWT-based authentication with access and refresh tokens
📝 Chirp Management - Create, read, update, and delete chirps (posts)
👥 User Management - User registration and profile updates
🔒 Admin Panel - Administrative endpoints for user metrics and database management
🪝 Webhook Integration - Support for external service webhooks (Polka integration)
🏥 Health Checks - System health monitoring endpoint

🛠️ Tech Stack

Language: Go
Database: PostgreSQL
Authentication: JWT (JSON Web Tokens)
HTTP Router: Go's built-in net/http

Dependencies
gogithub.com/google/uuid         // UUID generation
github.com/joho/godotenv      // Environment variable loading
github.com/golang-jwt/jwt/v5  // JWT token handling
🚀 Getting Started
Prerequisites

Go 1.21 or higher
PostgreSQL database
Git

Installation

1. Clone the repository
```
git clone https://github.com/yourusername/chirpy.git
cd chirpy
```

2. Install dependencies
```
go mod download
```

3. Set up environment variables
Create a .env file in the root directory:
```
DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key
POLKA_KEY=your-polka-webhook-key
```

4. Set up PostgreSQL database
```
CREATE DATABASE chirpy;
```

5. Run the application
```
go run main.go

```


The server will start on http://localhost:8080 (or your configured port).
📚 API Documentation
Base URL
```
http://localhost:8080
```
Authentication
Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```
Admin endpoints require an API key:
```
Authorization: ApiKey <your-api-key>
```

🏥 Health & Admin
Check System Health
```
GET /api/healthz
```
Response: "OK"
```
Get User Metrics (Admin)
```
GET /admin/metrics
Headers: Authorization: ApiKey <api_key>
Response:
json{
  "users": 5
}
Reset Database (Admin)
```
POST /admin/reset
```
Headers: Authorization: ApiKey <api_key>
Response: Empty (200 OK)

👤 User Management
Register New User
httpPOST /api/users
Request Body:
json{
  "email": "newuser@example.com",
  "password": "password123"
}
Response:
json{
  "id": 1,
  "email": "newuser@example.com"
}
Update User Profile
httpPUT /api/users
Headers: Authorization: Bearer <token>
Request Body:
json{
  "email": "updated@example.com",
  "password": "newpassword"
}
Response:
json{
  "id": 1,
  "email": "updated@example.com"
}

🔐 Authentication
User Login
httpPOST /api/login
Request Body:
json{
  "email": "user@example.com",
  "password": "password123"
}
Response:
json{
  "id": 1,
  "email": "user@example.com",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
Refresh Access Token
httpPOST /api/refresh
Request Body:
json{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
Response:
json{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
Revoke Refresh Token
httpPOST /api/revoke
Headers: Authorization: Bearer <token>
Response: Empty (200 OK)

🐦 Chirp Management
Create New Chirp
httpPOST /api/chirps
Headers: Authorization: Bearer <token>
Request Body:
json{
  "body": "This is my first chirp!"
}
Response:
json{
  "id": 1,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "body": "This is my first chirp!",
  "author_id": 2
}
Get All Chirps
httpGET /api/chirps?author_id=2
Query Parameters:

author_id (optional) - Filter chirps by author

Response:
json[
  {
    "id": 1,
    "body": "Chirp 1",
    "author_id": 2
  },
  {
    "id": 3,
    "body": "Chirp 2",
    "author_id": 5
  }
]
Get Single Chirp
httpGET /api/chirps/{chirpID}
Response:
json{
  "id": 1,
  "body": "This is a chirp!",
  "author_id": 2
}
Delete Chirp
httpDELETE /api/chirps/{chirpID}
Headers: Authorization: Bearer <token>
Response: Empty (200 OK)

🪝 Webhooks
Polka Webhook Handler
httpPOST /api/polka/webhooks
Headers: X-Polka-Signature: <signature>
Request Body:
json{
  "event": "user.upgraded",
  "data": {
    "chirpy_user_id": 123
  }
}
Response: Empty (204 No Content)

📊 Status Codes
CodeDescription200OK201Created204No Content400Bad Request401Unauthorized403Forbidden404Not Found500Internal Server Error
🧪 Testing
The project includes package-specific tests. Run tests with:
bashgo test ./...
🔒 Security

JWT tokens expire after 1 hour
Refresh tokens provide secure token renewal
Admin endpoints require API key authentication
Webhook signatures are validated for security

🤝 Contributing
This is a learning project, but contributions are welcome! Please feel free to:

Fork the repository
Create a feature branch
Make your changes
Submit a pull request

📝 License
This project is open source and available under the MIT License.
🙏 Acknowledgments

Built as part of learning Go backend development
Inspired by Twitter's microblogging concept
Uses industry-standard authentication practices
