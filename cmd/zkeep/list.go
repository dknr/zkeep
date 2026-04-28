package main

import (
	"fmt"
	"os"
	"zkeep/internal/client"
	
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list <dataset>",
	Short: "List ZFS snapshots",
	Long:  `List all snapshots for the specified dataset.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(client.SocketPath)
		response, err := c.SendCommand(fmt.Sprintf("list %s", args[0]))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(response)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
