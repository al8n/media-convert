package commands

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	root *cobra.Command
)

// Execute will execute media file convert
func Execute() (err error) {
	return newRootCommand().Execute()
}

func newRootCommand() *cobra.Command {
	root = &cobra.Command{
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ffmpeg := exec.Command("ffmpeg", os.Args[1:]...)
			ffmpeg.Stderr = os.Stderr

			if err = ffmpeg.Start(); err != nil {
				return err
			}

			if err = ffmpeg.Wait(); err != nil {
				return err
			}

			return
		},
	}

	root.SetHelpFunc(
		func(cmd *cobra.Command, args []string) {
			ffmpeg := exec.Command("ffmpeg", "-h")
			ffmpeg.Stderr = os.Stderr
			ffmpeg.Stdout = os.Stdout

			ffmpeg.Start()
			ffmpeg.Wait()
		},
	)

	root.AddCommand(newAllCommand())
	root.AddCommand(newRemoveCommand())
	return root
}
