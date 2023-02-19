package cleaning

import "strings"
import "taskplantdep/model"

// CleanDepends is a bug turnaround.
// Some depends value are surrounded by [\" and \"] and some are not, remove them.
func CleanDepends(tasks []model.ExportedTask) []model.ExportedTask {
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

// CleanDescriptions cleans all descriptions.
func CleanDescriptions(tasks []model.ExportedTask) []model.ExportedTask {
	for i, task := range tasks {
		tasks[i].Description = model.CleanOneDescription(task.Description)
	}
	return tasks
}
