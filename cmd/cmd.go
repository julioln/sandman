package cmd

import (
	"fmt"

	"github.com/julioln/sandman/build"
	"github.com/julioln/sandman/config"
	"github.com/julioln/sandman/constants"
	"github.com/julioln/sandman/podman"
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
		Version: constants.VERSION,
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
		Use:     "start [container_image]",
		Short:   "Start a sandboxed container",
		Long:    "Start a sandboxed container",
		Aliases: []string{"s", "run", "r"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			run.CmdExecute(Socket, Verbose, Keep, args)
		},
	}

	scaffoldCmd = &cobra.Command{
		Use:     "sample",
		Short:   "Prints a sample configuration file",
		Long:    "Print a sample configuration file",
		Aliases: []string{"scaffold"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(config.Scaffold())
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(scaffoldCmd)

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose mode (log debug). Defaults to false")
	rootCmd.PersistentFlags().StringVarP(&Socket, "socket", "", "", fmt.Sprintf("Specify podman socket. Defaults to %s", podman.DefaultSocket()))
	runCmd.Flags().BoolVarP(&Keep, "keep", "k", false, "Keep container after exit (omit --rm)")
	buildCmd.Flags().BoolVarP(&Layers, "layers", "l", false, "Use layers for building (default docker behavior)")
}
