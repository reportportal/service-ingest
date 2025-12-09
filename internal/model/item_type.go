package model

type ItemType string

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
