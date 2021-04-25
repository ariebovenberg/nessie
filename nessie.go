package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Names struct {
	Short  string
	Medium string
	Long   string
}

type Station struct {
	Names    Names
	Synonyms []string
	UIC      string
}

func cachedStationsRaw() []byte {
	usr, _ := user.Current()
	os.Mkdir(path.Join(usr.HomeDir, ".nessie"), 0755)
	dat, err := ioutil.ReadFile(path.Join(usr.HomeDir, ".nessie/stations.json"))
	if err != nil {
		return nil
	}
	return dat
}

func writeStationsCache(content []byte) {
	usr, _ := user.Current()
	os.Mkdir(path.Join(usr.HomeDir, ".nessie"), 0755)
	err := ioutil.WriteFile(path.Join(usr.HomeDir, ".nessie/stations.json"), content, 0644)
	check(err)
}

func getRequestBody(path string) []byte {
	req, err := http.NewRequest("GET", "https://gateway.apiportal.ns.nl/reisinformatie-api/api/v2"+path, nil)
	check(err)
	req.Header.Add("Ocp-Apim-Subscription-Key", os.Getenv("NESSIE_KEY"))
	resp, err := http.DefaultClient.Do(req)
	check(err)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		check(err)
		return bodyBytes
	}
	panic("Non-200 status code")
}

func fetchRawStations() []byte {
	return getRequestBody("/stations")
}

func getDepartures(uicCode string) []byte {
	return getRequestBody("/departures?uicCode=" + uicCode)
}

func stationsRaw() []byte {
	cached := cachedStationsRaw()
	if cached == nil {
		stations := fetchRawStations()
		writeStationsCache(stations)
		return stations
	}
	return cached
}

func parseStations(bytes []byte) map[string]Station {
	var body map[string]interface{}
	json.Unmarshal(bytes, &body)
	json := body["payload"].([]interface{})
	var stations = make(map[string]Station)
	for _, data := range json {
		aStation := parseOneStation(data.(map[string]interface{}))
		stations[strings.ToLower(aStation.Names.Short)] = aStation
		stations[strings.ToLower(aStation.Names.Medium)] = aStation
		stations[strings.ToLower(aStation.Names.Long)] = aStation
		for _, s := range aStation.Synonyms {
			stations[strings.ToLower(s)] = aStation
		}
	}
	return stations
}

func parseOneStation(json map[string]interface{}) Station {
	synonymsRaw := json["synoniemen"].([]interface{})
	synonyms := make([]string, len(synonymsRaw))
	for i, s := range synonyms {
		synonyms[i] = s
	}
	names := json["namen"].(map[string]interface{})
	return Station{
		Names: Names{
			Short:  names["kort"].(string),
			Medium: names["middel"].(string),
			Long:   names["lang"].(string),
		},
		Synonyms: synonyms,
		UIC:      json["UICCode"].(string)}
}

func main() {
	stationSearchPtr := flag.String("station", "", "The station")
	flag.Parse()

	stationIndex := parseStations(stationsRaw())
	station, prs := stationIndex[strings.ToLower(*stationSearchPtr)]
	if !prs {
		print("Station not found!")
		return
	}
	fmt.Printf("UIC Code: %s\n", station.UIC)
	fmt.Printf("Full name: %s\n", station.Names.Long)
	print("Under construction...")
}
