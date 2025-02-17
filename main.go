package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "miniETL",
	Short: "A simple ETL tool",
	Long:  `miniETL is a command line tool to extract, transform, and load data.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("miniETL is running...")
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the ETL pipeline",
	Long:  `Run the ETL pipeline using the specified configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set log level
		switch logLevel {
		case "debug":
			InfoLogger.SetOutput(os.Stdout)
			ErrorLogger.SetOutput(os.Stderr)
		case "info":
			InfoLogger.SetOutput(os.Stdout)
			ErrorLogger.SetOutput(os.Stderr)
		case "warning":
			InfoLogger.SetOutput(io.Discard)
			ErrorLogger.SetOutput(os.Stderr)
		case "error":
			InfoLogger.SetOutput(io.Discard)
			ErrorLogger.SetOutput(os.Stderr)
		default:
			fmt.Println("Invalid log level:", logLevel)
			os.Exit(1)
		}

		config, err := LoadConfig("./config.yaml")
		if err != nil {
			log.Fatal(err)
		}

		data, err := Extract(config)
		if err != nil {
			log.Fatal(err)
		}

		transformedData, err := Transform(data, config.Transformations)
		if err != nil {
			log.Fatal(err)
		}

		err = Load(config, transformedData)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("ETL process completed successfully!")
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the ETL pipeline configuration file",
	Long:  `Validate the ETL pipeline configuration file to ensure it is valid.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := LoadConfig("./config.yaml")
		if err != nil {
			fmt.Println("Configuration file is invalid:", err)
			os.Exit(1)
		}

		fmt.Println("Configuration file is valid.")
	},
}

var logLevel string

func init() {
	runCmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, warning, error)")
}

func main() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(validateCmd)

	var previewCmd = &cobra.Command{
		Use:   "preview",
		Short: "Preview data at each stage of the pipeline",
		Long:  `Preview data at each stage of the pipeline (extract, transform, load).`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := LoadConfig("./config.yaml")
			if err != nil {
				log.Fatal(err)
			}

			data, err := Extract(config)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Extracted Data:")
			fmt.Println(data)

			transformedData, err := Transform(data, config.Transformations)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Transformed Data:")
			fmt.Println(transformedData)

			// Load
			err = Load(config, transformedData)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Load process completed.")
		},
	}

	rootCmd.AddCommand(previewCmd)

	var manageCmd = &cobra.Command{
		Use:   "manage",
		Short: "Manage and schedule ETL pipeline executions",
		Long:  `Manage and schedule ETL pipeline executions. This feature is not yet implemented.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ETL pipeline management and scheduling is not yet implemented.")
		},
	}

	rootCmd.AddCommand(manageCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
