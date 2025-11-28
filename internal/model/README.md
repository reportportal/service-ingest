# Model Layer

This module contains the domain models (entities) for the service-ingest application.
Models represent the core business objects and contain domain logic.

## Responsibilities

- Define domain entities and their structure
- Implement domain validation logic
- Provide business methods on entities
- Define value objects and enums

## Structure

```text
model/
├── launch.go        # Launch domain model
├── test.go          # Test item domain model
├── step.go          # Nested step domain model
├── log.go           # Log domain model
└── errors.go        # Domain-specific errors (optional)
```

## Example Model

```go
package model

import (
    "errors"
    "time"
)

type EntityStatus string

const (
    EntityStatusActive   EntityStatus = "active"
	EntityStatusInactive EntityStatus = "inactive"
)

type Entity struct {
    ID        string       `json:"id"`
    Name      string       `json:"name"`
    Status    EntityStatus `json:"status"`
    CreatedAt time.Time    `json:"created_at"`
    UpdatedAt time.Time    `json:"updated_at"`
}

func (a *Entity) Validate() error {
    if a.Name == "" {
        return errors.New("name is required")
    }
    if len(a.Name) < 3 {
        return errors.New("name must be at least 3 characters")
    }
    return nil
}

func (a *Entity) Activate() {
    a.Status = EntityStatusActive
    a.UpdatedAt = time.Now()
}

func (a *Entity) Deactivate() {
    a.Status = EntityStatusInactive
    a.UpdatedAt = time.Now()
}

func (a *Entity) IsActive() bool {
    return a.Status == EntityStatusActive
}
```

## Struct Tags

Use struct tags for different purposes:

- `json:` - JSON serialization/deserialization
- `validate:` - Validation rules (if using validation library)

```go
package model

type Entity struct {
    ID   string `json:"id" validate:"required,uuid"`
    Name string `json:"name" validate:"required,min=3"`
}
```

## Domain Errors

Define domain-specific errors when needed:

```go
package model

import "errors"

var (
    ErrAgentNotFound    = errors.New("agent not found")
    ErrInvalidAgent     = errors.New("invalid agent")
    ErrDuplicateAgent   = errors.New("agent already exists")
)
```

## Guidelines

Models should:

- Contain domain logic and validation
- Be independent of external frameworks
- Use standard library types when possible
- Have clear, descriptive field names

Models should NOT:

- Know about HTTP, databases, or other infrastructure
- Contain business workflow logic (use service layer)
- Have dependencies on other layers
