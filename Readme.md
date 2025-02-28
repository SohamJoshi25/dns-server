# DNS Resolver CLI

A Command Line Interface (CLI) tool written in Go that functions as an Iterative DNS Resolver, allowing users to override default DNS answers using PostgreSQL. This project leverages the Cobra library to provide various commands for managing custom DNS records and starting the DNS server.

## Features

- **Iterative DNS Lookup**: Automatically performs iterative DNS resolution by querying root servers, TLD servers, and authoritative servers recursively. Also Caches the Itterative Resolution Rsponse.
- **Custom DNS Records**: Override DNS answers with custom entries stored in a PostgreSQL database.
- **Record Management**: Add, list, and remove DNS records from the database via CLI commands.
- **DNS Server**: Start a DNS server that listens on `127.0.0.1`.
- **Modular Architecture**: Separation of concerns with distinct modules for CLI, database interactions, and resolver logic.

## Installation

### Prerequisites

- Go (version 1.21 or higher)
- Docker (for running PostgreSQL)

### Setting Up PostgreSQL with Docker

To utilize the custom DNS records feature, ensure you have a running PostgreSQL instance. You can set this up using Docker:

```bash
docker run --name my-postgres-container   -e POSTGRES_USER=postgres   -e POSTGRES_PASSWORD=mypassword   -e POSTGRES_DB=postgres   -p 5432:5432   -d postgres:latest
```

This command pulls the latest PostgreSQL image from Docker Hub and runs it in a container with the specified environment variables. [Learn more about using PostgreSQL with Docker](https://www.docker.com/blog/how-to-use-the-postgres-docker-official-image/).

### Cloning the Repository

```bash
git clone https://github.com/SohamJoshi25/dns-server.git
cd dns-server
```

### Building the Project

```bash
go build -o dns-server
```

## Usage

### Adding a Custom DNS Record

```bash
./dns-server add --domain example.com --answer 1.2.3.4 --type A
```

### Listing All DNS Records

```bash
./dns-server list
```

### Removing a DNS Record by ID

```bash
./dns-server remove 1
```

### Starting the DNS Server

```bash
./dns-server start
```

By default, the server listens on `127.0.0.1:53`.
&nbsp;
![Commands Usage](https://raw.githubusercontent.com/SohamJoshi25/dns-server/refs/heads/main/docs/images/image.png)
&nbsp;
![Database Schema and Records](https://raw.githubusercontent.com/SohamJoshi25/dns-server/refs/heads/main/docs/images/db.png)
&nbsp;
![DNS Lookup](https://raw.githubusercontent.com/SohamJoshi25/dns-server/refs/heads/main/docs/images/lookup.png)
&nbsp;
# You can see my server looksup the AAAA record from internet but because A and TXT record was preset in Database, it fetched from there

## Configuration

Database connection settings can be configured via environment variables:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=mypassword
export DB_NAME=postgres
```

## Project Structure

The project follows a modular architecture:

```
dns-server/
├── cmd/
│   ├── add.go
│   ├── delete.go
│   ├── dns.go
│   ├── list.go
│   └── root.go
├── internal/
│   ├── dnsdb/
│   │   └── db.go
│   ├── dnslookup/
│   │   ├── constants.go
│   │   ├── lookup.go
│   │   └── types.go
│   └── dnsproxy/
│       └── dnsproxy.go
├── docs/
│   └── images/
│       ├── db.png
│       ├── image.png
│       └── lookup.png
├── .gitignore
├── dns-server.exe
├── go.mod
├── go.sum
├── main.go
├── Readme.md
└── Tood.md
```

- `cmd/`: Contains CLI command implementations using Cobra.
- `internal/`:
  - `dnsdb/`: Handles database interactions.
  - `dnslookup/`: Manages DNS resolution logic.
  - `dnsproxy/`: Implements DNS proxy functionalities.
- `docs/`: Documentation and related images.
- `main.go`: The entry point of the application.

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License.
