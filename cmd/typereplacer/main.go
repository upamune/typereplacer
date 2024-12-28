package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"

	"github.com/upamune/typereplacer/internal/config"
	"github.com/upamune/typereplacer/internal/rewriter"
)

var (
	version = "local"
	commit  = "none"
	date    = "unknown"
)

// CLI defines command-line interface with kong.
// - --config= path/to/config.yaml
// - PackageArg: the directory containing .go files to rewrite
type CLI struct {
	Config     string `kong:"required,help='Path to YAML config file'"`
	PackageArg string `kong:"arg,help='Directory (or package path) to rewrite'"`

	Version kong.VersionFlag `short:"v" help:"Show version and exit."`
}

func main() {
	if err := runMain(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// runMain is the core logic:
//  1. Parse CLI
//  2. Load YAML config
//  3. Rewrite all .go files in the specified package directory, changing struct fields' types
//     according to config.Structs
func runMain(args []string) error {
	var cli CLI
	parser, err := kong.New(
		&cli,
		kong.Name("typereplacer"),
		kong.Description("Rewrite specified struct fields' types based on a config."),
		kong.Vars{
			"version": fmt.Sprintf("%s (%s)", version, commit),
		},
	)
	if err != nil {
		return err
	}

	_, err = parser.Parse(args)
	if err != nil {
		return err
	}

	if cli.PackageArg == "" {
		return fmt.Errorf("no target directory specified")
	}

	// 1) Load config
	cfg, err := config.LoadConfig(cli.Config)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 2) Rewrite
	if err := rewriter.RewritePackage(cli.PackageArg, cfg); err != nil {
		return fmt.Errorf("rewrite error: %w", err)
	}

	return nil
}
