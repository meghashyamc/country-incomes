package countrydata

import (
	"bytes"
	"errors"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
)

type PPPDetails struct {
	ISO         string
	CurrentYear string
}

func ProjectAmount(countryFromISO, countryToISO string, amount int) (int, float64, error) {
	t := template.Must(template.New("projectAmount").Parse(os.Getenv("PARITY_URL")))
	urlBufCountryFrom := &bytes.Buffer{}
	pppDetails := PPPDetails{ISO: strings.ToLower(countryFromISO), CurrentYear: getCurrentYearAsStr()}

	if err := t.Execute(urlBufCountryFrom, pppDetails); err != nil {
		log.WithFields(log.Fields{"country_iso": countryFromISO, "err": err.Error()}).Error("could not form URL to get purchasing power parity data")
		return 0, 0.0, err
	}

	urlBufCountryTo := &bytes.Buffer{}
	pppDetails.ISO = strings.ToLower(countryToISO)

	if err := t.Execute(urlBufCountryTo, pppDetails); err != nil {
		log.WithFields(log.Fields{"country_iso": countryFromISO, "err": err.Error()}).Error("could not form URL to get purchasing power parity data")
		return 0, 0.0, err
	}

	countryFromParityDataList, err := makeRequestAndGetData(urlBufCountryFrom.String())
	if err != nil {
		return 0, 0.0, err
	}

	countryToParityDataList, err := makeRequestAndGetData(urlBufCountryTo.String())
	if err != nil {
		return 0, 0.0, err
	}

	countryFromDataMap, countryToDataMap, err := getCountryDataMapFromParityDataList(countryFromISO, countryToISO, countryFromParityDataList, countryToParityDataList)
	if err != nil {
		return 0, 0.0, err
	}
	valueKey := "value"

	parityFactor := countryToDataMap[valueKey].(float64) / countryFromDataMap[valueKey].(float64)

	return int(float64(amount) * parityFactor), parityFactor, nil

}

func getCountryDataMapFromParityDataList(countryFromISO, countryToISO string, countryFromParityDataList, countryToParityDataList []interface{}) (map[string]interface{}, map[string]interface{}, error) {

	errNoPPPData := "no purchasing power parity data is available for this country"
	if len(countryFromParityDataList) == 0 {
		log.WithFields(log.Fields{"country_iso": countryFromISO}).Error(errNoPPPData)
		return nil, nil, errors.New(errNoPPPData)
	}

	if len(countryToParityDataList) == 0 {
		log.WithFields(log.Fields{"country_iso": countryToISO}).Error(errNoPPPData)
		return nil, nil, errors.New(errNoPPPData)
	}

	errUnexpectedFormat := "received purchasing power parity data in an unexpected format"
	countryFromDataMap, ok := countryFromParityDataList[0].(map[string]interface{})
	if !ok {
		log.Error(errUnexpectedFormat)
		return nil, nil, errors.New(errUnexpectedFormat)
	}
	countryToDataMap, ok := countryToParityDataList[0].(map[string]interface{})
	if !ok {
		log.Error(errUnexpectedFormat)
		return nil, nil, errors.New(errUnexpectedFormat)
	}
	return countryFromDataMap, countryToDataMap, nil
}

func getCurrentYearAsStr() string {

	return strconv.Itoa(time.Now().Year())
}
