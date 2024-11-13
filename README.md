# Chirpy API Documentation

## Overview

Chirpy is a microblogging platform API that allows users to post short messages called "chirps". It features user authentication, content moderation, and premium user features.

## Base URL
```
http://localhost:8080
```

## Authentication

Chirpy uses JWT-based authentication with refresh tokens. Most endpoints require authentication via a Bearer token.

### Authentication Flow
1. Create a user account or login to receive both access and refresh tokens
2. Use the access token in the Authorization header for protected endpoints
3. When the access token expires, use the refresh token to get a new access token
4. Tokens can be revoked for security

### Token Usage
Include the token in the Authorization header:
```
Authorization: Bearer <your-token>
```

## Endpoints

### User Management

#### Create User
```http
POST /api/users
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "your-password"
}
```

**Response** (201 Created)
```json
{
    "id": "uuid",
    "email": "user@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "is_chirpy_red": false
}
```

#### Login
```http
POST /api/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "your-password"
}
```

**Response** (200 OK)
```json
{
    "id": "uuid",
    "email": "user@example.com",
    "token": "jwt-access-token",
    "refresh_token": "refresh-token",
    "is_chirpy_red": false
}
```

#### Update User
```http
PUT /api/users
Authorization: Bearer <access-token>
Content-Type: application/json

{
    "email": "new-email@example.com",
    "password": "new-password"
}
```

**Response** (200 OK)
```json
{
    "id": "uuid",
    "email": "new-email@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "is_chirpy_red": false
}
```

### Token Management

#### Refresh Access Token
```http
POST /api/refresh
Authorization: Bearer <refresh-token>
```

**Response** (200 OK)
```json
{
    "token": "new-jwt-access-token"
}
```

#### Revoke Refresh Token
```http
POST /api/revoke
Authorization: Bearer <refresh-token>
```

**Response** (204 No Content)

### Chirps

#### Create Chirp
```http
POST /api/chirps
Authorization: Bearer <access-token>
Content-Type: application/json

{
    "body": "Your chirp message"
}
```

**Notes:**
- Maximum length: 140 characters
- Automatic content moderation filters certain words
- Filtered words: "kerfuffle", "sharbert", "fornax" (replaced with "****")

**Response** (201 Created)
```json
{
    "id": "uuid",
    "body": "Your chirp message",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "user_id": "user-uuid"
}
```

#### Get All Chirps
```http
GET /api/chirps?author_id={uuid}&sort={asc|desc}
```

**Query Parameters:**
- `author_id` (optional): Filter chirps by author UUID
- `sort` (optional): Sort order ("asc" or "desc", defaults to "asc")

**Response** (200 OK)
```json
[
    {
        "id": "uuid",
        "body": "Chirp message",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "user_id": "user-uuid"
    }
]
```

#### Get Single Chirp
```http
GET /api/chirps/{chirpID}
```

**Response** (200 OK)
```json
{
    "id": "uuid",
    "body": "Chirp message",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "user_id": "user-uuid"
}
```

#### Delete Chirp
```http
DELETE /api/chirps/{chirpID}
Authorization: Bearer <access-token>
```

**Notes:**
- Users can only delete their own chirps
- Requires authentication
- Returns 204 on success

**Response** (204 No Content)

### Premium Features

#### Upgrade to Chirpy Red
```http
POST /api/polka/webhooks
Content-Type: application/json
Authorization: ApiKey <polka-api-key>
```

### Admin Endpoints

#### View Metrics
```http
GET /admin/metrics
```

**Response** (200 OK)
```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited {count} times!</p>
  </body>
</html>
```

#### Reset Metrics (Development Only)
```http
POST /admin/reset
```

**Notes:**
- Only available in development environment
- Resets visit counter
- Clears user database

**Response** (200 OK)

### Health Check

#### API Health Check
```http
GET /api/healthz
```

**Response** (200 OK)
```
OK
```

## Error Responses

All error responses follow this format:
```json
{
    "error": "Error message description"
}
```

Common HTTP Status Codes:
- 400: Bad Request (invalid input)
- 401: Unauthorized (missing or invalid token)
- 403: Forbidden (insufficient permissions)
- 404: Not Found
- 500: Internal Server Error

## Rate Limiting
Currently, no rate limiting is implemented.

## Development Tools Used
- Database: PostgreSQL
- Query Generation: sqlc
- Migrations: goose
- Authentication: JWT with refresh tokens
- UUID Generation: google/uuid

