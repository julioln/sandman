package cmd

import (
	"fmt"

	"github.com/julioln/sandman/build"
	"github.com/julioln/sandman/run"

	"github.com/spf13/cobra"
)

var (
	Verbose bool   = false
	Keep    bool   = false
	Layers  bool   = false
	Socket  string = ""

	rootCmd = &cobra.Command{
		Use:     "sandman",
		Short:   "sandman: Sandboxes with Podman",
		Long:    "sandman: Build and run sandboxes with Podman",
		Version: "2.0",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if Verbose {
				fmt.Println("Args: ", args)
			}
		},
	}

	buildCmd = &cobra.Command{
		Use:     "build [container_name]",
		Short:   "Build an image",
		Long:    "Build an image",
		Aliases: []string{"b"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			build.CmdExecute(Socket, Verbose, Layers, args)
		},
	}

	runCmd = &cobra.Command{
		Use:     "run [container_image]",
		Short:   "Run a sandboxed container",
		Long:    "Run a sandboxed container",
		Aliases: []string{"r"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			run.CmdExecute(Socket, Verbose, Keep, args)
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(runCmd)

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose mode (debug)")
	rootCmd.PersistentFlags().StringVarP(&Socket, "socket", "", "", "Specify podman socket")
	runCmd.Flags().BoolVarP(&Keep, "keep", "k", false, "Keep container after exit (omit --rm)")
	buildCmd.Flags().BoolVarP(&Layers, "layers", "l", false, "Use layers for building (default docker behavior)")
}
