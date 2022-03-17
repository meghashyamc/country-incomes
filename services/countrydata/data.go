package countrydata

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"text/template"

	"github.com/meghashyamc/country-incomes/services/cache"
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

func GetGDPPerCapitaPPP(countryISO string) (int, error) {
	currentYear := getCurrentYearAsStr()
	return getGDPPerCapitaPPPFromAPI(countryISO, currentYear)

}

func getGDPPerCapitaPPPFromAPI(countryISO, currentYear string) (int, error) {
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

	return int(gdpPerCapitaPPP), nil

}
