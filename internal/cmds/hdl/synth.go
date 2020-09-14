package hdl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/nj-designs/njd-cli/internal/njd"
	"github.com/spf13/cobra"
)

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

	filteredSrcFiles = append(filteredSrcFiles, "top.v")
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

func sanityCheckProject(project *hdlProject) error {

	// .name field
	if len(project.Name) == 0 {
		log.Fatalf("Missing 'name' field in project file")
	}

	// .fpga field
	if len(project.FPGA) == 0 {
		log.Fatalf("Missing 'fpga' field in project file")
	}
	if isSupportedFPGA(project.FPGA) == false {
		log.Fatalf("Unsupported fpga '%s', supported fpgas: %s", project.FPGA, strings.Join(supportedFPGAList, ", "))
	}

	// .top field
	if len(project.Top) == 0 {
		log.Fatalf("Missing 'top' field in project file")
	}

	return nil
}

// buildYOSYSSynthCmdFile returns the filepath of a YOSYS script to synthesise HDL project
func buildYOSYSSynthCmdFile(project *hdlProject) (string, error) {
	var sb strings.Builder

	fmt.Fprintf(&sb, "# Generated by njd-cli hdl synth\n")

	// Design files
	srcFiles, err := discoverSrcFiles(project)
	if err != nil {
		return "", err
	}
	fmt.Fprintf(&sb, "\n# Design files\n")
	for _, srcFile := range srcFiles {
		fmt.Fprintf(&sb, "read_verilog %s\n", srcFile)
	}

	// Synth
	fmt.Fprintf(&sb, "\n# Synth\n")
	fmt.Fprintf(&sb, "synth_%s -top %s", project.FPGA, project.Top)
	//   Add JSON output
	if len(SynthFlags.JSON) > 0 {
		fmt.Fprintf(&sb, " -json %s", SynthFlags.JSON)
	}
	fmt.Fprintf(&sb, "\n")

	// show
	if SynthFlags.Show {
		fmt.Fprintf(&sb, "show top")
	}

	fmt.Fprintf(&sb, "\n")

	f, err := ioutil.TempFile("/tmp", "yosys-cmds-*.txt")
	if err != nil {
		return "", err
	}
	defer f.Close()

	f.WriteString(sb.String())

	return f.Name(), nil
}

// Run the HDL synth command
func Run(cmd *cobra.Command, args []string) {
	var err error
	var absProjectDir string

	_, err = exec.LookPath(yosysExe)
	if err != nil {
		log.Fatal(err)
	}

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

	if err = sanityCheckProject(&project); err != nil {
		log.Fatalf("Invalid project file : %v", err)
	}

	synthCmdFile, err := buildYOSYSSynthCmdFile(&project)
	if err != nil {
		log.Fatalf("Failed to build cmd file: %v", err)
	}

	// Build actual yosys command line
	var yoySysArgs []string
	// var sb strings.Builder
	if SynthFlags.Verbose == false {
		yoySysArgs = append(yoySysArgs, "-q")
	}
	yoySysArgs = append(yoySysArgs, fmt.Sprintf("-s %s", synthCmdFile))

	// TODO(njohn) : Is there a better way to do this?
	p := strings.Split(strings.Join(yoySysArgs, " "), " ")
	yosysCmd := exec.Command(yosysExe, p...)

	allOP, err := yosysCmd.CombinedOutput()
	if len(allOP) > 0 {
		fmt.Printf("%s", allOP)
	}
	if err != nil {
		log.Fatalf("Synth failed : %v", err)
	}

	// fmt.Printf("\n\nSynth success!!!\n")
	fmt.Print(`
##############################################################
#                     Synth success!!!                       #
##############################################################


`)

	if len(SynthFlags.JSON) > 0 {
		fmt.Printf("  json output: %s\n", SynthFlags.JSON)
	}
}
