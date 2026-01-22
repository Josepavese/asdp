package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	root := ".." // Run from tools/ directory to reach root via ..

	// 1. Get current version from tools/engine/domain/version.go
	versionFilePath := filepath.Join(root, "tools/engine/domain/version.go")
	versionContent, err := ioutil.ReadFile(versionFilePath)
	if err != nil {
		fmt.Printf("Error reading version file: %v\n", err)
		os.Exit(1)
	}

	versionRegex := regexp.MustCompile(`const Version = "([^"]+)"`)
	matches := versionRegex.FindStringSubmatch(string(versionContent))
	if len(matches) < 2 {
		fmt.Println("Error: could not find Version in version.go")
		os.Exit(1)
	}
	newVersion := matches[1]

	fmt.Printf("Syncing version v%s across the project...\n", newVersion)

	// 2. Sync README.md
	syncFile(filepath.Join(root, "README.md"), `\(v\d+\.\d+\.\d+\)`, "(v"+newVersion+")")
	syncFile(filepath.Join(root, "README.md"), `Status\*\*: Core Implementation Complete \(v\d+\.\d+\.\d+\)`, "Status**: Core Implementation Complete (v"+newVersion+")")

	// 3. Sync codetree.md
	syncFile(filepath.Join(root, "codetree.md"), `"asdp_version": "\d+\.\d+\.\d+"`, `"asdp_version": "`+newVersion+`"`)

	// 4. Sync all codespec.md and codemodel.md files
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "venv" || info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		name := info.Name()
		if name == "codespec.md" || name == "codemodel.md" {
			syncFile(path, `"ASDPVersion": "\d+\.\d+\.\d+"`, `"ASDPVersion": "`+newVersion+`"`)
			syncFile(path, `"asdp_version": "\d+\.\d+\.\d+"`, `"asdp_version": "`+newVersion+`"`)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking tree: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Version sync complete.")
}

func syncFile(path string, pattern string, replacement string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Printf("Error reading %s: %v\n", path, err)
		return
	}

	re := regexp.MustCompile(pattern)
	if !re.Match(content) {
		return
	}

	newContent := re.ReplaceAllString(string(content), replacement)
	if string(content) == newContent {
		return
	}

	err = ioutil.WriteFile(path, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("Error writing %s: %v\n", path, err)
		return
	}

	fmt.Printf("Updated %s\n", path)
}
