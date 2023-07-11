package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	// Read the YAML content from a file
	filePath := "/tmp/test1.18.0/tetrate-istio.yaml"
	yamlContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read the file: %s\n", err)
		return
	}

	// Convert the YAML content to string
	text := string(yamlContent)
	// Split the text into separate YAML blocks
	blocks := strings.Split(text, "\n---\n")


	// Iterate over the blocks and check for blocks with only empty lines and comments
	for i, block := range blocks[1:] {
		trimmedBlock := strings.TrimSpace(block)
		lines := strings.Split(trimmedBlock, "\n")
		isEmptyBlock := true

		// Check if the block consists only of empty lines and lines starting with "#"
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "#") {
				isEmptyBlock = false
				break
			}
		}

		// Print the block if it consists only of empty lines and comments
		if isEmptyBlock {
			startLine := strings.Index(text, block)
			fmt.Printf("Invalid YAML structure in block %d (starting line: %d):\n", i+1, countLines(text[:startLine])+1)
			fmt.Println("---")
			fmt.Println(block)
		}
	}
}

// countLines counts the number of lines in a given text
func countLines(text string) int {
	return strings.Count(text, "\n")
}
