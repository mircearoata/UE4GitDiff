package unrealengine

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func discoverUProjectFromStartPath(startPath string) (string, error) {
	for currentPath := startPath; filepath.Dir(currentPath) != currentPath; currentPath = filepath.Dir(currentPath) {
		matches, err := filepath.Glob(filepath.Join(currentPath, "*.uproject"))
		if err == nil && len(matches) > 0 {
			return matches[0], nil
		}
	}
	return "", errors.New("no uproject found")
}

func DiscoverUProject() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "failed to get current working directory")
	}

	uproject, err := discoverUProjectFromStartPath(cwd)
	if err == nil {
		return uproject, nil
	}

	gitDir, ok := os.LookupEnv("GIT_DIR")
	if !ok {
		return "", errors.New("GIT_DIR environment variable not set")
	}
	uproject, err = discoverUProjectFromStartPath(gitDir)
	if err == nil {
		return uproject, nil
	}

	return "", errors.Wrap(err, "failed to discover uproject")
}
