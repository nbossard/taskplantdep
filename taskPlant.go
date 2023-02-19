package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"taskplantdep/cleaning"
	"taskplantdep/model"
)

func main() {

	filter := "-PERSO status:pending"
	filteredTasks := retrieveTasks(filter)
	allTasks := retrieveTasks("")
	fmt.Printf("Found %d filtered tasks\n", len(filteredTasks))
	fmt.Printf("Found %d all tasks\n", len(allTasks))

	var depsLines []string
	// objectLines is a map uuid -> object line
	objectPlantUML := make(map[string]string)
	neededObjPlantUML := make(map[string]string)

	filteredTasks = cleaning.CleanDescriptions(filteredTasks)
	filteredTasks = cleaning.CleanDepends(filteredTasks)

	// generate object lines for all tasks
	// e.g : object "fin formation" as 09a3937e99a540cba226fa0fa59399e4
	for _, task := range allTasks {
		objectPlantUML[task.UUID] = fmt.Sprintf("object \"%d: %s\" as %s #%s", task.ID, task.Description, task.GetUUIDCleaned(), task.GetColorUrgency())
		// display urgency with only two digits
		objectPlantUML[task.UUID] += fmt.Sprintf("\n%s : urgency = %.2f", task.GetUUIDCleaned(), task.Urgency)
		// display project if not empty
		if task.Project != "" {
			objectPlantUML[task.UUID] += fmt.Sprintf("\n%s : project = %s", task.GetUUIDCleaned(), task.Project)
		}
		if task.Due != "" {
			objectPlantUML[task.UUID] += fmt.Sprintf("\n%s : due = %s", task.GetUUIDCleaned(), task.Due)
		}
	}

	// generate dependency lines
	for _, task := range filteredTasks {
		if len(task.Depends) > 0 {
			for _, dep := range task.Depends {
				depsLines = append(depsLines, fmt.Sprintf("%s <-- %s", model.MakeOneUUIDCompatible(dep), task.GetUUIDCleaned()))
				neededObjPlantUML[dep] = objectPlantUML[dep]
				neededObjPlantUML[task.UUID] = objectPlantUML[task.UUID]
			}
		}
	}

	// generate object lines for filtered tasks
	var objectLines []string
	for _, task := range neededObjPlantUML {
		objectLines = append(objectLines, task)
	}

	puml := fmt.Sprintf("@startuml\n\n%s\n\n%s\n\n@enduml", strings.Join(objectLines, "\n"), strings.Join(depsLines, "\n"))

	err := ioutil.WriteFile("dependencies.puml", []byte(puml), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Dependencies written to dependencies.puml")
}

func retrieveTasks(filter string) []model.ExportedTask {
	var cmd *exec.Cmd
	if filter == "" {
		cmd = exec.Command("task", "export")
	} else {
		cmd = exec.Command("task", filter, "export")
	}

	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var tasks []model.ExportedTask

	err = json.Unmarshal(output, &tasks)
	if err != nil {
		log.Fatal(err)
	}
	return tasks
}
