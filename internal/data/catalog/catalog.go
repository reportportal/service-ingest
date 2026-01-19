package catalog

import (
	"fmt"
	"path/filepath"
)

// BuildPath constructs the partitioned path for Parquet files
// Format: project_id=<uuid>/launch_id=<uuid>/events/entity=<type>/event_date=<date>/batch_id=<ts>-<seq>/
func BuildPath(projectID, launchID, entityType, eventDate, batchID string) string {
	return filepath.Join(
		fmt.Sprintf("project_id=%s", projectID),
		fmt.Sprintf("launch_uuid=%s", launchID),
		"events",
		fmt.Sprintf("entity=%s", entityType),
		fmt.Sprintf("event_date=%s", eventDate),
		fmt.Sprintf("batch_id=%s", batchID),
	)
}
