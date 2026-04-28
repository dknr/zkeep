package main

import (
	"fmt"
	"os"
	
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zkeep",
	Short: "ZFS snapshot management client",
	Long:  `zkeep is a client for managing ZFS snapshots via the zkeep daemon.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
