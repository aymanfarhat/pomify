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

package internal

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/urfave/cli/v2"
)

func ScanJars(c *cli.Context) error {
	jarsPath := c.String("jars")
	outputPath := c.String("output")

	report := make([]ReportRow, 0)

	jars, err := listJarFiles(jarsPath)

	if (err != nil) {
		fmt.Printf("Error listing jar files: %s\n", err)
		return err
	}

	for _, jar := range jars {
		jarFilename := getFilename(jar)
		checksum, err := getFileSha1Checksum(jar)
		if err != nil {
			fmt.Printf("Error calculating checksum for %s: %s\n", jar, err)
			continue
		}

		dep, err := searchMaven(checksum)
		if err != nil {
			fmt.Printf("Error searching Maven Central for %s: %s\n", jar, err)
			continue
		}

		if dep != nil {
			fmt.Printf("\u2713 Found %s on Maven central\n", jar)
			report = append(report, ReportRow{
				JarFilename:    jarFilename,
				GroupId:        dep.GroupId,
				ArtifactId:     dep.ArtifactId,
				Version:        dep.Version,
				OnMavenCentral: true,
				FileChecksum:   checksum,
				LocalFilepath: jar,
			})
		} else {
			// For now, just add the jar to the report with default values
			// TODO: Drill down into the jar metadata to get more accurate values
			fmt.Printf("\u2715 Couldn't find %s on Maven cental\n", jar)
			report = append(report, ReportRow{
				JarFilename:    jarFilename,
				GroupId: jarFilename,
				ArtifactId: jarFilename,
				Version: "1.0.0",
				OnMavenCentral: false,
				FileChecksum:   checksum,
				LocalFilepath: jar,
			})
		}
	}

	reportFile, err := os.OpenFile(outputPath + "/pomify-report.csv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Printf("Error opening report file: %s\n", err)
		return err
	}

	defer reportFile.Close()

	err = gocsv.MarshalFile(&report, reportFile)
	if err != nil {
		fmt.Printf("Error writing report to CSV: %s\n", err)
		return err
	}

	return nil
}

func loadReportFile(reportPath string) ([]ReportRow, error) {
	reportFile, err := os.Open(reportPath)
	if err != nil {
		return nil, fmt.Errorf("error opening report file: %s", err)
	}

	defer reportFile.Close()

	report := make([]ReportRow, 0)

	if err := gocsv.UnmarshalFile(reportFile, &report); err != nil {
		return nil, fmt.Errorf("error reading report file: %s", err)
	}

	return report, nil
}

func GenDepXml(c *cli.Context) error {
	reportPath := c.String("report")
	outputPath := c.String("output")

	report, err := loadReportFile(reportPath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	mvnCentralDeps := make([]Dependency, 0)
	privateDeps := make([]Dependency, 0)

	for _, row := range report {
		if row.OnMavenCentral {
			mvnCentralDeps = append(mvnCentralDeps, Dependency{
				GroupId:    row.GroupId,
				ArtifactId: row.ArtifactId,
				Version:    row.Version,
			})
		} else {
			privateDeps = append(privateDeps, Dependency{
				GroupId:    row.GroupId,
				ArtifactId: row.ArtifactId,
				Version:    row.Version,
			})
		}
	}

	err = writeDependenciesToXML(mvnCentralDeps, outputPath + "/maven-central-deps.xml")
	if err != nil {
		fmt.Printf("Error writing Maven Central dependencies to XML: %s\n", err)
		return err
	}

	err = writeDependenciesToXML(privateDeps, outputPath + "/private-deps.xml")
	if err != nil {
		fmt.Printf("Error writing unknown dependencies to XML: %s\n", err)
		return err
	}

	return nil
}

func PushJars(c *cli.Context) error {
	reportPath := c.String("report")
	mvnRepo := c.String("mvn_repo")
	report, err := loadReportFile(reportPath)

	if err != nil {
		fmt.Println(err)
		return err
	}

	importCommands := []string{}

	for _, row := range report {
		if !row.OnMavenCentral {
			importCommand := fmt.Sprintf(
                `mvn deploy:deploy-file \
                    -Dfile="%s" \
                    -DgroupId="%s" \
                    -DartifactId="%s" \
                    -Dversion="%s" \
                    -Dpackaging="jar" \
                    -DgeneratePom=true \
                    -Durl="%s" \
                    -DcreateChecksum=true`, row.LocalFilepath, row.GroupId, row.ArtifactId, row.Version, mvnRepo)

            importCommands = append(importCommands, importCommand)
		}
	}

	writeStringsToFile(importCommands, "output/import-commands.sh")

	return nil
}