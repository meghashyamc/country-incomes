package countrydata

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/meghashyamc/country-incomes/services/db"
	log "github.com/sirupsen/logrus"
)

type PPPDetails struct {
	ISO         string
	CurrentYear int
}

func ProjectAmount(countryFromISO, countryToISO string, amount int) (int, float64, error) {

	currentYear := getCurrentYear()
	countryFromISO = strings.ToLower(countryFromISO)
	countryToISO = strings.ToLower(countryToISO)
	parityFactor, err := getParityFactorFromDB(countryFromISO, countryToISO, currentYear)
	if err == nil && parityFactor != 0 {
		return int(float64(amount) * parityFactor), parityFactor, nil
	}
	log.Info("getting parity factor from API")
	parityFactor, err = getParityFactorFromAPI(countryFromISO, countryToISO, currentYear)
	if err != nil {
		return 0, 0.0, err
	}

	return int(float64(amount) * parityFactor), parityFactor, nil

}

func getParityFactorFromAPI(countryFromISO, countryToISO string, currentYear int) (float64, error) {

	t := template.Must(template.New("parityFactor").Parse(os.Getenv("PARITY_URL")))
	urlBufCountryFrom := &bytes.Buffer{}
	pppDetails := PPPDetails{ISO: countryFromISO, CurrentYear: currentYear}

	if err := t.Execute(urlBufCountryFrom, pppDetails); err != nil {
		log.WithFields(log.Fields{"country_iso": countryFromISO, "err": err.Error()}).Error("could not form URL to get purchasing power parity data")
		return 0.0, err
	}

	urlBufCountryTo := &bytes.Buffer{}
	pppDetails.ISO = countryToISO

	if err := t.Execute(urlBufCountryTo, pppDetails); err != nil {
		log.WithFields(log.Fields{"country_iso": countryToISO, "err": err.Error()}).Error("could not form URL to get purchasing power parity data")
		return 0.0, err
	}
	countryFromParityDataList, err := makeRequestAndGetData(urlBufCountryFrom.String())
	if err != nil {
		return 0.0, err
	}

	countryToParityDataList, err := makeRequestAndGetData(urlBufCountryTo.String())
	if err != nil {
		return 0.0, err
	}

	countryFromDataMap, err := getCountryDataMapFromParityDataList(countryFromISO, countryFromParityDataList)
	if err != nil {
		return 0.0, err
	}

	countryToDataMap, err := getCountryDataMapFromParityDataList(countryToISO, countryToParityDataList)
	if err != nil {
		return 0.0, err
	}
	lcuForCountryTo, ok := countryToDataMap[valueKey].(float64)
	errNoDataTo := "no data is available for the country to which an amount needs to be projected"

	if !ok {
		log.WithFields(log.Fields{"country_iso": countryToISO}).Info(errNoDataTo)
		return 0.0, errors.New(errNoDataTo)
	}
	errNoDataFrom := "no data is available for the country from which an amount needs to be projected"
	lcuForCountryFrom, ok := countryFromDataMap[valueKey].(float64)
	if !ok {
		log.WithFields(log.Fields{"country_iso": countryFromISO}).Info(errNoDataFrom)
		return 0.0, errors.New(errNoDataTo)
	}

	parityFactor := lcuForCountryTo / lcuForCountryFrom
	insertParityFactorInDB(countryFromISO, countryToISO, forceIntFromString(countryFromDataMap[yearKey].(string)), parityFactor)
	return parityFactor, nil

}

func insertParityFactorInDB(countryFromISO, countryToISO string, year int, parityFactor float64) {
	dbClient, err := db.Get()
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info("could not reach DB to store parity data for future reference")
		return
	}

	if err := dbClient.InsertParityFactor(countryFromISO, countryToISO, year, parityFactor); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info("could not store parity data in DB for future reference")
		return
	}

}
func getParityFactorFromDB(countryFromISO, countryToISO string, currentYear int) (float64, error) {
	dbClient, err := db.Get()
	if err != nil {
		return 0.0, err
	}

	return dbClient.GetParityFactor(countryFromISO, countryToISO, currentYear)
}

func getCountryDataMapFromParityDataList(countryISO string, countryParityDataList []interface{}) (map[string]interface{}, error) {

	errNoPPPData := "no purchasing power parity data is available for this country"
	if len(countryParityDataList) == 0 {
		log.WithFields(log.Fields{"country_iso": countryISO}).Error(errNoPPPData)
		return nil, errors.New(errNoPPPData)
	}

	errUnexpectedFormat := "received purchasing power parity data in an unexpected format"
	countryDataMap, ok := countryParityDataList[0].(map[string]interface{})
	if !ok {
		log.Error(errUnexpectedFormat)
		return nil, errors.New(errUnexpectedFormat)
	}

	return countryDataMap, nil
}

func getCurrentYear() int {

	return time.Now().Year()
}
