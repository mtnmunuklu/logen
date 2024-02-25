package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mtnmunuklu/logen/sigma"
	"github.com/mtnmunuklu/logen/sigma/sevaluator"
)

var (
	filePath      string
	configPath    string
	fileContent   string
	configContent string
	showHelp      bool
	outputPath    string
	version       bool
	caseSensitive bool
)

// Set up the command-line flags
func init() {
	flag.StringVar(&filePath, "filepath", "", "Name or path of the file or directory to read")
	flag.StringVar(&configPath, "config", "", "Path to the configuration file")
	flag.StringVar(&fileContent, "filecontent", "", "Base64-encoded content of the file or directory to read")
	flag.StringVar(&configContent, "configcontent", "", "Base64-encoded content of the configuration file")
	flag.BoolVar(&showHelp, "help", false, "Show usage")
	flag.StringVar(&outputPath, "output", "", "Output directory for writing files")
	flag.BoolVar(&version, "version", false, "Show version information")
	flag.BoolVar(&caseSensitive, "cs", false, "Case sensitive mode")
	flag.Parse()

	// If the version flag is provided, print version information and exit
	if version {
		fmt.Println("Logen version 1.0.0")
		os.Exit(1)
	}

	// If the help flag is provided, print usage information and exit
	if showHelp {
		printUsage()
		os.Exit(1)
	}

	// Check if filepath and configpath are provided as command-line arguments
	if flag.NArg() > 0 {
		filePath = flag.Arg(0)
	}
	if flag.NArg() > 1 {
		configPath = flag.Arg(1)
	}

	// Check if both filecontent and configcontent are provided
	if (filePath == "" && fileContent == "") || (configPath == "" && configContent == "") {
		fmt.Println("Please provide either file paths or file contents, and either config path or config content.")
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: logen -filepath <path> -config <path> [flags]")
	fmt.Println("Flags:")
	flag.PrintDefaults()
	fmt.Println("Example:")
	fmt.Println("  logen -filepath /path/to/file -config /path/to/config")
}

func main() {
	// Read the contents of the file(s) specified by the filepath flag or filecontent flag
	fileContents := make(map[string][]byte)
	var err error

	// Check if file paths are provided
	if filePath != "" {
		// Check if the filepath is a directory
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("Error getting file/directory info:", err)
			return
		}

		if fileInfo.IsDir() {
			// filePath is a directory, so walk the directory to read all the files inside it
			filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					fmt.Println("Error accessing file:", err)
					return nil
				}
				if !info.IsDir() {
					// read file content
					content, err := os.ReadFile(path)
					if err != nil {
						fmt.Println("Error reading file:", err)
						return nil
					}
					fileContents[path] = content
				}
				return nil
			})
		} else {
			// filePath is a file, so read its contents
			fileContents[filePath], err = os.ReadFile(filePath)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}
		}
	} else if fileContent != "" {
		// Check if the filecontent is a directory
		lines := strings.Split(fileContent, "\n")
		if len(lines) > 1 {
			// fileContent is a directory, so read all lines as separate files
			for _, line := range lines {
				// decode base64 content
				decodedContent, err := base64.StdEncoding.DecodeString(line)
				if err != nil {
					fmt.Println("Error decoding base64 content:", err)
					return
				}
				fileContents[line] = decodedContent
			}
		} else {
			// fileContent is a file, so read its content
			// decode base64 content
			decodedContent, err := base64.StdEncoding.DecodeString(fileContent)
			if err != nil {
				fmt.Println("Error decoding base64 content:", err)
				return
			}
			fileContents["filecontent"] = decodedContent
		}
	}

	// Read the contents of the configuration file or use configcontent
	var configContents []byte
	if configPath != "" {
		configContents, err = os.ReadFile(configPath)
		if err != nil {
			fmt.Println("Error reading configuration file:", err)
			return
		}
	} else if configContent != "" {
		// decode base64 content
		decodedContent, err := base64.StdEncoding.DecodeString(configContent)
		if err != nil {
			fmt.Println("Error decoding base64 content:", err)
			return
		}
		configContents = decodedContent
	}

	// Loop over each file and parse its contents as a Sigma rule
	for _, fileContent := range fileContents {
		sigmaRule, err := sigma.ParseRule(fileContent)
		if err != nil {
			fmt.Println("Error parsing rule:", err)
			continue
		}

		// Parse the configuration file as a Sigma config
		config, err := sigma.ParseConfig(configContents)
		if err != nil {
			fmt.Println("Error parsing config:", err)
			continue
		}

		var sr *sevaluator.RuleEvaluator

		if caseSensitive {
			// Evaluate the Sigma rule against the config using case sensitive mode
			sr = sevaluator.ForRule(sigmaRule, sevaluator.WithConfig(config), sevaluator.CaseSensitive)
		} else {
			// Evaluate the Sigma rule against the config
			sr = sevaluator.ForRule(sigmaRule, sevaluator.WithConfig(config))
		}

		ctx := context.Background()
		result, err := sr.Alters(ctx)
		if err != nil {
			fmt.Println("Error converting rule:", err)
			continue
		}

		var output string

		// Print the results of the query
		var builder strings.Builder
		for _, queryResult := range result.QueryResults {
			builder.WriteString(queryResult + "\n")
		}
		output = builder.String()

		// Check if outputPath is provided
		if outputPath != "" {
			// Create the output file path using the Name field from the rule
			outputFilePath := filepath.Join(outputPath, fmt.Sprintf("%s.log", sigmaRule.Title))

			// Write the output string to the output file
			err := os.WriteFile(outputFilePath, []byte(output), 0644)
			if err != nil {
				fmt.Println("Error writing output to file:", err)
				continue
			}

			fmt.Printf("Output for rule '%s' written to file: %s\n", sigmaRule.Title, outputFilePath)
		} else {
			fmt.Printf("%s", output)
		}
	}
}
