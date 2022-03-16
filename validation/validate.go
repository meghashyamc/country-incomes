package validation

import (
	"errors"
	"strings"

	"github.com/meghashyamc/country-incomes/models"
	"github.com/meghashyamc/country-incomes/services/countrydata"

	log "github.com/sirupsen/logrus"
)

const amountRules = "amount to project must be a positive integer greater than one"

func ValidateCustomProjectIncome(countryFrom, countryTo *string, amount *int) (*models.CustomProjectIncomeResult, error) {

	errParamsNotSpecified := "from (country), to (country) and amount to project must be specified"
	if countryFrom == nil || countryTo == nil || amount == nil || *countryFrom == "" || *countryTo == "" || *amount == 0 {
		log.Error(errParamsNotSpecified)
		return nil, errors.New(errParamsNotSpecified)
	}

	cleanedCountryFrom := cleanString(*countryFrom)
	cleanedCountryTo := cleanString(*countryTo)

	countryFromISO, err := countrydata.GetISO(cleanedCountryFrom)
	if err != nil {
		return nil, err
	}
	countryToISO, err := countrydata.GetISO(cleanedCountryTo)
	if err != nil {
		return nil, err
	}

	if err := validateAmount(*amount); err != nil {
		return nil, err
	}
	return &models.CustomProjectIncomeResult{CountryFrom: cleanedCountryFrom, CountryTo: cleanedCountryTo, CountryFromISO: countryFromISO, CountryToISO: countryToISO, AmountToProject: *amount}, nil
}

func cleanString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
func validateAmount(amount int) error {
	if amount <= 0 {
		log.Error(amountRules)
		return errors.New(amountRules)
	}
	return nil
}
