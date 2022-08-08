package unrealengine

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

type UProject struct {
	EngineAssociation string `json:"EngineAssociation"`
}

func getEngineAssociationFromProject(projectPath string) (string, error) {
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

func GetEngineRootFromProject(projectPath string) (string, error) {
	engineAssociation, err := getEngineAssociationFromProject(projectPath)
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
