package main

import (
	"fmt"
	"os"
	"strings"
	"zkeep/internal/client"
	
	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot <dataset@tag>",
	Short: "Create a ZFS snapshot",
	Long:  `Create a snapshot of the specified dataset with the given tag.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		snapshot := args[0]
		if !strings.Contains(snapshot, "@") {
			fmt.Fprintln(os.Stderr, "Error: snapshot name must include '@' (e.g., tank/dataset@daily)")
			os.Exit(1)
		}
		
		c := client.New(client.SocketPath)
		response, err := c.SendCommand(fmt.Sprintf("snapshot %s", snapshot))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(response)
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
}
