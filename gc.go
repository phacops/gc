package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

var (
	EPO_POST_DATA = []byte{10, 45, 10, 7, 101, 120, 112, 114, 101, 115, 115, 18, 5, 100, 101, 95, 68, 69, 26, 7, 87, 105, 110, 100, 111, 119, 115, 34, 18, 54, 48, 49, 32, 83, 101, 114, 118, 105, 99, 101, 32, 80, 97, 99, 107, 32, 49, 18, 10, 8, 140, 180, 147, 184, 14, 18, 0, 24, 0, 24, 28, 34, 0}
)

func GetEPOFile(c *cli.Context) {
	if c.NArg() == 0 {
		panic(errors.New("you need to give the path of the watch"))
	}

	watchPath := c.Args()[0]
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
		if _, err := os.Stat(filepath.Join(watchPath, folder)); err == nil {
			err = ioutil.WriteFile(filepath.Join(watchPath, folder, "EPO.BIN"), epoBin.Bytes(), 0644)

			if err != nil {
				panic(err)
			}

			break
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "gc"
	app.Usage = "Interact with Garmin Connect"
	app.Commands = []cli.Command{
		{
			Name:    "epo",
			Aliases: []string{"e"},
			Usage:   "download an updated EPO file",
			Action:  GetEPOFile,
		},
	}

	app.Run(os.Args)
}
