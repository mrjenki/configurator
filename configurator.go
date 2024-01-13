package configmodule

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "sync"
    "syscall"
	"fmt"
)

// Config represents your application's configuration as key-value pairs.
type Config map[string]interface{}

var (
    config     Config
    configLock sync.RWMutex
    filePath   string // Store the filePath for configuration updates.
)

// InitConfig initializes the configuration with the given filePath and creates the file if it doesn't exist.
func InitConfig(filePath string, defaultConfig Config) error {
     if err := ensureConfigFile(filePath, defaultConfig); err != nil {
        return err
    }

    // Store the filePath for configuration updates.
    filePath = filePath

    if err := readConfigFile(filePath); err != nil {
        return err
    }
	


    return nil
}


// GetConfig returns the current configuration.
func GetConfig() Config {
	readConfigFile(filePath)
    configLock.RLock()
    defer configLock.RUnlock()
    return config
}

// UpdateConfigToFile updates the configuration in the file.
func UpdateConfigToFile(updatedConfig Config) error {
    configLock.Lock()
    config = updatedConfig
    configLock.Unlock()

    return writeConfigToFile(filePath, updatedConfig)
}

func ensureConfigFile(filePath string, defaultConfig Config) error {
    _, err := os.Stat(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            // Create the file if it doesn't exist and write the default configuration.
            if err := writeDefaultConfig(filePath, defaultConfig); err != nil {
                return err
            }
        } else {
            return err
        }
    }
    return nil
}

func writeDefaultConfig(filePath string, defaultConfig Config) error {
   
    data, err := json.MarshalIndent(defaultConfig, "", "    ")
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(filePath, data, 0644)
    if err != nil {
        return err
    }

    return nil
}

func writeConfigToFile(filePath string, updatedConfig Config) error {
    // Open the file with write access and create it if it doesn't exist.
    file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    // Acquire an exclusive lock on the file for writing.
    if err := fileLock(file); err != nil {
        return err
    }
    defer fileUnlock(file)

    // Serialize the updated configuration.
    data, err := json.MarshalIndent(updatedConfig, "", "    ")
    if err != nil {
        return err
    }

    // Write the updated configuration to the file atomically.
    tmpFile := filePath + ".tmp" // Temporary file.
    if err := ioutil.WriteFile(tmpFile, data, 0644); err != nil {
        return err
    }

    // Rename the temporary file to replace the original file atomically.
    if err := os.Rename(tmpFile, filePath); err != nil {
        // Cleanup the temporary file if the rename fails.
        os.Remove(tmpFile)
        return err
    }

    return nil
}

// fileLock acquires an exclusive lock on the file.
func fileLock(file *os.File) error {
    return syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
}

// fileUnlock releases the file lock.
func fileUnlock(file *os.File) error {
    return syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}
// UpdateKey updates the value of an existing key in the configuration.
func UpdateKey(key string, value interface{}) error {
    configLock.Lock()
    defer configLock.Unlock()

    // Check if the key exists in the configuration.
    if _, exists := config[key]; !exists {
        return fmt.Errorf("key '%s' does not exist in the configuration", key)
    }

    // Update the value of the key.
    config[key] = value

    // Update the configuration in the file.
    return writeConfigToFile(filePath, config)
}
// DeleteKey deletes a key-value pair from the configuration.
func DeleteKey(key string) error {
	configLock.Lock()
	defer configLock.Unlock()

	// Check if the key exists in the configuration.
	if _, exists := config[key]; !exists {
		return fmt.Errorf("key '%s' does not exist in the configuration", key)
	}

	// Delete the key-value pair from the configuration.
	delete(config, key)

	// Update the configuration in the file.
	return writeConfigToFile(filePath, config)
}
//HasKey checks if a key exists in the configuration.
func HasKey(key string) bool {
	configLock.RLock()
	defer configLock.RUnlock()

	_, exists := config[key]
	return exists
}
// AddKey adds a new key-value pair to the configuration.
func AddKey(key string, value interface{}) error {
    configLock.Lock()
    defer configLock.Unlock()

    // Check if the key already exists in the configuration.
    if _, exists := config[key]; exists {
        return fmt.Errorf("key '%s' already exists in the configuration", key)
    }

    // Add the new key-value pair to the configuration.
    config[key] = value

    // Update the configuration in the file.
    return writeConfigToFile(filePath, config)
}
func readConfigFile(filePath string) error {
    // check if the file exists
	_, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	// open the file with read-only access
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// acquire a shared lock on the file for reading
	if err := fileLock(file); err != nil {
		return err
	}
	defer fileUnlock(file)
	// read the file content to variable data
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	// unmarshal the data to a temporary configuration variable
	var tempConfig Config
	if err := json.Unmarshal(data, &tempConfig); err != nil {
		return err
	}
	//write data to console
	// fmt.Println(string(tempConfig))
    

    // Update the in-memory configuration in a thread-safe way.
    configLock.Lock()
    config = tempConfig
    configLock.Unlock()

    return nil
}