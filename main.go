package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Defaults struct {
		Proxy string `yaml:"proxy"`
		Port  int    `yaml:"port"`
	} `yaml:"defaults"`
	Environments map[string]struct {
		Proxy string `yaml:"proxy"`
		Port  int    `yaml:"port"`
	} `yaml:"environments"`
}

func main() {
	var (
		envFlag     string
		verboseFlag bool
	)
	flag.StringVar(&envFlag, "env", "", "Environment to use")
	flag.BoolVar(&verboseFlag, "verbose", false, "Enable verbose output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-env <environment>] [-verbose] [--] <command> [args...]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Remove the "--" separator if present
	if args[0] == "--" {
		args = args[1:]
	}

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No command specified")
		os.Exit(1)
	}

	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	proxy := config.Defaults.Proxy
	port := config.Defaults.Port

	if envFlag != "" {
		if env, ok := config.Environments[envFlag]; ok {
			if env.Proxy != "" {
				proxy = env.Proxy
			}
			if env.Port != 0 {
				port = env.Port
			}
		} else {
			fmt.Fprintf(os.Stderr, "Environment '%s' not found in config\n", envFlag)
			os.Exit(1)
		}
	}

	proxyURL := fmt.Sprintf("%s:%d", proxy, port)

	// Check if the command exists
	cmdPath, err := exec.LookPath(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "htproxctl: Error: Command '%s' not found in PATH\n", args[0])
		fmt.Fprintf(os.Stderr, "Please make sure the command is installed and accessible.\n")
		os.Exit(1)
	}

	cmd := exec.Command(cmdPath, args[1:]...)

	// Set environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("HTTP_PROXY=%s", proxyURL))
	env = append(env, fmt.Sprintf("HTTPS_PROXY=%s", proxyURL))
	cmd.Env = env

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Only output debug information if verbose flag is set
	if verboseFlag {
		fmt.Fprintf(os.Stderr, "htproxctl: Executing with HTTP_PROXY=%s, HTTPS_PROXY=%s\n", proxyURL, proxyURL)
		fmt.Fprintf(os.Stderr, "Command: %s\n", strings.Join(args, " "))
	}

	err = cmd.Run()
	if err != nil {
		if verboseFlag {
			fmt.Fprintf(os.Stderr, "htproxctl: Error executing command: %v\n", err)
		}
		os.Exit(1)
	}
}

func loadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "htproxctl.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
