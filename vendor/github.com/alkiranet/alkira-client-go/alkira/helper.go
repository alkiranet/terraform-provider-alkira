// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"fmt"
	"log"
	"os"
)

// logf a simple log wrapper to log based on ENV var
func logf(level string, message string, v ...interface{}) {
	logLevel := os.Getenv("TF_LOG")

	if logLevel == level {
		format := fmt.Sprintf("[%s] %s", level, message)
		log.Printf(format, v...)
	}
}
