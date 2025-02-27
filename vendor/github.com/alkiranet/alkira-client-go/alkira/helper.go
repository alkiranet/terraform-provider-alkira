// Copyright (C) 2022-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"fmt"
	"log"
	"os"
)

// SegmentNameToZone is a type used to explain a rather copmlex and
// nested API mapping of segments (using name as a key) to another map
// of zones to groups. The hope is that this naming is clearer and
// less error prone to use than a map of interface{}
type SegmentNameToZone map[string]OuterZoneToGroups

// ZoneToGroups is a type used to explain a mapping of zone names to
// an associated array of group names
type ZoneToGroups map[string][]string

// OuterZoneToGroups is an object that includes a segment ID as an
// identifier and a map of zone names to Groups
type OuterZoneToGroups struct {
	SegmentId     int          `json:"segmentId"`
	ZonesToGroups ZoneToGroups `json:"zonesToGroups"`
}

// logf a simple log wrapper to log based on ENV var
func logf(level string, message string, v ...interface{}) {
	logLevel := os.Getenv("TF_LOG")

	if logLevel == level {
		format := fmt.Sprintf("[%s] %s", level, message)
		log.Printf(format, v...)
	}
}
