# Go ETL Tool

## Description

This is a command-line ETL tool written in Go. It allows you to extract data from CSV and TXT files, transform the data using a simple expression language, and load the transformed data into JSON, CSV, TXT, or SQLite databases.

## Configuration

The tool is configured using a YAML or TOML file. The configuration file specifies the input file, output file, and transformations to apply.

## Usage

```
go run . --config config.yaml
```

You can also build the application and run the binary:

```
go build -o etl-tool
./etl-tool --config config.yaml
```

## Example Configuration

```yaml
input:
  path: input.csv
  type: csv
output:
  path: output.json
  type: json
transformations:
  - combined: header1 + " - " + header2
  - header1_upper: upper(header1)
  - total_chars: len(header1) + len(header2)