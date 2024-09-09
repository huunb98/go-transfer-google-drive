# Go Transfer Owner

## Description

Go Transfer Owner is a tool for transferring ownership of files in Google Drive. It uses the Google Drive API to check ownership then transfer it.


## Project Requirements

This project requires the following dependencies to be installed:

* Go version: 1.22.5 or higher
* Mysql: Database for server
* make: The make command is used to build and manage the project.
* swag: The swag command is used to generate API documentation.


## Docker Requirements

To build and run the project using Docker, you will need to have Docker installed on your system. You can install Docker by following the instructions on the official Docker website.

Once you have Docker installed, you can build the project by running the following command:

```bash
docker network create backend
docker compose -f docker-compose.mysql.yml up --build
docker compose -f docker-compose.yml up --build
```

## Usage
To use Go Transfer Owner, you can run the following command:

```bash
go mod tidy
make run
```