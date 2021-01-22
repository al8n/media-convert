package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/al8n/media-convert/config"
	"github.com/karrick/godirwalk"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var all *cobra.Command

type runner struct {
	mu sync.Mutex
	sync.WaitGroup

	jobs        int
	workingpool chan *ctr
	waitingPool chan *ctr
}

type ctr struct {
	name string
	cmd  *exec.Cmd
}

func newAllCommand() *cobra.Command {
	all = &cobra.Command{
		Use:   "all",
		Short: "convert all files in the specify directory",
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			if config.GetAllConfig().InputFormat == "" {
				return errors.New("input format is required")
			}

			if config.GetAllConfig().OutputFormat == "" {
				return errors.New("output format is required")
			}

			var r = &runner{
				workingpool: make(chan *ctr, int(config.GetAllConfig().WorkPoolSize)),
			}

			if err = r.prepare(); err != nil {
				logrus.Error(err)
				return err
			}

			if err = r.execute(); err != nil {
				logrus.Error(err)
				return err
			}
			r.Wait()
			return nil
		},
	}

	var cfg = config.GetAllConfig()

	all.Flags().StringVarP(&cfg.Path, "path", "p", ".", "specify the aim directory")
	all.Flags().StringVarP(&cfg.InputFormat, "input-format", "i", "", "specify the input file format")
	all.Flags().StringVarP(&cfg.OutputFormat, "output-format", "o", "", "specify the output file format")
	all.Flags().BoolVar(&cfg.WithProcessBar, "progress-bar", true, "show convert progress bar")
	all.Flags().UintVarP(&cfg.WorkPoolSize, "workpool-size", "s", 10, "allow how many files to convert simultaneously")

	all.Flags().BoolVarP(&cfg.RemoveSourses, "remove-sourses", "r", false, "remove input files after successfully convert")
	return all
}

func (r *runner) prepare() (err error) {

	var jobs []*ctr

	err = godirwalk.Walk(config.GetAllConfig().Path, &godirwalk.Options{
		Callback: func(name string, de *godirwalk.Dirent) error {
			if strings.Contains(name, ".git") {
				return godirwalk.SkipThis
			}

			var (
				filename, outputFile string
				ext                  string
			)

			ext = filepath.Ext(name)
			filename = strings.TrimSuffix(name, ext)

			if ext != "."+config.GetAllConfig().InputFormat {
				return nil
			}

			outputFile = filename + "." + config.GetAllConfig().OutputFormat

			job := &ctr{
				name: filename,
				cmd:  exec.Command("ffmpeg", "-i", name, outputFile),
			}

			jobs = append(jobs, job)
			r.jobs++
			return nil
		},
		Unsorted: true,
	})

	r.waitingPool = make(chan *ctr, len(jobs))
	for _, w := range jobs {
		r.waitingPool <- w
	}

	return
}

func (r *runner) execute() (err error) {
	for {
		if len(r.waitingPool) <= 0 {
			break
		}

		job := <-r.waitingPool
		r.workingpool <- job
		r.Add(1)
		go func(c *ctr) {
			err = c.cmd.Start()
			if err != nil {
				logrus.Error(err)
				return
			}

			logrus.Infof("Start convert %s\n", c.name)
			if err = c.cmd.Wait(); err != nil {
				logrus.Error(err)
				return
			}
			logrus.Infof("Finish convert %s\n", c.name)

			if config.GetAllConfig().RemoveSourses {
				if err = os.Remove(fmt.Sprintf("%s.%s", c.name, config.GetAllConfig().InputFormat)); err != nil {
					logrus.Error(err)
				}
			}

			r.Done()
			<-r.workingpool

		}(job)
	}
	return
}
