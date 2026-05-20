CREATE TABLE IF NOT EXISTS users (
    id            SERIAL       PRIMARY KEY,
    username      TEXT         UNIQUE NOT NULL,
    password_hash TEXT         NOT NULL,
    created_at    TIMESTAMPTZ  DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS categories (
    id      SERIAL       PRIMARY KEY,
    user_id INTEGER      REFERENCES users(id) ON DELETE CASCADE,
    name    VARCHAR(50)  NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
    id           SERIAL       PRIMARY KEY,
    user_id      INT          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_text    TEXT         NOT NULL,
    priority     SMALLINT     NOT NULL DEFAULT 2 CHECK (priority BETWEEN 1 AND 3),
    deadline     DATE,
    is_completed BOOLEAN      NOT NULL DEFAULT FALSE,
    position     INT          NOT NULL DEFAULT 0,
    category_id  INTEGER      REFERENCES categories(id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ  DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_user_position ON tasks(user_id, position);
