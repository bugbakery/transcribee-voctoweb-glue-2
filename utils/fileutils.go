package utils

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func ReadCsv(app core.App, file string, delimiter rune) (data [][]string, header []string, err error) {
	data = [][]string{}

	// initialize the filesystem
	fsys, err := app.NewFilesystem()
	if err != nil {
		return [][]string{}, []string{}, err
	}
	defer fsys.Close()

	// retrieve a file reader for the avatar key
	r, err := fsys.GetFile(file)
	if err != nil {
		return [][]string{}, []string{}, err
	}
	defer r.Close()
	reader := csv.NewReader(r)
	reader.Comma = delimiter // Set the custom delimiter

	// Read the header first
	header, err = reader.Read()
	if err != nil {
		return [][]string{}, []string{}, err
	}
	header[0] = strings.Trim(header[0], "\ufeff") // Trim leading whitespace (e.g., BOM)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break // End of file, exit the loop
		}
		if err != nil {
			return [][]string{}, []string{}, err
		}

		data = append(data, line)
	}

	return data, header, nil
}

func WriteCsv(app core.App, fileName string, data [][]string) (file *filesystem.File, err error) {
	// Initialize the filesystem
	fsys, err := app.NewFilesystem()
	if err != nil {
		return nil, err
	}
	defer fsys.Close()

	// Create a buffer to hold CSV data
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// Write each record, calls flush internally
	err = writer.WriteAll(data)
	if err != nil {
		return nil, err
	}

	f, err := filesystem.NewFileFromBytes(buf.Bytes(), fileName)
	if err != nil {
		return nil, err
	}

	return f, nil
}