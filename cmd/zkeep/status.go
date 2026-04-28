package main

import (
	"fmt"
	"os"
	"zkeep/internal/client"
	
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check daemon status",
	Long:  `Check if the zkeep daemon is running and responsive.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(client.SocketPath)
		response, err := c.SendCommand("status")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(response)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
