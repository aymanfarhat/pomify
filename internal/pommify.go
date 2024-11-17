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
	"text/template"

	"github.com/fatih/color"
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
		}

		dep, err := searchMaven(checksum)
		if err != nil {
			fmt.Printf("Error searching Maven Central for %s: %s\n", jar, err)
		}

		if dep != nil {
			color.Green("Found %s on Maven central", jar)
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
			color.Red("Couldn't find %s on Maven central, generating values...", jar)
			dep, err := jarDependencyFromManifest(jar, "META-INF/MANIFEST.MF", jarFilename)
			if err != nil {
				fmt.Printf("Error generating dependency for %s: %s\n", jar, err)
				report = append(report, ReportRow{
					JarFilename:    jarFilename,
					GroupId:        jarFilename,
					ArtifactId:     jarFilename,
					Version:        "1.0.0",
					OnMavenCentral: false,
					FileChecksum:   checksum,
					LocalFilepath: jar,
				})
			} else {
				report = append(report, ReportRow{
					JarFilename:    jarFilename,
					GroupId:        dep.GroupId,
					ArtifactId:     dep.ArtifactId,
					Version:        dep.Version,
					OnMavenCentral: false,
					FileChecksum:   checksum,
					LocalFilepath: jar,
				})
			}
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

	// Define a template for the bash script with a loop
	bashScriptTemplate := `#!/bin/bash

mvn_repo="{{.MvnRepo}}"

# Array to hold JAR details
jars=(
{{- range .Report}}
  "{{.LocalFilepath}} {{.GroupId}} {{.ArtifactId}} {{.Version}}"
{{- end}}
)

# Loop through the array and execute Maven deploy command
for jar in "${jars[@]}"; do
  local_filepath=$(echo "$jar" | awk '{print $1}')
  group_id=$(echo "$jar" | awk '{print $2}')
  artifact_id=$(echo "$jar" | awk '{print $3}')
  version=$(echo "$jar" | awk '{print $4}')

  mvn deploy:deploy-file \
    -Dfile="$local_filepath" \
    -DgroupId="$group_id" \
    -DartifactId="$artifact_id" \
    -Dversion="$version" \
    -Dpackaging="jar" \
    -DgeneratePom=true \
    -Durl="$mvn_repo" \
    -DcreateChecksum=true
done
`

	// Create a template object
	tmpl, err := template.New("bashScript").Parse(bashScriptTemplate)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Create the output file
	file, err := os.Create("output/import-commands.sh")
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer file.Close()

	// Execute the template with the report data
	err = tmpl.Execute(file, struct {
		Report  []ReportRow // Assuming your report data is in this struct
		MvnRepo string
	}{
		Report:  report,
		MvnRepo: mvnRepo,
	})
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}