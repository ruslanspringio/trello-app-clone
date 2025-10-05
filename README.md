# Trello Clone API

![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/postgresql-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![Prometheus](https://img.shields.io/badge/prometheus-%23E6522C.svg?style=for-the-badge&logo=prometheus&logoColor=white)
![Grafana](https://img.shields.io/badge/grafana-%23F46800.svg?style=for-the-badge&logo=grafana&logoColor=white)

A feature-rich, scalable, and observable backend for a Trello-like project management application. This project is built with Go and demonstrates a modern, production-ready architecture using a suite of powerful technologies.

## ✨ Features

- **Full Project Management Core**: Complete CRUD functionality for Boards, Lists, and Cards.
- **Real-time Collaboration**: Instant updates across all connected clients using **WebSockets**.
- **User Authentication**: Secure JWT-based authentication for user registration and login.
- **High Performance**: **Redis** caching for frequently accessed data to reduce database load.
- **Observability**:
    - **Metrics**: Instrumented with **Prometheus** for real-time monitoring of application health (RPS, latency, errors).
    - **Structured Logging**: Ready for integration with centralized logging systems.
- **Containerized**: Fully containerized with **Docker** and orchestrated with **Docker Compose** for easy setup and deployment.
- **Collaborative Workspaces**: Invite users to boards to work together.
- **Interactive API Documentation**: **Swagger (OpenAPI)** documentation available for easy testing and API exploration.

## 🛠️ Tech Stack

- **Language**: Go
- **Framework**: Gin
- **Database**: PostgreSQL
- **In-Memory Store**: Redis (for caching)
- **Real-time Communication**: Gorilla WebSocket
- **Monitoring**: Prometheus & Grafana
- **Containerization**: Docker & Docker Compose

## 🚀 Getting Started

The entire application stack can be launched with a single command using Docker Compose.

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop/) installed and running.
- [Docker Compose](https://docs.docker.com/compose/install/) (usually included with Docker Desktop).

### Installation & Launch

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/your-repo-name.git
    cd your-repo-name
    ```

2.  **Set up environment variables:**
    The application reads its configuration from environment variables. For convenience, you can create a `.env` file in the root of the project. A `docker-compose.yml` file will automatically pick it up.

    Create a `.env` file and add the following, customizing if needed:
    ```env
    # Application Port
    PORT=8080

    # PostgreSQL Configuration
    DB_HOST=db
    DB_PORT=5432
    DB_USER=notes_user
    DB_PASSWORD=notes_password
    DB_NAME=notes_db
    DB_SSLMODE=disable

    # Redis Configuration
    REDIS_ADDR=cache:6379

    # JWT Secret Key (use a long, random string)
    JWT_SECRET_KEY=your_super_secret_key_for_jwt_that_is_very_long
    ```

3.  **Run the entire stack:**
    ```bash
    docker-compose up -d --build
    ```
    This command will:
    - Build the Go application Docker image.
    - Start all services in the background (`-d`): App, PostgreSQL, Redis, Prometheus, and Grafana.

4.  **Verify that everything is running:**
    You can check the status of all containers with:
    ```bash
    docker-compose ps
    ```
    All services should have a `running` or `up` status.

## 🖥️ Available Services

Once the stack is running, the following services will be available on your `localhost`:

- **Trello Clone API**: `http://localhost:8080`
- **Swagger API Docs**: `http://localhost:8080/swagger/index.html`
- **Prometheus**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000` (Login: `admin` / `admin`)

##  API Usage

The best way to explore and test the API is through the interactive **Swagger documentation** available at `http://localhost:8080/swagger/index.html`.

### Authentication Flow

1.  **Register a new user:**
    - `POST /api/users/register`
2.  **Login to get a JWT token:**
    - `POST /api/users/login`
3.  **Authorize your requests:**
    - In Swagger UI, click the "Authorize" button and enter `Bearer <your_token>`.
    - In Postman or other clients, add the `Authorization` header with the value `Bearer <your_token>`.

All endpoints except for registration and login are protected and require this token.

## 📈 Monitoring

A pre-configured monitoring stack is included.

1.  **Prometheus**:
    - Go to `http://localhost:9090`.
    - Navigate to `Status -> Targets`. You should see the `trello-app` job with a state of `UP`.
2.  **Grafana**:
    - Go to `http://localhost:3000` and log in.
    - **Configure Prometheus Data Source:**
        - Go to `Configuration (gear icon) -> Data Sources -> Add data source`.
        - Select `Prometheus`.
        - Set the URL to `http://prometheus:9090`.
        - Click `Save & test`.
    - You can now create dashboards using metrics like `http_requests_total` and `http_request_duration_seconds`.

## 📂 Project Structure

The project follows a clean, layered architecture to separate concerns:

- **`cmd/main.go`**: The entry point of the application, responsible for initialization and startup.
- **`internal/handlers`**: Contains HTTP handlers (Gin) and middleware.
- **`internal/service`**: Contains the core business logic.
- **`internal/repository`**: Responsible for data access and communication with the database.
- **`internal/models`**: Defines the core data structures.
- **`internal/ws`**: Manages WebSocket connections and real-time communication.
- **`internal/metrics`**: Defines and registers Prometheus metrics.
- **`docs/`**: Contains auto-generated Swagger documentation.

---