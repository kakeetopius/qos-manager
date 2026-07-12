package cmd

import (
	"errors"
	"fmt"

	"github.com/kakeetopius/qosm/internal/core/nft"
	"github.com/kakeetopius/qosm/internal/service"
	"github.com/spf13/cobra"
)

func ServiceRuleAddCmd() *cobra.Command {
	var priority string
	ruleAddCmd := cobra.Command{
		Use:   "add service...",
		Short: "Add a QoS rule(s) that matches a service i.e protocol and port.",
		Example: `  qosm rules service add tcp/443 --priority high
  qosm rules service add tcp/80 udp/53 tcp/22 --priority high`,
		Aliases: []string{"a"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			qosManager, err := getQosManager(nft.NFTOpts{
				CreateTableIfNotExists: true,
			})
			if err != nil {
				return err
			}
			defer qosManager.Close()

			for _, serv := range args {
				_, err := qosManager.AddServiceRule(serv, priority)
				if err != nil {
					return err
				}
				fmt.Printf("Service rule for %v added successfully\n", serv)
			}

			return nil
		},
	}

	ruleAddCmd.Flags().StringVarP(&priority, "priority", "p", "", "Priority for the given services.")
	ruleAddCmd.MarkFlagRequired("priority")

	return &ruleAddCmd
}

func ServiceRuleDeleteCmd() *cobra.Command {
	ruleDeleteCmd := cobra.Command{
		Use:   "delete service...",
		Short: "Delete a QoS rule(s) that matches a service",
		Example: `  qosm rules service delete tcp/443
  qosm rules service delete tcp/80 udp/53 tcp/22`,
		Aliases: []string{"d"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			toDelete := make([]service.Service, 0, len(args))
			for _, serviceSpec := range args {
				serv, err := service.ServiceFromString(serviceSpec)
				if err != nil {
					return err
				}
				toDelete = append(toDelete, serv)
			}

			qosManager, err := getQosManager(nft.NFTOpts{
				CreateTableIfNotExists: false,
			})
			if err != nil && !errors.Is(err, nft.ErrTableNotFound) {
				return err
			}
			defer qosManager.Close()

			for _, serv := range toDelete {
				err := qosManager.DeleteServiceRule(serv)
				if err != nil {
					return err
				}
				fmt.Printf("Service rule for %v deleted successfully\n", serv)
			}

			return nil
		},
	}

	return &ruleDeleteCmd
}

func ServiceRuleListCmd() *cobra.Command {
	ruleListCmd := cobra.Command{
		Use:     "list",
		Short:   "List all QoS rules that match services",
		Aliases: []string{"l"},
		RunE: func(cmd *cobra.Command, args []string) error {
			qosManager, err := getQosManager(nft.NFTOpts{
				CreateTableIfNotExists: false,
			})
			if err != nil && !errors.Is(err, nft.ErrTableNotFound) {
				return err
			}
			defer qosManager.Close()

			highPrio, err := qosManager.GetHighPriorityServices()
			if err != nil {
				return err
			}
			lowPrio, err := qosManager.GetLowPriorityServices()
			if err != nil {
				return err
			}

			return printRules(highPrio, lowPrio)
		},
	}

	return &ruleListCmd
}
