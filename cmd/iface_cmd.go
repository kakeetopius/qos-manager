package cmd

import (
	"errors"
	"fmt"

	"github.com/kakeetopius/qosm/internal/core/htb"
	"github.com/kakeetopius/qosm/internal/core/nft"
	"github.com/spf13/cobra"
)

func IfaceCmd() *cobra.Command {
	ifaceCmd := cobra.Command{
		Use:     "iface",
		Short:   "Manage traffic control settings for an interface.",
		Aliases: []string{"i"},
	}

	ifaceCmd.AddCommand(
		IfaceEnableCmd(),
		IfaceDisableCmd(),
		IfaceListCmd(),
		IfaceStats(),
	)
	return &ifaceCmd
}

func IfaceEnableCmd() *cobra.Command {
	ifaceEnableCmd := cobra.Command{
		Use:     "enable interfaces...",
		Short:   "Enable the htb qdisc on an interface(s)",
		Aliases: []string{"e"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			qosManager, err := getQosManager(nft.NFTOpts{
				CreateTableIfNotExists: true,
			})
			if err != nil {
				return err
			}
			defer qosManager.Close()

			for _, iface := range args {
				err = qosManager.EnableTcOnInterface(iface, nil, &htb.ClassPercentages{
					HighPrioClass: 50,
					DefaultClass:  40,
					LowPrioClass:  10,
				})
				if err != nil {
					return fmt.Errorf(" Interface %v -> %w", iface, err)
				}
				fmt.Printf("Successfully enabled HTB qdisc on interface: %v\n", iface)
			}

			return nil
		},
	}

	return &ifaceEnableCmd
}

func IfaceDisableCmd() *cobra.Command {
	ifaceDisableCmd := cobra.Command{
		Use:     "disable interfaces...",
		Short:   "Disable the htb qdisc from an interface(s)",
		Aliases: []string{"d"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			qosManager, err := getQosManager(nft.NFTOpts{
				CreateTableIfNotExists: false,
			})
			if err != nil && !errors.Is(err, nft.ErrTableNotFound) {
				return err
			}
			defer qosManager.Close()

			for _, iface := range args {
				err = qosManager.DisableTcOnInterface(iface)
				if err != nil {
					return fmt.Errorf(" Interface %v -> %w", iface, err)
				}
				fmt.Printf("Successfully disabled the HTB qdisc on interface: %v\n", iface)
			}

			return nil
		},
	}

	return &ifaceDisableCmd
}

func IfaceListCmd() *cobra.Command {
	ifacelistCmd := cobra.Command{
		Use:     "list",
		Short:   "List htb enabled interfaces.",
		Aliases: []string{"l"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			qosManager, err := getQosManager(nft.NFTOpts{
				CreateTableIfNotExists: false,
			})
			if err != nil && !errors.Is(err, nft.ErrTableNotFound) {
				return err
			}
			enabledIfaces := qosManager.EnabledInterfaces()
			if len(enabledIfaces) == 0 {
				fmt.Println("No htb enabled interfaces.")
				return nil
			}

			HeadingPrinter.Println("Enabled Interfaces")
			return printIfaces(enabledIfaces)
		},
	}

	return &ifacelistCmd
}

func IfaceStats() *cobra.Command {
	ifacelistCmd := cobra.Command{
		Use:     "stats",
		Short:   "Get stats for an interface",
		Aliases: []string{"s"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			qosManager, err := getQosManager(nft.NFTOpts{
				CreateTableIfNotExists: false,
			})
			if err != nil && !errors.Is(err, nft.ErrTableNotFound) {
				return err
			}
			ifaceStats, err := qosManager.GetIfaceStats(args[0])
			if err != nil {
				return err
			}
			return printIfaceStats(&ifaceStats)
		},
	}

	return &ifacelistCmd
}
