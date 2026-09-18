package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "hier",
	Short: "data deployment tool",
	Long:  `support data lookup`,
}

func init() {
	rootCmd.PersistentFlags().StringP("data-dir", "d", "", "data directory")
	rootCmd.PersistentFlags().StringP("loglevel", "l", "info", "log level")
}
