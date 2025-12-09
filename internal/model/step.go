package model

import "time"

type NestedStep struct {
	UUID        string     `json:"uuid"`
	LaunchUUID  string     `json:"launch_uuid"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      ItemStatus `json:"status"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Attributes  Attributes `json:"attributes"`
	Parameters  Parameters `json:"parameters"`
	CodeRef     string     `json:"code_ref"`
	TestCaseId  string     `json:"test_case_id"`
	ParentUUID  string     `json:"parent_uuid"`
}
