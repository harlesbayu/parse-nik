package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	// fill with nik
	parseNIK := parseNIK("")
	fmt.Println(parseNIK)
}

type ParseKTP struct {
	Province string
	City     string
	Dob      string
	Gender   string
	Serial   string
}

type provinceAndAdress struct {
	province string
	city     string
}

func getProvinceAndCity(provinceId string, cityId string) provinceAndAdress {
	jsonFileProvince, _ := os.Open("./internal/shared/data_kemendagri/province.json")
	defer jsonFileProvince.Close()

	type Province struct {
		Id      string `json:"id"`
		CpsName string `json:"cpsName"`
	}

	type Provinces struct {
		Provinces []Province `json:"Provinces"`
	}

	byteValueProvince, _ := ioutil.ReadAll(jsonFileProvince)
	var provinces Provinces
	json.Unmarshal(byteValueProvince, &provinces)

	var dataProvince string
	for i := 0; i < len(provinces.Provinces); i++ {
		if provinceId == provinces.Provinces[i].Id {
			dataProvince = provinces.Provinces[i].CpsName
			break
		}
	}

	jsonFileCity, _ := os.Open("./internal/shared/data_kemendagri/city.json")
	defer jsonFileCity.Close()

	type City struct {
		Id         string `json:"id"`
		CpsName    string `json:"cpsName"`
		ProvinceId string `json:"provinceId"`
	}

	type Cities struct {
		Cities []City `json:"Cities"`
	}

	byteValueCity, _ := ioutil.ReadAll(jsonFileCity)
	var cities Cities
	json.Unmarshal(byteValueCity, &cities)

	var dataCity string
	for i := 0; i < len(cities.Cities); i++ {
		if cities.Cities[i].Id == cityId && cities.Cities[i].ProvinceId == provinceId {
			dataCity = cities.Cities[i].CpsName
		}
	}

	return provinceAndAdress{
		province: dataProvince,
		city:     dataCity,
	}
}

func parseNIK(nik string) ParseKTP {
	numProvince := nik[0:2]
	numCity := nik[2:4]
	numGender, _ := strconv.ParseInt(nik[6:8], 10, 8)
	numDob := nik[6:12]
	numSerial := nik[12:16]

	getDay, _ := strconv.ParseInt(numDob[0:2], 10, 8)
	getMonth := numDob[2:4]
	getYear := numDob[4:6]
	startYear, _ := strconv.ParseInt(getYear[0:1], 10, 8)

	if startYear > 4 {
		getYear = "19" + getYear
	} else {
		getYear = "20" + getYear
	}

	respGetProvinceAndCity := getProvinceAndCity(numProvince, numCity)

	var newGetDay string
	if getDay >= 40 {
		getDay = getDay - 40
		if getDay < 10 {
			newGetDay = "0" + strconv.FormatInt(getDay, 10)
		} else {
			newGetDay = strconv.FormatInt(getDay, 10)
		}
	} else if getDay < 10 {
		newGetDay = "0" + strconv.FormatInt(getDay, 10)
	}

	gender := "male"
	if numGender >= 40 {
		gender = "female"
	}

	return ParseKTP{
		Province: respGetProvinceAndCity.province,
		City:     respGetProvinceAndCity.city,
		Dob:      fmt.Sprintf("%s-%s-%s", getYear, getMonth, newGetDay),
		Gender:   gender,
		Serial:   numSerial,
	}
}
