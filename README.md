# typereplacer (Struct Field Type Rewriter)

This CLI tool reads a YAML `Config` that specifies:

- **Imports**: Packages to analyze (but we only use this for analysis or loading)
- **Structs**: A list of struct definitions, each with multiple fields that specify 
  which field name to rewrite and what the **new type** should be.

Then, we parse the target Go package/directory provided on the CLI and replace 
the specified struct field **types** in all `.go` files.

## Usage

```bash
$ typereplacer --config=./myconfig.yaml ./path/to/pkg
```

Where myconfig.yaml might look like:

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
