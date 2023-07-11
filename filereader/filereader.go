package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/karuppiah7890/go-jsonschema-generator"

	"gopkg.in/yaml.v2"
)

func main() {
	istioVersion := "1.18.0"
	reposPath := "/home/petr/go/src/github.com/tetrateio/PartnerEng/packager/aws/helmcharts"
	valuesFilePath := reposPath + "/" + istioVersion

	// Read values.yaml file
	values := make(map[string]interface{})
	valuesFileData, err := ioutil.ReadFile(valuesFilePath + "/values.yaml")
	if err != nil {
		fmt.Printf("Error when reading file '%s': %v", valuesFilePath, err)
		return
	}

	// Unmarshal values.yaml into a map
	err = yaml.Unmarshal(valuesFileData, &values)
	if err != nil {
		fmt.Println("Error unmarshaling YAML:", err)
		return
	}

	// Generate JSON schema document from values
	jsonDocument := generateJSONSchemaDocument(&values)

	// Convert JSON schema document to JSON bytes
	jsonBytes, err := json.Marshal(jsonDocument)
	if err != nil {
		fmt.Println("Error marshaling JSON schema:", err)
		return
	}

	// Write the original JSON schema to a file
	err = ioutil.WriteFile("/tmp/nochanges.json", jsonBytes, 0644)
	if err != nil {
		fmt.Println("Error writing original JSON schema to file:", err)
		return
	}

	// Convert JSON schema to JSON data map
	jsonData := make(map[string]interface{})
	err = json.Unmarshal(jsonBytes, &jsonData)
	if err != nil {
		fmt.Println("Error unmarshaling JSON schema:", err)
		return
	}

	// Remove sections from the JSON data
	sectionsToRemove := []string{
		"image",
		"imagePullSecrets",
		"livenessProbe",
		"readinessProbe",
		"startupProbe",
		"podDisruptionBudget",
		"serviceAccount",
		"priorityClassName",
		"podSecurityContext",
		"securityContext",
		"updateStrategy",
		"nameOverride",
		"podSecurityPolicy",
		"extraVolumeMounts",
		"extraVolumes",
	}
	removeSections(jsonData, sectionsToRemove)

	// Add the new section to jsonData
	jsonData["properties"].(map[string]interface{})["pilot"].(map[string]interface{})["properties"].(map[string]interface{})["tolerations"] = map[string]interface{}{
		"type": "array",
	}

	// Convert the modified JSON data back to JSON bytes
	modifiedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling modified JSON:", err)
		return
	}

	tmpFile := "/tmp/um.json"

	// Write the modified JSON to the temporary file
	err = ioutil.WriteFile(tmpFile, modifiedJSON, 0644)
	if err != nil {
		fmt.Println("Error writing modified JSON to file:", err)
		return
	}

	fmt.Println("Properties removed successfully!")
}

// generateJSONSchemaDocument generates a JSON schema document from the given values
func generateJSONSchemaDocument(values *map[string]interface{}) *jsonschema.Document {
	jsonDocument := &jsonschema.Document{}
	jsonDocument.ReadDeep(values)
	return jsonDocument
}

// removeSections recursively removes sections from the jsonData based on their names
func removeSections(jsonData map[string]interface{}, sectionNames []string) {
	for key, value := range jsonData {
		if nestedObj, ok := value.(map[string]interface{}); ok {
			removeSections(nestedObj, sectionNames)
		}
		if contains(sectionNames, key) {
			delete(jsonData, key)
		}
	}
}

// contains checks if a string is present in a string slice
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
