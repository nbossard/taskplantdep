package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Run taskwarrior dependency command
	cmd := exec.Command("task", "dependency")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running task dependency command:", err)
		return
	}

	// Parse output to generate PlantUML code
	var lines []string
	lines = append(lines, "@startuml")
	lines = append(lines, "skinparam monochrome true")
	lines = append(lines, "skinparam defaultFontName Courier")
	lines = append(lines, "skinparam ArrowFontName Courier")

	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " depends on ")
		if len(parts) != 2 {
			continue
		}
		parent := strings.TrimSpace(parts[1])
		child := strings.TrimSpace(parts[0])
		lines = append(lines, fmt.Sprintf("%s --> %s", parent, child))
	}

	lines = append(lines, "@enduml")

	// Write PlantUML code to file
	file, err := os.Create("dependencies.puml")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	for _, line := range lines {
		fmt.Fprintln(file, line)
	}

	fmt.Println("Task dependencies exported to dependencies.puml")
}
