package cmd

import (
	"fmt"

	"github.com/kakeetopius/qosm/internal/core/nft"
	"github.com/spf13/cobra"
)

func RestoreCmd() *cobra.Command {
	restoreCmd := cobra.Command{
		Use:   "restore",
		Short: "Restore all traffic control rules and interface settings according to the state stored in the database.",
		Long: `Restore all traffic control rules and interface settings according to the state stored in the database.
Useful when the system was rebooted or the QoS rules and interface qdisc settings were altered externally without 
using qosm and are no longer in sync with what qosm expects.`,
		Args:    cobra.NoArgs,
		Aliases: []string{"res"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRestore()
		},
	}

	return &restoreCmd
}

func runRestore() error {
	qosManager, err := getQosManager(nft.NFTOpts{
		CreateTableIfNotExists: true,
	})
	if err != nil {
		return err
	}
	defer qosManager.Close()

	err = qosManager.RestoreRules()
	if err != nil {
		return err
	}

	err = qosManager.RestoreInterfaceStates()
	if err != nil {
		return err
	}

	fmt.Println("Restore Done Successfully")

	return nil
}
