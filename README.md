# Mailing List Microservice

A simple mailing list management microservice built with Go that provides both JSON REST API and gRPC interfaces for managing email subscriptions.

## Features

- **Dual API Support**: Both JSON REST API and gRPC API
- **SQLite Database**: Lightweight database for email storage
- **Email Management**: Create, read, update, delete, and batch retrieve email entries
- **Opt-out Support**: Soft delete functionality with opt-out flags
- **Pagination**: Batch retrieval with page-based pagination

## Project Structure

```
mailinglist-ms/
├── server/          # Main server application
├── client/          # gRPC client example
├── jsonapi/         # JSON REST API handlers
├── grpcapi/         # gRPC API implementation
├── maildb/          # Database operations
├── proto/           # Protocol buffer definitions
└── go.mod           # Go module file
```

## Prerequisites

- Go 1.24+
- Protocol Buffers compiler (protoc)
- SQLite3

## Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd mailinglist-ms
```

2. Install dependencies:

```bash
go mod tidy
```

3. Generate protobuf files (if needed):

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/mail.proto
```

## Usage

### Running the Server

Start both JSON and gRPC servers:

```bash
go run server/server.go
```

The server will start with default settings:

- JSON API: `localhost:8080`
- gRPC API: `localhost:8082`
- Database: `list.db`

### Environment Variables

You can configure the server using environment variables:

```bash
export MAILINGLIST_DB="custom.db"
export MAILINGLIST_BIND_JSON=":8080"
export MAILINGLIST_BIND_GRPC=":8082"
```

### JSON REST API

#### Endpoints

- `POST /email/create` - Create a new email entry
- `GET /email/get` - Get an email entry
- `PUT /email/update` - Update an email entry
- `POST /email/delete` - Soft delete (opt-out) an email
- `GET /email/get_batch` - Get emails in batches

#### Examples

**Create Email:**

```bash
curl -X POST http://localhost:8080/email/create \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'
```

**Get Email:**

```bash
curl -X GET http://localhost:8080/email/get \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'
```

**Get Email Batch:**

```bash
curl -X GET http://localhost:8080/email/get_batch \
  -H "Content-Type: application/json" \
  -d '{"page":1,"count":10}'
```

### gRPC API

#### Running the Client

The included client demonstrates gRPC usage:

```bash
go run client/client.go
```

## Database Schema

The SQLite database contains a single `emails` table:

```sql
CREATE TABLE emails (
    id INTEGER PRIMARY KEY,
    email TEXT UNIQUE,
    confirmed_at INTEGER,
    opt_out INTEGER
);
```

## API Reference

### EmailEntry Structure

```json
{
  "id": 1,
  "email": "user@example.com",
  "confirmed_at": 1640995200,
  "opt_out": false
}
```

### gRPC Service Methods

- `CreateEmail(CreateEmailRequest) returns (EmailResponse)`
- `GetEmail(GetEmailRequest) returns (EmailResponse)`
- `UpdateEmail(UpdateEmailRequest) returns (EmailResponse)`
- `DeleteEmail(DeleteEmailRequest) returns (EmailResponse)`
- `GetEmailBatch(GetEmailBatchRequest) returns (GetEmailBatchResponse)`

## Development

### Regenerating Protobuf Files

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/mail.proto
```

### Building

```bash
go build -o mailinglist-server server/server.go
go build -o mailinglist-client client/client.go
```
