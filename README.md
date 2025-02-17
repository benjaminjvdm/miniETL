# Go ETL Tool

## Description

This is a command-line ETL tool written in Go. It allows you to extract data from CSV, JSON, and TXT files, transform the data using user-defined Go functions or a built-in expression language, and load the transformed data into CSV, JSON, TXT, or SQLite databases.

## Installation

```bash
go get github.com/yourusername/miniETL
cd $GOPATH/github.com/yourusername/miniETL
go build
```

## Usage

```bash
./miniETL run -config config.yaml
```

### Commands

*   `run`: Executes the ETL pipeline.
    *   `-config`: Specifies the configuration file to use.
    *   `-log-level`: Specifies the log level (debug, info, warning, error).
*   `validate`: Validates the ETL pipeline configuration file.
*   `preview`: Previews data at each stage of the pipeline (extract, transform, load).
*   `manage`: Manages and schedules ETL pipeline executions (not yet implemented).

## Configuration

The tool is configured using a YAML file. The configuration file specifies the input file, output file, and transformations to apply.

### Example Configuration

```yaml
input:
  path: input.csv
  type: csv
output:
  path: output.db
  type: sqlite
transformations:
  - field: name
    action: uppercase
  - field: date
    action: date_format
    params:
      format: "YYYY-MM-DD"

```

## Extension

To add support for new data sources or destinations, you need to implement the `Extractor` or `Loader` interface, respectively. You can also add custom transformation functions by defining them in Go code and registering them with the `Transform` function.