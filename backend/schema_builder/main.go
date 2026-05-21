package main

import (
	"os"
	"log"
	"strings"
)

// Cathedral

func main() {
	//get into the correct directory
	dirPath := "../sql/schema"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("Couldn't open the directory: %s", err)
	}

	results := [][]string{}

	for _, file := range files {
		//read file
		fileName := file.Name()
		fullFileName := dirPath + "/" + fileName
		byteContent, err := os.ReadFile(fullFileName) 
		if err != nil {
			log.Fatalf("Could not read the file %s: %s", fullFileName, err)
		}

		//get content as a string
		stringContent := string(byteContent)
		trimmedStringContent := strings.TrimPrefix(stringContent, "-- +goose Up\n")

		linesContent := strings.Split(trimmedStringContent, "--")

		// split by semi colon for each change
		trimmedChanges := strings.TrimSpace(linesContent[0])
		changes := strings.Split(trimmedChanges, ";")

		// change loop
		for _, change := range changes {
			if len(change) == 0 {
				continue
			}

			//split it by line
			trimmed := cleanLine(change)
			lines := strings.Split(trimmed, "\n")

			//checks first line to decide method of parsing
			firstLineSplit := strings.Split(lines[0], " ")

			switch firstLineSplit[0] {
				case "CREATE":
					textLines, err := parseCreate(lines)
					if err != nil {
						log.Fatalf("Couldn't parse the lines of this create statement: %s", err)
					}
					results = append(results, textLines)

				case "ALTER":
					alteredResults, err := parseAlter(lines, results)
					if err != nil {
						log.Fatalf("Couldn't parse the lines of this alter statement: %s", err)
					}
					results = alteredResults
			}
		}
	}

	var finalString strings.Builder
	for _, model := range results {
		innerJoin := strings.Join(model, "\n")
		finalString .WriteString(innerJoin + "\n\n")
	}

	errWriting := os.WriteFile("../sql/schema.txt", []byte(finalString.String()), 0644)
	if errWriting != nil {
		log.Fatalf("Error while writing to file: %s", errWriting)
	}
}