package models

import (
	"os"
	"strconv"

	"github.com/iancoleman/strcase"
	"github.com/olekukonko/tablewriter"
)

type CustomProjectIncomeResult struct {
	CountryFrom          string
	CountryFromISO       string
	CountryTo            string
	CountryToISO         string
	AmountToProject      int
	AmountProjected      int
	MultiplicationFactor float64
}

func (cp *CustomProjectIncomeResult) Print() {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Projected From", "Amount", "Projected To", "Amount After Projection", "Mult. Factor"})

	row := []string{strcase.ToCamel(cp.CountryFrom), strconv.Itoa(cp.AmountToProject), strcase.ToCamel(cp.CountryTo), strconv.Itoa(cp.AmountProjected), strconv.FormatFloat(cp.MultiplicationFactor, 'f', 4, 64)}
	table.Append(row)

	table.Render()
}
