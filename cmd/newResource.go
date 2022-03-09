/**
 * Copyright © 2021 Mirco Veltri <github@mircoveltri.me>
 *
 * Use of this source code is governed by Apache 2.0 license
 * that can be found in the LICENSE file.
 */

// Package cmd ...
package cmd

import (
	"errors"
	"fmt"

	"github.com/sveltinio/sveltin/common"
	"github.com/sveltinio/sveltin/config"
	"github.com/sveltinio/sveltin/helpers"
	"github.com/sveltinio/sveltin/helpers/factory"
	"github.com/sveltinio/sveltin/resources"
	"github.com/sveltinio/sveltin/sveltinlib/composer"
	"github.com/sveltinio/sveltin/sveltinlib/logger"
	"github.com/sveltinio/sveltin/sveltinlib/sveltinerr"
	"github.com/sveltinio/sveltin/utils"

	"github.com/spf13/cobra"
)

//=============================================================================

var newResourceCmd = &cobra.Command{
	Use:     "resource [name]",
	Aliases: []string{"r"},
	Short:   "Command to create new resources",
	Long: resources.GetAsciiArt() + `
Command to create new resources.

What is a "resource" for Sveltin?
A resource is a way to group, serve and expose your content.

This command:

- Create a <resource_name> folder within "content" folder, so that you can add new content for the resource
- Scaffold a GET endpoint for the resource within "src/routes/api/<api_version>/<resource_name> so that you resource is ready to be consumed by the REST endpoint
- Add the resource as route within the "src/routes" folder, creating its own folder
- Scaffold index.svelte component to list all the content belongs to a resource
- Scaffold [slug].svelte component
	`,
	Run: RunNewResourceCmd,
}

// RunNewResourceCmd is the actual work function.
func RunNewResourceCmd(cmd *cobra.Command, args []string) {
	listLogger := log.WithList()

	resourceName, err := promptResourceName(args)
	utils.ExitIfError(err)

	// GET FOLDER: content folder
	contentFolder := fsManager.GetFolder(CONTENT)

	// NEW FOLDER: content/<resource_name>. Here is where the "new content" command saves files
	listLogger.Append(logger.LevelInfo, "Creating the content folder for your resource")
	resourceContentFolder := composer.NewFolder(resourceName)
	contentFolder.Add(resourceContentFolder)

	// GET FOLDER: src/routes/api/<api_version> folder
	apiFolder := fsManager.GetFolder(API)

	// NEW FOLDER: src/routes/api/<version>/<resource_name>
	resourceAPIFolder := composer.NewFolder(resourceName)

	// NEW FILE: src/routes/api/<resource_name>/published.json.ts
	listLogger.Append(logger.LevelInfo, "Creating the API endpoint for the resource")
	apiFile := &composer.File{
		Name:       conf.GetAPIFilename(),
		TemplateId: API,
		TemplateData: &config.TemplateData{
			Name:   resourceName,
			Config: &conf,
		},
	}
	resourceAPIFolder.Add(apiFile)
	apiFolder.Add(resourceAPIFolder)

	// GET FOLDER: src/lib folder
	libFolder := fsManager.GetFolder(LIB)

	// NEW FOLDER: /src/lib/<resource_name>
	resourceLibFolder := composer.NewFolder(resourceName)

	// NEW FILE: src/lib/<resource_name>/get<resource_name>.js
	listLogger.Append(logger.LevelInfo, "Creating the lib file for the resource")
	libFile := &composer.File{
		Name:       utils.ToLibFile(resourceName),
		TemplateId: LIB,
		TemplateData: &config.TemplateData{
			Name:   resourceName,
			Config: &conf,
		},
	}
	resourceLibFolder.Add(libFile)
	libFolder.Add(resourceLibFolder)

	// GET FOLDER: src/routes folder
	routesFolder := fsManager.GetFolder(ROUTES)

	// NEW FOLDER: src/routes/<resource_name>
	resourceRoutesFolder := composer.NewFolder(resourceName)

	// NEW FILE: <resource_name>/index.svelte and [slug].svelte
	for _, item := range []string{INDEX, SLUG} {
		listLogger.Append(logger.LevelInfo, fmt.Sprintf("Creating the %s component for the resource", item))
		f := &composer.File{
			Name:       helpers.GetResourceRouteFilename(item, &conf),
			TemplateId: item,
			TemplateData: &config.TemplateData{
				Name:   resourceName,
				Config: &conf,
			},
		}
		resourceRoutesFolder.Add(f)
	}
	routesFolder.Add(resourceRoutesFolder)

	// SET FOLDER STRUCTURE
	projectFolder := fsManager.GetFolder(ROOT)
	projectFolder.Add(contentFolder)
	projectFolder.Add(apiFolder)
	projectFolder.Add(libFolder)
	projectFolder.Add(routesFolder)

	// CREATE FOLDER STRUCTURE
	sfs := factory.NewResourceArtifact(&resources.SveltinFS, AppFs)
	err = projectFolder.Create(sfs)
	utils.ExitIfError(err)

	// LOG TO STDOUT
	listLogger.Info(fmt.Sprintf("New '%s' as resource will be created", resourceName))
	log.Success("Done")

	// NEXT STEPS
	listLogger = log.WithList()
	listLogger.Append(logger.LevelSuccess, "Resource ready to be used. Start by adding content to it.")
	listLogger.Append(logger.LevelImportant, fmt.Sprintf("Eg: sveltin new content %s/getting-started", resourceName))
	listLogger.Info("Next Steps")
}

func init() {
	newCmd.AddCommand(newResourceCmd)
}

//=============================================================================

func promptResourceName(inputs []string) (string, error) {
	var name string
	switch numOfArgs := len(inputs); {
	case numOfArgs < 1:
		resourceNamePromptContent := config.PromptContent{
			ErrorMsg: "Please, provide a name for the resource.",
			Label:    "What's the resource name? (e.g. posts, portfolio ...)",
		}
		name = common.PromptGetInput(resourceNamePromptContent)
		return utils.ToSlug(name), nil
	case numOfArgs == 1:
		name = inputs[0]
		return utils.ToSlug(name), nil
	default:
		err := errors.New("something went wrong: value not valid")
		return "", sveltinerr.NewDefaultError(err)
	}
}
