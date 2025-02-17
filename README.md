# Go ETL Tool

## Description

This is a command-line ETL tool written in Go. It allows you to extract data from CSV, JSON, and TXT files, transform the data using user-defined Go functions or a built-in expression language, and load the transformed data into CSV, JSON, TXT, or SQLite databases.

## Configuration

The tool is configured using a YAML or TOML file. The configuration file specifies the input file, output file, and transformations to apply.

## Usage

```
go run main.go -config config.yaml
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
  - field1: "expression1"
  - field2: "expression2"