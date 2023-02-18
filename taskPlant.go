package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

type exportedTask struct {
	ID          int      `json:"id"`
	Description string   `json:"description"`
	Due         string   `json:"due"`
	End         string   `json:"end"`
	Entry       string   `json:"entry"`
	Modified    string   `json:"modified"`
	Project     string   `json:"project"`
	Status      string   `json:"status"`
	UUID        string   `json:"uuid"`
	Wait        string   `json:"wait"`
	Depends     []string `json:"depends"`
	Urgency     float64  `json:"urgency"`
	Tags        []string `json:"tags"`
}

func main() {
	cmd := exec.Command("task", "export")

	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var tasks []exportedTask

	err = json.Unmarshal(output, &tasks)
	if err != nil {
		log.Fatal(err)
	}

	var outputLines []string

	tasks = makeUUIDsObjectCompatible(tasks)
	tasks = cleanDescriptions(tasks)
	tasks = cleanDepends(tasks)

	// generate object lines
	// e.g : object "fin formation début" as 09a3937e99a540cba226fa0fa59399e4
	for _, task := range tasks {
		outputLines = append(outputLines, fmt.Sprintf("object \"%s\" as %s", task.Description, task.UUID))
	}

	// generate dependency lines
	for _, task := range tasks {
		if len(task.Depends) > 0 {
			for _, dep := range task.Depends {
				outputLines = append(outputLines, fmt.Sprintf("%s --> %s", dep, task.UUID))
			}
		}
	}

	puml := fmt.Sprintf("@startuml\n\n%s\n\n@enduml", strings.Join(outputLines, "\n"))
	err = ioutil.WriteFile("dependencies.puml", []byte(puml), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Dependencies written to dependencies.puml")
}

// cleanDepends is a bug turnaround.
// Some depends value are surrounded by [\" and \"] and some are not, remove them.
func cleanDepends(tasks []exportedTask) []exportedTask {
	for i, task := range tasks {
		for j, dep := range task.Depends {
			if strings.HasSuffix(dep, "]") {
				tasks[i].Depends[j] = strings.TrimRight(dep, "]")
			}
		}
	}
	for i, task := range tasks {
		for j, dep := range task.Depends {
			if strings.HasSuffix(dep, "\"") {
				tasks[i].Depends[j] = strings.TrimRight(dep, "\"")
			}
		}
	}
	for i, task := range tasks {
		for j, dep := range task.Depends {
			if strings.HasPrefix(dep, "[") {
				tasks[i].Depends[j] = strings.TrimLeft(dep, "[\"")
			}
		}
	}
	for i, task := range tasks {
		for j, dep := range task.Depends {
			if strings.HasPrefix(dep, "\"") {
				tasks[i].Depends[j] = strings.TrimLeft(dep, "\"")
			}
		}
	}
	return tasks
}

// makeUUIDsObjectCompatible transforms UUIDs to be compatible with PlantUML objects.
// remove all dashes from UUIDs
// remove all dashes from UUIDS in depends
// remove all carriage returns from descriptions
// replace all " with '
func makeUUIDsObjectCompatible(tasks []exportedTask) []exportedTask {
	for i, task := range tasks {
		tasks[i].UUID = strings.Replace(task.UUID, "-", "", -1)
		for j, dep := range task.Depends {
			tasks[i].Depends[j] = strings.Replace(dep, "-", "", -1)
		}
	}
	return tasks
}

// cleanDescriptions cleans all descriptions.
func cleanDescriptions(tasks []exportedTask) []exportedTask {
	for i, task := range tasks {
		tasks[i].Description = cleanOneDescription(task.Description)
	}
	return tasks
}

// cleanOneDescription cleans one description.
// Removes all carriage returns from descriptions.
// Also replaces all " with '.
func cleanOneDescription(parDescription string) string {
	parDescription = strings.Replace(parDescription, "\r", " ", -1)
	parDescription = strings.Replace(parDescription, "\n", " ", -1)
	parDescription = strings.Replace(parDescription, "\"", "'", -1)
	return parDescription
}
