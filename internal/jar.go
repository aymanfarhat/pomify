// A set of functions to run operations on a JAR file such as extracting
// the manifest file, searching for dependencies, and extracting the Sha1 hash
package internal

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"
)

// Dependency represents a Maven dependency.
func splitLines(s string) []string {
	trimmedS := strings.TrimSpace(s)

	if trimmedS == "" {
		return nil
	}

	delimiter := "\n"
	if strings.Contains(s, "\r\n") {
		delimiter = "\r\n"
	}

	return strings.Split(trimmedS, delimiter)
}
// parseManifest parses a manifest file and returns a map of key-value pairs.
// It returns an error if the manifest file is empty or if a line is invalid.
func parseManifest(manifest string) (map[string]string, error) {
	curr := make(map[string]string)

	lines := splitLines(manifest)
	if lines == nil {
		return nil, errors.New("empty manifest file")
	}

	key := ""
	value := ""
	for n, line := range lines {
		cleanLine := strings.ReplaceAll(line, "\x00", "")
		if cleanLine != "" {
			if unicode.IsSpace(rune(cleanLine[0])) {
				value += strings.Replace(cleanLine, " ", "",1)
				curr[key] = value
			} else {
				parts := strings.SplitN(cleanLine, ":", 2)

				if len(parts) != 2 {
					fmt.Println(cleanLine)
					return nil, fmt.Errorf("invalid manifest line %d", n+1)
				}

				key = strings.TrimSpace(parts[0])
				value = strings.TrimSpace(parts[1])
				curr[key] = value
			}
		}
	}
	return curr, nil
}

// extractJarManifest extracts the manifest file from a JAR file.
// It returns the manifest file as a byte slice and an error if the file cannot be opened.
func extractJarManifest(jarFilePath string, manifestPath string) ([]byte, error) {
	jar, err := zip.OpenReader(jarFilePath)
	if err != nil {
		return nil, err
	}
	defer jar.Close()

	manifestFile, err := jar.Open(manifestPath)
	if err != nil {
		return nil, err
	}
	defer manifestFile.Close()

	data, err := io.ReadAll(manifestFile)
	if err != nil {
		return nil, err
	}

	return data, err
}

// A valid groupId should not be empty, no spaces and following reverse domain name notation.
// Apply a regex to validate this
func validateGroupId(groupId string) bool {
	pattern := `^[a-zA-Z_][a-zA-Z0-9_-]*(\.[a-zA-Z_][a-zA-Z0-9_-]*)*$`
	matched, err := regexp.MatchString(pattern, groupId)
	if err != nil {
	  return false // Handle the error appropriately
	}
	return matched
  }

// A valid artifactId should not be empty, no spaces
func validateArtifactId(artifactId string) bool {
	pattern := `^[a-zA-Z0-9_-]+$`
	matched, err := regexp.MatchString(pattern, artifactId)
	if err != nil {
	  return false // Handle the error appropriately
	}
	return matched
  }

// Attempts to estimate the values of groupId, artifactId, and version from the manifest file of a JAR file.
func jarDependencyFromManifest(jarFilepath string, manifestPath string, jarFilename string) (Dependency, error) {
	defaultDependency := Dependency{
		GroupId:        jarFilename,
		ArtifactId:     jarFilename,
		Version:        "1.0.0",
	}

	rawManifest, err := extractJarManifest(jarFilepath, manifestPath)
	manifest, _ := parseManifest(string(rawManifest))

	if err != nil {
		return defaultDependency, err
	}
	
	artifactId := manifest["Implementation-Title"]
	groupId := manifest["Implementation-Vendor"]
	version := manifest["Implementation-Version"]

	if validateGroupId(groupId) {
		defaultDependency.GroupId = groupId
	}

	if validateArtifactId(artifactId) {
		defaultDependency.ArtifactId = artifactId
	}

	if version != "" {
		defaultDependency.Version = version
	}

	return defaultDependency, nil
}