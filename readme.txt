ToDoList - приложение для управления списком задач.

Функционал:
- Добавление новых задач;
- Просмотр списка всех задач;
- Выделение задачи выполненной;
- Удаление задачи

![ToDoList Preview](frontend/public/Screenshot.png)

Технологический стек:

**Backend:** 
- Go 1.21+
- PostgreSQL 16 (реляционная БД)

**Frontend:**
- React 18+
- TypeScript 4.8.5

Установка и запуск.

1. Клонирование репозитория
git clone https://github.com/sudolicious/todolist.git
cd todolist

2. Запуск PostgreSQL в Docker
docker run --name postgres \
  -e POSTGRES_USER=your_username \
  -e POSTGRES_PASSWORD=your_password \
  -e POSTGRES_DB=your_database \
  -p 5432:5432 \
  -d postgres:16

3. Бэкенд
cd backend
cp .env.example .env  # заполните переменные для БД
go mod download
go run main.go

4. Фронтенд
cd frontend
npm install
npm run build
serve -s build

Структура проекта:
todolist/
├── backend/            # Go-бэкенд
│   ├── go.mod          # Модули Go
│   ├── go.sum          # Зависимости
│   ├── main.go         # Точка входа
│   └── migrations/     # Миграции БД
│
├── frontend/           # React-фронтенд
│   ├── build/          # Собранный проект (после npm run build)
│   ├── node_modules/   # Зависимости npm
│   ├── public/         # Статические файлы
│   ├── src/            # Исходники
│   ├── package.json    # Зависимости
│   └── tsconfig.json   # Настройки TypeScript
│
├── openapi.yml         # OpenAPI спецификация
└── README.md           # Этот файл

