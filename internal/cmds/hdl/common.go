package hdl

// SynthFlagsT cotains
type SynthFlagsT struct {
	Verbose bool
	Show    bool
	JSON    string
}

// SynthFlags cotains
var SynthFlags SynthFlagsT

// Expected name of HDL project file.
// It's location is the root dir of the project
const projectFileName string = "hdl-project.json"

const yosysExe string = "yosys"

// Supported FPGA
var supportedFPGAList = []string{
	"ice40",
	"intel",
}

// hdlProject contains unmarshalled contents of project file
type hdlProject struct {
	Name    string   `json:"name"`
	FPGA    string   `json:"fpga"`
	Top     string   `json:"top"`
	Modules []string `json:"modules"`
}

func isSupportedFPGA(fpga string) bool {
	for _, f := range supportedFPGAList {
		if f == fpga {
			return true
		}
	}
	return false
}
