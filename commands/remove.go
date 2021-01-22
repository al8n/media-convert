package commands

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	ext    string
	remove *cobra.Command
	path   string
	force  bool
)

func newRemoveCommand() *cobra.Command {
	remove = &cobra.Command{
		Use:   "remove",
		Short: "remove files in specify directory according to extension",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if ext == "" {
				return errors.New("file extension is required")
			}

			return execute()
		},
	}

	remove.Flags().StringVarP(&path, "path", "p", ".", "specify the directory")
	remove.Flags().StringVarP(&ext, "extension", "e", "", "specify the file extension")
	remove.Flags().BoolVarP(&force, "force", "f", false, "force remove a file")
	return remove
}

func execute() (err error) {

	err = godirwalk.Walk(path, &godirwalk.Options{
		Callback: func(name string, de *godirwalk.Dirent) error {
			if strings.Contains(name, ".git") {
				return godirwalk.SkipThis
			}

			if filepath.Ext(name) != "."+ext {
				return nil
			}

			if err = os.Remove(name); err != nil {
				logrus.Error(err)
			}
			return nil
		},
		Unsorted: true,
	})

	return
}
