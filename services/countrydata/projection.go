package countrydata

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
)

type Country struct {
	ISO string
}

func ProjectAmount(countryFromISO, countryToISO string, amount int) (int, float64, error) {
	t := template.Must(template.New("projectAmount").Parse(os.Getenv("PARITY_URL")))
	urlBufCountryFrom := &bytes.Buffer{}
	iso := Country{ISO: strings.ToLower(countryFromISO)}

	if err := t.Execute(urlBufCountryFrom, iso); err != nil {
		log.WithFields(log.Fields{"country_iso": countryFromISO, "err": err.Error()}).Error("could not form URL to get purchasing power parity data")
		return 0, 0.0, err
	}

	urlBufCountryTo := &bytes.Buffer{}
	iso = Country{ISO: strings.ToLower(countryToISO)}

	if err := t.Execute(urlBufCountryTo, iso); err != nil {
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

	if len(countryFromParityDataList) == 0 {
		log.WithFields(log.Fields{"country_iso": countryFromISO}).Error("no purchasing power parity data is available for the country from which you want to project an amount")
		return 0, 0.0, err
	}

	if len(countryToParityDataList) == 0 {
		log.WithFields(log.Fields{"country_iso": countryToISO}).Error("no purchasing power parity data is available for the country to which you want to project an amount")
		return 0, 0.0, err
	}

	errUnexpectedFormat := "received purchasing power parity data in an unexpected format"
	countryFromDataMap, ok := countryFromParityDataList[0].(map[string]interface{})
	if !ok {
		log.Error(errUnexpectedFormat)
		return 0, 0.0, err
	}
	countryToDataMap, ok := countryToParityDataList[0].(map[string]interface{})
	if !ok {
		log.Error(errUnexpectedFormat)
		return 0, 0.0, err
	}
	valueKey := "value"

	parityFactor := countryToDataMap[valueKey].(float64) / countryFromDataMap[valueKey].(float64)

	return int(float64(amount) * parityFactor), parityFactor, nil

}
