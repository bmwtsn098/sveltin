/**
 * Copyright © 2021 Mirco Veltri <github@mircoveltri.me>
 *
 * Use of this source code is governed by Apache 2.0 license
 * that can be found in the LICENSE file.
 */

// Package cmd ...
package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/sveltinio/sveltin/helpers"
	"github.com/sveltinio/sveltin/resources"
	"github.com/sveltinio/sveltin/utils"
)

//=============================================================================

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "Update the dependencies from the `package.json` file",
	Long: resources.GetAsciiArt() + `
Update all dependencies from the package.json file.

It wraps (npm|pnpm|yarn) update.
`,
	Run: RunUpdateCmd,
}

// RunUpdateCmd is the actual work function.
func RunUpdateCmd(cmd *cobra.Command, args []string) {
	textLogger.Reset()
	textLogger.SetTitle("Update dependencies Sveltin project")
	textLogger.SetContent("* Updating dependencies")

	pathToPkgFile := filepath.Join(pathMaker.GetRootFolder(), "package.json")
	npmClient, err := utils.RetrievePackageManagerFromPkgJson(AppFs, pathToPkgFile)
	utils.CheckIfError(err)

	// LOG TO STDOUT
	utils.PrettyPrinter(textLogger).Print()

	err = helpers.RunPMCommand(npmClient.Name, "update", "", nil, false)
	utils.CheckIfError(err)
}

func init() {
	rootCmd.AddCommand(updateCmd)
}