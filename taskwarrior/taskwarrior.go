package taskwarrior

import (
	"encoding/json"
	"log"
	"os/exec"
	"taskplantdep/model"
)

func RetrieveTasks(filter string) []model.ExportedTask {
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
