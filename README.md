# Pomify

A simple CLI tool for classifying, organizing and managing jar dependencies into a Maven project.

## Features

- [x] Scan a directory for jar files and check if they're on Maven Central or need a custom repository
- [x] Generate a status report of the jar files found in the directory
- [x] Generate pom.xml dependency blocks with the dependencies found on Maven Central and those that need a custom repository
- [x] Generate an mvn import command to migrate jar files into a custom Maven repository


## Usage

### Scan
Scan a directory of jar files and output a CSV report

USAGE:
   pomify scan [command options]

OPTIONS:
   --jars value, -j value    Path to the directory containing the jar files
   --output value, -o value  Path to the output CSV file

```bash
pomify scan --jars /path/to/jars --output /path/to/output.csv
```

### Generate

Generate a pom.xml dependencies block from a CSV report for each jar category

USAGE:
   pomify generate [command options]

OPTIONS:
   --report value, -r value  Path to the Pomify CSV report file
   --output value, -o value  Path to the output directory of XML files

```bash
pomify generate --report /path/to/output.csv --output /path/to/output
```

