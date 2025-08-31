ToDoList - A task management application.

Features:

    Add new tasks

    View the list of all tasks

    Mark tasks as completed

    Delete tasks

Preview: https://github.com/sudolicious/todolist/blob/main/frontend/public/Screenshot.png?raw=true

Technology Stack:

    PostgreSQL 16

    Backend: Go 1.21+

    Frontend: React 18+, TypeScript 4.8.5

    CI/CD: Azure DevOps Pipelines

Installation and Setup:

    Clone the repository
    git clone https://github.com/sudolicious/todolist.git
    cd todolist

    Run PostgreSQL in Docker
    docker run --name postgres
    -e POSTGRES_USER=your_username
    -e POSTGRES_PASSWORD=your_password
    -e POSTGRES_DB=your_database
    -p 5432:5432
    -d postgres:16

    Start the backend
    cd backend
    cp .env.example .env
    go mod download
    go run main.go
    API: http://localhost:8080/api/tasks

    Run tests
    Unit tests:
    go test -v ./...

Integration tests:
go test -v -tags=integration ./test_integration/ -timeout 30s

    Start the frontend
    cd frontend
    npm install
    npm run build
    npm install -g serve
    serve -s build
    Frontend: http://localhost:3000

Docker Setup:

Backend Docker:
cd backend
docker build -t todolist-backend .
docker run -p 8080:8080 --env-file .env todolist-backend

Frontend Docker:
cd frontend
docker build -t todolist-frontend .
docker run -p 3000:80 todolist-frontend

CI/CD Pipeline:
The project includes Azure DevOps CI/CD pipeline configuration (azure-pipelines.yml) that automates versioning, security scanning, Docker building, integration testing, and deployment.

Integration Tests:
Integration tests are located in the test_integration/ directory. They verify database connectivity, API functionality, and service integration.

Run integration tests:
cd backend
go test -v -tags=integration ./test_integration/ -timeout 30s

Project Structure:
todolist/
├── backend/
│ ├── go.mod
│ ├── go.sum
│ ├── main.go
│ ├── Dockerfile
│ ├── test_integration/
│ └── migrations/
├── frontend/
│ ├── Dockerfile
│ ├── build/
│ ├── node_modules/
│ ├── public/
│ ├── src/
│ ├── package.json
│ └── tsconfig.json
├── azure-pipelines.yml
├── openapi.yml
└── README.md
