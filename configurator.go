package configmodule

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Config represents your application's configuration as key-value pairs.
type Config map[string]string

var (
	config Config
)

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Data struct {
	Items []Item `json:"item"`
}

type JSONData struct {
	Data Data `json:"data"`
}

func parseJSON(jsonStr string) error {
	var jsonData JSONData

	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return err
	}

	config = make(Config)
	for _, item := range jsonData.Data.Items {
		config[item.Key] = item.Value
	}

	return nil
}

// InitConfig initializes the configuration with the given filePath and creates the file if it doesn't exist.
func InitConfig(file_Path string, defaultConfig Config) error {
	return nil
}

// GetConfig returns the current configuration.
func GetConfig() Config {
	err := readConfigFile()
	if err != nil {
		fmt.Println(err)
	}

	return config
}

func readConfigFile() error {
	core_host := os.Getenv("CORE_HOST")
	endpoint := "http://" + core_host + "/api/buffer-configuration?populate=*"
	auth := "Bearer " + os.Getenv("CORE_TOKEN")
	// Create a new request using http
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	// add content type to the request
	req.Header.Add("Content-Type", "application/json")
	// add authorization header to the request
	req.Header.Add("Authorization", auth)
	// Send http request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	// Close the response body
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// Parse the response body
	err = parseJSON(string(body))
	if err != nil {
		return err
	}

	return nil
}
func HasKey(key string) bool {

	_, exists := config[key]
	return exists
}
