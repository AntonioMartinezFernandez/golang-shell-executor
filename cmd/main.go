package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	ctx := context.Background()

	shellCmdOutput, err := executeShellCmd(
		ctx,
		750*time.Millisecond,
		"curl",
		"-X", "GET", "https://pokeapi.co/api/v2/pokemon/",
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	parsedCmdOutput, err := parseShellCmdOutput(shellCmdOutput)
	if err != nil {
		fmt.Println(err)
		return
	}

	printPokemons(parsedCmdOutput)
}

func printPokemons(parsedCmdOutput *ShellCommandOutput) {
	fmt.Println("POKEMONS:")
	for i, pokemon := range parsedCmdOutput.Results {
		fmt.Printf(
			"%d. %s\n",
			i+1,
			cases.Title(language.English, cases.Compact).String(pokemon.Name),
		)
	}
}

func executeShellCmd(
	ctx context.Context,
	timeout time.Duration,
	binPath string,
	args ...string,
) ([]byte, error) {
	// Set timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(ctx, binPath, args...)

	// Run command
	cmdOutput, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}

	// Return command output
	return cmdOutput, nil
}

func parseShellCmdOutput(output []byte) (*ShellCommandOutput, error) {
	var shellCmdOutput ShellCommandOutput
	if err := json.Unmarshal(output, &shellCmdOutput); err != nil {
		return nil, fmt.Errorf("failed to parse output: %w", err)
	}
	return &shellCmdOutput, nil
}

type ShellCommandOutput struct {
	Count    int64       `json:"count"`
	Next     string      `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []Result    `json:"results"`
}

type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
