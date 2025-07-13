# Article Processing Microservice

## Overview

This project is a **gRPC-based microservice** for processing articles, extracting tags, and storing them in MongoDB. It is designed for scalability, concurrency, and easy deployment using Docker and Docker Compose.

### Key Features
- **gRPC API** for article processing and tag extraction.
- **Tag extraction** using custom logic (see `tagextractor/`).
- **MongoDB integration** for persistent storage of articles and tags.
- **Concurrent processing** for batch article requests.
- **Dockerized** for easy deployment and local development.
- **Graceful shutdown** for safe server termination.

---

## Project Structure

```
.
├── main.go                # Entry point, starts gRPC server
├── proto/                 # Protobuf definitions and generated code
├── server/                # gRPC server implementation
├── tagextractor/          # Tag extraction logic
├── database/              # MongoDB connection and queries
├── utils/                 # Utility functions
├── Dockerfile             # Docker build instructions
├── docker-compose.yml     # Multi-container orchestration
├── Makefile               # Simple build/run/test automation
├── go.mod, go.sum         # Go dependencies
└── readMe.md              # Project documentation
```

---

## gRPC API

- **ProcessSingleArticle**: Extract tags from a single article and store it.
- **ProcessArticles**: Bidirectional streaming for batch processing.
- **GetTopTags**: Retrieve the most frequent tags across all articles.

See `proto/article.proto` for full API details.

---

## Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)
- (For local dev) [Go 1.20+](https://go.dev/dl/)

---

## Quick Start (with Docker)

1. **Clone the repository:**
   ```sh
   git clone git@github.com:DKeshavarz/article-processing-microservice.git
   cd article-processing-microservice
   ```

2. **Build and start the service and MongoDB:**
   ```sh
   docker-compose up --build
   ```
   - This will:
     - Build your Go gRPC server
     - Start MongoDB (on port 27017)
     - Start your server (on port 50051)
     - Persist MongoDB data in a Docker volume

3. **Stop all services:**
   ```sh
   docker-compose down
   ```

---

## Local Development (without Docker)

1. **Install dependencies:**
   ```sh
   make deps
   ```

2. **Start MongoDB (if not already running):**
   - You can use Docker:
     ```sh
     docker run -d --name mongoDp -p 27017:27017 mongo:latest
     ```
   - Or install MongoDB locally.

3. **Run the server:**
   ```sh
   make run
   ```

4. **Run tests:**
   ```sh
   make test
   ```

---

## Configuration

- **MongoDB URI**: The app reads the `MONGODB_URI` environment variable.  
  - In Docker Compose, it’s set to `mongodb://mongoDp:27017`
  - Locally, defaults to `mongodb://localhost:27017`

---

## API Usage

- Use any gRPC client (e.g., [grpcurl](https://github.com/fullstorydev/grpcurl), Postman, or your own code) to interact with the service on port `50051`.
- See `proto/article.proto` for message and service definitions.

---

## Graceful Shutdown

- The server handles `SIGINT`/`SIGTERM` and will:
  - Stop accepting new requests
  - Finish in-flight requests
  - Close the MongoDB connection
  - Exit cleanly

---

## Cleaning Up

- To remove all containers, networks, and volumes:
  ```sh
  docker-compose down -v
  ```

---

## Troubleshooting

- If you see `UNAVAILABLE: failed to connect to all addresses`, check:
  - Both containers are running (`docker-compose ps`)
  - The app is connecting to the correct MongoDB URI
  - Ports are not blocked by other processes

