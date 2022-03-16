package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "incomes",
	Short: "get insights about incomes in a specific country",
	Long:  "get insights about incomes in a specific country",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("error executing root command")
		os.Exit(1)
	}
}

func init() {

	setupCustomProjectIncomeCmd()
	// setupGetProjectedStatsCmd()
	rootCmd.AddCommand(customProjectIncomeCmd)
	// rootCmd.AddCommand(getProjectedStatsCmd)

}
