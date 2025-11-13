# 🚀 Distributed Push Notification Service API

## Overview
A robust and scalable backend service engineered with **FastAPI** for managing and dispatching push notifications via **Firebase Cloud Messaging (FCM)**. This system leverages **RabbitMQ** for resilient asynchronous message processing, incorporating a sophisticated retry mechanism and a Dead-Letter Queue (DLQ) to ensure reliable delivery. Notification states are meticulously persisted using **SQLModel** with an **SQLite** database.

## Features
-   **FastAPI**: Utilizes the high-performance, asynchronous Python web framework for building efficient and responsive API endpoints.
-   **Firebase Cloud Messaging (FCM)**: Integrates seamlessly with FCM to facilitate reliable and real-time push notification delivery to various client devices.
-   **RabbitMQ**: Implements a robust message queuing system that decouples notification requests from the actual sending process, enhancing API responsiveness and scalability.
-   **Asynchronous Processing**: Handles the intricate logic of notification dispatch in the background, ensuring non-blocking operations and improved system throughput.
-   **Retry Mechanism**: Automatically re-attempts failed notification deliveries up to a configurable maximum number of retries, increasing delivery assurance.
-   **Dead-Letter Queue (DLQ)**: Provides a dedicated queue for messages that fail after exhausting all retries or encounter unrecoverable processing errors, preventing data loss and enabling post-mortem analysis.
-   **SQLModel (SQLite)**: Persists comprehensive details of push notification requests and their delivery statuses within a lightweight, local SQLite database.
-   **Docker & Docker Compose**: Offers a streamlined development and deployment experience, enabling quick spin-up of the API, RabbitMQ broker, and consumer worker with minimal setup.
-   **Structured Logging**: Employs a custom logging solution with file rotation for comprehensive monitoring, debugging, and operational insights.

## Getting Started

### Prerequisites
Before you begin, ensure your system meets the following requirements:
-   **Docker**: [Install Docker](https://www.docker.com/get-started) to run the services in isolated containers.
-   **Docker Compose**: [Install Docker Compose](https://docs.docker.com/compose/install/) for defining and running multi-container Docker applications.
-   **Python 3.8+**: (Optional, if running Python components directly without Docker).
-   **Firebase Project**: Access to a Firebase project and its associated service account key.

### Installation

1.  **Clone the Repository**:
    ```bash
    git clone git@github.com:brainox/hng-group55-distributed-notification-system.git
    cd hng-group55-distributed-notification-system
    ```

2.  **Firebase Service Account Key**:
    -   🔐 Obtain your Firebase service account key JSON file. Navigate to your Firebase project settings (`Project settings` > `Service accounts` > `Generate new private key`).
    -   Rename this downloaded file to `distributed-systems-3349d-firebase-adminsdk-fbsvc-b84f851127.json` and place it in the root directory of the cloned project. If you prefer a different filename, update the `credential_path` argument in `consumer.py` and `service/senders.py` accordingly.

3.  **Build and Run with Docker Compose**:
    🚀 This project is designed for easy deployment using Docker Compose, which orchestrates RabbitMQ, the FastAPI application, and the consumer worker.
    ```bash
    docker-compose up --build
    ```
    This command will:
    -   Build the necessary Docker image for the application.
    -   Initialize and start a RabbitMQ container.
    -   Launch the `push_api` service, hosting the FastAPI application.
    -   Start the `push_worker` service, which consumes messages from RabbitMQ.

### Environment Variables
The project primarily relies on internal configurations and Docker Compose for its dependencies:

-   **Firebase Service Account Key Path**: The path to your Firebase service account JSON file. By default, the system expects `distributed-systems-3349d-firebase-adminsdk-fbsvc-b84f851127.json` to be present in the project's root directory.
-   **RabbitMQ Host**: Configured as `localhost` in the Python scripts (`consumer.py`, `service/queue.py`) for establishing Pika connections. When running within the Docker Compose network, the `rabbitmq` service is accessible via the `rabbitmq` hostname.
-   **RabbitMQ Default User**: `guest` (configured in `docker-compose.yml` for the RabbitMQ service).
-   **RabbitMQ Default Password**: `guest` (configured in `docker-compose.yml` for the RabbitMQ service).

## Usage
Once all services are up and running via `docker-compose up`, the FastAPI application will be accessible at `http://localhost:8000`.

To initiate a push notification, simply send a `POST` request to the `/send` endpoint with the appropriate payload. The message will be queued, processed asynchronously by the consumer worker, and its delivery status updated in the database.

Example using `curl`:
```bash
curl -X POST "http://localhost:8000/send" \
     -H "Content-Type: application/json" \
     -d '{
           "id": 1,
           "title": "Welcome to Our App!",
           "body": "Thanks for joining. Explore our new features!",
           "token": "YOUR_FCM_DEVICE_TOKEN_HERE",
           "image": "https://example.com/welcome.png",
           "url": "https://yourapp.com/onboarding",
           "data": {
             "campaign_id": "onboarding_flow",
             "user_level": "new"
           }
         }'
```
Remember to replace `YOUR_FCM_DEVICE_TOKEN_HERE` with a valid FCM device token to successfully receive the notification.

## API Documentation

### Base URL
`http://localhost:8000`

### Endpoints

#### GET /
**Overview**: Checks the root endpoint of the API and provides foundational service information.

**Request**:
No payload required.

**Response**:
```json
{
  "success": true,
  "data": {
    "title": "Push Notification Service",
    "version": "1.0.0",
    "description": "A service to send push notifications via FCM."
  },
  "error": null,
  "message": "RootAPI",
  "meta": null
}
```

**Errors**:
-   `500 Internal Server Error`: An unexpected server error occurred during processing.

#### GET /health
**Overview**: Performs a health check on the API to ascertain its operational status.

**Request**:
No payload required.

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "ok"
  },
  "error": null,
  "message": "Health Check OK",
  "meta": null
}
```

**Errors**:
-   `500 Internal Server Error`: An unexpected server error occurred during the health check.

#### POST /send
**Overview**: Enqueues a new push notification message for asynchronous delivery through FCM. The message details, including its initial `pending` status, are persisted in the database.

**Request**:
```json
{
  "id": 1,
  "title": "string",
  "body": "string",
  "token": "string",
  "image": "string",
  "url": "string",
  "data": {
    "key": "value"
  }
}
```
**Required fields**:
-   `id`: `int` - A unique integer identifier for the push message.
-   `title`: `str` - The title text of the push notification.
-   `body`: `str` - The main content body of the push notification.
-   `token`: `str` - The Firebase Cloud Messaging device token of the intended recipient.

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Welcome to Our App!",
    "body": "Thanks for joining. Explore our new features!",
    "token": "FCM_DEVICE_TOKEN",
    "status": "pending",
    "retry_count": 0,
    "created_at": "2023-10-27T10:00:00.000000",
    "updated_at": "2023-10-27T10:00:00.000000"
  },
  "error": null,
  "message": "Health Check OK",
  "meta": null
}
```
*(Note: The `data` field in the response contains the complete `PushMessage` object created and stored in the database.)*

**Errors**:
-   `422 Unprocessable Entity`: Occurs if required fields (`id`, `title`, `body`, `token`) are missing or invalid.
-   `500 Internal Server Error`:
    -   A `PushMessage` with the provided `id` already exists in the database.
    -   The database session is unavailable.
    -   Failure to create a `PushMessage` instance.
    -   Other unforeseen server errors during the queuing or database interaction process.

#### GET /notifications
**Overview**: Retrieves a comprehensive list of all stored push notification messages, including their current statuses and metadata.

**Request**:
No payload required.

**Response**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "Welcome to Our App!",
      "body": "Thanks for joining.",
      "token": "FCM_TOKEN_1",
      "status": "success",
      "retry_count": 0,
      "created_at": "2023-10-27T10:00:00.000000",
      "updated_at": "2023-10-27T10:05:00.000000"
    },
    {
      "id": 2,
      "title": "System Update",
      "body": "Important notice about upcoming maintenance.",
      "token": "FCM_TOKEN_2",
      "status": "failed",
      "retry_count": 3,
      "created_at": "2023-10-26T14:30:00.000000",
      "updated_at": "2023-10-26T14:45:00.000000"
    }
  ],
  "error": null,
  "message": "Notifications retrieved successfully",
  "meta": null
}
```
*(Note: The `data` field in the response contains an array of `PushMessage` objects.)*

**Errors**:
-   `500 Internal Server Error`: An unexpected server error occurred during database retrieval.

#### GET /notifications/{notification_id}
**Overview**: Retrieves the details of a single push notification message using its unique identifier.

**Request**:
No payload required.
**Path Parameters**:
-   `notification_id`: `int` - The unique ID of the specific notification to retrieve.

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Welcome to Our App!",
    "body": "Thanks for joining.",
    "token": "FCM_TOKEN_1",
    "status": "success",
    "retry_count": 0,
    "created_at": "2023-10-27T10:00:00.000000",
    "updated_at": "2023-10-27T10:05:00.000000"
  },
  "error": null,
  "message": "Notification retrieved successfully",
  "meta": null
}
```
*(Note: The `data` field in the response contains a single `PushMessage` object.)*

**Errors**:
-   `404 Not Found`: No notification could be found corresponding to the provided `notification_id`.
-   `500 Internal Server Error`: An unexpected server error occurred during database retrieval.

## Technologies Used
| Technology         | Description                                                        | Link                                                       |
| :----------------- | :----------------------------------------------------------------- | :--------------------------------------------------------- |
| Python             | The primary programming language used for developing the service.    | [python.org](https://www.python.org/)                      |
| FastAPI            | A modern, fast (high-performance) web framework for building APIs. | [fastapi.tiangolo.com](https://fastapi.tiangolo.com/)      |
| SQLModel           | A library for interacting with SQL databases, built on SQLAlchemy. | [sqlmodel.tiangolo.com](https://sqlmodel.tiangolo.com/)    |
| SQLite             | A lightweight, serverless, and self-contained relational database engine. | [sqlite.org](https://www.sqlite.org/index.html)            |
| RabbitMQ           | A robust and widely-used open-source message broker for asynchronous processing. | [rabbitmq.com](https://www.rabbitmq.com/)                  |
| Pika               | The official Python client library for RabbitMQ.                   | [pika.readthedocs.io](https://pika.readthedocs.io/en/stable/) |
| Firebase Admin SDK | Enables server-side integration with Firebase services, specifically FCM for sending messages. | [firebase.google.com/docs/admin/setup](https://firebase.google.com/docs/admin/setup) |
| Docker             | A platform for developing, shipping, and running applications in containers. | [docker.com](https://www.docker.com/)                      |
| Uvicorn            | An ASGI (Asynchronous Server Gateway Interface) server, used to serve FastAPI applications. | [www.uvicorn.org](https://www.uvicorn.org/)                |

## License
This project is licensed under the MIT License.

<!-- ## Author Info
-   **Developer**: [Your Name Here]
-   **LinkedIn**: [Your LinkedIn Profile Link]
-   **Twitter**: [Your Twitter Handle Link] -->

## Badges
[![Python 3.10](https://img.shields.io/badge/Python-3.10%2B-blue?logo=python&logoColor=white)](https://www.python.org/)
[![FastAPI](https://img.shields.io/badge/FastAPI-0.121.1-009688?logo=fastapi)](https://fastapi.tiangolo.com/)
[![SQLModel](https://img.shields.io/badge/SQLModel-0.0.18-orange?logo=sqlmodel)](https://sqlmodel.tiangolo.com/)
[![RabbitMQ](https://img.shields.io/badge/RabbitMQ-4.0.9-ff6600?logo=rabbitmq&logoColor=white)](https://www.rabbitmq.com/)
[![Docker](https://img.shields.io/badge/Docker-24.0.5-blue?logo=docker)](https://www.docker.com/)

[![Readme was generated by Dokugen](https://img.shields.io/badge/Readme%20was%20generated%20by-Dokugen-brightgreen)](https://www.npmjs.com/package/dokugen)