package cmd

import (
	"github.com/kakeetopius/qosm/internal/daemon"
	"github.com/spf13/cobra"
)

func DaemonCmd() *cobra.Command {
	daemonCmd := cobra.Command{
		Use:     "daemon",
		Aliases: []string{"d"},
		Short:   "Manage the qos daemon",
		Long: `Manage the qos daemon

Some operations performed by qosm require root privileges, but running the entire qosm process as root is not always desirable for example, 
when running the web server. 
To address this, a small privileged daemon can be started with 'sudo qosm daemon run' and then any subsquent usages of qosm will send all privileged 
operations to the daemon hence no longer requiring to run with 'sudo' as long the daemon is running. Only the daemon requires sudo.

Note: subsequent usages of qosm will require using the --daemon-mode or -d flag e.g
  qosm rule service add tcp/80 udp/53 --daemon-mode
  qosm web run -d
`,
	}

	daemonCmd.AddCommand(runDaemonCmd())
	return &daemonCmd
}

func runDaemonCmd() *cobra.Command {
	runCmd := cobra.Command{
		Use:   "run",
		Short: "run the daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := daemon.New(daemon.Options{
				SocketPath: appConfig.GetString("daemon.sock"),
				Debug:      debug,
			})
			if err != nil {
				return err
			}

			return d.Run()
		},
	}

	return &runCmd
}
