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

```bash
pomify scan --jars /path/to/jars --output /path/to/output.csv
```

### Generate

Generate a pom.xml dependencies block from a CSV report for each jar category

```bash
pomify generate --report /path/to/output.csv --output /path/to/output
```


### Push

Generate an mvn import command to migrate jar files into a custom Maven repository

```bash
pomi push --report /path/to/output.csv --mvn_repo http://localhost:8081/repository/maven-releases
```

