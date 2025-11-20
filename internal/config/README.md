# Configuration

This module handles application configuration loading from environment variables,
config files, or other sources.

## Responsibilities

- Load configuration from environment variables
- Provide default values
- Validate configuration
- Expose configuration as a typed struct

## Structure

```text
config/
└── config.go        # Configuration loading and validation
```

## Best Practices

- Never commit secrets or credentials
- Provide sensible defaults for local development
- Validate configuration on startup
- Use typed configuration structs
- Document all environment variables
- Consider using a library like `viper` for complex config needs
