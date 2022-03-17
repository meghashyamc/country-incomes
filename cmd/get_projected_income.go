package cmd

import (
	"fmt"
	"time"

	"github.com/meghashyamc/country-incomes/services/countrydata"
	"github.com/meghashyamc/country-incomes/validation"
	"github.com/spf13/cobra"
)

const (
	usISO        = "us"
	monthsInYear = 12
)

var (
	countryFromForGDPProjection *string
	countryToForGDPProjection   *string
)

var getProjectedPerCapitaIncomeCmd = &cobra.Command{
	Use:   "averageincome",
	Short: "stats related to per capita income in a country from the perspective of another country",
	Long:  "stats related to per capita income in a country from the perspective of another country",
	Run:   getProjectedPerCapitaIncome,
}

func getProjectedPerCapitaIncome(cmd *cobra.Command, args []string) {
	st := time.Now()

	getProjectedPerCapitaIncomeResult, err := validation.ValidateGetProjectedPerCapitaIncome(countryFromForGDPProjection, countryToForGDPProjection)
	if err != nil {
		return
	}

	gdpPerCapitaPPP, err := countrydata.GetGDPPerCapitaPPP(getProjectedPerCapitaIncomeResult.CountryFromISO)
	if err != nil {
		return
	}
	projectedAmount, _, err := countrydata.ProjectAmount(usISO, getProjectedPerCapitaIncomeResult.CountryToISO, gdpPerCapitaPPP)
	if err != nil {
		return
	}

	getProjectedPerCapitaIncomeResult.GDPPerCapitaAnnual = projectedAmount
	getProjectedPerCapitaIncomeResult.GDPPerCapitaMonthly = projectedAmount / monthsInYear

	getProjectedPerCapitaIncomeResult.Print()
	fmt.Println(time.Since(st).Truncate(100000))

}

func setupGetProjectedPerCapitaIncomeCmd() {
	countryFromForGDPProjection = getProjectedPerCapitaIncomeCmd.Flags().StringP("from", "f", "", "the country from which to project an amount")
	countryToForGDPProjection = getProjectedPerCapitaIncomeCmd.Flags().StringP("to", "t", "", "the country to which an amount should be projected")

}
