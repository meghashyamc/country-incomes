package cmd

import (
	"github.com/spf13/cobra"
)

var getProjectedStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "stats related to per capita income in a country from the perspective of another country",
	Long:  "stats related to per capita income in a country from the perspective of another country",
	Run:   getProjectedStats,
}

func getProjectedStats(cmd *cobra.Command, args []string) {

}
