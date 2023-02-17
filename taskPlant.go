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

	var dependencies []string

	for _, task := range tasks {
		if len(task.Depends) > 0 {
			for _, dep := range task.Depends {
				dependencies = append(dependencies, fmt.Sprintf("%s --> %s : %s", dep, task.UUID, task.Description))
			}
		}
	}

	puml := fmt.Sprintf("@startuml\n\n%s\n\n@enduml", strings.Join(dependencies, "\n"))
	err = ioutil.WriteFile("dependencies.puml", []byte(puml), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Dependencies written to dependencies.puml")
}
