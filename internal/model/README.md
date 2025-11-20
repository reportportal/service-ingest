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
├── agent.go         # Agent domain model
├── launch.go        # Launch domain model
├── test.go          # Test item domain model
└── errors.go        # Domain-specific errors (optional)
```

## Example Model

```go
package model

import (
    "errors"
    "time"
)

type AgentStatus string

const (
    AgentStatusActive   AgentStatus = "active"
    AgentStatusInactive AgentStatus = "inactive"
)

type Agent struct {
    ID        string      `json:"id" db:"id"`
    Name      string      `json:"name" db:"name"`
    Version   string      `json:"version" db:"version"`
    Status    AgentStatus `json:"status" db:"status"`
    CreatedAt time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}

// Validate performs domain validation
func (a *Agent) Validate() error {
    if a.Name == "" {
        return errors.New("agent name is required")
    }
    if len(a.Name) < 3 {
        return errors.New("agent name must be at least 3 characters")
    }
    if a.Version == "" {
        return errors.New("agent version is required")
    }
    return nil
}

// Activate sets the agent status to active
func (a *Agent) Activate() {
    a.Status = AgentStatusActive
    a.UpdatedAt = time.Now()
}

// Deactivate sets the agent status to inactive
func (a *Agent) Deactivate() {
    a.Status = AgentStatusInactive
    a.UpdatedAt = time.Now()
}

// IsActive checks if agent is currently active
func (a *Agent) IsActive() bool {
    return a.Status == AgentStatusActive
}
```

## Struct Tags

Use struct tags for different purposes:

- `json:` - JSON serialization/deserialization
- `db:` - Database column mapping
- `validate:` - Validation rules (if using validation library)

```go
type Entity struct {
    ID   string `json:"id" db:"id" validate:"required,uuid"`
    Name string `json:"name" db:"name" validate:"required,min=3"`
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
