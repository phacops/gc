package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/phacops/garminconnect"
)

const (
	CONFIG_FILE = "${XDG_CONFIG_HOME}/gc/config"
)

var (
	EPO_POST_DATA = []byte{10, 45, 10, 7, 101, 120, 112, 114, 101, 115, 115, 18, 5, 100, 101, 95, 68, 69, 26, 7, 87, 105, 110, 100, 111, 119, 115, 34, 18, 54, 48, 49, 32, 83, 101, 114, 118, 105, 99, 101, 32, 80, 97, 99, 107, 32, 49, 18, 10, 8, 140, 180, 147, 184, 14, 18, 0, 24, 0, 24, 28, 34, 0}

	config Config
)

type Config struct {
	GarminConnectUsername string `json:"gc_username"`
	GarminConnectPassword string `json:"gc_password"`
	WatchDir              string `json:"watch_dir"`
}

func GetEPOFile(c *cli.Context) {
	client := &http.Client{}
	request, _ := http.NewRequest("POST", "http://omt.garmin.com/Rce/ProtobufApi/EphemerisService/GetEphemerisData", bytes.NewBuffer(EPO_POST_DATA))

	request.Header.Set("Garmin-Client-Name", "CoreService")
	request.Header.Set("Content-Type", "application/octet-stream")

	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	var epoBin bytes.Buffer

	for i := 0; i < 28; i++ {
		skip := make([]byte, 3)
		append := make([]byte, 2304)

		response.Body.Read(skip)
		response.Body.Read(append)

		epoBin.Write(append)
	}

	possibleFolders := []string{
		"GARMIN/REMOTESW",
		"GARMIN/GPS",
	}

	for _, folder := range possibleFolders {
		if _, err := os.Stat(filepath.Join(config.WatchDir, folder)); err == nil {
			err = ioutil.WriteFile(filepath.Join(config.WatchDir, folder, "EPO.BIN"), epoBin.Bytes(), 0644)

			if err != nil {
				panic(err)
			}

			break
		}
	}
}

func SyncActivities(c *cli.Context) {
	client, err := garminconnect.NewClient()

	if err != nil {
		panic(err)
	}

	err = client.Auth(config.GarminConnectUsername, config.GarminConnectPassword)

	if err != nil {
		panic(err)
	}

	err = filepath.Walk(filepath.Join(config.WatchDir, "GARMIN/ACTIVITY"), func(path string, info os.FileInfo, _ error) error {
		stat, err := os.Stat(path)

		if err != nil {
			return err
		}

		if !stat.IsDir() {
			fmt.Printf("syncing %s...", filepath.Base(path))
			upload, err := client.UploadActivity(path)

			if err != nil {
				panic(err)
			}

			if len(upload.DetailedImportResult.Successes) > 0 {
				fmt.Printf(" success.\n")
			} else if len(upload.DetailedImportResult.Failures) > 0 {
				fmt.Printf(" failure (%s).\n", upload.DetailedImportResult.Failures[0].Messages[0].Content)
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func SyncWorkouts(c *cli.Context) {
	client, err := garminconnect.NewClient()

	if err != nil {
		panic(err)
	}

	err = client.Auth(config.GarminConnectUsername, config.GarminConnectPassword)

	if err != nil {
		panic(err)
	}

	messages, err := client.Messages()

	if err != nil {
		panic(err)
	}

	for _, message := range messages {
		if message.DeviceXmlDataType == garminconnect.WORKOUT_FILE_TYPE {
			fmt.Printf(`downloading "%s" to the watch...`, message.Metadata.MessageName)

			fitFile, err := os.Create(filepath.Join(config.WatchDir, fmt.Sprintf("GARMIN/WORKOUTS/%d.FIT", message.Metadata.Id)))

			if err != nil {
				panic(err)
			}

			defer fitFile.Close()

			response, err := http.Get(garminconnect.GARMIN_CONNECT_URL + "/" + message.Metadata.MessageUrl)

			if err != nil {
				panic(err)
			}

			defer response.Body.Close()

			_, err = io.Copy(fitFile, response.Body)

			if err != nil {
				panic(err)
			}

			err = client.DeleteMessage(message.Id)

			if err != nil {
				fmt.Sprintf(" Error (%s)\n", err.Error())
			} else {
				fmt.Sprintf(" Success\n")
			}
		}
	}
}

func main() {
	configFile, err := os.Open(os.ExpandEnv(CONFIG_FILE))

	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(configFile).Decode(&config)

	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(config.WatchDir); err != nil {
		panic(errors.New("you must give a valid path for to the watch"))
	}

	app := cli.NewApp()
	app.Name = "gc"
	app.Usage = "Interact with Garmin Connect"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "username, u",
			Value: "",
			Usage: "Garmin Connect username",
		},
		cli.StringFlag{
			Name:  "password, p",
			Value: "",
			Usage: "Garmin Connect password",
		},
		cli.StringFlag{
			Name:  "dir, d",
			Value: "/mnt",
			Usage: "Watch root",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "sync your watch with Garmin Connect",
			Action: func(c *cli.Context) {
				GetEPOFile(c)
				SyncActivities(c)
				SyncWorkouts(c)
			},
			Subcommands: []cli.Command{
				{
					Name:    "epo",
					Aliases: []string{"e"},
					Usage:   "download an updated EPO file",
					Action:  GetEPOFile,
				},
				{
					Name:    "activities",
					Aliases: []string{"a"},
					Usage:   "sync activities",
					Action:  SyncActivities,
				},
				{
					Name:    "workouts",
					Aliases: []string{"w"},
					Usage:   "sync workouts",
					Action:  SyncWorkouts,
				},
			},
		},
	}

	app.Run(os.Args)
}
