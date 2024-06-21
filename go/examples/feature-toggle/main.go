package main

import (
	"context"
	"log"
	"time"

	dasher "github.com/flohansen/dasher/go"
)

var (
	Debug = dasher.NewFeature(dasher.FeatureParams{
		Name:        "DEBUG_ENABLED",
		Description: "If the logging should be more verbose",
	})

	GoTemplates = dasher.NewFeature(dasher.FeatureParams{
		Name:        "USE_GO_TEMPLATES",
		Description: "If the rendering module should use Go Templates instead of Handlebars",
	})
)

func init() {
	dasher.MustRegister(Debug)
	dasher.MustRegister(GoTemplates)
}

func main() {
	go dasher.Connect(context.Background(), ":50051")

	for {
		if Debug.Enabled {
			log.Println("wait for one second...")
		}

		if GoTemplates.Enabled {
			log.Println("rendering with Go Templates")
		} else {
			log.Println("rendering with Handlebars")
		}

		time.Sleep(1 * time.Second)
	}
}
