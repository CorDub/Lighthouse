package main

import (
	"strings"
	"log"
	"fmt"
)


func parseCreate(lines []string) ([]string, error) {
	// check what is being created
	topLineSplit := strings.Split(lines[0], " ")
	object := topLineSplit[1]
	tableName := topLineSplit[2]

	if object == "TABLE" {
		// parsing from after the title and skipping the ); line
		parsedLines, err := parseCreateTable(lines[1:len(lines)-1], tableName)
		if err != nil {
			log.Printf("Couldn't parse this create table migration: %s", err)
			return []string{}, err
		}

		return parsedLines, nil

	} else {
		return []string{}, fmt.Errorf("Object to be parsed not recognized")
	}
}


func parseCreateTable(lines []string, name string) ([]string, error) {
	// starting the res array with the name of the table capitalized
	parsedLines := []string{strings.ToUpper(name)}

	for _, line := range lines {
		finalLine := "\t"
		// remove final comma
		trimmedLine := strings.Trim(line, ",")

		lineSplit := strings.Split(trimmedLine, " ")
		extraInfo := strings.Join(lineSplit, " ")
		finalLine += extraInfo
		parsedLines = append(parsedLines, finalLine)
	}

	return parsedLines, nil
}