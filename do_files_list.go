package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Smartling/api-sdk-go"
	hierr "github.com/reconquest/hierr-go"
)

func doFilesList(
	client *smartling.Client,
	config Config,
	args map[string]interface{},
) error {
	var (
		short   = args["--short"].(bool)
		project = args["--project"].(string)
		uri, _  = args["<uri>"].(string)
	)

	if args["--format"] == nil {
		args["--format"] = defaultFilesListFormat
	}

	format, err := compileFormat(args["--format"].(string))
	if err != nil {
		return err
	}

	files, err := globFilesRemote(client, project, uri)
	if err != nil {
		return err
	}

	table := NewTableWriter(os.Stdout)

	for _, file := range files {
		if short {
			fmt.Fprintf(table, "%s\n", file.FileURI)
		} else {
			row, err := format.Execute(file)
			if err != nil {
				return err
			}

			_, err = io.WriteString(table, row)
			if err != nil {
				return hierr.Errorf(
					err,
					"unable to write row to output table",
				)
			}
		}
	}

	err = RenderTable(table)
	if err != nil {
		return err
	}

	return nil
}
