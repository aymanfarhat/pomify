# Pomify

A simple CLI tool for classifying, organizing and managing jar dependencies into a Maven project.

## Features

- [x] Scan a directory for jar files and check if they're on Maven Central or need a custom repository
- [x] Generate a status report of the jar files found in the directory
- [x] Generate pom.xml dependency blocks with the dependencies found on Maven Central and those that need a custom repository
- [x] Generate an mvn import command to migrate jar files into a custom Maven repository

