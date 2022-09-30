package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mircearoata/UEGitDiff/unrealengine"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Runs UE4Editor diffing the two assets",
	Long:  `Usage: UEGitDiff diff --old="old/file/path.uasset" --new="new/file/path.uasset"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "failed to get current working directory")
		}

		uproject, err := unrealengine.DiscoverUProject()
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

		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().String("old", "", "old file")
	diffCmd.Flags().String("new", "", "new file")
}
