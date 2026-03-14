<?php

function requireAuth(): int
{
    if (empty($_SESSION['user_id'])) {
        http_response_code(401);
        echo json_encode(['error' => 'Unauthorized']);
        exit;
    }
    return (int)$_SESSION['user_id'];
}

function handleRegister(PDO $databaseConnection): void
{
    $data = json_decode(file_get_contents('php://input'), true) ?? [];
    $username = trim($data['username'] ?? '');
    $password = $data['password'] ?? '';

    if (mb_strlen($username) < 3 || strlen($password) < 8) {
        http_response_code(422);
        echo json_encode(['error' => 'Username must be at least 3 characters and password at least 8 characters long']);
        return;
    }

    $stmt = $databaseConnection->prepare('SELECT id FROM users WHERE username = ?');
    $stmt->execute([$username]);

    if ($stmt->fetch()) {
        http_response_code(409);
        echo json_encode(['error' => 'This username is already taken']);
        return;
    }

    $hash = password_hash($password, PASSWORD_BCRYPT);
    $stmt = $databaseConnection->prepare('INSERT INTO users (username, password_hash) VALUES (?, ?)');
    $stmt->execute([$username, $hash]);

    $userId = (int)$databaseConnection->lastInsertId();
    $_SESSION['user_id'] = $userId;
    $_SESSION['username'] = $username;

    http_response_code(201);
    echo json_encode(['id' => $userId, 'username' => $username]);
}

function handleLogin(PDO $databaseConnection): void
{
    $data = json_decode(file_get_contents('php://input'), true) ?? [];
    $username = trim($data['username'] ?? '');
    $password = $data['password'] ?? '';

    $stmt = $databaseConnection->prepare('SELECT id, password_hash FROM users WHERE username = ?');
    $stmt->execute([$username]);
    $user = $stmt->fetch();

    if (!$user || !password_verify($password, $user['password_hash'])) {
        http_response_code(401);
        echo json_encode(['error' => 'Invalid credentials']);
        return;
    }

    $_SESSION['user_id']  = $user['id'];
    $_SESSION['username'] = $username;

    echo json_encode(['id' => (int)$user['id'], 'username' => $username]);
}

function handleLogout(): void
{
    session_destroy();
    echo json_encode(['message' => 'Logged out']);
}

function handleMe(): void
{
    if (empty($_SESSION['user_id'])) {
        http_response_code(401);
        echo json_encode(['error' => 'Unauthorized']);
        return;
    }

    echo json_encode([
        'id' => (int)$_SESSION['user_id'],
        'username' => $_SESSION['username'],
    ]);
}

function deleteAccount(PDO $databaseConnection): void
{
    $userId = requireAuth();

    $stmt = $databaseConnection->prepare('DELETE FROM users WHERE id = ?');
    $stmt->execute([$userId]);
    session_destroy();

    echo json_encode(['message' => 'Account deleted successfully']);
}
