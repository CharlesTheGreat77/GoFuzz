package utils

import (
	"bufio"
	"os"
)

func ReadFile(path string) ([]string, error) {
	var contents []string
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		contents = append(contents, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return contents, nil
}
