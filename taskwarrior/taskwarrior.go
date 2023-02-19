package taskwarrior

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"taskplantdep/model"
)

func RetrieveTasks(parFilter string) []model.ExportedTask {
	var cmd *exec.Cmd
	if parFilter == "" {
		cmd = exec.Command("task", "export")
	} else {
		// parFilter can contain multiple filters separated by space
		// "export" must be the last argument
		if strings.Contains(parFilter, " ") {
			parFilters := strings.Split(parFilter, " ")
			cmd = exec.Command("task", parFilters[0], parFilters[1], "export")
		} else {
			cmd = exec.Command("task", parFilter, "export")
		}
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
