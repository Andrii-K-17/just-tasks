# Just Tasks

A full-stack task manager with user authentication, priorities, and deadlines.

---

## stack

| layer     | technology                               |
|-----------|------------------------------------------|
| frontend  | vue 3 · typescript · vite · tailwind css |
| state     | pinia · vue router 4                     |
| backend   | php 8.2 · pdo · sessions                 |
| database  | mysql 8.0                                |
| infra     | docker compose                           |

---

## features

- register · login · logout · delete account
- create · toggle · inline edit · delete tasks
- priority levels: low / medium / high
- deadline tracking with overdue warnings
- live search · filter by status
- completion statistics with progress ring
- animated transitions
- sql cascade delete on account removal

---

## getting started
```bash
git clone https://github.com/Andrii-K-17/just-tasks.git
cd just-tasks

cp .env.example .env

docker-compose up -d --build
```

open `http://localhost:5173`

---

## project structure
```
just-tasks/
├── docker-compose.yml
├── init.sql
│
├── backend/
│   ├── Dockerfile
│   ├── 000-default.conf
│   ├── .htaccess
│   ├── config.php          # pdo connection
│   ├── index.php           # api router
│   └── handlers/
│       ├── auth.php        # register · login · logout · delete account
│       └── tasks.php       # crud — list · create · update · delete
│
└── frontend/
    ├── index.html
    ├── vite.config.ts
    ├── package.json
    └── src/
        ├── main.ts
        ├── App.vue
        ├── style.css
        ├── types/
        │   └── index.ts        # task · user interfaces
        ├── api/
        │   ├── auth.ts         # fetch wrappers for auth endpoints
        │   └── tasks.ts        # fetch wrappers for task endpoints
        ├── stores/
        │   ├── useAuthStore.ts # auth state · session init
        │   └── useTaskStore.ts # tasks · filters · search · stats
        ├── router/
        │   └── index.ts        # routes · auth guard
        ├── views/
        │   ├── LandingPage.vue
        │   ├── LoginView.vue
        │   └── DashboardView.vue
        └── components/
            ├── TaskForm.vue    # add task — text · priority · deadline
            ├── TaskItem.vue    # inline edit · toggle · delete
            ├── SearchBar.vue   # live search with clear button
            └── StatsModal.vue  # completion ring · priority breakdown
```

---

## api
```
POST   /api/register        create account
POST   /api/login           authenticate
POST   /api/logout          end session
GET    /api/me              current user

GET    /api/tasks           list tasks
POST   /api/tasks           create task
PUT    /api/tasks/:id       update task
DELETE /api/tasks/:id       delete task

DELETE /api/account         delete account + all tasks
```

---

## environment
```bash
DB_HOST=db
DB_NAME=todo_db
DB_USER=appuser
DB_PASSWORD=your_password
```

---

## license

mit
