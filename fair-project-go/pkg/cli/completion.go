package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/posener/complete"
)

// CompletionCommand generates shell completion scripts
func CompletionCommand(shell, outFile string) error {
	if shell == "" {
		return fmt.Errorf("shell type is required")
	}

	// Normalize shell name
	shell = strings.ToLower(shell)

	// Get the binary name
	binaryName := filepath.Base(os.Args[0])

	// Create the completion command
	var script string
	switch shell {
	case "bash":
		script = generateBashCompletion(binaryName)
	case "zsh":
		script = generateZshCompletion(binaryName)
	case "fish":
		script = generateFishCompletion(binaryName)
	default:
		return fmt.Errorf("unsupported shell type: %s (supported: bash, zsh, fish)", shell)
	}

	// Write to file or stdout
	if outFile != "" {
		if err := os.WriteFile(outFile, []byte(script), 0644); err != nil {
			return fmt.Errorf("failed to write completion script to %s: %w", outFile, err)
		}
		fmt.Printf("Completion script written to %s\n", outFile)
		fmt.Printf("To install, run:\n")
		switch shell {
		case "bash":
			fmt.Printf("  echo \"source %s\" >> ~/.bashrc\n", outFile)
		case "zsh":
			fmt.Printf("  echo \"source %s\" >> ~/.zshrc\n", outFile)
		case "fish":
			fmt.Printf("  echo \"source %s\" >> ~/.config/fish/config.fish\n", outFile)
		}
	} else {
		fmt.Println(script)
		fmt.Printf("\nTo install, add the above to your shell configuration file or run:\n")
		switch shell {
		case "bash":
			fmt.Printf("  %s completion bash > ~/.%s-completion.bash && echo \"source ~/.%s-completion.bash\" >> ~/.bashrc\n",
				binaryName, binaryName, binaryName)
		case "zsh":
			fmt.Printf("  %s completion zsh > ~/.%s-completion.zsh && echo \"source ~/.%s-completion.zsh\" >> ~/.zshrc\n",
				binaryName, binaryName, binaryName)
		case "fish":
			fmt.Printf("  %s completion fish > ~/.config/fish/%s-completion.fish\n", binaryName, binaryName)
		}
	}

	return nil
}

// generateBashCompletion generates a bash completion script
func generateBashCompletion(binaryName string) string {
	// Since BashComplete is not available in the current version of the package,
	// we'll use a simple template for bash completion
	return fmt.Sprintf(`#!/bin/bash

_%s_completions() {
  COMPREPLY=()
  local word="${COMP_WORDS[COMP_CWORD]}"
  local completions="$(COMP_LINE="${COMP_LINE}" COMP_POINT="${COMP_POINT}" %s __complete)"
  COMPREPLY=( $(compgen -W "$completions" -- "$word") )
}

complete -F _%s_completions %s
`, binaryName, binaryName, binaryName, binaryName)
}

// generateZshCompletion generates a zsh completion script
func generateZshCompletion(binaryName string) string {
	// Since ZshComplete is not available in the current version of the package,
	// we'll use a simple template for zsh completion
	return fmt.Sprintf(`#compdef %s

_%s() {
  local -a completions
  completions=("${(@f)$(COMP_LINE="${words[*]}" COMP_POINT=$#words %s __complete)}")
  _describe 'completions' completions
}

compdef _%s %s
`, binaryName, binaryName, binaryName, binaryName, binaryName)
}

// generateFishCompletion generates a fish completion script
func generateFishCompletion(binaryName string) string {
	// Since FishComplete is not available in the current version of the package,
	// we'll use a simple template for fish completion
	return fmt.Sprintf(`function __fish_%s_complete
  set -l cl (commandline --tokenize --current-process)
  set -l tokens (commandline --tokenize --cut-at-cursor --current-process)
  %s __complete $tokens | tr '\n' ' '
end

complete -f -c %s -a '(__fish_%s_complete)'
`, binaryName, binaryName, binaryName, binaryName)
}

// CreateCompletionSpec creates a completion spec for the CLI
func CreateCompletionSpec(binaryName string) *complete.Command {
	// Create a completion command
	cmd := &complete.Command{
		Sub: complete.Commands{
			"create-class": {
				Flags: complete.Flags{
					"-name": complete.PredictAnything,
				},
			},
			"list-classes": {},
			"import-projects": {
				Flags: complete.Flags{
					"-class": predictClasses(),
					"-file":  complete.PredictFiles("*"),
				},
			},
			"import-students": {
				Flags: complete.Flags{
					"-class": predictClasses(),
					"-file":  complete.PredictFiles("*"),
				},
			},
			"list-projects": {
				Flags: complete.Flags{
					"-class": predictClasses(),
				},
			},
			"list-students": {
				Flags: complete.Flags{
					"-class": predictClasses(),
				},
			},
			"run-assignment": {
				Flags: complete.Flags{
					"-class": predictClasses(),
				},
			},
			"list-assignments": {
				Flags: complete.Flags{
					"-class": predictClasses(),
				},
			},
			"show-assignment": {
				Flags: complete.Flags{
					"-class": predictClasses(),
					"-id":    predictAssignmentIDs(),
				},
			},
			"export-assignment": {
				Flags: complete.Flags{
					"-class":  predictClasses(),
					"-id":     predictAssignmentIDs(),
					"-output": complete.PredictFiles("*"),
				},
			},
			"create-superadmin": {
				Flags: complete.Flags{
					"-username": complete.PredictAnything,
					"-password": complete.PredictAnything,
					"-key":      complete.PredictAnything,
				},
			},
			"change-password": {
				Flags: complete.Flags{
					"-username": predictUsers(),
					"-password": complete.PredictAnything,
				},
			},
			"list-users": {},
			"generate-env": {
				Flags: complete.Flags{
					"-output": complete.PredictFiles("*"),
				},
			},
			"completion": {
				Flags: complete.Flags{
					"-shell":  complete.PredictSet("bash", "zsh", "fish"),
					"-output": complete.PredictFiles("*"),
				},
			},
			"help": {},
		},
		Flags: complete.Flags{
			"-v":       complete.PredictNothing,
			"-version": complete.PredictNothing,
		},
		GlobalFlags: complete.Flags{
			"-v":       complete.PredictNothing,
			"-version": complete.PredictNothing,
		},
	}

	return cmd
}

// predictClasses predicts class names
func predictClasses() complete.Predictor {
	return complete.PredictFunc(func(args complete.Args) []string {
		// Try to list classes
		classes, err := GetClassTerms()
		if err != nil {
			return nil
		}
		return classes
	})
}

// predictAssignmentIDs predicts assignment IDs
func predictAssignmentIDs() complete.Predictor {
	return complete.PredictFunc(func(args complete.Args) []string {
		// Try to get the class name from the command line
		className := ""
		for i, arg := range args.Completed {
			if arg == "-class" && i+1 < len(args.Completed) {
				className = args.Completed[i+1]
				break
			}
		}

		if className == "" {
			return nil
		}

		// Try to list assignments for the class
		assignments, err := GetAssignments(className)
		if err != nil {
			return nil
		}
		return assignments
	})
}

// predictUsers predicts user names
func predictUsers() complete.Predictor {
	return complete.PredictFunc(func(args complete.Args) []string {
		// Try to list users
		usernames, err := GetUsernames()
		if err != nil {
			return nil
		}
		return usernames
	})
}

// InstallCompletion installs the completion script for the current shell
func InstallCompletion() error {
	// Get the shell from the environment
	shell := os.Getenv("SHELL")
	if shell == "" {
		return fmt.Errorf("could not determine shell type from environment")
	}

	// Extract the shell name from the path
	shell = filepath.Base(shell)

	// Normalize shell name
	shell = strings.ToLower(shell)

	// Map shell name to supported shell type
	switch {
	case shell == "bash":
		shell = "bash"
	case shell == "zsh":
		shell = "zsh"
	case shell == "fish":
		shell = "fish"
	default:
		return fmt.Errorf("unsupported shell type: %s (supported: bash, zsh, fish)", shell)
	}

	// Get the binary name
	binaryName := filepath.Base(os.Args[0])

	// Generate the completion script
	var script string
	switch shell {
	case "bash":
		script = generateBashCompletion(binaryName)
	case "zsh":
		script = generateZshCompletion(binaryName)
	case "fish":
		script = generateFishCompletion(binaryName)
	}

	// Determine the output file
	var outFile string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	switch shell {
	case "bash":
		outFile = filepath.Join(homeDir, fmt.Sprintf(".%s-completion.bash", binaryName))
	case "zsh":
		outFile = filepath.Join(homeDir, fmt.Sprintf(".%s-completion.zsh", binaryName))
	case "fish":
		fishDir := filepath.Join(homeDir, ".config", "fish")
		if err := os.MkdirAll(fishDir, 0755); err != nil {
			return fmt.Errorf("failed to create fish config directory: %w", err)
		}
		outFile = filepath.Join(fishDir, fmt.Sprintf("%s-completion.fish", binaryName))
	}

	// Write the completion script
	if err := os.WriteFile(outFile, []byte(script), 0644); err != nil {
		return fmt.Errorf("failed to write completion script to %s: %w", outFile, err)
	}

	fmt.Printf("Completion script written to %s\n", outFile)

	// Add the source line to the shell configuration file
	var configFile string
	var sourceLine string
	switch shell {
	case "bash":
		configFile = filepath.Join(homeDir, ".bashrc")
		sourceLine = fmt.Sprintf("source %s", outFile)
	case "zsh":
		configFile = filepath.Join(homeDir, ".zshrc")
		sourceLine = fmt.Sprintf("source %s", outFile)
	case "fish":
		configFile = filepath.Join(homeDir, ".config", "fish", "config.fish")
		sourceLine = fmt.Sprintf("source %s", outFile)
	}

	// Check if the source line already exists in the config file
	if _, err := os.Stat(configFile); err == nil {
		content, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to read shell config file %s: %w", configFile, err)
		}

		if strings.Contains(string(content), sourceLine) {
			fmt.Printf("Completion already installed in %s\n", configFile)
			return nil
		}
	}

	// Append the source line to the config file
	f, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open shell config file %s: %w", configFile, err)
	}
	defer f.Close()

	if _, err := f.WriteString("\n# Shell completion for " + binaryName + "\n" + sourceLine + "\n"); err != nil {
		return fmt.Errorf("failed to write to shell config file %s: %w", configFile, err)
	}

	fmt.Printf("Completion installed in %s\n", configFile)
	fmt.Printf("Restart your shell or run 'source %s' to enable completion\n", configFile)

	return nil
}
