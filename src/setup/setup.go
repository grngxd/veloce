package setup

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

func SetupAppDir() map[string]interface{} {
	var appdir string
	var err error

	if runtime.GOOS != "windows" {
		sudoUser := os.Getenv("SUDO_USER")
		if sudoUser != "" {
			usr, err := user.Lookup(sudoUser)
			if err != nil {
				panic(err)
			}
			appdir = filepath.Join(usr.HomeDir, ".config", ".veloce")
		} else {
			appdir, err = os.UserConfigDir()
			if err != nil {
				panic(err)
			}
			appdir = filepath.Join(appdir, ".veloce")
		}
	} else {
		appdir, err = os.UserConfigDir()
		if err != nil {
			panic(err)
		}
		appdir = filepath.Join(appdir, ".veloce")
	}

	configDir := filepath.Join(appdir, "config")
	configFile := filepath.Join(configDir, "config.json")

	os.MkdirAll(configDir, fs.ModePerm)

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		defaultConfig := map[string]interface{}{
			"port":   "6134",
			"public": "public",
			"ignore": []string{
				".env",
			},
		}

		configJSON, err := json.Marshal(defaultConfig)
		if err != nil {
			panic(err)
		}

		if err := os.WriteFile(configFile, configJSON, fs.ModePerm); err != nil {
			panic(err)
		}
	}

	fmt.Println("Config file created at " + configFile)

	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	config := map[string]interface{}{}
	if err := json.Unmarshal(configBytes, &config); err != nil {
		panic(err)
	}

	publicDir := filepath.Join(appdir, config["public"].(string))
	os.MkdirAll(publicDir, fs.ModePerm)

	return config
}
