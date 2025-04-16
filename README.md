# Receipt Processor Web Service

A Go-based web service built to fulfill the requirements of the Fetch Rewards coding assessment. It processes retail receipt data, calculates reward points based on defined rules, and provides access to those points via a RESTful API.

## Overview

This application provides two primary RESTful API endpoints:

1.  `/receipts/process` (POST): Accepts receipt details in JSON format, calculates points according to specific rules, stores the receipt data (in memory), and returns a unique ID for retrieval.
2.  `/receipts/{id}/points` (GET): Retrieves the calculated points for a previously processed receipt using its unique ID.

The service utilizes in-memory storage (i.e., data does not persist across restarts), includes structured logging, and is designed to be run easily using Docker or the Go toolchain directly.

## Prerequisites

To build and run this application, you will need:

*   **Go:** Version 1.20 or later.
*   **Docker:** For building and running the containerized application.
*   **Docker Desktop:** Preferred if you are on Windows or macOS for easier Docker management.

## Running the Application

Choose one of the following methods to run the service.

### Option 1: Running with Docker 

1.  **Build the Docker Image:**
    From the project's root directory, run:
    ```bash
    docker build -t receipt-processor .
    ```

2.  **Run the Docker Container:**
    To start the container and map the application's port (8080) to your host machine:
    ```bash
    docker run -d -p 8080:8080 --name receipt-processor-test receipt-processor
    ```
    *   The API will be accessible at `http://localhost:8080`.

3.  **Accessing Logs with Volume Mount:**
    If you want logs persisted outside the container:
    ```bash
    # Please verify 'logs' dir exists first: mkdir -p logs (macOS/Linux) or mkdir logs (Windows)
    # Linux/macOS:
    docker run -d -p 8080:8080 -v "$(pwd)/logs":/app/logs --name receipt-processor-test-logs receipt-processor
    # Windows:
    docker run -d -p 8080:8080 -v "%cd%\logs":/app/logs --name receipt-processor-test-logs receipt-processor
    ```

### Option 2: Running Locally (Without Docker)

1.  **Clone the Repository:**
    ```bash
    git clone https://github.com/saurabhd96/receipt-processor.git
    cd receipt-processor
    ```

2.  **Download Dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Run the Application:**
    ```bash
    go run main.go
    ```
    *   The API will be accessible at `http://localhost:8080`. Logs appear in a `logs` sub-directory. Press `Ctrl+C` to stop.

## Points Calculation Rules Summary

Points are calculated based on the following implemented rules:

1.  **One point for every alphanumeric character** in the retailer name.
2.  **50 points** if the total is a round dollar amount with no cents (e.g. `15.00`).
3.  **25 points** if the total is a multiple of `0.25`.
4.  **5 points for every two items** on the receipt.
5.  If the **trimmed length** of the item description is a multiple of 3, multiply the price by `0.2` and **round up** to the nearest integer. The result is the number of points earned for that item.
6.  **6 points** if the day in the purchase date is odd.
7.  **10 points** if the time of purchase is **after 2:00pm (14:00) and before 4:00pm (16:00)**.

## Project Structure

```
.
├── Dockerfile           # Defines the Docker image build process
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── main.go              # Application entry point, server setup, routing
├── README.md            # This file: Project documentation and instructions
├── handlers/            # Request handlers for API endpoints
│   └── receipts.go      # Handles /receipts/process and /receipts/{id}/points
├── logging/             # Custom structured logging implementation
│   └── logger.go        # Sets up and provides logging functions
├── models/              # Data structures, business logic, storage
│   ├── errors.go        # Custom error types and validation helpers
│   └── receipt.go       # Receipt/Item structs, points calculation, in-memory store
└── logs/                # Directory for log files (created automatically at runtime)
```

## Troubleshooting

*   **Port Conflict:** If you get an error like `bind: address already in use` or "port is already allocated", ensure no other application is using port 8080. For mapping to a different host port with Docker:
    `docker run -p <other_host_port>:8080 receipt-processor`.
*   **Docker Daemon:** Ensure the Docker engine/daemon is running before executing `docker` commands. Check Docker Desktop status or run `systemctl status docker` (Linux).
*   **Log Volume Mounts (Docker):** If logs aren't appearing in your host `logs` directory:
    *   Verify the `logs` directory exists on your host *before* running the `docker run -v ...` command.
    *   Double-check the host path (`$(pwd)/logs` or `%cd%\logs`) is correct for your OS and current directory.
    *   Check file permissions on the host `logs` directory.
*   **Invalid Input (400 Bad Request):** Carefully check that your JSON payload matches the required structure. Ensure values adhere to the expected formats: `YYYY-MM-DD` for `purchaseDate`, `HH:MM` (24-hour) for `purchaseTime`, and `"X.XX"` (string with two decimal places) for `total` and item `price`.