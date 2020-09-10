/*
Copyright © 2020 Neil Johnson <nj.designs@protonmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const projectFileName string = "hdl-project.json"

type hdlProject struct {
	Top     string   `json:"top"`
	Modules []string `json:"modules"`
}

func readHDLProjectFile(project *hdlProject) error {

	contents, err := ioutil.ReadFile(projectFileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(contents, project)
	if err != nil {
		return err
	}
	return nil
}

func discoverSrcFiles(project *hdlProject) ([]string, error) {
	var err error
	var filteredSrcFiles []string

	filteredSrcFiles = append(filteredSrcFiles, project.Top)
	for _, moduleDir := range project.Modules {
		var srcFiles []os.FileInfo
		srcFiles, err = ioutil.ReadDir(moduleDir)
		if err != nil {
			return nil, err
		}
		for _, srcFile := range srcFiles {
			srcFileName := srcFile.Name()
			// Currently only verilog files
			if strings.HasSuffix(srcFileName, ".v") == false {
				continue
			}
			// Ignore test bench files
			if strings.HasSuffix(srcFileName, "_tb.v") {
				continue
			}
			filteredSrcFiles = append(filteredSrcFiles, path.Join(moduleDir, srcFileName))
		}
	}

	return filteredSrcFiles, nil
}

func runHDLSynth(cmd *cobra.Command, args []string) {
	var err error
	var absProjectDir string

	// Get absolute project directory
	if absProjectDir, err = filepath.Abs(args[0]); err != nil {
		log.Fatalf("Failed to get abs project dir: %v", err)
	}

	// Switch to project dir
	if err = os.Chdir(absProjectDir); err != nil {
		log.Fatalf("Failed to switch to project dir: %v", err)
	}

	// Does project file exist
	if err = isFileReadable(projectFileName); err != nil {
		log.Fatalf("Can't read project file in %s : %v", args[0], err)
	}

	// Project file contents
	var project hdlProject

	if err = readHDLProjectFile(&project); err != nil {
		log.Fatalf("Can't read project file in %s : %v", projectFileName, err)
	}

	srcFiles, err := discoverSrcFiles(&project)
	if err != nil {
		log.Fatalf("Failed to discovet src files: %v", err)
	}

	fmt.Println(srcFiles)
}

func init() {
	hdlCmd.AddCommand(hdlSynthCmd)

	// hdlSynthCmd.Flags().BoolVarP(&recurse, "recurse", "r", false, "Clone & init sub modules")
}

// cloneCmd represents the clone command
var hdlSynthCmd = &cobra.Command{
	Use:   "synth <project dir>",
	Short: "Synthesise given HDL project",
	Long: `Run synthesis step on given HDL project.

A HDL project is a directory containing a project file 'hdl-project.json'

Example: njd-cli hdl synth .

Synthesise project in current directory

`,

	Run:  runHDLSynth,
	Args: cobra.ExactArgs(1),
}
