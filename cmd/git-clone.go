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
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var recurse bool //Recursive clone

func runGitClone(cmd *cobra.Command, args []string) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	for _, repoURL := range args {
		var h, p string

		if strings.HasPrefix(repoURL, "git@") {
			// Assume repoURL is something like "git@github.com:thingsforhackers/CopyMe84.git"
			parts := strings.Split(repoURL, ":")
			if len(parts) != 2 {
				log.Fatalf("runGitClone() : Don't know how to parse '%s'", repoURL)
			}
			h = strings.Split(parts[0], "@")[1]
			p = strings.TrimSuffix(parts[1], ".git")
		} else {
			// Assume repoURL is somethng like "https://github.com/thingsforhackers/Arduino-Make.git"
			u, err := url.Parse(repoURL)
			if err != nil {
				log.Fatalf("runGitClone() : failed to parse url '%s'", repoURL)
			}
			h = u.Host
			p = strings.TrimSuffix(u.Path, ".git")
		}

		absCloneDir := path.Join(userHomeDir, h, p)

		if _, err := os.Stat(absCloneDir); !os.IsNotExist(err) {
			fmt.Printf("'%s' already exists, skipping\n", absCloneDir)
			continue
		}

		fmt.Printf("Cloning '%s' into '%s' recurse:%v\n", repoURL, absCloneDir, recurse)

		var recurseSubmodules = git.NoRecurseSubmodules
		if recurse {
			recurseSubmodules = git.DefaultSubmoduleRecursionDepth
		}
		_, err := git.PlainClone(absCloneDir, false, &git.CloneOptions{
			URL:               repoURL,
			RecurseSubmodules: recurseSubmodules,
		})
		if err != nil {
			log.Fatalln(err)
		}

	}

}

func init() {
	gitCmd.AddCommand(gitCloneCmd)

	gitCloneCmd.Flags().BoolVarP(&recurse, "recurse", "r", false, "Clone & init sub modules")
}

// cloneCmd represents the clone command
var gitCloneCmd = &cobra.Command{
	Use:   "clone [REPO]...",
	Short: "Clone one or more git repos to usual location",
	Long: `Git clone one or more repos

Example: njd-cli git clone git@github.com:thingsforhackers/CopyMe84.git

Repo will be cloned in to $HOME/github.com/thingsforhackers/CopyMe84

`,

	Run:  runGitClone,
	Args: cobra.MinimumNArgs(1),
}
