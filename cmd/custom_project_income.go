package cmd

import (
	"fmt"
	"time"

	"github.com/meghashyamc/country-incomes/services/countrydata"
	"github.com/meghashyamc/country-incomes/validation"
	"github.com/spf13/cobra"
)

var (
	countryFromCustomProject, countryToCustomProject *string
	amount                                           *int
)

var customProjectIncomeCmd = &cobra.Command{
	Use:   "project",
	Short: "project an amount in one country to another country",
	Long:  "project an amount in one country to another country (for eg.what amount in rupees would $20 in the US represent in India?)",
	Run:   projectAmount,
}

func projectAmount(cmd *cobra.Command, args []string) {
	st := time.Now()
	customProjectIncomeResult, err := validation.ValidateCustomProjectIncome(countryFromCustomProject, countryToCustomProject, amount)
	if err != nil {
		return
	}

	projectedAmount, parityFactor, err := countrydata.ProjectAmount(customProjectIncomeResult.CountryFromISO, customProjectIncomeResult.CountryToISO, customProjectIncomeResult.AmountToProject)
	if err != nil {
		return
	}
	customProjectIncomeResult.AmountProjected = projectedAmount
	customProjectIncomeResult.MultiplicationFactor = parityFactor
	customProjectIncomeResult.Print()
	fmt.Println(time.Since(st).Truncate(100000))
}

func setupCustomProjectIncomeCmd() {
	countryFromCustomProject = customProjectIncomeCmd.Flags().StringP("from", "f", "", "the country from which to project an amount")
	countryToCustomProject = customProjectIncomeCmd.Flags().StringP("to", "t", "", "the country to which an amount should be projected")
	amount = customProjectIncomeCmd.Flags().IntP("amount", "a", 0, "the amount to be projected")

}
