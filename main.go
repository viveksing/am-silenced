package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
)

type Silences []struct {
	SilenceName string    `yaml:"silencename"`
	StartTime   string    `yaml:"starttime"`
	Duration    string    `yaml:"duration"`
	Matchers    []Matcher `yaml:"matchers"`
}

type Matcher struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	IsRegex bool   `json:"isRegex"`
}

type SilencePayload struct {
	Matchers  []Matcher `json:"matchers"`
	StartsAt  string    `json:"startsAt"`
	EndsAt    string    `json:"endsAt"`
	CreatedBy string    `json:"createdBy"`
	Comment   string    `json:"comment"`
}

type ProcessConfig struct {
	URL string `yaml:"url"`
}

func generate(url string, silencename string, duration string, matchers []Matcher) func() {
	return func() {
		stopdur, _ := time.ParseDuration(duration)
		starttime := time.Now().UTC().Format("2006-01-02T15:04:05.070Z")
		stoptime := time.Now().Add(stopdur).UTC().Format("2006-01-02T15:04:05.070Z")
		fmt.Printf("Triggering Silence %s will run from %s for %s duration \n", silencename, starttime, stoptime)

		fmt.Println("URL:>", url)

		var data = SilencePayload{Matchers: matchers}

		data.StartsAt = starttime
		data.EndsAt = stoptime
		data.Comment = "Created by daemon"
		data.CreatedBy = "Silence Daemon"

		var jsonStr, _ = json.Marshal(data)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))

	}
}

func main() {
	fmt.Println("Starting Application")

	processconfig := ProcessConfig{}
	processrf, _ := ioutil.ReadFile("process.yaml")
	yaml.Unmarshal(processrf, &processconfig)
	config := Silences{}

	data, _ := ioutil.ReadFile("config.yaml")

	yaml.Unmarshal(data, &config)
	c := cron.New()

	for _, i := range config {
		c.AddFunc(i.StartTime, generate(processconfig.URL, i.SilenceName, i.Duration, i.Matchers))
	}

	c.Start()

	for {
		time.Sleep(time.Minute * 1)
	}

}
