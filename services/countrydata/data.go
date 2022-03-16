package countrydata

import (
	"errors"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GetISO(country string) (string, error) {

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

			if strings.ToLower(countryName.(string)) == country {
				return countryISOCode.(string), nil
			}
		}
	}

	errNoSuchCountry := "could not validate country"
	log.WithFields(log.Fields{"country_name": country}).Info(errNoSuchCountry)

	return "", errors.New(errNoSuchCountry)
}
