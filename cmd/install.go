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
	"github.com/sveltinio/sveltin/sveltinlib/logger"
	"github.com/sveltinio/sveltin/utils"
)

//=============================================================================

var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"i", "init"},
	Short:   "Get all the dependencies from the `package.json` file",
	Long: resources.GetAsciiArt() + `
Initialize the Sveltin project getting all dependencies from the package.json file.

It wraps (npm|pnpm|yarn) install.
`,
	Run: RunInstallCmd,
}

// RunInstallCmd is the actual work function.
func RunInstallCmd(cmd *cobra.Command, args []string) {
	listLogger := log.WithList()
	listLogger.Append(logger.LevelInfo, "Getting dependencies")
	listLogger.Info("Prepare Sveltin project")

	pathToPkgFile := filepath.Join(pathMaker.GetRootFolder(), "package.json")
	npmClient, err := utils.RetrievePackageManagerFromPkgJson(AppFs, pathToPkgFile)
	utils.ExitIfError(err)

	err = helpers.RunPMCommand(npmClient.Name, "install", "", nil, false)
	utils.ExitIfError(err)
	log.Success("Done")
}

func init() {
	rootCmd.AddCommand(installCmd)
}
