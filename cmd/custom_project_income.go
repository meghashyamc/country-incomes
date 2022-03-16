package cmd

import (
	"github.com/meghashyamc/country-incomes/services/countrydata"
	"github.com/meghashyamc/country-incomes/validation"
	"github.com/spf13/cobra"
)

var (
	countryFrom, countryTo *string
	amount                 *int
)

var customProjectIncomeCmd = &cobra.Command{
	Use:   "project",
	Short: "project an amount in one country to another country",
	Long:  "project an amount in one country to another country (for eg.what amount in rupees would $20 in the US represent in India?)",
	Run:   projectAmount,
}

func projectAmount(cmd *cobra.Command, args []string) {

	customProjectIncomeResult, err := validation.ValidateCustomProjectIncome(countryFrom, countryTo, amount)
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
}

func setupCustomProjectIncomeCmd() {
	countryFrom = customProjectIncomeCmd.Flags().StringP("from", "f", "", "the country from which to project an amount")
	countryTo = customProjectIncomeCmd.Flags().StringP("to", "t", "", "the country to which an amount should be projected")
	amount = customProjectIncomeCmd.Flags().IntP("amount", "a", 0, "the amount to be projected")

}
