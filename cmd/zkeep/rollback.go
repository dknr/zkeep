package main

import (
	"fmt"
	"os"
	"zkeep/internal/client"
	
	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback <snapshot>",
	Short: "Rollback to a ZFS snapshot",
	Long:  `Rollback the dataset to the specified snapshot.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(client.SocketPath)
		response, err := c.SendCommand(fmt.Sprintf("rollback %s", args[0]))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(response)
	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
}
