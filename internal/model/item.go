package model

import "time"

const (
	ItemTypeSuite        ItemType = "SUITE"
	ItemTypeStory        ItemType = "STORY"
	ItemTypeTest         ItemType = "TEST"
	ItemTypeScenario     ItemType = "SCENARIO"
	ItemTypeStep         ItemType = "STEP"
	ItemTypeBeforeClass  ItemType = "BEFORE_CLASS"
	ItemTypeBeforeGroups ItemType = "BEFORE_GROUPS"
	ItemTypeBeforeMethod ItemType = "BEFORE_METHOD"
	ItemTypeBeforeSuite  ItemType = "BEFORE_SUITE"
	ItemTypeBeforeTest   ItemType = "BEFORE_TEST"
	ItemTypeAfterClass   ItemType = "AFTER_CLASS"
	ItemTypeAfterGroups  ItemType = "AFTER_GROUPS"
	ItemTypeAfterMethod  ItemType = "AFTER_METHOD"
	ItemTypeAfterSuite   ItemType = "AFTER_SUITE"
	ItemTypeAfterTest    ItemType = "AFTER_TEST"
)

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

type Item struct {
	UUID        string     `json:"uuid"`
	LaunchUUID  string     `json:"launch_uuid"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        ItemType   `json:"type"`
	Status      ItemStatus `json:"status"`
	StartTime   time.Time  `json:"start_time"`
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

type ItemType string

type ItemStatus string

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
