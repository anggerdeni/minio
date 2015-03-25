/*
 * Minimalist Object Storage, (C) 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"os"

	"github.com/minio-io/cli"
	"github.com/minio-io/minio/pkg/server"
	"github.com/minio-io/minio/pkg/utils/log"
)

// commitID is automatically set by git. Settings are controlled
// through .gitattributes
const commitID = "$Id$"

var flags = []cli.Flag{
	cli.StringFlag{
		Name:  "domain,d",
		Value: "",
		Usage: "domain used for routing incoming API requests",
	},
	cli.StringFlag{
		Name:  "api-address,a",
		Value: ":9000",
		Usage: "address for incoming API requests",
	},
	cli.StringFlag{
		Name:  "web-address,w",
		Value: ":9001",
		Usage: "address for incoming Management UI requests",
	},
	cli.StringFlag{
		Name:  "cert,c",
		Hide:  true,
		Value: "",
		Usage: "cert.pem",
	},
	cli.StringFlag{
		Name:  "key,k",
		Hide:  true,
		Value: "",
		Usage: "key.pem",
	},
	cli.StringFlag{
		Name:  "driver-type,t",
		Value: "donut",
		Usage: "valid entries: file,inmemory,donut",
	},
}

func getDriverType(input string) server.DriverType {
	switch {
	case input == "file":
		return server.File
	case input == "memory":
		return server.Memory
	case input == "donut":
		return server.Donut
	default:
		{
			log.Println("Unknown driver type:", input)
			log.Println("Choosing default driver type as 'file'..")
			return server.File
		}
	}
}

func runCmd(c *cli.Context) {
	driverTypeStr := c.String("driver-type")
	domain := c.String("domain")
	apiaddress := c.String("api-address")
	webaddress := c.String("web-address")
	certFile := c.String("cert")
	keyFile := c.String("key")
	if (certFile != "" && keyFile == "") || (certFile == "" && keyFile != "") {
		log.Fatal("Both certificate and key must be provided to enable https")
	}
	tls := (certFile != "" && keyFile != "")
	driverType := getDriverType(driverTypeStr)
	var serverConfigs []server.Config
	apiServerConfig := server.Config{
		Domain:   domain,
		Address:  apiaddress,
		TLS:      tls,
		CertFile: certFile,
		KeyFile:  keyFile,
		APIType: server.MinioAPI{
			DriverType: driverType,
		},
	}
	webUIServerConfig := server.Config{
		Domain:   domain,
		Address:  webaddress,
		TLS:      false,
		CertFile: "",
		KeyFile:  "",
		APIType: server.Web{
			Websocket: false,
		},
	}
	serverConfigs = append(serverConfigs, apiServerConfig)
	serverConfigs = append(serverConfigs, webUIServerConfig)
	server.Start(serverConfigs)
}

func main() {
	app := cli.NewApp()
	app.Name = "minio"
	app.Version = "0.1.0"
	app.Author = "Minio.io"
	app.Usage = "Minimalist Object Storage"
	app.EnableBashCompletion = true
	app.Flags = flags
	app.Action = runCmd
	app.Run(os.Args)
}
