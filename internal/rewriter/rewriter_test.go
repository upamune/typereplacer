package rewriter_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/upamune/typereplacer/internal/config"
	"github.com/upamune/typereplacer/internal/rewriter"
)

func TestRewritePackage_MultiCases(t *testing.T) {
	tests := []struct {
		name       string
		sourceFile string
		goldenFile string
		cfg        *config.Config
	}{
		{
			name:       "Case1_SingleFieldStringToSlice",
			sourceFile: "struct1.go",
			goldenFile: "struct1.golden",
			cfg: &config.Config{
				Imports: []string{"fmt"},
				Structs: []config.Struct{
					{
						Name: "MyStruct",
						Fields: []config.Field{
							{Name: "Value", Type: "[]string"},
						},
					},
				},
			},
		},
		{
			name:       "Case2_MultipleStructButOnlyOneFieldChange",
			sourceFile: "struct2.go",
			goldenFile: "struct2.golden",
			cfg: &config.Config{
				Imports: []string{},
				Structs: []config.Struct{
					{
						Name: "Task",
						Fields: []config.Field{
							{Name: "Title", Type: "int"},
						},
					},
				},
			},
		},
		{
			name:       "Case3_ChangeBoolToString",
			sourceFile: "struct3.go",
			goldenFile: "struct3.golden",
			cfg: &config.Config{
				Imports: []string{"strings"},
				Structs: []config.Struct{
					{
						Name: "Complex",
						Fields: []config.Field{
							{Name: "Next", Type: "string"},
						},
					},
				},
			},
		},
		{
			name:       "Case4_AlreadyHasImport",
			sourceFile: "struct4.go",
			goldenFile: "struct4.golden",
			cfg: &config.Config{
				Imports: []string{"database/sql"},
				Structs: []config.Struct{
					{
						Name: "MyData",
						Fields: []config.Field{
							{Name: "Name", Type: "sql.Null[string]"},
						},
					},
				},
			},
		},
		{
			name:       "Case5_NoImportBlockAddSql",
			sourceFile: "struct5.go",
			goldenFile: "struct5.golden",
			cfg: &config.Config{
				Imports: []string{"database/sql"},
				Structs: []config.Struct{
					{
						Name: "Info",
						Fields: []config.Field{
							{Name: "Text", Type: "sql.Null[string]"},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			srcFilePath := filepath.Join("testdata", tc.sourceFile)
			dstFilePath := filepath.Join(tmpDir, tc.sourceFile)
			if err := copyFile(dstFilePath, srcFilePath); err != nil {
				t.Fatalf("failed to copy source file: %v", err)
			}

			if err := rewriter.RewritePackage(tmpDir, tc.cfg); err != nil {
				t.Fatalf("RewritePackage error: %v", err)
			}

			gotBytes, err := os.ReadFile(dstFilePath)
			if err != nil {
				t.Fatalf("failed to read rewritten file: %v", err)
			}

			wantBytes, err := os.ReadFile(filepath.Join("testdata", tc.goldenFile))
			if err != nil {
				t.Fatalf("failed to read golden file: %v", err)
			}

			wantStr := string(wantBytes)
			gotStr := string(gotBytes)

			if diff := cmp.Diff(wantStr, gotStr); diff != "" {
				t.Errorf("-want +got:\n%s", diff)
			}
		})
	}
}

// copyFile は単純なファイルコピー用ヘルパー関数
func copyFile(dst, src string) error {
	sf, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create dst: %w", err)
	}
	defer df.Close()

	if _, err := io.Copy(df, sf); err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	return nil
}
