package model

import "time"

type Log struct {
	UUID       string
	ItemUUID   string
	LaunchUUID string
	Timestamp  time.Time
	Level      LogLevel
	Message    string
	File       LogFile
}

type LogFile struct {
	Name string
}
