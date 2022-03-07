/**
 * Copyright © 2021 Mirco Veltri <github@mircoveltri.me>
 *
 * Use of this source code is governed by Apache 2.0 license
 * that can be found in the LICENSE file.
 */

// Package utils ...
package utils

import (
	"os/exec"

	"github.com/spf13/afero"
	"github.com/sveltinio/sveltin/sveltinlib/npmc"
	"github.com/sveltinio/sveltin/sveltinlib/sveltinerr"
)

// GetInstalledNPMClientList returns the list of installed npmClient as slice of NPMClient.
func GetInstalledNPMClientList() []npmc.NPMClient {
	valid := []string{"npm", "yarn", "pnpm"}
	npmClientList := []npmc.NPMClient{}
	for _, pm := range valid {
		valid, version := GetNPMClientInfo(pm)
		if valid {
			a := npmc.NPMClient{
				Name:    pm,
				Version: version,
			}
			npmClientList = append(npmClientList, a)
		}
	}
	return npmClientList
}

// GetNPMClientNames returns the list of installed npmClient as slice of strings.
func GetNPMClientNames(items []npmc.NPMClient) []string {
	npmClientNames := []string{}
	for _, v := range items {
		npmClientNames = append(npmClientNames, v.Name)
	}
	return npmClientNames
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command(name, "-v")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// GetNPMClientInfo returns true and npm client version as string.
func GetNPMClientInfo(name string) (bool, string) {
	out, err := exec.Command(name, "-v").Output()
	if err != nil {
		return false, ""
	}
	return true, string(out)
}

func filter(in []npmc.NPMClient, name string) npmc.NPMClient {
	var out npmc.NPMClient
	for _, each := range in {
		if each.Name == name {
			out = each
		}
	}
	return out
}

// GetSelectedNPMClient returns the selected NPMClient struct out of the available ones.
func GetSelectedNPMClient(in []npmc.NPMClient, name string) npmc.NPMClient {
	return filter(in, name)
}

func RetrieveProjectName(appFS afero.Fs, pathToPkgJson string) (string, error) {
	pkgFileContent, err := afero.ReadFile(appFS, pathToPkgJson)
	ExitIfError(err)
	pkgParsed := npmc.Parse(pkgFileContent)
	if pkgParsed.Name != "" {
		return pkgParsed.Name, nil
	}

	return "", sveltinerr.NewProjectNameNotFoundError()
}

// RetrievePackageManagerFromPkgJson returns NPMClient struct parsing the package.json file.
func RetrievePackageManagerFromPkgJson(appFS afero.Fs, pathToPkgJson string) (npmc.NPMClient, error) {
	pkgFileContent, err := afero.ReadFile(appFS, pathToPkgJson)
	ExitIfError(err)
	pkgParsed := npmc.Parse(pkgFileContent)
	if pkgParsed.PackageManager != "" {
		pmInfoString := npmc.NPMClientInfoStr(pkgParsed.PackageManager)
		return pmInfoString.ToNPMClient(), nil
	}
	return npmc.NPMClient{}, sveltinerr.NewPackageManagerKeyNotFoundOnPackageJSONFile()
}
