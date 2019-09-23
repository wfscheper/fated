// Copyright © 2017 Walter Scheper <walter.scheper@gmal.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build mage

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

const (
	targetName  = "fated"
	packageName = "github.com/wfscheper/" + targetName
)

var (
	// Default is the default mage target
	Default = Build

	// allow override via GOEXE=...
	goexe   = "go"
	ldflags = "-X $PACKAGE/cmd.Version=$VERSION -X $PACKAGE/cmd.BuildDate='$BUILD_DATE' -X $PACKAGE/cmd.Commit=$COMMIT"

	// commands
	gofmt        = sh.RunCmd(goexe, "fmt")
	gotest       = sh.RunCmd(goexe, "test", "-timeout", "15s")
	goveralls    = filepath.Join("tools", "goveralls")
	golangcilint = filepath.Join("tools", "golangci-lint")
)

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// Force use of go modules
	os.Setenv("GO111MODULE", "on")
}

// Build builds the fated binary
func Build(ctx context.Context) error {
	mg.CtxDeps(ctx, Fmt, Lint)
	fmt.Println("building " + targetName + "…")
	return sh.RunWith(buildEnvs(), goexe, "build", "-v", "-tags", buildTags(), "-ldflags", ldflags, "-o",
		filepath.Join("bin", targetName), ".")
}

// Fmt runs go fmt
func Fmt(ctx context.Context) error {
	fmt.Println("running go fmt…")
	return gofmt("./...")
}

// Lint runs golangci-lint
func Lint(ctx context.Context) error {
	mg.CtxDeps(ctx, getGolangciLint)
	fmt.Println("running golagnci-lint…")
	return sh.Run(golangcilint, "run")
}

func getGolangciLint(ctx context.Context) error {
	if rebuild, err := target.Path(golangcilint); err != nil {
		return err
	} else if rebuild {
		fmt.Println("getting golangci-lint…")
		return sh.RunWith(toolsEnv(), goexe, "install", "github.com/golangci/golangci-lint/cmd/golangci-lint")
	}
	return nil
}

// Test runs the test suite
func Test(ctx context.Context) error {
	mg.CtxDeps(ctx, Fmt, Lint)
	return runTest()
}

// TestRace runs the test suite with race detection
func TestRace(ctx context.Context) error {
	mg.CtxDeps(ctx, Fmt, Lint)
	return runTest("-race")
}

// TestShort runs only tests marked as short
func TestShort(ctx context.Context) error {
	mg.CtxDeps(ctx, Fmt, Lint)
	return runTest("-short")
}

// Benchmark runs the benchmark suite
func Benchmark(ctx context.Context) error {
	mg.CtxDeps(ctx, Fmt, Lint)
	return runTest("-run=__absolutelynothing__", "-bench")
}

func runTest(testType ...string) error {
	var space string
	if len(testType) > 1 {
		space = " "
	}
	fmt.Printf("running go test%s%s…\n", space, strings.Join(testType, " "))
	testType = append(testType, "./...")
	return gotest(testType...)
}

// Coverage generates coverage reports
func Coverage(ctx context.Context) error {
	mg.CtxDeps(ctx, Fmt, Lint)
	sh.Run("mkdir", "-p", "coverage")
	mode := os.Getenv("COVERAGE_MODE")
	if mode == "" {
		mode = "atomic"
	}
	if err := runTest("-cover", "-covermode", mode, "-coverprofile=coverage/cover.out"); err != nil {
		return err
	}
	if err := sh.Run(goexe, "tool", "cover", "-html=coverage/cover.out", "-o", "coverage/index.html"); err != nil {
		return err
	}
	return nil
}

// Coveralls uploads coverage report
func Coveralls(ctx context.Context) error {
	// only do something if within travis
	if os.Getenv("TRAVIS_HOME") == "" {
		return nil
	}
	mg.CtxDeps(ctx, getGoveralls)
	fmt.Println("running goveralls…")
	return sh.Run(goveralls, "-coverprofile=coverage/cover.out", "-service=travis-ci")
}

func getGoveralls(ctx context.Context) error {
	if rebuild, err := target.Path(goveralls); err != nil {
		return err
	} else if rebuild {
		fmt.Println("getting goveralls…")
		return sh.RunWith(toolsEnv(), goexe, "install", "github.com/mattn/goveralls")
	}
	return nil
}

func Clean() error {
	return sh.Run("rm", "-rf", "bin/", "dist/", "tools/", "coverage/")
}

func buildEnvs() map[string]string {
	commit, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return map[string]string{
		"BUILD_DATE": time.Now().Format(time.RFC3339),
		"COMMIT":     commit,
		"PACKAGE":    packageName,
	}
}

func buildTags() string {
	if tags := os.Getenv("BUILD_TAGS"); tags != "" {
		return tags
	}
	return "none"
}

func toolsEnv() map[string]string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	tools := filepath.Join(cwd, "tools")
	path := strings.Join([]string{
		tools,
		os.Getenv("PATH"),
	}, string(os.PathListSeparator))
	return map[string]string{
		"GOBIN": tools,
		"PATH":  path,
	}
}
