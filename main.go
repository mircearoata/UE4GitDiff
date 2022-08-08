package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"
)

type UProject struct {
	EngineAssociation string `json:"EngineAssociation"`
}

func getEngineAssociation(projectPath string) (string, error) {
	projectData, err := os.ReadFile(projectPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to read project file")
	}

	var uproject UProject
	if err := json.Unmarshal(projectData, &uproject); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal project file")
	}

	return uproject.EngineAssociation, nil
}

func getEngineRoot(projectPath string) (string, error) {
	engineAssociation, err := getEngineAssociation(projectPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to get engine association")
	}

	engines, err := listEngines()
	if err != nil {
		return "", errors.Wrap(err, "failed to list engines")
	}

	engineRoot, ok := engines[engineAssociation]
	if !ok {
		return "", errors.New("failed to find engine")
	}

	return engineRoot, nil
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	uproject := path.Join(cwd, "FactoryGame.uproject")

	oldFile := os.Args[1]
	newFile := path.Join(cwd, os.Args[2])

	fmt.Println(oldFile, newFile)

	engineRoot, err := getEngineRoot(uproject)
	if err != nil {
		panic(err)
	}

	enginePath := path.Join(engineRoot, "Engine", "Binaries", "Win64", "UE4Editor.exe")

	cmd := exec.Command(enginePath, uproject, "-diff", oldFile, newFile)

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}
