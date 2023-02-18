package main

import (
	"testing"
)

func TestKeepOnlyTasksWithoutTag(t *testing.T) {
	var tasks []Task
	tasks = append(tasks, Task{Description: "task 1", Tags: []string{"PERSO"}})

	filteredTasks := keepOnlyTasksWithoutTag(tasks, "PERSO")

	if len(filteredTasks) != 0 {
		t.Error("testKeepOnlyTasksWithoutTag failed on one tag PERSO")
	}

	tasks = tasks[:0]
	tasks = append(tasks, Task{Description: "task 1", Tags: []string{"PERSO IKEA"}})

	filteredTasks = keepOnlyTasksWithoutTag(tasks, "PERSO")

	if len(filteredTasks) != 0 {
		t.Error("testKeepOnlyTasksWithoutTag failed on one tag two words PERSO IKEA")
	}

	tasks = tasks[:0]
	tasks = append(tasks, Task{Description: "task 1", Tags: []string{"IKEA"}})

	filteredTasks = keepOnlyTasksWithoutTag(tasks, "PERSO")

	if len(filteredTasks) != 1 {
		t.Error("testKeepOnlyTasksWithoutTag failed on one tag IKEA")
	}

	tasks = tasks[:0]
	tasks = append(tasks, Task{Description: "task 1", Tags: []string{"IKEA", "PERSO"}})

	filteredTasks = keepOnlyTasksWithoutTag(tasks, "PERSO")

	if len(filteredTasks) != 0 {
		t.Error("testKeepOnlyTasksWithoutTag failed on two tags IKEA PERSO")
	}

	tasks = tasks[:0]
	tasks = append(tasks, Task{Description: "task 1", Tags: []string{}})

	filteredTasks = keepOnlyTasksWithoutTag(tasks, "PERSO")

	if len(filteredTasks) != 1 {
		t.Error("testKeepOnlyTasksWithoutTag failed on no tags")
	}
}
