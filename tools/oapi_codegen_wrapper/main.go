package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config map[string]interface{}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: oapi-gen <config.yaml> <openapi.yaml>")
		os.Exit(1)
	}

	configPath := os.Args[1]
	openapiPath := os.Args[2]

	configDir := filepath.Dir(configPath)

	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var configs map[string]Config
	if err := yaml.Unmarshal(configData, &configs); err != nil {
		log.Fatalf("Error parsing config YAML: %v", err)
	}

	tempDir, err := os.MkdirTemp("", "oapi-configs")
	if err != nil {
		log.Fatalf("Error creating temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Process each top-level prop of config.yaml (module1_server, module1_models, etc)
	for name, config := range configs {
		if outputPath, ok := config["output"].(string); ok {
			if !filepath.IsAbs(outputPath) {
				outputPath = strings.TrimPrefix(outputPath, "./")
				absoluteOutputPath := filepath.Join(configDir, outputPath)
				config["output"] = absoluteOutputPath
				fmt.Printf("Output path for %s: %s\n", name, absoluteOutputPath)
			}
		}

		if outputPath, ok := config["output"].(string); ok {
			outputDir := filepath.Dir(outputPath)
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				log.Fatalf("Error creating output directory %s: %v", outputDir, err)
			}
		}

		tempConfigPath := filepath.Join(tempDir, name+".yaml")

		configBytes, err := yaml.Marshal(config)
		if err != nil {
			log.Fatalf("Error marshaling config for %s: %v", name, err)
		}

		if err := os.WriteFile(tempConfigPath, configBytes, 0644); err != nil {
			log.Fatalf("Error writing temp config for %s: %v", name, err)
		}

		cmd := exec.Command("go", "tool", "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen",
			"-config", tempConfigPath, openapiPath)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Printf("Generating code for %s...\n", name)
		if err := cmd.Run(); err != nil {
			log.Fatalf("Error running oapi-codegen for %s: %v", name, err)
		}
		fmt.Printf("Successfully generated code for %s\n\n", name)
	}

	fmt.Println("All code generation completed successfully!")
}
