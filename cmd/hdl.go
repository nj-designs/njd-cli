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
	"github.com/nj-designs/njd-cli/internal/cmds/hdl"
	"github.com/spf13/cobra"
)

// hdlCmd represents the hdl command
var hdlCmd = &cobra.Command{
	Use:   "hdl",
	Short: "HDL related commands",
}

// hdlSynthCmd represents the hld synth command
var hdlSynthCmd = &cobra.Command{
	Use:   "synth <project dir>",
	Short: "Synthesise given HDL project",
	Long: `Run synthesis step on given HDL project.

A HDL project is a directory containing a project file 'hdl-project.json'

Example: njd-cli hdl synth .

Synthesise project in current directory

`,

	Run:  hdl.Run,
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(hdlCmd)

	// hdl synth sub command
	hdlSynthCmd.Flags().StringVar(&hdl.SynthFlags.JSON, "json", "", "Save output as json file")
	hdlCmd.AddCommand(hdlSynthCmd)
}
