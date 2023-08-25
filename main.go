package main

import (
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var err error

	var startPath string
	if len(os.Args) == 2 {
		startPath, err = filepath.Abs(os.Args[1])
	} else {
		startPath, err = os.Getwd()
	}
	if err != nil {
		log.Fatal(err)
	}

	m := &model{
		path: startPath,
	}

	err = m.list()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tea.NewProgram(m, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
