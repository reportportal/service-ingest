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
```

## Dependencies

Data layer depends on:

- **model** - for domain models
- Database drivers (postgres, redis, etc.)

Data layer should NOT:

- Contain business logic (belongs in service layer)
- Know about HTTP or presentation details
- Validate business rules (use service layer)
