package countrydata

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"text/template"

	"github.com/meghashyamc/country-incomes/services/cache"
	"github.com/meghashyamc/country-incomes/services/db"
	log "github.com/sirupsen/logrus"
)

const isoKey = "iso"

func GetISO(country string) (string, error) {
	iso, err := cache.ReadHash(isoKey, country)
	if err == nil && iso != "" {
		return iso, nil
	}
	isoDataList, err := makeRequestAndGetData(os.Getenv("ISO_URL"))
	if err != nil {
		return "", err
	}

	for _, isoData := range isoDataList {

		isoDataMap, ok := isoData.(map[string]interface{})
		if !ok {
			errDataReceivedFormat := "unexpected format of country details received from data source"
			log.WithFields(log.Fields{"data_received": isoData}).Error(errDataReceivedFormat)
			return "", errors.New(errDataReceivedFormat)
		}
		countryName := isoDataMap["name"]
		countryISOCode := isoDataMap["iso2Code"]

		if countryName != nil && countryISOCode != nil {

			if strings.Contains(strings.ToLower(countryName.(string)), country) {
				cache.WriteHash(isoKey, country, countryISOCode.(string))
				return countryISOCode.(string), nil
			}
		}
	}

	errNoSuchCountry := "could not validate country"
	log.WithFields(log.Fields{"country_name": country}).Info(errNoSuchCountry)

	return "", errors.New(errNoSuchCountry)
}

func GetGDPPerCapitaPPP(countryName, countryISO string) (int, error) {
	currentYear := getCurrentYear()
	countryISO = strings.ToLower(countryISO)
	gdpPerCapitaPPP, err := getGDPPerCapitaPPPFromDB(countryISO, currentYear)
	if err == nil && gdpPerCapitaPPP > 0 {
		return gdpPerCapitaPPP, nil
	}
	log.Info("getting GDP Per capita from API")

	return getGDPPerCapitaPPPFromAPI(countryName, countryISO, currentYear)

}

func insertGDPPerCapitaPPPInDB(countryName, countryISO string, year, gdpPerCapitaPPP int) {
	dbClient, err := db.Get()
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info("could not reach DB to store parity data for future reference")
		return
	}

	if err := dbClient.InsertGDPPerCapitaPPP(countryName, countryISO, year, gdpPerCapitaPPP); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info("could not store GDP Per capita (PPP) data in DB for future reference")
		return
	}
}
func getGDPPerCapitaPPPFromDB(countryISO string, currentYear int) (int, error) {

	dbClient, err := db.Get()
	if err != nil {
		return 0.0, err
	}

	return dbClient.GetGDPPerCapitaPPP(countryISO, currentYear)

}

func getGDPPerCapitaPPPFromAPI(countryName, countryISO string, currentYear int) (int, error) {
	t := template.Must(template.New("gdpPerCapitaPPP").Parse(os.Getenv("GDPPPP_URL")))
	urlBuf := &bytes.Buffer{}
	pppDetails := PPPDetails{ISO: strings.ToLower(countryISO), CurrentYear: currentYear}

	if err := t.Execute(urlBuf, pppDetails); err != nil {
		log.WithFields(log.Fields{"country_iso": countryISO, "err": err.Error()}).Error("could not form URL to get purchasing power parity data")
		return 0.0, err
	}

	gdpPerCapitaPPPDataList, err := makeRequestAndGetData(urlBuf.String())
	if err != nil {
		return 0.0, err
	}

	countryDataMap, err := getCountryDataMapFromParityDataList(countryISO, gdpPerCapitaPPPDataList)
	if err != nil {
		return 0.0, err
	}

	gdpPerCapitaPPP, ok := countryDataMap[valueKey].(float64)
	errNoData := "no GDP PPP data is available for the country whose average income needs to be projected"
	if !ok {
		log.WithFields(log.Fields{"country_iso": countryISO}).Info(errNoData)
		return 0, errors.New(errNoData)
	}

	insertGDPPerCapitaPPPInDB(countryName, countryISO, forceIntFromString(countryDataMap[yearKey].(string)), int(gdpPerCapitaPPP))

	return int(gdpPerCapitaPPP), nil

}
