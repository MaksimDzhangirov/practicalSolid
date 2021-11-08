package moreExamples

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

type DataReader func(source string) ([]byte, error)

func ReadFromFile(fileName string) ([]byte, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ReadFromLink(link string) ([]byte, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return data, nil
}

func GetCities(reader DataReader, source string) ([]City, error) {
	data, err := reader(source)
	if err != nil {
		return nil, err
	}

	var cities []City
	err = yaml.Unmarshal(data, &cities)
	if err != nil {
		return nil, err
	}

	return cities, nil
}
