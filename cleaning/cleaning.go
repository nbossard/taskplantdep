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
