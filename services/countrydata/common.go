package countrydata

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/meghashyamc/country-incomes/services/httprequest"
	log "github.com/sirupsen/logrus"
)

func makeRequestAndGetData(url string) ([]interface{}, error) {
	resp, err := httprequest.Make(nil, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not read response after successfully getting ISO codes")
		return nil, err
	}
	responseDataList := []interface{}{}
	if err := json.Unmarshal(respBodyBytes, &responseDataList); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not unmarshal response after successfully getting ISO codes")
		return nil, err
	}
	if len(responseDataList) < 2 {
		errDataReceivedLength := "unexpected length of country details received from data source"
		log.WithFields(log.Fields{"length": len(responseDataList), "response_data_list": responseDataList, "url": url}).Error(errDataReceivedLength)
		return nil, errors.New(errDataReceivedLength)
	}

	actualDataList, ok := responseDataList[1].([]interface{})
	if !ok {
		errDataReceivedFormat := "unexpected format of country details received from data source"
		log.WithFields(log.Fields{"data_received": responseDataList[1]}).Error(errDataReceivedFormat)
		return nil, errors.New(errDataReceivedFormat)
	}

	return actualDataList, nil
}
