package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func loadEnv(fileName string) error {
	file, err := os.Open(fileName)

	if err != nil {
		return fmt.Errorf("fail to load environment: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)

		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		os.Setenv(key, val)
	}

	return scanner.Err()
}
