package main

import (
	"testing"
)

func TestCleanOneDescription(t *testing.T) {
	task := "Ceci est une description"
	res := cleanOneDescription(task)
	if res != task {
		t.Error("testCleanOneDescription failed")
	}
}
