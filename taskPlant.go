package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"taskplantdep/cleaning"
	"taskplantdep/model"
	"taskplantdep/taskwarrior"
)

func main() {

	// retrieve parameters from cammand line call
	// search for parameter "filter"
	defaultFilter := "status:pending"
	defaultFilename := "dependencies.puml"
	filter, outputFilename, dispStatus := parseParameters(defaultFilter, defaultFilename)

	// retrieve all tasks from taskwarrior
	filteredTasks := taskwarrior.RetrieveTasks(filter)
	allTasks := taskwarrior.RetrieveTasks("")
	fmt.Printf("Found %d filtered tasks with filter \"%s\"\n", len(filteredTasks), filter)
	fmt.Printf("Found %d tasks in total\n", len(allTasks))

	var depsLines []string
	// objectLines is a map uuid -> object line
	objectPlantUML := make(map[string]string)
	neededObjPlantUML := make(map[string]string)

	filteredTasks = cleaning.CleanDescriptions(filteredTasks)
	filteredTasks = cleaning.CleanDepends(filteredTasks)

	// generate object lines for all tasks
	// e.g : object "fin formation" as 09a3937e99a540cba226fa0fa59399e4
	generateObjectLines(allTasks, objectPlantUML, dispStatus)

	// generate dependency lines
	depsLines = generateDependencyLines(filteredTasks, objectPlantUML, neededObjPlantUML)

	// generate object lines for filtered tasks
	var objectLines []string
	for _, task := range neededObjPlantUML {
		objectLines = append(objectLines, task)
	}
	fmt.Printf("Found %d objects concerned by dependencies\n", len(objectLines))

	puml := fmt.Sprintf("@startuml\n\n%s\n\n%s\n\n@enduml", strings.Join(objectLines, "\n"), strings.Join(depsLines, "\n"))

	err := ioutil.WriteFile(outputFilename, []byte(puml), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Finished writing to \"./" + outputFilename + "\"")
	fmt.Println("FINISHED !")
}

// parseParameters parses the parameters from the command line.
func parseParameters(parDefaultFilter string, parDefaultFilename string) (string, string, bool) {
	filterPtr := flag.String("filter", parDefaultFilter, "filter to apply to taskwarrior")
	filterOutputPtr := flag.String("output", parDefaultFilename, "output file")
	filterDispStatus := flag.Bool("dispstatus", false, "display status of tasks")
	flag.Parse()
	filter := *filterPtr
	outputFilename := *filterOutputPtr
	dispStatus := *filterDispStatus
	return filter, outputFilename, dispStatus
}

func generateObjectLines(allTasks []model.ExportedTask, objectPlantUML map[string]string, dispStatus bool) {
	for _, task := range allTasks {
		objectPlantUML[task.UUID] = fmt.Sprintf(
			"object \"%d: %s\" as %s #%s",
			task.ID,
			task.GetDescriptionCleaned(),
			task.GetUUIDCleaned(),
			task.GetColorUrgency(),
		)
		// display urgency with only two digits
		objectPlantUML[task.UUID] += fmt.Sprintf("\n%s : urgency = %.2f", task.GetUUIDCleaned(), task.Urgency)
		// display project if not empty
		if task.Project != "" {
			objectPlantUML[task.UUID] += fmt.Sprintf("\n%s : project = %s", task.GetUUIDCleaned(), task.Project)
		}
		if task.Due != "" {
			objectPlantUML[task.UUID] += fmt.Sprintf("\n%s : due = %s", task.GetUUIDCleaned(), task.Due)
		}
		if dispStatus {
			objectPlantUML[task.UUID] += fmt.Sprintf("\n%s : status = %s", task.GetUUIDCleaned(), task.Status)
		}
	}
}

func generateDependencyLines(
	filteredTasks []model.ExportedTask,
	objectPlantUML map[string]string,
	neededObjPlantUML map[string]string,
) (depsLines []string) {
	nbrDeps := 0
	for _, task := range filteredTasks {
		if len(task.Depends) > 0 {
			for _, dep := range task.Depends {
				depsLines = append(depsLines, fmt.Sprintf("%s <-- %s", model.MakeOneUUIDCompatible(dep), task.GetUUIDCleaned()))
				nbrDeps++
				neededObjPlantUML[dep] = objectPlantUML[dep]
				neededObjPlantUML[task.UUID] = objectPlantUML[task.UUID]
			}
		}
	}
	fmt.Printf("Found %d dependencies\n", nbrDeps)
	return depsLines
}
