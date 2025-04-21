# Go Microservice

[![Build Status](https://github.com/YasCay/go-microservice/actions/workflows/go.yml/badge.svg)](https://github.com/YasCay/go-microservice/actions/workflows/go.yml)
[![SonarCloud](https://sonarcloud.io/api/project_badges/measure?project=YasCay_go-microservice&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=YasCay_go-microservice)

A RESTful API microservice for product management built with Go.

## Features

- RESTful API for product management
- CRUD operations for products (Create, Read, Update, Delete)
- PostgreSQL database integration
- Environment variable configuration
- Database seeding functionality

## Prerequisites

- Go 1.24+
- PostgreSQL database
- Environment variables set up (see Configuration section)

## Installation

```bash
# Clone the repository
git clone https://github.com/YasCay/go-microservice.git
cd go-microservice

# Install dependencies
go mod download
```

## Configuration

Create a `.env` file in the project root with the following variables:

```
APP_DB_USERNAME=your_username
APP_DB_PASSWORD=your_password
APP_DB_NAME=your_dbname
APP_DB_HOST=localhost
APP_DB_PORT=5432
```

## Usage

### Running the Application

```bash
# Start the server
go run *.go

# Start the server with database seeding
go run *.go -seed -count 10
```

The server will start on port 8010.

### API Endpoints

- `GET /products` - Get all products
- `GET /product/{id}` - Get a specific product
- `POST /product` - Create a new product
- `PUT /product/{id}` - Update a product
- `DELETE /product/{id}` - Delete a product
- `DELETE /products` - Delete all products

## Testing

```bash
go test -v
```

## License

[MIT](LICENSE)