package main

import (
	"strings"
	"log"
	"fmt"
)

func parseAlter(lines []string, results [][]string) ([][]string, error) {
	//check what is being created
	cleanTopLine := cleanLine(lines[0])
	topLineSplit := strings.Split(cleanTopLine, " ")
	object := topLineSplit[1]
	tableName := topLineSplit[2]

	if object == "TABLE" {
		// parsing from after the title
		parsedLines, err := parseAlterTable(lines[1:], results, tableName)
		if err != nil {
			log.Printf("Couldn't parse this alter table migration: %s", err)
			return [][]string{}, err
		}

		// create modified results
		alteredResults, err := alterModel(results, tableName, parsedLines)
		if err != nil {
			log.Fatalf("Could not create altered results: \n%s\n", err)
		}

		return alteredResults, nil
	} else {
		return [][]string{}, fmt.Errorf("Object to be parsed not recognized")
	} 
}


func parseAlterTable(lines []string, results [][]string, name string) ([]string, error) {
	model, err := getModelToAlter(results, name)
	if err != nil {
		log.Fatalf("Could not find a model to alter: %s", err)
	}

	trimmedFirstLine := strings.TrimSpace(lines[0])
	alteration := strings.Split(trimmedFirstLine, " ")
	preciseAlteration := alteration[0] + " " + alteration[1]

	for _, line := range lines {
		switch preciseAlteration {
			case "ADD COLUMN":
				model = addColumn(line, model)

			case "ALTER COLUMN":
				model, err = updateColumn(line, model)
				if err != nil {
					log.Fatalf("Could not alter column")
				}
		}
	}

	return model, nil
}


func getModelToAlter(results [][]string, name string) ([]string, error) {
	for _, model := range results {
		if strings.ToLower(model[0]) == name {
			return model, nil
		}
	}
	return []string{}, fmt.Errorf("Could not find a match for %s within current results", name)
}


func alterModel(results [][]string, name string, alteredModel []string) ([][]string, error) {
	alteredResults := [][]string{}

	altered := false
	for _, model := range results {
		if strings.ToLower(model[0]) == name {
			alteredResults = append(alteredResults, alteredModel)
			altered = true
		} else {
			alteredResults = append(alteredResults, model)
		}
	}

	if altered {
		return alteredResults, nil
	}

	return [][]string{}, fmt.Errorf("Could not find a match for %s within current results", name)
}


func addColumn(line string, model []string) ([]string) {
	lineSplit := strings.Split(line, " ")
	lineToAdd := strings.Join(lineSplit[2:], " ")
	tabbed := "\t\t" + lineToAdd
	updatedModel := append(model, tabbed)
	return updatedModel
}


func updateColumn(line string, model []string) ([]string, error) {
	cleanedLine := cleanLine(line)
	lineSplit := strings.Split(cleanedLine, " ")
	target := lineSplit[2]
	lineToAlter, err := getLineToAlter(model, target)
	if err != nil {
		log.Fatalf("Couldn't update the column: %s", err)
	}

	action := lineSplit[3]
	keywords := lineSplit[4:]

	switch action {
		case "DROP":
			alteredLine := dropKeyword(lineToAlter, keywords)
			updatedModel, err := alterLine(alteredLine, target, model)
			if err != nil {
				log.Fatalf("Could not update the column: %s", err)
			}
			return updatedModel, nil
		
		case "TYPE":
			alteredLine := changeType(lineToAlter, keywords[0])
			updatedModel, err := alterLine(alteredLine, target, model)
			if err != nil {
				log.Fatalf("Could not update the column: %s", err)
			}
			return updatedModel, nil
	}
	
	return []string{}, fmt.Errorf("Could not update the column")
}


func getLineToAlter(model []string, target string) (string, error) {
	for _, line := range model[1:] {
		trimmedLine := strings.TrimSpace(line)
		lineSplit := strings.Split(trimmedLine, " ")
		if lineSplit[0] == target {
			return line, nil
		}
	}
	return "", fmt.Errorf("Target wasn't encountered in the model")
}


func dropKeyword(lineToAlter string, keywords []string) string {
	alteredLine := ""

	cleanLineToAlter := cleanLine(lineToAlter)
	lineToAlterSplit := strings.Split(cleanLineToAlter, " ")
	defaut := false
	for _, word := range lineToAlterSplit {
		matched := false
		if defaut {
			defaut = false
			continue
		}
		
		for _, keyword := range keywords {
			if word == keyword {
				matched = true
				if word == "DEFAULT" {
					defaut = true
				}
			}
		}
		if matched {
			continue
		}

		alteredLine += word + " "
	}

	return strings.TrimSpace(alteredLine)
}


func changeType(lineToAlter, newType string) string {
	cleanLineToAlter := cleanLine(lineToAlter)
	lineToAlterSplit := strings.Split(cleanLineToAlter, " ")
	lineToAlterSplit[1] = newType
	reformed := strings.Join(lineToAlterSplit, " ")
	return reformed
}


func alterLine(alteredLine, target string, model []string) ([]string, error) {
	updatedModel := []string{}
	updated := false

	for _, line := range model {
		cleanedLine := cleanLine(line)
		lineSplit := strings.Split(cleanedLine, " ")
		if lineSplit[0] == target {
			updatedModel = append(updatedModel, "\t\t" + alteredLine)
			updated = true
		} else {
			updatedModel = append(updatedModel, line)
		}
	}

	if updated {
		return updatedModel, nil
	}

	return []string{}, fmt.Errorf("Could not find the line to alter")
}
