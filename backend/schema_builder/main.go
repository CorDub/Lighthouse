package main

import (
	"os"
	"log"
	"fmt"
	"strings"
)

// Brick functions

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


func parseAlterTable(lines []string, results [][]string, name string) ([]string, error) {
	model, err := getModelToAlter(results, name)
	if err != nil {
		log.Fatalf("Could not find a model to alter: %s", err)
	}
	//TBF
}


func getModelToAlter(results [][]string, name string) ([]string, error) {
	for _, model := range results {
		if strings.ToLower(model[0]) == name {
			return model, nil
		}

		return []string{}, fmt.Errorf("Could not find a match for %s within current results", name)
	}
} 


// Mortar functions


func parseCreate(lines []string) ([]string, error) {
	// check what is being created
	topLineSplit := strings.Split(lines[0], " ")
	object := topLineSplit[1]
	tableName := topLineSplit[2]

	if object == "table" {
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


func parseAlter(lines []string, results [][]string) ([]string, error) {
	//check what is being created
	topLineSplit := strings.Split(lines[0], " ")
	object := topLineSplit[1]
	tableName := topLineSplit[2]

	if object == "table" {
		// parsing from after the title
		parsedLines, err := parseAlterTable(lines[1:], results, tableName)
		if err != nil {
			log.Printf("Couldn't parse this alter table migration: %s", err)
			return []string{}, err
		}

		return parsedLines, nil
	} else {
		return []string{}, fmt.Errorf("Object to be parsed not recognized")
	} 
}

// Bricks and mortar

func main() {
	//get into the correct directory
	files, err := os.ReadDir("../sql/schema")
	if err != nil {
		log.Fatalf("Couldn't open the directory: %s", err)
	}

	results := [][]string{}

	for _, file := range files {
		//read file
		fileName := file.Name()
		fmt.Println(fileName)
		byteContent, err := os.ReadFile(fileName) 
		if err != nil {
			log.Fatalf("Could not read the file %s: %s", fileName, err)
		}

		//get content as a string
		stringContent := string(byteContent)
		linesContent := strings.Split(stringContent, ";")[1:]

		//split it by line
		lines := strings.Split(linesContent[0], "\n")

		//checks first line to decide method of parsing
		firstLineSplit := strings.Split(lines[0], " ")

		if firstLineSplit[0] == "CREATE" {
			textLines, err := parseCreate(lines)
			if err != nil {
				log.Fatalf("Couldn't parse the lines of this create statement")
			}
			results = append(results, textLines)

		} else if firstLineSplit[0] == "ALTER" {
			textLines, err := parseAlter(lines, results)
			if err != nil {
				log.Fatalf("Couldn't parse the lines of this create statement")
			}
			results = append(results, textLines)
		}

		
	}
}