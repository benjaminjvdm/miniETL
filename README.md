# Go-ETL: A Simple & Powerful ETL Tool

![Go Version](https://img.shields.io/badge/go-1.18+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)

Go-ETL is a flexible, command-line tool written in Go for performing simple Extract, Transform, and Load (ETL) tasks. It's designed to be easy to configure and fast to run, helping you move and reshape data with minimal effort.

Whether you need to migrate data between formats, clean up a dataset, or create new columns from existing data, Go-ETL provides a straightforward solution.

## Key Features

- **Multiple Data Sources:** Extract data from CSV and TXT (JSONL) files.
- **Flexible Output Options:** Load data into JSON, CSV, TXT (JSONL), or a SQLite database.
- **Powerful Transformations:**
    - Create new columns using data from existing columns.
    - Use a simple expression language to combine, concatenate, and modify data.
    - Built-in functions like `upper()` for case conversion and `len()` for string length.
- **Easy Configuration:** Define your entire ETL pipeline in a single, human-readable YAML or TOML file.
- **Fast & Efficient:** Built with Go for high performance and low resource usage.
- **Cross-Platform:** Runs on any system where you can compile a Go program (Windows, macOS, Linux).

## Getting Started

### Prerequisites

- You need to have [Go](https://golang.org/doc/install) (version 1.18 or newer) installed on your system.

### Installation & Usage

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/go-etl.git
    cd go-etl
    ```

2.  **Run the tool:**
    The easiest way to run the tool is directly using `go run`. You can point it to a configuration file using the `--config` flag.
    ```bash
    go run . --config config.yaml
    ```

3.  **Build a binary (Optional):**
    For more permanent use, you can build an executable.
    ```bash
    go build -o etl-tool
    ./etl-tool --config config.yaml
    ```

## Configuration

Go-ETL is controlled by a configuration file, typically named `config.yaml` or `config.toml`. Here is a breakdown of the structure:

```yaml
# --- Input Configuration ---
input:
  path: path/to/your/input.csv # Path to the source file
  type: csv                      # Type of the source file

# --- Output Configuration ---
output:
  path: path/to/your/output.json # Path for the destination file
  type: json                     # Type of the destination file

# --- Transformation Rules ---
transformations:
  # Each item is a new column to be created.
  # The key is the new column name, and the value is the expression.
  - new_column_name: expression
```

### Supported Formats

| Format | Input | Output | Notes |
| :--- | :---: | :---: | :--- |
| `csv` | ✓ | ✓ | |
| `json` | | ✓ | |
| `txt` | ✓ | ✓ | Assumes JSONL format (one JSON object per line). |
| `sqlite`| | ✓ | The table name will be the output filename. |


## Transformation Expressions

The `transformations` section allows you to create new columns based on existing data. The expression language is simple and powerful.

**Syntax:** `new_column_name: term + term + ...`

An expression is composed of one or more "terms" joined by the `+` operator. A term can be:

- **A column name:** e.g., `header1`
- **A string literal:** e.g., `" - "` (must be in double quotes)
- **A function call:** e.g., `upper(header1)`

**Behavior of `+`:**
- If the first term is a number (e.g., from `len()`), `+` will perform **addition**.
- Otherwise, `+` will perform **string concatenation**.

### Available Functions

- `upper(column_name)`: Converts the text in `column_name` to uppercase.
- `len(column_name)`: Returns the length of the string in `column_name`.

### Example

Given an input row: `{"first_name": "John", "last_name": "Doe"}`

The following transformations:
```yaml
transformations:
  - full_name: first_name + " " + last_name
  - upper_last: upper(last_name)
  - name_length: len(first_name) + len(last_name)
```

Will produce the new columns:
- `full_name`: "John Doe"
- `upper_last`: "DOE"
- `name_length`: 8 (4 + 4)

## Contributing

Contributions are welcome! If you find a bug, have a feature request, or want to contribute code, please feel free to:

1.  [Open an issue](https://github.com/your-username/go-etl/issues) to discuss the change.
2.  Fork the repository and create a new branch.
3.  Make your changes and ensure all tests pass (`go test ./...`).
4.  Submit a pull request.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.