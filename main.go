package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "syscall"

    "github.com/goccy/go-yaml"
)

type Node struct {
	Cmd string `yaml:"cmd,omitempty"`
	Run string `yaml:"run,omitempty"`
	Args []*Node `yaml:"args,omitempty"`
}

func main() {
	root, err := ParseConfig()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("No arguments provided.")
		return
	}

	var command []string
	pointer := root
	for _, arg := range args {
		found := false
		for _, possible := range pointer.Args {
			if possible.Cmd == arg {
				pointer = possible
				command = append(command, possible.Run)
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("Could not find command: %v", arg)
		}
	}
	binary, err := exec.LookPath(command[0])
	if err != nil {
		log.Fatal(err)
	}

	env := os.Environ()
	execErr := syscall.Exec(binary, command, env)
	
	if execErr != nil {
		log.Fatalf("Failed to start process: %v", execErr)
	}
}

func ParseConfig() (*Node, error) {
	homedir, err := os.UserHomeDir()
    if err != nil {
		return nil, err
    }

	data, err := os.ReadFile(homedir + "/.config/runner/config.yaml")
    if err != nil {
		return nil, err
    }

    var root Node
    if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, err
    }

	return &root, nil
}

