package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Adds UE4GitDiff as a difftool to the global gitconfig",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := exec.LookPath("git")
		if err != nil {
			return errors.Wrap(err, "failed to find git")
		}

		ex, err := os.Executable()
		if err != nil {
			return errors.Wrap(err, "failed to get current executable path")
		}

		difftoolCmd := fmt.Sprintf(`%s diff --old="$LOCAL" --new="$REMOTE"`, filepath.ToSlash(ex))

		err = exec.Command("git", "config", "--global", "difftool.ue4.cmd", difftoolCmd).Run()
		if err != nil {
			return errors.Wrap(err, "failed to add difftool")
		}

		mergetoolCmd := fmt.Sprintf(`%s merge --local="$LOCAL" --remote="$REMOTE" --base="$BASE" --result="$MERGED"`, filepath.ToSlash(ex))

		err = exec.Command("git", "config", "--global", "mergetool.ue4.cmd", mergetoolCmd).Run()
		if err != nil {
			return errors.Wrap(err, "failed to add mergetool")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
