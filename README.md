# Ingest service

> This is a proof of concept for validating the hypothesis of writing agent data to Parquet files for subsequent analytics. The primary goal is to build a lightweight, self-contained service that can be run independently in any environment.

Ingest service for ReportPortal agents data

## Tech stack

- Chi - Go HTTP router.
- BadgerDB - Fast key-value database for buffering data.
- Parquet-go - Library for reading and writing Parquet files in Go.