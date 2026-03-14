<?php

$databaseConnection = null;

try {
    $connectionString = sprintf(
        'mysql:host=%s;dbname=%s;charset=utf8mb4',
        getenv('DB_HOST') ?: 'db',
        getenv('DB_NAME') ?: 'todo_db'
    );

    $databaseConnection = new PDO(
        $connectionString,
        getenv('DB_USER')     ?: 'appuser',
        getenv('DB_PASSWORD') ?: '',
        [
            PDO::ATTR_ERRMODE            => PDO::ERRMODE_EXCEPTION,
            PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
            PDO::ATTR_EMULATE_PREPARES   => false,
        ]
    );
} catch (PDOException $e) {
    http_response_code(500);
    echo json_encode(['error' => 'Database connection failed']);
    exit;
}