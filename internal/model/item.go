package model

import "time"

type Item struct {
	UUID        string     `json:"uuid"`
	LaunchUUID  string     `json:"launch_uuid"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        ItemType   `json:"type"`
	Status      ItemStatus `json:"status"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Attributes  Attributes `json:"attributes"`
	Parameters  Parameters `json:"parameters"`
	CodeRef     string     `json:"code_ref"`
	TestCaseId  string     `json:"test_case_id"`
	ParentUUID  string     `json:"parent_uuid"`
	IsRetry     bool       `json:"retry"`
	RetryOf     string     `json:"retry_of"`
	Issue       Issue      `json:"issue"`
}

type Attribute struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsSystem bool   `json:"isSystem"`
}

type Attributes []Attribute

type Parameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Parameters []Parameter

type Issue struct {
	Type string `json:"type"`
}
