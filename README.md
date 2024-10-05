# Go User Caching API

This project is a Go-based API that demonstrates user profile management with caching using Redis and PostgreSQL for persistent data storage. The API supports CRUD operations and rate limiting for efficient handling of requests.

## Features

- Fetch user profile with Redis caching to improve read performance
- Paginated list of users with caching support
- Update user profile and automatically invalidate cache
- Search for users by name or email, caching search results
- Rate limiting to restrict the number of API requests
- PostgreSQL as the primary database
- Redis as the caching layer

## Prerequisites

Ensure you have the following installed:

- Go (1.18+)
- PostgreSQL (12+)
- Redis (5+)
- Docker (optional, if running PostgreSQL and Redis using containers)

## Project Setup

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/go-user-caching.git
cd go-user-caching
```

### 2. 

