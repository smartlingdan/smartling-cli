package main

import (
	"os"
	"path/filepath"

	"github.com/Smartling/api-sdk-go"
	"github.com/gobwas/glob"
	hierr "github.com/reconquest/hierr-go"
)

func globFilesRemote(
	client *smartling.Client,
	project string,
	uri string,
) ([]smartling.File, error) {
	if uri == "" {
		uri = "**"
	}

	pattern, err := glob.Compile(uri, '/')
	if err != nil {
		return nil, NewError(
			err,
			"Search file URI is malformed. Check out help for more "+
				"information about search patterns.",
		)
	}

	request := smartling.FilesListRequest{}

	files, err := client.ListAllFiles(project, request)
	if err != nil {
		return nil, hierr.Errorf(
			err,
			`unable to list files in project "%s"`,
			project,
		)
	}

	result := []smartling.File{}

	for _, file := range files {
		if pattern.Match(file.FileURI) {
			result = append(result, file)
		}
	}

	return result, nil
}

func globFilesLocally(directory string, mask string) ([]string, error) {
	pattern, err := glob.Compile(mask, '/')
	if err != nil {
		return nil, NewError(
			err,
			"Search file pattern is malformed. Check out help for more "+
				"information about search patterns.",
		)
	}

	var result []string

	err = filepath.Walk(
		directory,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if pattern.Match(path) {
				result = append(result, path)
			}

			return nil
		},
	)

	if err != nil {
		return nil, hierr.Errorf(
			err,
			`unable to walk down files in dir "%s"`,
			directory,
		)
	}

	return result, nil
}
