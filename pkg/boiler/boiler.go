package boiler

import (
	"bufio"
	"io/fs"
)

func ReadJsonFromDisk(fs fs.FS, filename string) (string, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var jsonString string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		jsonString += scanner.Text()
	}

	return jsonString, nil
}
