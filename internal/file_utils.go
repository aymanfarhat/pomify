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
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// getFileSha1Checksum calculates the SHA1 checksum of the file
// at the given file path. It returns the checksum as a string
// and an error if the file cannot be opened or read.
func getFileSha1Checksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Error opening file: %s", err)
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("Error copying file to hash: %s", err)
	}

	hashInBytes := hash.Sum(nil)
	return fmt.Sprintf("%x", hashInBytes), nil
}

// writeStringsToFile writes the given strings to a file at the
// given file path. It returns an error if the file cannot be
// created or if the strings cannot be written.
func writeStringsToFile(strings []string, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, str := range strings {
		_, err = fmt.Fprintln(file, str)
		if err != nil {
			return err
		}
	}

	return nil
}

// copyFile copies the file at the source path to the destination path.
// It returns the number of bytes copied and an error if the file cannot
// be opened or created.
func copyFile(src, dst string) (int64, error) {
    sourceFileStat, err := os.Stat(src)
    if err != nil {
    	return 0, err
    }
    if !sourceFileStat.Mode().IsRegular() {
        return 0, fmt.Errorf("%s is not a regular file", src)
    }
    source, err := os.Open(src)
    if err != nil {
        return 0, err
    }
    defer source.Close()
    destination, err := os.Create(dst)
    if err != nil {
        return 0, err
    }
    defer destination.Close()
    nBytes, err := io.Copy(destination, source)
    return nBytes, err
}

// writeDependenciesToXML writes the dependencies to an XML file
// at the given file path. It returns an error if the file cannot
// be created or if the encoding fails.
// The XML file will have the following format:
// <dependencies>
//     <dependency>
//         <groupId>...</groupId>
//         <artifactId>...</artifactId>
//         <version>...</version>
//     </dependency>
//     ...
// </dependencies>
func writeDependenciesToXML(dependencies []Dependency, filePath string) error {
	xmlData := Dependencies{Dependencies: dependencies}
	file, err := os.Create(filePath)
	if err != nil {
			return err
	}
	defer file.Close()

	file.WriteString(xml.Header)

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "    ")

	err = encoder.Encode(xmlData)
	if err != nil {
			return err
	}

	return nil
}

// getFilename returns the filename from the given path.
func getFilename(path string) string {
	filename := filepath.Base(path)
	return filename
}

// listJarFiles returns a list of all the JAR files in the given directory.
// It returns an error if the directory cannot be read.
func listJarFiles(jarDir string) ([]string, error) {
	var jarFiles []string
	err := filepath.Walk(jarDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
					return fmt.Errorf("directory %s not found", jarDir)
			}

			if !info.IsDir() && strings.HasSuffix(info.Name(), ".jar") {
					jarFiles = append(jarFiles, path)
			}

			return nil
	})

	if err != nil {
			return nil, err
	}

	return jarFiles, nil
}