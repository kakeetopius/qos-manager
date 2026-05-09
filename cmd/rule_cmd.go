package cmd

import (
	"fmt"

	"github.com/kakeetopius/qosm/internal/core/tc"
	"github.com/kakeetopius/qosm/internal/core/util"
	"github.com/spf13/cobra"
)

var iface string

func RuleCmd() *cobra.Command {
	var flush bool
	ruleCmd := cobra.Command{
		Use:     "rule",
		Short:   "Manage the traffic control rules.",
		Aliases: []string{"r"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if flush {
				if iface == "" {
					return fmt.Errorf("please provide the interface")
				}
				return tc.FlushRules(iface)
			}

			return nil
		},
	}

	ruleCmd.PersistentFlags().StringVarP(&iface, "iface", "i", "", "The network interface to use.")
	ruleCmd.Flags().BoolVar(&flush, "flush", false, "Flush all rules.")

	ruleCmd.AddCommand(
		RuleAddCmd(),
	)
	return &ruleCmd
}

func RuleAddCmd() *cobra.Command {
	var priority string
	ruleAddCmd := cobra.Command{
		Use:     "add <targets>",
		Short:   "Add a QoS rule.",
		Aliases: []string{"a"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Adding rule for these targets: %v\n", args[0])
			targets, _, err := util.TargetsFromStringWithDNSLookup(args[0])
			if err != nil {
				return err
			}

			var tcPriority tc.Priority

			switch priority {
			case "high":
				tcPriority = tc.PRIORITYHIGH
			case "low":
				tcPriority = tc.PRIORITYLOW
			default:
				return fmt.Errorf("unknown priority %v", priority)
			}

			err = tc.AddRule(iface, targets[0], tcPriority)
			if err != nil {
				return err
			}
			return nil
		},
	}

	ruleAddCmd.Flags().StringVar(&priority, "priority", "", "Priority for the given targets.")
	ruleAddCmd.MarkFlagRequired("priority")

	return &ruleAddCmd
}

// func RuleDeleteCmd() *cobra.Command {
// 	var deleteAll bool
// 	ruleDelCmd := cobra.Command{
// 		Use:     "delete",
// 		Short:   "Add a QoS rule.",
// 		Aliases: []string{"d"},
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 		},
// 	}
//
// 	ruleDelCmd.Flags().BoolVarP(&deleteAll, "all", "a", false, "Delete all rules.")
// 	return &ruleDelCmd
// }
