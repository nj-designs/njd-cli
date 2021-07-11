/*
Copyright © 2021 Neil Johnson <nj.designs@protonmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/nj-designs/njd-cli/internal/cmds/tiled"
	"github.com/spf13/cobra"
)

// tiledCmd represents the tiled command
var tiledCmd = &cobra.Command{
	Use:   "tiled",
	Short: "A collection of tools to process Tiled mapeditor files",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run:  tiled.Run,
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(tiledCmd)
	tiledCmd.Flags().StringVar(&tiled.TilemapFileName, "tilemap", "", "Input tilemap file (required)")
	tiledCmd.MarkFlagRequired("tilemap")
	tiledCmd.Flags().StringVar(&tiled.OutputFileName, "output", "", "Output file (required)")
	tiledCmd.MarkFlagRequired("output")
}
