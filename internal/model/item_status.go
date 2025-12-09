package model

type ItemStatus string

const (
	ItemStatusPassed      ItemStatus = "PASSED"
	ItemStatusFailed      ItemStatus = "FAILED"
	ItemStatusSkipped     ItemStatus = "SKIPPED"
	ItemStatusStopped     ItemStatus = "STOPPED"
	ItemStatusInterrupted ItemStatus = "INTERRUPTED"
	ItemStatusCancelled   ItemStatus = "CANCELLED"
	ItemStatusInfo        ItemStatus = "INFO"
	ItemStatusWarn        ItemStatus = "WARN"
)
