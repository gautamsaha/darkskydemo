package main

import (
	"fmt"
	//s"github.com/mlbright/darksky"
	"encoding/json"
	"io"
	// "io/ioutil"
	"net/http"
	"strconv"
	"time"
	// "html/template"
	// "net/url"
	//"net"
	"math"
	"strings"
)

const (
	//
	BASEURL = "https://api.darksky.net/forecast"
)

type Flags struct {
	DarkSkyUnavailable string   `json:"darksky-unavailable,omitempty"`
	DarkSkyStations    []string `json:"darksky-stations,omitempty"`
	DataPointStations  []string `json:"datapoint-stations,omitempty"`
	ISDStations        []string `json:"isds-stations,omitempty"`
	LAMPStations       []string `json:"lamp-stations,omitempty"`
	MADISStations      []string `json:"madis-stations,omitempty"`
	METARStations      []string `json:"metars-stations,omitempty"`
	METNOLicense       string   `json:"metnol-license,omitempty"`
	Sources            []string `json:"sources,omitempty"`
	Units              string   `json:"units,omitempty"`
}

type DataPoint struct {
	Time                       int64   `json:"time,omitempty"`
	Summary                    string  `json:"summary,omitempty"`
	Icon                       string  `json:"icon,omitempty"`
	SunriseTime                int64   `json:"sunriseTime,omitempty"`
	SunsetTime                 int64   `json:"sunsetTime,omitempty"`
	PrecipIntensity            float64 `json:"precipIntensity,omitempty"`
	PrecipIntensityMax         float64 `json:"precipIntensityMax,omitempty"`
	PrecipIntensityMaxTime     int64   `json:"precipIntensityMaxTime,omitempty"`
	PrecipProbability          float64 `json:"precipProbability,omitempty"`
	PrecipType                 string  `json:"precipType,omitempty"`
	PrecipAccumulation         float64 `json:"precipAccumulation,omitempty"`
	Temperature                float64 `json:"temperature,omitempty"`
	TemperatureMin             float64 `json:"temperatureMin,omitempty"`
	TemperatureMinTime         int64   `json:"temperatureMinTime,omitempty"`
	TemperatureMax             float64 `json:"temperatureMax,omitempty"`
	TemperatureMaxTime         int64   `json:"temperatureMaxTime,omitempty"`
	ApparentTemperature        float64 `json:"apparentTemperature,omitempty"`
	ApparentTemperatureMin     float64 `json:"apparentTemperatureMin,omitempty"`
	ApparentTemperatureMinTime int64   `json:"apparentTemperatureMinTime,omitempty"`
	ApparentTemperatureMax     float64 `json:"apparentTemperatureMax,omitempty"`
	ApparentTemperatureMaxTime int64   `json:"apparentTemperatureMaxTime,omitempty"`
	NearestStormBearing        float64 `json:"nearestStormBearing,omitempty"`
	NearestStormDistance       float64 `json:"nearestStormDistance,omitempty"`
	DewPoint                   float64 `json:"dewPoint,omitempty"`
	WindSpeed                  float64 `json:"windSpeed,omitempty"`
	WindBearing                float64 `json:"windBearing,omitempty"`
	CloudCover                 float64 `json:"cloudCover,omitempty"`
	Humidity                   float64 `json:"humidity,omitempty"`
	Pressure                   float64 `json:"pressure,omitempty"`
	Visibility                 float64 `json:"visibility,omitempty"`
	Ozone                      float64 `json:"ozone,omitempty"`
	MoonPhase                  float64 `json:"moonPhase,omitempty"`
}

type DataBlock struct {
	Summary string      `json:"summary,omitempty"`
	Icon    string      `json:"icon,omitempty"`
	Data    []DataPoint `json:"data,omitempty"`
}

type alert struct {
	Title       string   `json:"title,omitempty"`
	Regions     []string `json:"regions,omitempty"`
	Severity    string   `json:"severity,omitempty"`
	Description string   `json:"description,omitempty"`
	Time        int64    `json:"time,omitempty"`
	Expires     float64  `json:"expires,omitempty"`
	URI         string   `json:"uri,omitempty"`
}

type Forecast struct {
	Latitude  float64   `json:"latitude,omitempty"`
	Longitude float64   `json:"longitude,omitempty"`
	Timezone  string    `json:"timezone,omitempty"`
	Offset    float64   `json:"offset,omitempty"`
	Currently DataPoint `json:"currently,omitempty"`
	Minutely  DataBlock `json:"minutely,omitempty"`
	Hourly    DataBlock `json:"hourly,omitempty"`
	Daily     DataBlock `json:"daily,omitempty"`
	Alerts    []alert   `json:"alerts,omitempty"`
	Flags     Flags     `json:"flags,omitempty"`
	APICalls  int       `json:"apicalls,omitempty"`
	Code      int       `json:"code,omitempty"`
}

type DemoData struct {
	NumberOfDays int        `json:"numberofdays,omitempty"`
	Days         []Forecast `json:"days,omitempty"`
}

type Units string

const (
	CA   Units = "ca"
	SI   Units = "si"
	US   Units = "us"
	UK   Units = "uk"
	AUTO Units = "auto"
)

type Lang string

const (
	Arabic             Lang = "ar"
	Azerbaijani        Lang = "az"
	Belarusian         Lang = "be"
	Bosnian            Lang = "bs"
	Catalan            Lang = "ca"
	Czech              Lang = "cs"
	German             Lang = "de"
	Greek              Lang = "el"
	English            Lang = "en"
	Spanish            Lang = "es"
	Estonian           Lang = "et"
	French             Lang = "fr"
	Croatian           Lang = "hr"
	Hungarian          Lang = "hu"
	Indonesian         Lang = "id"
	Italian            Lang = "it"
	Icelandic          Lang = "is"
	Cornish            Lang = "kw"
	Indonesia          Lang = "nb"
	Dutch              Lang = "nl"
	Polish             Lang = "pl"
	Portuguese         Lang = "pt"
	Russian            Lang = "ru"
	Slovak             Lang = "sk"
	Slovenian          Lang = "sl"
	Serbian            Lang = "sr"
	Swedish            Lang = "sv"
	Tetum              Lang = "te"
	Turkish            Lang = "tr"
	Ukrainian          Lang = "uk"
	IgpayAtinlay       Lang = "x-pig-latin"
	SimplifiedChinese  Lang = "zh"
	TraditionalChinese Lang = "zh-tw"
)

func Get(key string, lat string, long string, time string, units Units, lang Lang) (*Forecast, error) {
	res, err := GetResponse(key, lat, long, time, units, lang)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	f, err := FromJSON(res.Body)
	if err != nil {
		return nil, err
	}

	calls, _ := strconv.Atoi(res.Header.Get("X-Forecast-API-Calls"))
	f.APICalls = calls

	return f, nil
}

func FromJSON(reader io.Reader) (*Forecast, error) {
	var f Forecast
	if err := json.NewDecoder(reader).Decode(&f); err != nil {
		return nil, err
	}

	return &f, nil
}

func GetResponse(key string, lat string, long string, time string, units Units, lang Lang) (*http.Response, error) {
	coord := lat + "," + long

	var url string
	if time == "now" {
		url = BASEURL + "/" + key + "/" + coord + "?units=" + string(units) + "&lang=" + string(lang)
	} else {
		url = BASEURL + "/" + key + "/" + coord + "," + time + "?units=" + string(units) + "&lang=" + string(lang) + "&exclude=currently,minutely,hourly"
	}

	res, err := http.Get(url)
	if err != nil {
		return res, err
	}

	return res, nil
}

func handler(w http.ResponseWriter, r *http.Request) {

	lat := ""
	lng := ""
	qp, ok := r.URL.Query()["loc"]
	if !ok || len(qp) < 1 {
		fmt.Println("NO location")
	} else {
		loc := string(qp[0])
		fmt.Printf(" Loc=%s", string(loc))

		points := strings.Split(loc, ",")
		if len(points) == 2 {
			lat = points[0]
			lng = points[1]
		}
	}

	if ("" == lat) || ("" == lng) {
		lat = "40.7127"
		lng = "-74.0059"
	}

	// -------API params
	key := ""
	// lat := "40.7127"
	// lng := "-74.0059"

	// htmlout = createMockForecastInMem()s
	var htmlout DemoData
	htmlout.NumberOfDays = 7
	for i := 1; i <= htmlout.NumberOfDays; i++ {
		// htmlout.Days = append(htmlout.Days,createMockForecastWithDay(i))

		time_fetch := strconv.FormatInt(time.Now().AddDate(0, 0, 0-i).Unix(), 10)
		response, err := Get(key, lat, lng, time_fetch, CA, English)
		if err != nil {
			fmt.Printf("%s", err)
			// os.Exit(1)
		} else {
			htmlout.Days = append(htmlout.Days, *response)
		}
	}

	///----- Hand code HTML
	// fmt.Fprint(w,`<head><title>Dark Sky Demo</title><script src=\"timeline.js\"></script><link rel=\"stylesheet\" href=\"styles.min.css?version=1.0.7/></head>`)
	fmt.Fprintf(w, "<h1> Dark Sky Timeline : Past %s Days for location(%s.%s)</h1>", strconv.Itoa(htmlout.NumberOfDays), lat, lng)

	for _, _day := range htmlout.Days {
		fmt.Printf("%s\n", _day.Daily.Data[0].Summary)

		tm := time.Unix(_day.Daily.Data[0].Time, 0)
		tmpMin := strconv.FormatFloat(round(_day.Daily.Data[0].TemperatureMin), 'f', -1, 64)
		tmpMax := strconv.FormatFloat(round(_day.Daily.Data[0].TemperatureMax), 'f', -1, 64)

		fmt.Fprintf(w, "<li>  %s-%s-%s : Temp Min=%s, Temp Max=%s,Summary:%s </li>", strconv.Itoa(tm.Year()), tm.Month(), strconv.Itoa(tm.Day()), tmpMin, tmpMax, _day.Daily.Data[0].Summary)
	}
	///- Hand coded html

	// ///-------- Template HTML
	// tmpl := template.Must(template.ParseFiles("layout.html"))
	// tmpl.Execute(w, createMockForecastInMem())
	// ///-------- Template HTML

	// response := ""
	// 	response :=

	//     contents, err := json.Marshal(response)
	//     if err != nil {
	//         fmt.Printf("%s", err)
	//         // os.Exit(1)
	//     }
	// 	fmt.Printf("%s", string(time))
	// 	// fmt.Printf("%+v",contents)
	// 	fmt.Fprintf(w,"%s", string(contents))
	// // }
}

func main() {
	// http.HandleFunc("/v1", handler)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func createMockForecastWithDay(i int) Forecast {
	daynumber := fmt.Sprintf("Day%s", strconv.Itoa(i))
	fmt.Printf("createMockForecastWithDay: Day # is %s\n", daynumber)

	return Forecast{
		Latitude:  40.7127,
		Longitude: -74.0059,
		Timezone:  "America/New_York",
		Offset:    -5,
		Daily: DataBlock{
			Data: []DataPoint{
				DataPoint{
					Time:                       1518238800,
					Summary:                    daynumber,
					Icon:                       "rain",
					SunriseTime:                1518263837,
					SunsetTime:                 1518301555,
					PrecipIntensity:            0.3937,
					PrecipIntensityMax:         2.2809,
					PrecipIntensityMaxTime:     1518310800,
					PrecipProbability:          0.89,
					PrecipType:                 "rain",
					TemperatureMin:             2.19,
					TemperatureMinTime:         1518238800,
					TemperatureMax:             7.84,
					TemperatureMaxTime:         1518282000,
					ApparentTemperatureMin:     2.19,
					ApparentTemperatureMinTime: 1518238800,
					ApparentTemperatureMax:     7.84,
					ApparentTemperatureMaxTime: 1518282000,
					DewPoint:                   3.64,
					WindSpeed:                  0.48,
					WindBearing:                226,
					CloudCover:                 0.63,
					Humidity:                   0.89,
					Pressure:                   1023.3,
					Visibility:                 10.44,
					Ozone:                      323.17,
					MoonPhase:                  0.84,
				}, //DataPoint instance
			}, //Data
		}, //Daily
	}
}

func createMockForecastInMem() DemoData {

	return DemoData{
		NumberOfDays: 7,
		Days: []Forecast{
			Forecast{
				Latitude:  40.7127,
				Longitude: -74.0059,
				Timezone:  "America/New_York",
				Offset:    -5,
				Daily: DataBlock{
					Data: []DataPoint{
						DataPoint{
							Time:                       1518238800,
							Summary:                    "Day1",
							Icon:                       "rain",
							SunriseTime:                1518263837,
							SunsetTime:                 1518301555,
							PrecipIntensity:            0.3937,
							PrecipIntensityMax:         2.2809,
							PrecipIntensityMaxTime:     1518310800,
							PrecipProbability:          0.89,
							PrecipType:                 "rain",
							TemperatureMin:             2.19,
							TemperatureMinTime:         1518238800,
							TemperatureMax:             7.84,
							TemperatureMaxTime:         1518282000,
							ApparentTemperatureMin:     2.19,
							ApparentTemperatureMinTime: 1518238800,
							ApparentTemperatureMax:     7.84,
							ApparentTemperatureMaxTime: 1518282000,
							DewPoint:                   3.64,
							WindSpeed:                  0.48,
							WindBearing:                226,
							CloudCover:                 0.63,
							Humidity:                   0.89,
							Pressure:                   1023.3,
							Visibility:                 10.44,
							Ozone:                      323.17,
							MoonPhase:                  0.84,
						}, //DataPoint instance
					}, //Data
				}, //Daily
			}, //Forecast instance
			Forecast{
				Latitude:  40.7127,
				Longitude: -74.0059,
				Timezone:  "America/New_York",
				Offset:    -5,
				Daily: DataBlock{
					Data: []DataPoint{
						DataPoint{
							Time:                       1518238800,
							Summary:                    "Day2",
							Icon:                       "rain",
							SunriseTime:                1518263837,
							SunsetTime:                 1518301555,
							PrecipIntensity:            0.3937,
							PrecipIntensityMax:         2.2809,
							PrecipIntensityMaxTime:     1518310800,
							PrecipProbability:          0.89,
							PrecipType:                 "rain",
							TemperatureMin:             2.19,
							TemperatureMinTime:         1518238800,
							TemperatureMax:             7.84,
							TemperatureMaxTime:         1518282000,
							ApparentTemperatureMin:     2.19,
							ApparentTemperatureMinTime: 1518238800,
							ApparentTemperatureMax:     7.84,
							ApparentTemperatureMaxTime: 1518282000,
							DewPoint:                   3.64,
							WindSpeed:                  0.48,
							WindBearing:                226,
							CloudCover:                 0.63,
							Humidity:                   0.89,
							Pressure:                   1023.3,
							Visibility:                 10.44,
							Ozone:                      323.17,
							MoonPhase:                  0.84,
						}, //DataPoint instance
					}, //Data
				}, //Daily
			}, //Forecast instance
			Forecast{
				Latitude:  40.7127,
				Longitude: -74.0059,
				Timezone:  "America/New_York",
				Offset:    -5,
				Daily: DataBlock{
					Data: []DataPoint{
						DataPoint{
							Time:                       1518238800,
							Summary:                    "Day3",
							Icon:                       "rain",
							SunriseTime:                1518263837,
							SunsetTime:                 1518301555,
							PrecipIntensity:            0.3937,
							PrecipIntensityMax:         2.2809,
							PrecipIntensityMaxTime:     1518310800,
							PrecipProbability:          0.89,
							PrecipType:                 "rain",
							TemperatureMin:             2.19,
							TemperatureMinTime:         1518238800,
							TemperatureMax:             7.84,
							TemperatureMaxTime:         1518282000,
							ApparentTemperatureMin:     2.19,
							ApparentTemperatureMinTime: 1518238800,
							ApparentTemperatureMax:     7.84,
							ApparentTemperatureMaxTime: 1518282000,
							DewPoint:                   3.64,
							WindSpeed:                  0.48,
							WindBearing:                226,
							CloudCover:                 0.63,
							Humidity:                   0.89,
							Pressure:                   1023.3,
							Visibility:                 10.44,
							Ozone:                      323.17,
							MoonPhase:                  0.84,
						}, //DataPoint instance
					}, //Data
				}, //Daily
			}, //Forecast instance
		}, //Days
	}
}

func round(input float64) float64 {
	if math.IsNaN(input) {
		return math.NaN()
	}
	sign := 1.0
	if input < 0 {
		sign = -1
		input *= -1
	}
	_, decimal := math.Modf(input)
	var rounded float64
	if decimal >= 0.5 {
		rounded = math.Ceil(input)
	} else {
		rounded = math.Floor(input)
	}
	return rounded * sign
}

// func createMockForcastFromJSONString() {
// 	return `
// 	{
// 		"numberofdays": 7,
// 		"days": [
// 			{
// 				"latitude": 40.7127,
// 				"longitude": -74.0059,
// 				"timezone": "America/New_York",
// 				"offset": -5,
// 				"daily": {
// 					"data": [
// 						{
// 							"time": 1518238800,
// 							"summary": "Rain in the afternoon and overnight.",
// 							"icon": "rain",
// 							"sunriseTime": 1518263837,
// 							"sunsetTime": 1518301555,
// 							"precipIntensity": 0.3937,
// 							"precipIntensityMax": 2.2809,
// 							"precipIntensityMaxTime": 1518310800,
// 							"precipProbability": 0.89,
// 							"precipType": "rain",
// 							"temperatureMin": 2.19,
// 							"temperatureMinTime": 1518238800,
// 							"temperatureMax": 7.84,
// 							"temperatureMaxTime": 1518282000,
// 							"apparentTemperatureMin": 2.19,
// 							"apparentTemperatureMinTime": 1518238800,
// 							"apparentTemperatureMax": 7.84,
// 							"apparentTemperatureMaxTime": 1518282000,
// 							"dewPoint": 3.64,
// 							"windSpeed": 0.48,
// 							"windBearing": 226,
// 							"cloudCover": 0.63,
// 							"humidity": 0.89,
// 							"pressure": 1023.3,
// 							"visibility": 10.44,
// 							"ozone": 323.17,
// 							"moonPhase": 0.84
// 						}
// 					]
// 				}
// 			}
// 		]
// 	}`

// }
