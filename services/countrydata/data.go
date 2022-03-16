package countrydata

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/meghashyamc/country-incomes/services/httprequest"
	log "github.com/sirupsen/logrus"
)

func GetISO(country string) (string, error) {
	resp, err := httprequest.Make(nil, http.MethodGet, os.Getenv("ISOURL"), nil)
	if err != nil {
		return "", err
	}
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not read response after successfully getting ISO codes")
		return "", err
	}
	responseDataList := []interface{}{}
	if err := json.Unmarshal(respBodyBytes, &responseDataList); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not unmarshal response after successfully getting ISO codes")
		return "", err
	}
	if len(responseDataList) < 2 {
		errDataReceivedLength := "unexpected length of country details received from data source"
		log.Error(errDataReceivedLength)
		return "", errors.New(errDataReceivedLength)
	}

	isoDataList, ok := responseDataList[1].([]interface{})
	if !ok {
		errDataReceivedFormat := "unexpected format of country details received from data source"
		log.WithFields(log.Fields{"data_received": responseDataList[1]}).Error(errDataReceivedFormat)
		return "", errors.New(errDataReceivedFormat)
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
