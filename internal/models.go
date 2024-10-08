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
	"encoding/xml"
)

// Dependency represents a Maven dependency.
// It has a GroupId, ArtifactId, and Version.
// For example:
// <dependency>
//     <groupId>com.google.guava</groupId>
//     <artifactId>guava</artifactId>
//     <version>30.1-jre</version>
// </dependency>
type Dependency struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
}

// Dependencies represents a list of Maven dependencies.
// It has a list of Dependency objects.
// For example:
// <dependencies>
//     <dependency>
//         <groupId>com.google.guava</groupId>
//         <artifactId>guava</artifactId>
//         <version>30.1-jre</version>
//     </dependency>
//     ...
// </dependencies>
type Dependencies struct {
	XMLName     xml.Name     `xml:"dependencies"`
	Dependencies []Dependency `xml:"dependency"`
}

// ReportRow represents a row in the report.
type ReportRow struct {
	JarFilename string `csv:"JarFilename"`
	GroupId	 string `csv:"GroupID"`
	ArtifactId string `csv:"ArtifactID"`
	Version	 string `csv:"Version"`
	OnMavenCentral bool `csv:"OnMavenCentral"`
	FileChecksum string `csv:"FileChecksum"`
	LocalFilepath string `csv:"LocalFilepath"`
}