# Running Statistics API

This API allows you to manage and retrieve running statistics data.

## Prerequisites
- Go installed on your machine
- CockroachDB or PostgreSQL database setup
- Environment variables:
    - `PGHOST`: Hostname of the database
    - `PGPORT`: Port of the database
    - `PGDATABASE`: Name of the database
    - `PGUSER`: Database user
    - `PGPASSWORD`: Database password
    - `HTTP_PORT`: Port on which the server will run (default: 8080)

## Installation and Setup
1. Clone the repository
2. Set the required environment variables
3. Build and run the application by executing: `go run main.go`

## Endpoints
- `GET /`: Home route
- `GET /ping`: Check if the server is running
- `POST /running_statistics`: Add a new running statistic record
- `GET /get_running_statistics`: Retrieve all running statistics records

## API Payload
### Running Statistics (POST /running_statistics)
```json
{
  "date": "2022-01-01",
  "distance": "10 km",
  "time": "50 mins"
}

## Response format
{
  "message": "Record creation status",
  "record_created": true
}


## Acknowledgements
This API is built using the Echo framework for Go.