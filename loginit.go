// +build log

package main

import (
	"log"
)

func init() {
	var zerr error

	logger, zerr = zap.NewDevelopment()

	if zerr != nil {
		log.Fatalf("can't initialize zap logger: %v", zerr)
	}

	defer logger.Sync()
}
