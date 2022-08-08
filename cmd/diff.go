package cmd

import (
	"github.com/mircearoata/UEGitDiff/unrealengine"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
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

func discoverUProject() (string, error) {
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

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Runs UE4Editor diffing the two assets",
	Long:  `Usage: UEGitDiff diff --old="old/file/path.uasset" --new="new/file/path.uasset"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "failed to get current working directory")
		}

		uproject, err := discoverUProject()
		if err != nil {
			return errors.Wrap(err, "failed to discover uproject")
		}

		oldFile, err := cmd.Flags().GetString("old")
		if err != nil {
			return errors.Wrap(err, "failed to get old file")
		}
		newFile, err := cmd.Flags().GetString("new")
		if err != nil {
			return errors.Wrap(err, "failed to get new file")
		}

		if !filepath.IsAbs(oldFile) {
			oldFile = filepath.Join(cwd, oldFile)
		}
		if !filepath.IsAbs(newFile) {
			newFile = filepath.Join(cwd, newFile)
		}

		engineRoot, err := unrealengine.GetEngineRootFromProject(uproject)
		if err != nil {
			return errors.Wrap(err, "failed to get engine root")
		}

		enginePath := filepath.Join(engineRoot, "Engine", "Binaries", "Win64", "UE4Editor.exe")

		ueCmd := exec.Command(enginePath, uproject, "-diff", oldFile, newFile)

		err = ueCmd.Run()
		if err != nil {
			return errors.Wrap(err, "failed to run UE4Editor")
		}

		err = ueCmd.Wait()
		if err != nil {
			return errors.Wrap(err, "failed to wait for UE4Editor")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().String("old", "", "old file")
	diffCmd.Flags().String("new", "", "new file")
}
