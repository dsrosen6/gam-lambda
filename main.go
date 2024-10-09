package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/aws/aws-lambda-go/lambda"
)

const (
	gamBin     = "/opt/gam7/gam"
	rscDir     = "/tmp/resources"
	gamWorkDir = "GAMWork"
	gamCfgDir  = "GAMConfig"
)

type event struct {
	Org  string    `json:"org"`
	Cmds []command `json:"cmds"`
}

type command struct {
	Args []string `json:"args"`
}

type results struct {
	Results []cmdOutput `json:"results"`
}

type cmdOutput struct {
	Cmd     string `json:"cmd"`
	Success bool   `json:"success"`
	Out     string `json:"output"`
}

func main() {
	lambda.Start(runGamCommands)
}

func runGamCommands(ctx context.Context, event *event) (*results, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr,
		&slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	// If nil event, return error
	if event == nil {
		slog.Error("event is nil")
		return nil, fmt.Errorf("event is nil")
	}

	// Set up files
	if err := setUpFiles(); err != nil {
		slog.Error("could not set up files", "error", err)
		return nil, fmt.Errorf("setUpFiles: %w", err)
	}

	// Set org via runGam command - cannot run further commands without selecting org
	_, err := runGam("select", event.Org, "save")
	if err != nil {
		slog.Error("error selecting org", "error", err)
		return nil, fmt.Errorf("error selecting org: %w", err)
	}

	// Run all commands in event, do not stop on error
	var res results
	for _, cmd := range event.Cmds {
		out, err := runGam(cmd.Args...)
		success := true
		if err != nil {
			success = false
		}

		// Append output to results
		res.Results = append(res.Results, cmdOutput{
			Cmd:     commandToString(cmd.Args),
			Success: success,
			Out:     out,
		})
	}

	return &res, nil
}

// runGam runs a GAM command with the given arguments. GAM binary is provided by a Lambda layer
func runGam(args ...string) (string, error) {
	cmd := exec.Command(gamBin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("cmd.CombinedOutput: %w", err)
	}

	return string(out), nil
}

// setUpFiles handles creating necessary directories and copying files from Lambda layers into said directories
func setUpFiles() error {
	// Create Resources and GAMWork directories
	gw := fmt.Sprintf("%s/%s", rscDir, gamWorkDir)
	if err := makeDir(gw); err != nil {
		return fmt.Errorf("makeDir: %w", err)
	}

	// Copy GAM Config layer to resources directory
	if err := copyDir("/opt/GAMConfig", rscDir); err != nil {
		return fmt.Errorf("copyDir: %w", err)
	}

	return nil
}

// copyDir copies a directory from src to dst
func copyDir(src, dst string) error {
	cmd := exec.Command("cp", "-r", src, dst)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cmd.Run: %w", err)
	}

	return nil
}

// makeDir creates a directory at dirPath (and any necessary parent directories, same functionality as mkdir -p)
func makeDir(dirPath string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	return nil
}

// commandToString converts a command slice to a string - used for logging
func commandToString(cmd []string) string {
	var str string
	for _, c := range cmd {
		str += c + " "
	}

	return str
}
