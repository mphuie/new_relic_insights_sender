package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	NewRelicAccountId string `yaml:"NewRelicAccountId"`
	NewRelicInsertKey string `yaml:"NewRelicInsertKey"`
	EventValues map[string]interface{} `yaml:"EventValues"`
}

func main() {

	var sendValueKey string
	if len(os.Args) == 2 {
		sendValueKey = os.Args[1]
	} else {
		sendValueKey = "value"
	}

	configFile, err := ioutil.ReadFile("./nr_insights_sender.yaml")
	if err != nil {
		log.Fatal(err)
	}
	c := Config{}

	err = yaml.Unmarshal(configFile, &c)

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode() & os.ModeNamedPipe == 0 {

	} else {
		reader := bufio.NewReader(os.Stdin)
		pipeData, _ := ioutil.ReadAll(reader)
		val := strings.TrimSpace(string(pipeData))

		i, err := strconv.Atoi(val)

		if err != nil {
			c.EventValues[sendValueKey] = val
		} else {
			c.EventValues[sendValueKey] = i
		}
	}

	data := make([]map[string]interface{}, 0, 0)
	data = append(data, c.EventValues)

	dataString, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(dataString)); err != nil {
		log.Fatal(err)
	}

	_ = gz.Close()


	url := fmt.Sprintf("https://insights-collector.newrelic.com/v1/accounts/%s/events", c.NewRelicAccountId)
	fmt.Println(url)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b.Bytes()))

	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Insert-Key", c.NewRelicInsertKey)

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

	defer resp.Body.Close()



}
