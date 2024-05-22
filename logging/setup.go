package logging

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type Options struct {
	Labels map[string]string
}

func SetupWithOptions(options Options) error {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)

	if len(options.Labels) > 0 {
		log.AddHook(newLabelHook(options.Labels))
	}

	return nil
}
