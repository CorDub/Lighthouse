package main

import (
	"fmt"
	"strings"
)

func printV[T any](name string, variable T) {
	switch v := any(variable).(type) {
	case string:
		fmt.Printf("%s \n%q\n\n", name, v)
	default:
		fmt.Printf("%s \n%v\n\n", name, v)
	}
}

func cleanLine(line string) string {
	trimmed := strings.TrimSpace(line)
	commasRemoved := strings.ReplaceAll(trimmed, ",", "")
	semiColonRemoved := strings.ReplaceAll(commasRemoved, ";", "")
	return semiColonRemoved
}