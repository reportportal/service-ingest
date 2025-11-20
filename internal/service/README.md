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
├── agent.go         # Agent business logic
├── launch.go        # Launch business logic
└── test.go          # Test item business logic
```

## Example Service

```go
package service

import (
    "context"
    "errors"
    "github.com/reportportal/service-ingest/internal/data/postgres"
    "github.com/reportportal/service-ingest/internal/model"
)

type AgentService struct {
    agentRepo *postgres.AgentRepository
}

func NewAgentService(repo *postgres.AgentRepository) *AgentService {
    return &AgentService{agentRepo: repo}
}

func (s *AgentService) Register(ctx context.Context, agent *model.Agent) error {
    // Business validation
    if err := agent.Validate(); err != nil {
        return err
    }

    // Check for duplicates
    existing, _ := s.agentRepo.FindByName(ctx, agent.Name)
    if existing != nil {
        return errors.New("agent with this name already exists")
    }

    // Business logic: set initial status
    agent.Status = model.AgentStatusActive

    // Persist
    return s.agentRepo.Save(ctx, agent)
}

func (s *AgentService) GetByID(ctx context.Context, id string) (*model.Agent, error) {
    agent, err := s.agentRepo.FindByID(ctx, id)
    if err != nil {
        return nil, errors.New("agent not found")
    }
    return agent, nil
}
```

## When to Add Interfaces

Start with concrete types. Add interfaces when you need:

```go
// Add interface when you need mocking for tests
type agentRepository interface {
    Save(ctx context.Context, agent *model.Agent) error
    FindByID(ctx context.Context, id string) (*model.Agent, error)
    FindByName(ctx context.Context, name string) (*model.Agent, error)
}

type AgentService struct {
    agentRepo agentRepository  // now uses interface
}
```

## Dependencies

Services depend on:

- **data layer** - for persistence (repositories)
- **model** - for domain models

Services should NOT:

- Know about HTTP details (status codes, headers, etc.)
- Parse HTTP requests or format responses
- Access database directly (use repositories)
- Contain framework-specific code
