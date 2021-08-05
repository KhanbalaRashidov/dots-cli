package commands

import (
	"bytes"
	"fmt"
	"github.com/alvanrahimli/dots-cli/models"
	"github.com/alvanrahimli/dots-cli/utils"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

// Push command pushes package to already added registry
// Registries can be added using `dots-cli remote add` command
type Push struct {
}

func (p Push) GetArguments() []string {
	return []string{}
}

func (p Push) CheckRequirements() (bool, string) {
	return true, ""
}

func (p Push) ExecuteCommand(opts *models.Opts, config *models.AppConfig) models.CommandResult {
	// Check for token
	if config.AuthorToken == "" {
		return models.CommandResult{
			Code:    1,
			Message: "User is not signed in",
		}
	}

	manifest, err := utils.ReadManifestFile(opts.OutputDir)
	if err != nil {
		return models.CommandResult{
			Code:    1,
			Message: "Could not read manifest. Did you initialize package?",
		}
	}

	var selectedVersion string

	// If more than 1 version found
	if len(manifest.Versions) > 1 {
		fmt.Println("Several versions found. Type version of package you want to push")
		for _, version := range manifest.Versions {
			fmt.Printf("\n\t%s", version.ToString())
		}

		fmt.Print("\nEnter version: ")
		_, scanErr := fmt.Scanln(&selectedVersion)
		if scanErr != nil {
			return models.CommandResult{
				Code:    1,
				Message: "Could not read selected version",
			}
		}

		versionRe := regexp.MustCompile("^([0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3})")
		if !versionRe.MatchString(selectedVersion) {
			return models.CommandResult{
				Code:    1,
				Message: "Invalid version number entered",
			}
		}

		// Find package name by version:
		selectedVersion = strings.ReplaceAll(selectedVersion, ".", "_")
	} else {
		selectedVersion = manifest.Versions[0].ToFormattedString()
	}

	versFolder := path.Join(opts.OutputDir, ".vers")
	archiveName := fmt.Sprintf(models.TarballNameFormat, manifest.Name, selectedVersion)
	archiveFullPath := path.Join(versFolder, archiveName)
	_, statErr := os.Stat(archiveFullPath)
	if statErr != nil {
		return models.CommandResult{
			Code:    1,
			Message: fmt.Sprintf("Could not read archive file '%s'", archiveFullPath),
		}
	}

	// SEND POST REQUEST
	client := &http.Client{}
	values := map[string]io.Reader{
		"name":    strings.NewReader(manifest.Name),
		"version": strings.NewReader(strings.ReplaceAll(selectedVersion, "_", ".")),
		"archive": mustOpen(archiveFullPath),
	}
	uploadErr := Upload(client, models.AddPackageEndpoint, values, config)
	if uploadErr != nil {
		fmt.Println(uploadErr.Error())
		return models.CommandResult{
			Code:    1,
			Message: "Could not send request",
		}
	}

	return models.CommandResult{
		Code:    0,
		Message: "Package pushed successfully",
	}
}

func Upload(client *http.Client, url string, values map[string]io.Reader, config *models.AppConfig) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", config.AuthorToken)

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
