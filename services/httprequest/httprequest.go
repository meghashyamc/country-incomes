package httprequest

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	clientTimeout   = 30 * time.Second
	maxRetries      = 3
	timeBeforeRetry = 500 * time.Millisecond
)

func Make(headers map[string]string, method, url string, bodyBytes []byte) (*http.Response, error) {
	client := http.Client{Timeout: clientTimeout}
	var (
		retries  int
		response *http.Response
		body     io.Reader
	)
	for retries = 1; retries <= maxRetries; retries++ {
		if bodyBytes != nil {
			body = bytes.NewReader(bodyBytes)
		}

		request, err := http.NewRequest(method, url, body)
		if err != nil {
			log.WithFields(log.Fields{"err": err.Error(), "method": method, "url": url}).Error("failed to create new request")
			return nil, err
		}

		for headerKey, headerVal := range headers {

			request.Header.Set(headerKey, headerVal)
		}

		response, err = client.Do(request)
		if err != nil {
			log.WithFields(log.Fields{"err": err.Error(), "method": method, "url": url}).Error("failed to get a response")
			time.Sleep(timeBeforeRetry)
			continue
		}

		if response.StatusCode != http.StatusOK {
			respBodyBytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.WithFields(log.Fields{"err": err.Error(), "method": method, "url": url}).Error("failed to read gotten response")
				return nil, err
			}
			log.WithFields(log.Fields{"err": err.Error(), "response_code": response.StatusCode, "response_body": string(respBodyBytes), "method": method, "url": url, "retry_number": retries}).Error("failed to get an OK response")

			response.Body.Close()
			time.Sleep(timeBeforeRetry)
			continue
		}

		break
	}

	if retries > maxRetries {
		log.WithFields(log.Fields{"method": method, "url": url}).Error("failed to get an OK response despite multiple retries")

		return nil, errors.New("could not get an OK response despite multiple retries")
	}

	return response, nil
}
