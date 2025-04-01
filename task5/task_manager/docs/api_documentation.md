# Task Management API Documentation

## Overview

This REST API allows you to manage tasks with basic CRUD operations. Data is stored in memory.

## Endpoints

### 1. Get All Tasks

- **URL:** `/tasks`
- **Method:** `GET`
- **Response:**
  - **200 OK:** An array of tasks.

Example Response:

```json
[
  {
    "id": 1,
    "title": "Complete assignment",
    "description": "Finish the REST API assignment",
    "due_date": "2025-04-01",
    "status": "pending"
  }
]
```
