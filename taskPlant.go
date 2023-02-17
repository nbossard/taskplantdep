package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("Exporting task dependencies...")

	// Run task export command
	cmd := exec.Command("task", "export", "status:pending")
	output, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error running task export command:", err)
		return
	}
	defer output.Close()

	// Start command and scan output
	scanner := bufio.NewScanner(output)
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting task export command:", err)
		return
	}

	// Parse tasks to generate PlantUML code
	var tasks []string
	var dependencies []string
	for scanner.Scan() {
		task := scanner.Text()
		if strings.HasPrefix(task, "depends: ") {
			dependencies = append(dependencies, task)
		} else {
			tasks = append(tasks, task)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning task export output:", err)
		return
	}

	var lines []string
	lines = append(lines, "@startuml")
	lines = append(lines, "skinparam monochrome true")
	lines = append(lines, "skinparam defaultFontName Courier")
	lines = append(lines, "skinparam ArrowFontName Courier")

	// Add tasks to PlantUML code
	for _, task := range tasks {
		parts := strings.SplitN(task, "description:", 2)
		if len(parts) != 2 {
			continue
		}
		id := strings.TrimSpace(parts[0])
		desc := strings.TrimSpace(parts[1])
		lines = append(lines, fmt.Sprintf("task %s { %s }", id, desc))
	}

	// Add dependencies to PlantUML code
	for _, dep := range dependencies {
		parts := strings.SplitN(dep, "depends: ", 2)
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
