package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	configFile := flag.String("config", "config.yaml", "Path to the configuration file")
	help := flag.Bool("help", false, "Show help message and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	fmt.Println("Starting miniETL...")

	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Configuration: %+v\n", config)

	source, err := NewSource(config.Source)
	if err != nil {
		log.Fatalf("Failed to create source: %v", err)
	}

	recordChan, errChan := source.Read()

	for _, transformConfig := range config.Transform {
		transform, err := NewTransform(transformConfig)
		if err != nil {
			log.Fatalf("Failed to create transform: %v", err)
		}

		recordChan, errChan = transform.Apply(recordChan)
		go func() {
			for err := range errChan {
				log.Printf("Error: %v", err)
			}
		}()
	}

	destination, err := NewDestination(config.Destination)
	if err != nil {
		log.Fatalf("Failed to create destination: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	destErrChan := destination.Write(recordChan)

	go func() {
		defer wg.Done()
		for err := range destErrChan {
			log.Printf("Error: %v", err)
		}
	}()

	wg.Wait()

	log.Println("miniETL finished.")
}
