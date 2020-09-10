package hdl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/nj-designs/njd-cli/internal/njd"
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

// Run the HDL synth command
func Run(cmd *cobra.Command, args []string) {
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
	if err = njd.IsFileReadable(projectFileName); err != nil {
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
