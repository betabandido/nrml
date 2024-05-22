package logging

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type labelHook struct {
	labels map[string]string
}

func newLabelHook(labels map[string]string) *labelHook {
	return &labelHook{
		labels: labels,
	}
}

func (hook *labelHook) Fire(entry *log.Entry) error {
	if entry.Context != nil {
		for key, value := range hook.labels {
			labelKey := fmt.Sprintf("labels.%s", key)
			if _, exists := entry.Data[labelKey]; !exists {
				entry.Data[labelKey] = value
			}
		}
	}

	return nil
}

func (hook *labelHook) Levels() []log.Level {
	return log.AllLevels
}
