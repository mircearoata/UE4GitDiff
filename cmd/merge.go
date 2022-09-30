package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mircearoata/UEGitDiff/unrealengine"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Runs UE4Editor merging the assets",
	Long:  `Usage: UEGitDiff merge --remote="remote/file/path.uasset" --local="local/file/path.uasset"--base="base/file/path.uasset" --result="result/file/path.uasset"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "failed to get current working directory")
		}

		uproject, err := unrealengine.DiscoverUProject()
		if err != nil {
			return errors.Wrap(err, "failed to discover uproject")
		}

		remoteFile, err := cmd.Flags().GetString("remote")
		if err != nil {
			return errors.Wrap(err, "failed to get remote file")
		}
		localFile, err := cmd.Flags().GetString("local")
		if err != nil {
			return errors.Wrap(err, "failed to get local file")
		}
		baseFile, err := cmd.Flags().GetString("base")
		if err != nil {
			return errors.Wrap(err, "failed to get base file")
		}
		resultFile, err := cmd.Flags().GetString("result")
		if err != nil {
			return errors.Wrap(err, "failed to get result file")
		}

		if !filepath.IsAbs(remoteFile) {
			remoteFile = filepath.Join(cwd, remoteFile)
		}
		if !filepath.IsAbs(localFile) {
			localFile = filepath.Join(cwd, localFile)
		}
		if !filepath.IsAbs(baseFile) {
			baseFile = filepath.Join(cwd, baseFile)
		}
		if !filepath.IsAbs(resultFile) {
			resultFile = filepath.Join(cwd, resultFile)
		}

		engineRoot, err := unrealengine.GetEngineRootFromProject(uproject)
		if err != nil {
			return errors.Wrap(err, "failed to get engine root")
		}

		enginePath := filepath.Join(engineRoot, "Engine", "Binaries", "Win64", "UE4Editor.exe")

		ueCmd := exec.Command(enginePath, uproject, "-diff", remoteFile, localFile, baseFile, resultFile)

		err = ueCmd.Run()
		if err != nil {
			return errors.Wrap(err, "failed to run UE4Editor")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().String("remote", "", "remote file")
	mergeCmd.Flags().String("local", "", "local file")
	mergeCmd.Flags().String("base", "", "base file")
	mergeCmd.Flags().String("result", "", "result file")
}
