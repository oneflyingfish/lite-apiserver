package app

import (
	"LiteKube/cmd/lite-apiserver/app/options"
	"LiteKube/pkg/version"
	verflag "LiteKube/pkg/version/varflag"
	"fmt"
	"io"

	"LiteKube/pkg/util"

	"github.com/moby/term"
	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
)

var ComponentName = "lite-apiserver"

func NewServerCommand() *cobra.Command {
	opt := options.NewServerRunOption()

	cmd := &cobra.Command{
		Use:  ComponentName,
		Long: `The lite-apiserver is one simplified version of kube-apiserver, which is only service for one node and deal with pods.`,

		// stop printing usage when the command errors
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested() // --version=false/true/raw to print version

			klog.Infof("Versions: %+v\n", version.Get())

			// load config from disk-file and merge with flags
			if errs := opt.LoadConfig(); len(errs) != 0 {
				klog.Error("some error in your configs")
				return fmt.Errorf("some error in your configs")
			}

			// complete all default server options,current is none-function
			if err := opt.Complete(); err != nil {
				klog.Errorf("complete options error: %v", err)
				return err
			}

			// ready to run
			return Run(opt, util.SetupSignalHandler())
		},
		Args: func(cmd *cobra.Command, args []string) error { // Validate unresolved args
			for _, arg := range args {
				if len(arg) > 0 {
					klog.Errorf("%q does not support subcommands at this time but get %q", cmd.CommandPath(), args)
					return fmt.Errorf("%q does not support subcommands at this time but get %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	fs := cmd.Flags()
	namedFlagSets := opt.GetNamedFlagsSet()

	// Add the custom Flagset to cmd.Flags(), so the value will be parse
	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	// better print
	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := TerminalSize(cmd.OutOrStdout()) // terminal_width, terminal_height, error
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})

	return cmd

}

// start to run lite-apiserver
func Run(serverOptions *options.ServerRunOption, stopCh <-chan struct{}) error {
	return nil
}

func TerminalSize(w io.Writer) (int, int, error) {
	outFd, isTerminal := term.GetFdInfo(w)
	if !isTerminal {
		return 0, 0, fmt.Errorf("given writer is no terminal")
	}
	winsize, err := term.GetWinsize(outFd)
	if err != nil {
		return 0, 0, err
	}
	return int(winsize.Width), int(winsize.Height), nil
}
