# Data Layer

This module contains the data access layer for the service-ingest application.
It provides implementations for interacting with various data sources, including databases,
caches, file systems, and external APIs.

## Responsibilities

- Implement data persistence operations (CRUD)
- Execute database queries
- Handle database connections and transactions
- Manage caching operations
- Interact with external APIs for data
- Map between database models and domain models (if needed)

## Structure

```text
data/
├── postgres/
│   ├── agent.go         # Agent repository implementation
│   ├── launch.go        # Launch repository implementation
│   └── db.go            # Database connection and utilities
└── redis/
    └── cache.go         # Cache repository implementation
```

## Database Connection

```go
// data/postgres/db.go
package postgres

import (
    "database/sql"
    "time"
    _ "github.com/lib/pq"
)

func NewDB(dsn string, maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration) (*sql.DB, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(maxOpenConns)
    db.SetMaxIdleConns(maxIdleConns)
    db.SetConnMaxLifetime(connMaxLifetime)

    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}
```

## Cache Repository Example

```go
// data/redis/cache.go
package redis

import (
    "context"
    "time"
    "github.com/redis/go-redis/v9"
)

type CacheRepository struct {
    client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
    return &CacheRepository{client: client}
}

func (r *CacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *CacheRepository) Get(ctx context.Context, key string) (string, error) {
    return r.client.Get(ctx, key).Result()
}

func (r *CacheRepository) Delete(ctx context.Context, key string) error {
    return r.client.Del(ctx, key).Err()
}
```

## Dependencies

Data layer depends on:

- **model** - for domain models
- Database drivers (postgres, redis, etc.)

Data layer should NOT:

- Contain business logic (belongs in service layer)
- Know about HTTP or presentation details
- Validate business rules (use service layer)
