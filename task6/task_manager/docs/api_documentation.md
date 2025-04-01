# Task Management API Documentation

## Overview

This API allows you to manage tasks with CRUD operations. In this enhanced version, tasks are persisted in MongoDB so that data remains available across API restarts.

## MongoDB Integration

- **MongoDB URI:** `mongodb://localhost:27017`
- **Database:** `task_manager_db`
- **Collection:** `tasks`

Make sure your MongoDB instance is running before starting the API. You can run a local MongoDB instance or use a cloud service.

## Endpoints

### 1. Get All Tasks

- **URL:** `/tasks`
- **Method:** `GET`
- **Response:**
  - **200 OK:** An array of task objects.

Example Response:

```json
[
  {
    "id": "606d1c3f72f9a5d1f1e12345",
    "title": "Complete assignment",
    "description": "Finish the Task Management API assignment",
    "due_date": "2025-04-01",
    "status": "pending"
  }
]
```
