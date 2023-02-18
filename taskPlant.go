package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

type Task struct {
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

	var tasks []Task
	err = json.Unmarshal(output, &tasks)
	if err != nil {
		log.Fatal(err)
	}

	var outputLines []string

	tasks = keepOnlyTasksOfStatus(tasks, "pending")
	tasks = keepOnlyTasksWithoutTag(tasks, "PERSO")

	tasks = makeUUIDsObjectCompatible(tasks)

	tasks = cleanDepends(tasks)

	// write object lines
	// e.g : object "fin formation début" as 09a3937e99a540cba226fa0fa59399e4
	for _, task := range tasks {
		outputLines = append(outputLines, fmt.Sprintf("object \"%s\" as %s", task.Description, task.UUID))
	}

	// write dependency lines
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
func cleanDepends(tasks []Task) []Task {
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
func makeUUIDsObjectCompatible(tasks []Task) []Task {
	for i, task := range tasks {
		tasks[i].UUID = strings.Replace(task.UUID, "-", "", -1)
		tasks[i].Description = strings.Replace(task.Description, "\r", " ", -1)
		tasks[i].Description = strings.Replace(tasks[i].Description, "\n", " ", -1)
		tasks[i].Description = strings.Replace(tasks[i].Description, "\"", "'", -1)
		for j, dep := range task.Depends {
			tasks[i].Depends[j] = strings.Replace(dep, "-", "", -1)
		}
	}
	return tasks
}

// keepOnlyTasksOfStatus returns only tasks of the given status
func keepOnlyTasksOfStatus(tasks []Task, parStatus string) []Task {
	var filteredTasks []Task
	for _, task := range tasks {
		if task.Status == parStatus {
			filteredTasks = append(filteredTasks, task)
		}
	}
	return filteredTasks
}

// keepOnlyTasksWithoutTag returns only tasks without the given tag
// each task can or not have a tag
// parTag can appear as part of tags[Ø] (two or more words separated by a space)
// or be full tags[Ø] or be full tags[1]
func keepOnlyTasksWithoutTag(tasks []Task, parTag string) []Task {
	var filteredTasks []Task
	for _, task := range tasks {
		if len(task.Tags) == 0 {
			filteredTasks = append(filteredTasks, task)
		} else {
			foundTag := false
			for _, tag := range task.Tags {
				if strings.Contains(tag, parTag) {
					foundTag = true
				}
			}
			if !foundTag {
				filteredTasks = append(filteredTasks, task)
			}
		}
	}
	return filteredTasks
}
