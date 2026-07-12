package cmd

import (
	"errors"

	"github.com/kakeetopius/qosm/internal/core/nft"
	"github.com/spf13/cobra"
)

func StatsCmd() *cobra.Command {
	ifaceEnableCmd := cobra.Command{
		Use:     "stats",
		Short:   "Get QoS stats",
		Aliases: []string{"s"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			qosManager, err := getQosManager(nft.NFTOpts{
				CreateTableIfNotExists: false,
			})
			if err != nil && !errors.Is(err, nft.ErrTableNotFound) {
				return err
			}
			defer qosManager.Close()

			stats, err := qosManager.GetStats()
			if err != nil {
				return err
			}
			return printStats(&stats)
		},
	}

	return &ifaceEnableCmd
}
