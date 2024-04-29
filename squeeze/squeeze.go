package squeeze

import (
	"flag"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CountCandidateFiles(directory string) int {
	var count int
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			slog.Error("error walking directory", "error", err)
			return err
		}

		if !info.IsDir() && isCommittedInGit(path) {
			count++
		}

		return nil
	})
	if err != nil {
		slog.Error("error counting candidate files", "error", err)
		os.Exit(1)
	}

	return count
}

func RunSqueezer(directory string) {
	var excludeDirs arrayFlags
	flag.Var(&excludeDirs, "exclude-dir", "Directories to exclude (can be specified multiple times)")
	flag.Parse()

	slog.Debug("processing directory", "dir", directory)
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			slog.Error("error walking directory", "error", err)
			return err
		}

		if info.IsDir() {
			if info.Name() == ".git" || excludeDirs.contains(info.Name()) {
				slog.Debug("skipping directory", "dir", info.Name())
				return filepath.SkipDir
			}
			return nil
		}

		slog.Debug("processing file", "file", path)
		if isCommittedInGit(path) {
			squeezeNewlines(path)
		} else {
			slog.Debug("skipping file not committed in Git", "file", path)
		}

		return nil
	})
	if err != nil {
		slog.Error("error running squeezer", "error", err)
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
		slog.Error("error reading file", "file", file, "error", err)
		return
	}

	lines := strings.Split(string(content), "\n")
	var squeezed []string
	for i := 0; i < len(lines); i++ {
		if i > 0 && (lines[i] == "" || lines[i] == "\r") && (lines[i-1] == "" || lines[i-1] == "\r") {
			continue
		}
		squeezed = append(squeezed, lines[i])
	}

	err = os.WriteFile(file, []byte(strings.Join(squeezed, "\n")), 0o644)
	if err != nil {
		slog.Error("error writing file", "file", file, "error", err)
	} else {
		slog.Debug("file written successfully", "file", file)
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
