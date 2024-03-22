package configmodule

import (
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
	endpoint := "http://" + core_host + ":8080" + "/api/buffer-configuration"
	auth := "Bearer " + os.Getenv("CORE_TOKEN")
	// Create a new request using http
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	// add authorization header to the request
	req.Header.Add("Authorization: ", auth)
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
	// log the response body
	fmt.Println(string(body))

	return nil
}
