ToDoList - A task management application.

Features:

- Add new tasks;
- View the list of all tasks;
- Mark tasks as completed;
- Delete tasks.

Preview:
https://github.com/sudolicious/todolist/blob/main/frontend/public/Screenshot.png?raw=true

Technology Stack:
- PostgreSQL 16
- Backend: Go 1.21+
- Frontend:
- React 18+
- TypeScript 4.8.5

Installation and Setup.

    1. Clone the repository
    git clone https://github.com/sudolicious/todolist.git
    cd todolist

    2. Run PostgreSQL in Docker
    docker run --name postgres \
    -e POSTGRES_USER=your_username \
    -e POSTGRES_PASSWORD=your_password \
    -e POSTGRES_DB=your_database \
    -p 5432:5432 \
    -d postgres:16

    3. Start the backend
    cd backend
    cp .env.example .env # fill in the database variables
    go mod download
    go run main.go
    API: http://localhost:8080/api/tasks

    4. Start the frontend
    cd frontend
    npm install
    npm run build
    serve -s build
    Frontend: http://localhost:3000

Project Structure:
todolist/
├── backend/ # Go backend
│ ├── go.mod # Go modules
│ ├── go.sum # Dependencies
│ ├── main.go # Entry point
│ └── migrations/ # Database migrations
│
├── frontend/ # React frontend
│ ├── build/ # Build output
│ ├── node_modules/ # npm dependencies
│ ├── public/ # Static files
│ ├── src/ # Source files
│ ├── package.json # Dependencies
│ └── tsconfig.json # TypeScript settings
│
├── openapi.yml # OpenAPI specification
└── README.md # This file
