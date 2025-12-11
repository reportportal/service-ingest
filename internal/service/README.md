# Service Layer

This module contains the business logic for the service-ingest application.
Services orchestrate operations, enforce business rules, and coordinate between
handlers and the data layer.

## Responsibilities

- Implement business logic and rules
- Validate domain models
- Coordinate multiple repository operations
- Handle transactions
- Transform data between layers when needed
- Return domain errors

## Structure

```text
service/
├── launch.go        # Launch business logic
├── launch_repo.go   # Launch repository interface
├── item.go          # Test item business logic
├── item_repo.go     # Test item repository interface
├── log.go           # Test log business logic
└── log_repo.go      # Test log repository interface
```

## When to Add Interfaces

Start with concrete types. Add interfaces when you need.

## Dependencies

Services depend on:

- **data layer** - for persistence (repositories)
- **model** - for domain models

Services should NOT:

- Know about HTTP details (status codes, headers, etc.)
- Parse HTTP requests or format responses
- Access database directly (use repositories)
- Contain framework-specific code
