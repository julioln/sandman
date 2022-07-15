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
		Use:   "sandman",
		Short: "sandman: short",
		Long:  "sandman: long",
	}

	buildCmd = &cobra.Command{
		Use:     "build",
		Short:   "build: short",
		Long:    "build: long",
		Aliases: []string{"b"},
		Run: func(cmd *cobra.Command, args []string) {
			if Verbose {
				fmt.Println("Args: ", args)
			}
			build.CmdExecute(Socket, Verbose, Layers, args)
		},
	}

	runCmd = &cobra.Command{
		Use:     "run",
		Short:   "run: short",
		Long:    "run: long",
		Aliases: []string{"r"},
		Run: func(cmd *cobra.Command, args []string) {
			if Verbose {
				fmt.Println("Args: ", args)
			}
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

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose mode")
	rootCmd.PersistentFlags().StringVarP(&Socket, "socket", "", "", "Specify podman socket")
	runCmd.Flags().BoolVarP(&Keep, "keep", "k", false, "Keep container after exit (omit --rm)")
	buildCmd.Flags().BoolVarP(&Layers, "layers", "l", false, "Use layers for building (default docker behavior)")
}
