# TASKFLOW

## Overview

A full-featured task management system with:

- User authentication
- Task CRUD operations
- Comment functionality
- Statistics dashboard
- Priority-based task organization

## Features

✅ User authentication & authorization
✅ Task management with due dates & priorities
✅ Comment system for task discussions
✅ Comprehensive statistics dashboard
✅ Responsive frontend built with Next.js
✅ RESTful API backend in Go

## Tech Stack

**Frontend**:

- Next.js
- TypeScript
- Tailwind CSS

**Backend**:

- Go
- MongoDB
- JWT Authentication

## Getting Started

### Prerequisites

- Go 1.20+
- Node.js 18+
- MongoDB

### Installation

1. Clone the repository:

```bash
git clone https://github.com/YeabTesfaye/TaskFlow.git
```

2. Set up backend:

```bash
cd api
cp .env-example .env
# Edit .env with your configuration
go mod download
go run main.go
```

3. Set up frontend:

```bash
cd client
npm install
npm run dev
```

## API Documentation

### Authentication

`POST /api/auth/login` - User login
`POST /api/auth/register` - User registration

### Tasks

`GET /api/tasks` - List all tasks
`POST /api/tasks` - Create new task
`GET /api/tasks/:id` - Get task details

## Frontend Routes

- `/` - Home page
- `/dashboard` - Statistics dashboard
- `/task/[id]` - Task details
- `/new` - Create new task

## Contributing

Pull requests are welcome. For major changes, please open an issue first.

## License

[MIT](https://choosealicense.com/licenses/mit/)

```

Would you like me to add any additional sections or modify any part of this template?

```
