package model

type LaunchStatus string

const (
	LaunchStatusInProgress  LaunchStatus = "IN_PROGRESS"
	LaunchStatusPassed      LaunchStatus = "PASSED"
	LaunchStatusFailed      LaunchStatus = "FAILED"
	LaunchStatusStopped     LaunchStatus = "STOPPED"
	LaunchStatusSkipped     LaunchStatus = "SKIPPED"
	LaunchStatusInterrupted LaunchStatus = "INTERRUPTED"
	LaunchStatusCancelled   LaunchStatus = "CANCELLED"
	LaunchStatusInfo        LaunchStatus = "INFO"
	LaunchStatusWarn        LaunchStatus = "WARN"
)
