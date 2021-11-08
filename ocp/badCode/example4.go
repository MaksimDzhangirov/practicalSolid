package badCode

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

type City struct {
	name      string
	latitude  float64
	longitude float64
	country   string
}

func GetCities(sourceType string, source string) ([]City, error) {
	var data []byte
	var err error

	if sourceType == "file" {
		data, err = ioutil.ReadFile(source)
		if err != nil {
			return nil, err
		}
	} else if sourceType == "link" {
		resp, err := http.Get(source)
		if err != nil {
			return nil, err
		}

		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}

	var cities []City
	err = yaml.Unmarshal(data, &cities)
	if err != nil {
		return nil, err
	}

	return cities, nil
}
