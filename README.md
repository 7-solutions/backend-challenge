# Backend Golang Coding Test

## Run the project

### Prerequisites

Choose one of the following ways to run the application:

- **Docker Compose (recommended):** install Docker Desktop (or Docker Engine) with the Docker Compose plugin. Docker Compose builds the Go application and starts MongoDB for you.
- **Run Go locally:** install Go **1.26.4** and MongoDB **8.0**. Start MongoDB and make sure it is reachable at `mongodb://localhost:27017` (the default connection URL).

Git is useful for cloning the repository. No separate Node.js, `protoc`, or Go code-generation tools are needed to run the application; generated protobuf and Swagger files are already included.

Swagger documentation is generated during the Docker image build from the annotations in `app/cmd/server` and `app/layers/deliveries/http/swagger_docs`. To regenerate it locally after changing API annotations, run `make swag` from the repository root.

### Option A: Docker Compose

From the repository root, run:

```sh
docker compose up --build -d
```

This starts the API and MongoDB. The API is available at `http://localhost:8080`, and gRPC listens on port `50051`. Check the application log with:

```sh
docker compose logs -f app
```

Stop the services with `docker compose down`. To also delete the persisted MongoDB data, run `docker compose down -v`.

### Option B: Run the Go application locally

1. Start MongoDB 8.0 and confirm it is listening on `localhost:27017`.
2. From the repository root, download Go dependencies and start the server:

```sh
go mod download
go run ./app/cmd/server
```

The HTTP API listens on port `8080`; gRPC listens on port `50051`. You can also use `make dev`, which installs the repository's protobuf and Swagger generation tools before starting the server. These tools are only needed for that Make target or for regenerating code, not for a normal run.

### Configuration

The application reads configuration from environment variables. Defaults are shown below.

| Variable | Default | Description |
| --- | --- | --- |
| `MONGODB_URL` | `mongodb://localhost:27017` | MongoDB connection URL |
| `MONGODB_NAME` | `backend-challenge` | Database name |
| `MONGODB_USER` | empty | MongoDB username (optional) |
| `MONGODB_PASSWORD` | empty | MongoDB password (optional) |
| `MONGODB_REQUEST_TIMEOUT` | `10` | MongoDB timeout in seconds |
| `JWT_SECRET` | `secret` | HMAC signing key for JWTs; set a strong private value outside local development |
| `JWT_EXPIRE` | `3600` | JWT lifetime in seconds |
| `SERVICE_NAME` | `backend-challenge` | JWT issuer and audience |
| `BASE_URL` | `http://localhost:8080` | Service base URL |
| `DEBUG_MODE` | `false` | Enables debug mode and Swagger UI when `true` |

For example, to run locally with a custom JWT key and Swagger enabled:

```sh
export JWT_SECRET='replace-with-a-long-random-secret'
export DEBUG_MODE=true
go run ./app/cmd/server
```

Swagger UI is then available at `http://localhost:8080/api/swagger/index.html`.

### Quick API check

See [docs/JWT_GUIDE.md](docs/JWT_GUIDE.md) for how to configure the signing key, obtain a token, and call protected endpoints.

Check that the API is up:

```sh
curl http://localhost:8080/api/health
```

Create a user (registration does not require a token):

```sh
curl -X POST http://localhost:8080/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Jane Doe","email":"jane@example.com","password":"secret123"}'
```

Log in to receive a JWT:

```sh
curl -X POST http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"jane@example.com","password":"secret123"}'
```

Use the returned token as a Bearer token for the protected `/api/users` endpoints:

```sh
curl http://localhost:8080/api/users \
  -H 'Authorization: Bearer <JWT_TOKEN>'
```

User routes are `GET /api/users`, `POST /api/users`, `GET /api/users/:id`, `PUT /api/users/:id`, and `DELETE /api/users/:id`. Auth routes are `POST /api/auth/register` and `POST /api/auth/login`; health check is `GET /api/health`. `POST /api/users` remains as a public registration alias.

Successful requests use a JSON envelope. A registration response looks like this (the ID and timestamp are generated):

```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "id": "65f1a4d76a9c0d1234567890",
    "email": "jane@example.com",
    "name": "Jane Doe"
  },
  "meta": {},
  "timestamp": "2026-11-23T12:00:00Z"
}
```

Login returns the JWT string in `data`:

```json
{
  "code": 200,
  "message": "OK",
  "data": "<JWT_TOKEN>",
  "meta": {},
  "timestamp": "2026-11-23T12:00:00Z"
}
```

### API assumptions and design decisions

- Registration (`POST /api/auth/register` or the legacy `POST /api/users` alias) and login are public. All user listing, lookup, update, and delete routes require `Authorization: Bearer <JWT_TOKEN>`.
- User IDs are MongoDB ObjectIDs represented as 24-character hexadecimal strings. Passwords are bcrypt-hashed and never included in user responses.
- Email uniqueness is enforced by a unique MongoDB index as well as an application-level check. Duplicate registrations return a conflict response, including concurrent attempts.
- Updates require both `name` and `email`; password changes are outside this API's scope. User list filters accept optional `name` and `email` query parameters.
- JWTs are signed with HS256. Configure a private `JWT_SECRET` outside local development; tokens carry issuer, audience, subject, issued-at, not-before, and expiration claims.
- The API uses a consistent success/error envelope. User operations currently return HTTP 200 on success; registration does not use HTTP 201.

### Optional: seed a local admin account

With the Docker Compose MongoDB service running, run `make seed-admin`. This creates `admin@admin.com` with password `secret123` if the account does not already exist. Use it to log in at `/api/auth/login` and receive a JWT. This fixed credential is for local development only.

### Tests and coverage

Run the Go test suite, excluding generated protobuf packages under `app/gen` and generated Swagger packages under `app/docs`, with:

```sh
make test
```

Generate a coverage report with:

```sh
make coverage
```

The coverage command writes `coverage.out` and prints per-function coverage to the terminal.

---

## Overview

This coding test has two parts.

| Section | Focus | Submission Type |
| --- | --- | --- |
| User Management API | Build a Golang user management API with MongoDB and JWT authentication | Code implementation |
| Lottery Search System | Design a real-world lottery ticket search solution with wildcard matching | Design proposal only (no code) |

## User Management API

### Objective (User Management API)

Build a RESTful API in Golang to manage users, using MongoDB for persistence, JWT for authentication, and clean code practices.

### Requirements (User Management API)

#### 1. User Model

Define a user entity with the following fields:

- `ID` (auto-generated)
- `Name` (string)
- `Email` (string, unique)
- `Password` (hashed)
- `CreatedAt` (timestamp)

#### 2. Authentication

Implement:

- User registration
- User authentication that returns a JWT token

JWT requirements:

- Protect endpoints with JWT
- Validate tokens via middleware
- Sign tokens using HMAC (`HS256`) with a secret key

#### 3. User Operations

Implement the following operations:

- Create a new user
- Fetch a user by ID
- List all users
- Update a user's name or email
- Delete a user

#### 4. MongoDB Integration

- Use the official Go MongoDB driver
- Persist and retrieve user data from MongoDB

#### 5. Middleware

- Implement logging middleware to capture HTTP method, path, and execution time

#### 6. Concurrency Task

- Run a background goroutine every 10 seconds to log the total number of users in the database

#### 7. Testing

- Write unit tests using Go's standard `testing` package
- Mock MongoDB interactions where appropriate

### Bonus (Optional, User Management API)

- **Containerization**: Add Docker and `docker-compose` support for the API and MongoDB
- **Abstraction**: Use Go interfaces to abstract MongoDB operations for better testability
- **Validation**: Implement input validation (for example, required fields and email format)
- **Graceful Shutdown**: Handle system signals using `context.Context`
- **gRPC Support**:
  - Define a `.proto` file for `CreateUser` and `GetUser`
  - Implement a gRPC server (optionally secure with token metadata)
- **Hexagonal Architecture**:
  - Structure the project using ports and adapters
  - Separate domain, application, and infrastructure layers
  - Decouple business logic from frameworks and drivers

### Deliverables (User Management API)

Provide a Git repository containing:

- `README.md` with setup and execution instructions
- A guide explaining how to generate and use JWT tokens
- Sample API requests and responses
- Documentation of assumptions or design decisions

### Evaluation Criteria (User Management API)

- Code quality, structure, and readability
- Correctness and completeness of the REST API
- Security and implementation of JWT
- Proper usage and abstraction of MongoDB
- Test coverage and effective mocking
- Idiomatic Go usage
- Bonus implementations (gRPC, Docker, validation, architecture)

## Lottery Search System

### Objective (Lottery Search System)

Design a real-world solution to search a large dataset of lottery tickets using pattern matching with wildcard support.

> This section is a design exercise. Do not implement code.

### Requirements (Lottery Search System)

#### 1. Data Volume

- Handle a dataset of **1 million** lottery tickets
- Each ticket is a 6-digit number

#### 2. Search Pattern

- Support a 6-character search pattern containing digits and wildcards (`*`)
- Example patterns:

| Pattern | Matches |
| --- | --- |
| `****23` | Numbers ending in `23` |
| `1****5` | Numbers starting with `1` and ending with `5` |
| `123***` | Numbers starting with `123` |

#### 3. Result Distribution

- Constraint: the same search pattern should not return the same ticket to multiple users at the same time
- Propose a distribution mechanism so matching tickets are assigned without duplicate simultaneous selection

#### 4. Performance

- Ensure the search is performant for `1M+` records
- Propose an efficient approach for querying and allocation

#### 5. Real-World Design Proposal (No Code Required)

- Recommend the database/storage technology you would use in production and explain why
- Describe the algorithm and indexing strategy used for wildcard pattern matching
- Explain how you would prevent duplicate simultaneous results for the same pattern (for example, locking, reservation, or atomic allocation)
- No code implementation is required; provide a solution/design only

### Deliverables (Lottery Search System)

Submit a design document only (no code implementation) that includes:

- Proposed solution architecture, data structures, and algorithms
- Recommended production database/storage choice with justification (for example, query performance, concurrency handling, operational simplicity)
- Performance analysis summarizing efficiency and tradeoffs
- Concurrency/distribution strategy explaining how duplicate results are avoided for the same pattern

### Evaluation Criteria (Lottery Search System)

- Feasibility: the solution addresses the stated requirements
- Performance: the search approach is efficient for the target scale
- Correctness: the distribution constraint is handled correctly
- Real-world practicality: the database/storage and concurrency approach are appropriate for production use
- Creativity: thoughtful use of data structures and algorithms
