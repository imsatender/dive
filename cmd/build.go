package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wagoodman/dive/image"
	"github.com/wagoodman/dive/ui"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:                "build",
	Short:              "Build and analyze a docker image",
	Long:               `Build and analyze a docker image`,
	DisableFlagParsing: true,
	Run:                doBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

// doBuild implements the steps taken for the build command
func doBuild(cmd *cobra.Command, args []string) {
	iidfile, err := ioutil.TempFile("/tmp", "dive.*.iid")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(iidfile.Name())

	allArgs := append([]string{"--iidfile", iidfile.Name()}, args...)
	err = runDockerCmd("build", allArgs...)
	if err != nil {
		log.Fatal(err)
	}

	imageId, err := ioutil.ReadFile(iidfile.Name())
	if err != nil {
		log.Fatal(err)
	}

	manifest, refTrees := image.InitializeData(string(imageId))
	ui.Run(manifest, refTrees)
}

// runDockerCmd runs a given Docker command in the current tty
func runDockerCmd(cmdStr string, args ...string) error {

	allArgs := cleanArgs(append([]string{cmdStr}, args...))

	cmd := exec.Command("docker", allArgs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// cleanArgs trims the whitespace from the given set of strings.
func cleanArgs(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, strings.Trim(str, " "))
		}
	}
	return r
}
