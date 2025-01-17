package main

import (
	"bufio"
	"os"
	"strings"
)

// Specific funtions to deal with world.mt

func worldMtGetActiveModDirs() ([]string, error) {
	result := []string{}
	file, err := os.Open(cfg.worldMtFile)
	if err != nil {
		return result, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "load_mod") && !strings.Contains(line, "= false") {
			value := extractValue(line)
			if value != "" {
				result = append(result, value)
			}
		}
	}
	return result, nil
}

func worldMtGetOnlyOptions() string { //NOT slice!
	//TODO: Handle errors
	result := ""
	file, err := os.Open(cfg.worldMtFile)
	if err != nil {
		return result
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "load_mod") {
			result = result + line + "\n"
		}
	}

	return result
}

// Extract value from `key = value` line
func extractValue(line string) string {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}
