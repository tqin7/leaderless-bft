package util

import (
	log "github.com/sirupsen/logrus"
	"errors"
)

func StringInArray(target string, arr []string) bool {
	for _, element := range arr {
		if element == target {
			return true
		}
	}
	return false
}

func LogAndGetError(msg string) error {
	log.Error(msg)
	return errors.New(msg)
}