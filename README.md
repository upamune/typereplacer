# 🔄 typereplacer

> ✨ A Go struct field type rewriter

## 🎯 What it does

`typereplacer` is a CLI tool that helps you refactor Go struct field types across your entire codebase with ease. Simply define your type changes in a YAML config, and let the tool do the heavy lifting! 

### ✨ Key Features

- 🛠️ **Simple YAML Configuration** - Define your type changes in a clear, structured way
- 📦 **Bulk Updates** - Change multiple struct fields across many files at once

## 🚀 Getting Started

### Installation

```bash
go install github.com/upamune/typereplacer@latest
```

### Usage

```bash
typereplacer --config=./myconfig.yaml ./path/to/pkg
```

## 📝 Configuration

Create a YAML config file that specifies your desired type changes:

```yaml
imports:
  - "fmt"
  - "strings"
structs:
  - name: "MyStruct"
    fields:
      - name: "Value"
        type: "int"
  - name: "YourStruct"
    fields:
      - name: "Text"
        type: "string"
```

### Config Structure

- `imports`: List of packages to analyze (used for type resolution)
- `structs`: List of struct definitions to modify
  - `name`: The struct name to target
  - `fields`: List of fields to update
    - `name`: Field name to modify
    - `type`: New type to apply

## 📋 Example

Let's say you have this initial code:

```go
type User struct {
    ID        int
    Balance   float64
    CreatedAt string
    UpdatedAt string
}
```

And you want to change the time fields to use `time.Time` and the balance field to use `decimal.Decimal`. Create a config file `typeconfig.yaml`:

```yaml
imports:
  - "time"
  - "github.com/shopspring/decimal"
structs:
  - name: "User"
    fields:
      - name: "Balance"
        type: "decimal.Decimal"
      - name: "CreatedAt"
        type: "time.Time"
      - name: "UpdatedAt"
        type: "time.Time"
```

Run the command:

```bash
typereplacer --config=typeconfig.yaml ./...
```

The tool will update your code to:

```go
type User struct {
    ID        int
    Balance   decimal.Decimal
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
