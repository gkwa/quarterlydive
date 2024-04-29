package squeeze

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunSqueezer() {
	var excludeDirs arrayFlags
	flag.Var(&excludeDirs, "exclude-dir", "Directories to exclude (can be specified multiple times)")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Usage: squeeze-newlines [--exclude-dir dir] <directory>")
		os.Exit(1)
	}

	directory := flag.Args()[0]
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == ".git" || excludeDirs.contains(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if isCommittedInGit(path) {
			squeezeNewlines(path)
		}

		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
		os.Exit(1)
	}
}

func isCommittedInGit(file string) bool {
	cmd := exec.Command("git", "ls-files", "--error-unmatch", file)
	err := cmd.Run()
	return err == nil
}

func squeezeNewlines(file string) {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file, err)
		return
	}

	lines := strings.Split(string(content), "\n")
	var squeezed []string
	for i := 0; i < len(lines); i++ {
		if i > 0 && lines[i] == "" && lines[i-1] == "" {
			continue
		}
		squeezed = append(squeezed, lines[i])
	}

	err = os.WriteFile(file, []byte(strings.Join(squeezed, "\n")), 0o644)
	if err != nil {
		fmt.Printf("Error writing file %s: %v\n", file, err)
	}
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *arrayFlags) contains(value string) bool {
	for _, v := range *i {
		if v == value {
			return true
		}
	}
	return false
}
