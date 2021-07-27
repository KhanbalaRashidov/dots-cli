package commands

import (
	"dots/apphandlers"
	"dots/models"
	"dots/utils"
	"fmt"
	"os"
	"path"
	"strings"
)

// Add command when name is specified adds new application to package.
type Add struct {
	Options *models.Opts
}

func (a Add) GetArguments() []string {
	return []string{}
}

func (a Add) CheckRequirements() (bool, string) {
	if len(a.Options.Arguments) < 2 {
		return false, fmt.Sprintf("%s is not enough arguments for add command.", a.Options.Arguments)
	}

	return true, ""
}

func (a Add) ExecuteCommand(opts *models.Opts) models.CommandResult {
	a.Options = opts
	// Check if arguments satisfy required arguments for add command
	satisfiesRequirements, message := a.CheckRequirements()
	if !satisfiesRequirements {
		fmt.Printf("Add command can not work in this directory:\n\t%s\n", message)
		os.Exit(1)
	}

	// Read manifest
	manifest, manifestErr := utils.ReadManifestFile(opts.OutputDir)
	if manifestErr != nil {
		fmt.Printf("ERROR: %s\n", manifestErr.Error())
		os.Exit(1)
	}

	// Check if apps exist in package or not
	possibleAppNames := make([]string, 0)
	for _, appName := range opts.Arguments[1:] {
		appExists := utils.AppExistsInManifest(appName, &manifest)
		if appExists {
			fmt.Printf("App %s already exists in this package\n", appName)
			continue
		}

		possibleAppNames = append(possibleAppNames, appName)
	}

	// Exit app if all apps are in package
	if len(possibleAppNames) == 0 {
		fmt.Println("All packages are already in package")
		os.Exit(1)
	}

	addedApps := make([]string, 0)
	// Copy files to package
	for _, appName := range possibleAppNames {
		added, message := apphandlers.HandleApp(opts.OutputDir, appName)
		if added {
			addedApps = append(addedApps, appName)
			manifest.Apps = append(manifest.Apps, apphandlers.GetApp(appName))
		} else {
			fmt.Printf("ERROR: %s\n", message)
		}
	}

	// Write manifest to file
	manifestPath := path.Join(opts.OutputDir, "manifest.json")
	removeErr := os.Remove(manifestPath)
	if removeErr != nil {
		return models.CommandResult{
			Code:    1,
			Message: "Could not remove old manifest file",
		}
	}

	manifestWriteErr := WriteManifestFile(opts.OutputDir, &manifest)
	if manifestWriteErr != nil {
		return models.CommandResult{
			Code:    1,
			Message: "Could write updated manifest file",
		}
	}

	return models.CommandResult{
		Code: 0,
		Message: fmt.Sprintf("Following apps added to package: %s\n",
			strings.Join(addedApps, ", ")),
	}
}
