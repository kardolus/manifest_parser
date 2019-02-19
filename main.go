package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	USAGE   = "Usage:   go run main.go <path_to_buildpack/manifest.yml>\n"
	EXAMPLE = "Example: go run main.go ~/workspace/go-buildpack/manifest.yml\n"
)

func init() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, USAGE)
		fmt.Fprint(os.Stderr, EXAMPLE)
		os.Exit(0)
	}
}

func main() {
	flag.Parse()

	args := len(flag.Args())
	if args != 1 {
		flag.Usage()
	}

	fileName := flag.Arg(0)
	data := Manifest{}
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("Unable to read file")
	}
	if err := yaml.Unmarshal(file, &data); err != nil {
		log.Fatal("Unable to unmarshal", err.Error())
	}

	for i, dep := range data.Dependencies {
		if dep.Source != "" && dep.SourceSha256 == "" {
			cmd := exec.Command("sh", "getsha.sh", dep.Source)
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatalf("Unable to run getsha script: %s", err.Error())
			}
			data.Dependencies[i].SourceSha256 = strings.TrimSpace(string(output))
		}
	}

	output, err := yaml.Marshal(&data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := ioutil.WriteFile("new_manifest.yml", output, 0666); err != nil {
		log.Fatalf("error writing to file %s", err.Error())
	}
}

type Manifest struct {
	Language        string `yaml:"language"`
	DefaultVersions []struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"default_versions"`
	DependencyDeprecationDates []struct {
		VersionLine string `yaml:"version_line"`
		Name        string `yaml:"name"`
		Date        string `yaml:"date"`
		Link        string `yaml:"link"`
	} `yaml:"dependency_deprecation_dates"`
	Dependencies []struct {
		Name         string   `yaml:"name"`
		Version      string   `yaml:"version"`
		URI          string   `yaml:"uri"`
		Sha256       string   `yaml:"sha256"`
		CfStacks     []string `yaml:"cf_stacks"`
		Source       string   `yaml:"source,omitempty"`
		SourceSha256 string   `yaml:"source_sha256,omitempty"`
	} `yaml:"dependencies"`
	IncludeFiles []string `yaml:"include_files"`
}
