// +build mage

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

const (
	targetName  = "fated"
	packageName = "github.com/wfscheper/" + targetName
	version     = "v0.1.0"
)

var (
	// Default target to run when none is specified
	// If not set, running mage will list available targets
	Default = Build
	// allow override via GOEXE=...
	goexe   = "go"
	ldflags = "-X $PACKAGE/cmd.Version=$VERSION -X $PACKAGE/cmd.BuildDate=$BUILD_DATE -X $PACKAGE/cmd.Commit=$COMMIT"
)

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// Force go modules
	os.Setenv("GO111MODULE", "on")
}

func getEnv() map[string]string {
	commit, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return map[string]string{
		"BUILD_DATE": time.Now().Format(time.RFC3339),
		"COMMIT":     commit,
		"PACKAGE":    packageName,
		"VERSION":    version,
	}
}

func isGoLatest() bool {
	return strings.Contains(runtime.Version(), "1.11")
}

// Build fated binary
func Build(ctx context.Context) error {
	mg.CtxDeps(ctx, Fmt, Lint, Vet)
	if rebuild, err := target.Dir(filepath.Join("bin", targetName), "."); rebuild {
		fmt.Println("Building...")
		return sh.RunWith(getEnv(), goexe, "build", "-ldflags", ldflags, packageName)
	} else if err != nil {
		return err
	}
	return nil
}

// Run gofmt
func Fmt(ctx context.Context) error {
	if !isGoLatest() {
		return nil
	}
	fmt.Println("Formatting...")
	if err := sh.Run(goexe, "fmt", "./..."); err != nil {
		return fmt.Errorf("error running go fmt: %v", err)
	}
	return nil
}

// Run go lint
func Lint(ctx context.Context) error {
	fmt.Println("Linting...")
	if err := sh.Run("golint", "./..."); err != nil {
		return fmt.Errorf("error running golint: %v", err)
	}
	return nil
}

// Run vet
func Vet(ctx context.Context) error {
	fmt.Println("Vetting...")
	if err := sh.Run(goexe, "vet", "./..."); err != nil {
		return fmt.Errorf("error running go vet: %v", err)
	}
	return nil
}
