package rewriter

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/upamune/typereplacer/internal/config"
	"golang.org/x/tools/imports"
)

// RewritePackage walks through the specified directory,
// parsing each .go file, and rewrites the type of the fields
// listed in config.Structs.
func RewritePackage(pkgDir string, cfg *config.Config) error {
	var goFiles []string
	err := filepath.Walk(pkgDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory %q: %w", pkgDir, err)
	}

	for _, gf := range goFiles {
		if err := rewriteFile(gf, cfg); err != nil {
			return fmt.Errorf("rewriteFile %q error: %w", gf, err)
		}
	}
	return nil
}

// rewriteFile reads a single .go file, adds missing imports from cfg.Imports to
// the first import block found (or creates a new one below the package declaration)
// and rewrites struct fields' types as specified in cfg.Structs.
func rewriteFile(filePath string, cfg *config.Config) error {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		return err
	}

	existing := make(map[string]bool)
	for _, imp := range f.Imports {
		p := strings.Trim(imp.Path.Value, `"`)
		existing[p] = true
	}

	var importsToAdd []string
	for _, i := range cfg.Imports {
		if !existing[i] {
			importsToAdd = append(importsToAdd, i)
		}
	}

	// Add imports to the first import block found, or create a new one if none exist
	if len(importsToAdd) > 0 {
		var importBlock *ast.GenDecl
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.IMPORT {
				continue
			}
			importBlock = gd
			break
		}
		if importBlock != nil {
			for _, i := range importsToAdd {
				importBlock.Specs = append(importBlock.Specs, &ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("%q", i),
					},
				})
			}
		} else {
			var specs []ast.Spec
			for _, i := range importsToAdd {
				specs = append(specs, &ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("%q", i),
					},
				})
			}
			newDecl := &ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: specs,
			}
			f.Decls = append([]ast.Decl{newDecl}, f.Decls...)
		}
	}

	// Rewrite struct fields based on cfg.Structs
	ast.Inspect(f, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}
		var match *config.Struct
		for i := range cfg.Structs {
			if cfg.Structs[i].Name == ts.Name.Name {
				match = &cfg.Structs[i]
				break
			}
		}
		if match == nil || st.Fields == nil {
			return true
		}
		for _, field := range st.Fields.List {
			if len(field.Names) == 0 {
				continue
			}
			fName := field.Names[0].Name
			for _, fc := range match.Fields {
				if fc.Name == fName {
					newExpr, parseErr := parser.ParseExpr(fc.Type)
					if parseErr == nil {
						field.Type = newExpr
					}
				}
			}
		}
		return true
	})

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, f); err != nil {
		return err
	}

	res, err := imports.Process(filePath, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("formatting failed: %w", err)
	}

	return os.WriteFile(filePath, res, 0644)
}
