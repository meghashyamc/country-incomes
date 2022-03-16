package models

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

type CustomProjectIncomeResult struct {
	CountryFrom     string
	CountryFromISO  string
	CountryTo       string
	CountryToISO    string
	AmountToProject int
	AmountProjected int
}

func (cp *CustomProjectIncomeResult) Print() {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Projected From", "Amount Projected", "Projected To", "Amount After Projection"})

	row := []string{cp.CountryFrom, strconv.Itoa(cp.AmountToProject), cp.CountryTo, strconv.Itoa(cp.AmountProjected)}
	table.Append(row)

	table.Render()
}
