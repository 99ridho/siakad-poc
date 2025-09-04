# Course Offering Management Technical Documentation

## Role

Admin, Koorprodi

## Endpoints

### GET /academic/course-offering

**Query Params:**

- page (default = 1)
- page_size (default = 10)

**Expected Success Responses Format (200):**

```
{
    "status": "success",
    "data": [
        {
            "id": "01f436c6-f6ae-4552-8184-5a6cd1a9f116",
            "course_name": "Pemrograman Dasar",
            "course_code": "151000"
            "section_code": "151011",
            "capacity": 50,
            "start_time": "2025-09-04T18:51:52Z",
            "end_time": "2025-09-04T18:51:52Z"
        }
    ],
    "paging": {
        "page": 1,
        "page_size": 10,
        "total_records": 100,
        "total_pages": 10
    }
}
```

### POST /academic/course-offering

**Example payload:**

```
{
    "course_id": "01f436c6-f6ae-4552-8184-5a6cd1a9f116",
    "semester_id": "01f436c6-f6ae-4552-8184-5a6cd1a9f116",
    "section_code": "10000",
    "capacity": 40,
    "start_time": "2025-09-04T18:51:52Z"
}
```

Validation:

- All attributes must be present
- Respect the unique constraint on DB (throw error if DB operation fails)

**Expected success response format (200):**

```
{
    "status": "success",
    "data": {
        "id": "01f436c6-f6ae-4552-8184-5a6cd1a9f116"
    }
}
```

**Response Error**

- When validation fails (HTTP 400)

### PUT /academic/course-offering/{id}

**Example payload:**

```
{
    "course_id": "01f436c6-f6ae-4552-8184-5a6cd1a9f116",
    "semester_id": "01f436c6-f6ae-4552-8184-5a6cd1a9f116",
    "capacity": 40,
    "section_code": "10000",
    "start_time": "2025-09-04T18:51:52Z"
}
```

Validation:

- All attributes must be present
- Respect the unique constraint on DB (throw error if DB operation fails)

**Expected success response format (200):**

```
{
    "status": "success",
    "data": {
        "id": "01f436c6-f6ae-4552-8184-5a6cd1a9f116"
    }
}
```

**Response Error**

- When not found (HTTP 404)
- When validation fails (HTTP 400)

### DELETE /academic/course-offering/{id}

**Expected success response:**

```
No content (HTTP code 204)
```

**Response Error**

- When not found (HTTP 404)
