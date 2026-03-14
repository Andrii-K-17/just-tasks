<?php

function getTasks(PDO $databaseConnection): void
{
    $userId = requireAuth();

    $stmt = $databaseConnection->prepare(
        'SELECT id, task_text, priority, deadline, is_completed, created_at
         FROM tasks
         WHERE user_id = ?
         ORDER BY id DESC'
    );
    $stmt->execute([$userId]);

    $tasks = array_map(
        static fn($task) => [
            ...$task,
            'id' => (int)$task['id'],
            'priority' => (int) $task['priority'],
            'is_completed' => (bool)$task['is_completed'],
        ],
        $stmt->fetchAll()
    );

    echo json_encode($tasks);
}

function createTask(PDO $databaseConnection): void
{
    $userId = requireAuth();
    $data = json_decode(file_get_contents('php://input'), true) ?? [];

    $text = trim($data['task_text'] ?? '');
    $priority = (int) ($data['priority'] ?? 2);
    $deadline = ($data['deadline'] ?? '') ?: null;

    if (!$text) {
        http_response_code(422);
        echo json_encode(['error' => 'The task text field is required']);
        return;
    }

    if (!in_array($priority, [1, 2, 3], true)) {
        $priority = 2;
    }

    $stmt = $databaseConnection->prepare(
        'INSERT INTO tasks (user_id, task_text, priority, deadline) VALUES (?, ?, ?, ?)'
    );
    $stmt->execute([$userId, $text, $priority, $deadline]);

    http_response_code(201);
    echo json_encode([
        'id'           => (int)$databaseConnection->lastInsertId(),
        'task_text'    => $text,
        'priority'     => $priority,
        'deadline'     => $deadline,
        'is_completed' => false,
        'created_at'   => date('Y-m-d H:i:s'),
    ]);
}

function updateTask(PDO $databaseConnection, ?int $id): void
{
    $userId = requireAuth();

    if (!$id) {
        http_response_code(400);
        echo json_encode(['error' => 'Task ID required']);
        return;
    }

    $data = json_decode(file_get_contents('php://input'), true) ?? [];
    $fields = [];
    $values = [];

    if (array_key_exists('task_text', $data)) {
        $text = trim($data['task_text']);
        if (!$text) {
            http_response_code(422);
            echo json_encode(['error' => 'The task text field is required']);
            return;
        }
        $fields[] = 'task_text = ?';
        $values[] = $text;
    }

    if (array_key_exists('is_completed', $data)) {
        $fields[] = 'is_completed = ?';
        $values[] = (int)$data['is_completed'];
    }

    if (array_key_exists('priority', $data)) {
        $priority = (int) $data['priority'];
        $fields[] = 'priority = ?';
        $values[] = in_array($priority, [1, 2, 3], true) ? $priority : 2;
    }

    if (array_key_exists('deadline', $data)) {
        $fields[] = 'deadline = ?';
        $values[] = ($data['deadline'] ?: null);
    }

    if (!$fields) {
        http_response_code(422);
        echo json_encode(['error' => 'No fields to update']);
        return;
    }

    $values[] = $id;
    $values[] = $userId;

    $stmt = $databaseConnection->prepare(
        'UPDATE tasks SET ' . implode(', ', $fields) . ' WHERE id = ? AND user_id = ?'
    );
    $stmt->execute($values);

    echo json_encode(['updated' => $stmt->rowCount() > 0]);
}

function deleteTask(PDO $databaseConnection, ?int $id): void
{
    $userId = requireAuth();

    if (!$id) {
        http_response_code(400);
        echo json_encode(['error' => 'Task ID required']);
        return;
    }

    $stmt = $databaseConnection->prepare('DELETE FROM tasks WHERE id = ? AND user_id = ?');
    $stmt->execute([$id, $userId]);

    echo json_encode(['deleted' => $stmt->rowCount() > 0]);
}
