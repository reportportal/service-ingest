package handler

import "time"

type SaveLogRQ struct {
	UUID       string    `json:"uuid,omitempty" verify:"omitempty,uuid"`
	ItemUUID   string    `json:"itemUuid,omitempty" verify:"omitempty,uuid"`
	LaunchUUID string    `json:"launchUuid,omitempty" verify:"required,uuid"`
	Timestamp  time.Time `json:"time,omitempty" verify:"required"`
	Level      string    `json:"level,omitempty" verify:"omitempty,oneof=error warn info debug trace fatal unknown"`
	Message    string    `json:"message,omitempty"`
	File       LogFile   `json:"file,omitempty"`
}

type LogFile struct {
	Name string `json:"name,omitempty"`
}

type SaveLogRS struct {
	Responses []LogResponse `json:"responses"`
}

type LogResponse struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	StackTrace string `json:"stackTrace"`
}
