# Task Management API Documentation

## Overview

This API manages tasks and now includes user authentication and authorization using JWT. Users can register and log in; upon login, a JWT token is returned. Protected endpoints require a valid token. Only admin users (the first registered user becomes admin) can create, update, or delete tasks and promote users.

## Endpoints

### User Endpoints

#### 1. Register a New User

- **URL:** `/register`
- **Method:** `POST`
- **Request Payload:**

```json
{
  "username": "johndoe",
  "password": "your_password"
}
```
