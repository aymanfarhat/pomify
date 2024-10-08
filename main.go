/*
* Copyright 2024 Google LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* 	https://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package main

import (
	"log"
	"os"

	"github.com/google/pomify-jars/internal"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "pomify",
		Usage:                "Create Java pom.xml dependency definitions from a list of unknown jar files",
		Version:              "0.0.1",
		EnableBashCompletion: true,
		Suggest:              true,
		Commands: []*cli.Command{
			{
				Name:    "scan",
				Usage:   "Scan a directory of jar files and output a CSV report",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "jars",
						Aliases:  []string{"j"},
						Usage:    "Path to the directory containing the jar files",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "Path to the output CSV file",
						Required: true,
					},
				},
				Action: internal.ScanJars,
			},
			{
				Name:    "generate",
				Usage:   "Generate a pom.xml dependencies block from a CSV report for each jar category",
				Aliases: []string{"g"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "report",
						Aliases:  []string{"r"},
						Usage:    "Path to the Pomify CSV report file",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "Path to the output directory of XML files",

						Required: true,
					},
				},
				Action: internal.GenDepXml,
			},
			{
				Name:    "push",
				Usage:   "Push the private jars to a custom Maven repository",
				Aliases: []string{"p"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "report",
						Aliases:  []string{"r"},
						Usage:    "Path to the Pomify CSV report file",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "mvn_repo",
						Aliases:  []string{"m"},
						Usage:    "URL of the custom Maven repository",
						Required: true,
					},
				},
				Action: internal.PushJars,
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
