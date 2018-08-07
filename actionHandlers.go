package main

import (
	"os"
	"fmt"
	"path/filepath"
)

func findActionHandler(kind string) func(chan<- bool, []Dep, Perform) int {
	switch kind {
	case "git":
		return gitActionHandler
	case "secret":
		return secretesActionHandler
	case "docker-compose":
		return dockerComposeHandler
	default:
		return defaultActionHandler
	}
}


/// *** Action Handlers *** ///

func defaultActionHandler(complete chan<- bool, deps []Dep, perform Perform) int {
	n := 0
	for _, dep := range deps {
		n++
		go defaultAction(complete, dep, perform)
	}
	return n
}

func dockerComposeHandler(complete chan<- bool, deps []Dep, perform Perform) int {
	fmt.Println("Called dockerComposeHandler")
	override := ""

	action := perform.Action[0]
	perform.Action = perform.Action[1:]

	for _, dep := range deps {
		override += " -f " + filepath.Clean(prepDockerComposeAction(dep, perform))
	}

	//go dockerComposeAction(complete, perform, override)
	dockerComposeAction(complete, perform, action, override)

	return 0
}

func gitActionHandler(complete chan<- bool, deps []Dep, perform Perform) int {
	n := 0
	for _, dep := range deps {
		n++
		switch perform.Action[0] {
		case "clone": // Clone breaks if the parent dirs aren't already there.
			there, _ := exists(dep.Location)
			if !there {
				if perform.DryRun {
					fmt.Println("DryRun, would have performed a: `mkdir -p ", dep.Location, "`")
				} else {
					// mkdir -p <location>
					os.MkdirAll(dep.Location, os.ModePerm)
				}
			}

			go defaultAction(complete, dep, perform)

		default:
			go defaultAction(complete, dep, perform)
		}
	}
	return n
}


// TODO: make a special handler for secrets
func secretesActionHandler(complete chan<- bool, deps []Dep, perform Perform) int {
	n := 0
	for _, dep := range deps {
		n++
		switch perform.Action[0] {
		//case "get": // Or something along these lines
		// Right now, everything falls through to default.
		default:
			go defaultAction(complete, dep, perform)
		}
	}
	return n
}


/// ***  Helpers *** ///

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}