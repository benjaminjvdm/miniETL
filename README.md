# miniETL 🚀

A simple and extensible ETL (Extract, Transform, Load) application written in Go.

## Features ✨

*   **Reads data from CSV files 📁**
*   **Transforms data using:**
    *   Field renaming ✏️
    *   Data type conversion 🗂️
    *   Filtering 🔎
    *   Field calculation ➕
*   **Writes transformed data to CSV files 📤**
*   **Configuration via YAML files ⚙️**
*   **Help text and usage examples ℹ️**
*   **Configuration validation ✅**

## Usage 💻

1.  **Build the application:**

    ```bash
    go build -o miniETL main.go config.go csv_source.go source.go rename_transform.go convert_transform.go filter_transform.go calculate_transform.go transform.go csv_destination.go destination.go
    ```
2.  **Run the application:**

    ```bash
    ./miniETL -config config.yaml
    ```

## Configuration ⚙️

The ETL process is configured using a YAML file. See `config.yaml` for an example.

## Transformations 🔀

The following transformations are supported:

*   **rename:** Renames a field.
*   **convert:** Converts a field to a different data type (e.g., int, float, string).
*   **filter:** Filters records based on a condition.
*   **calculate:** Calculates a new field based on an expression.

## Contributing 🤝

Contributions are welcome! Please submit a pull request.

## License 📝

MIT License