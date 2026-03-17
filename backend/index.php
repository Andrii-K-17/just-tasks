<?php

session_start();

header('Content-Type: application/json');
header('Access-Control-Allow-Origin: http://localhost:5173');
header('Access-Control-Allow-Credentials: true');
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type');

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(204);
    exit;
}

require_once __DIR__ . '/config.php';
require_once __DIR__ . '/handlers/auth.php';
require_once __DIR__ . '/handlers/tasks.php';

$method = $_SERVER['REQUEST_METHOD'];
$path = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);
$parts = array_values(array_filter(explode('/', $path)));

// URL: /api/{resource}/{id?}
$resource = $parts[1] ?? '';
$id = isset($parts[2]) ? (int)$parts[2] : null;

match (true) {
    $resource === 'register' && $method === 'POST'   => handleRegister($databaseConnection),
    $resource === 'login'    && $method === 'POST'   => handleLogin($databaseConnection),
    $resource === 'logout'   && $method === 'POST'   => handleLogout(),
    $resource === 'me'       && $method === 'GET'    => handleMe(),
    $resource === 'tasks'    && $method === 'GET'    => getTasks($databaseConnection),
    $resource === 'tasks'    && $method === 'POST'   => createTask($databaseConnection),
    $resource === 'tasks'    && $method === 'PUT'    => updateTask($databaseConnection, $id),
    $resource === 'tasks'    && $method === 'DELETE' => deleteTask($databaseConnection, $id),
    $resource === 'tasks'    && $method === 'PUT' && $parts[2] === 'reorder' => reorderTasks($pdo),
    $resource === 'account'  && $method === 'DELETE' => deleteAccount($databaseConnection),
    default => (static function () {
        http_response_code(404);
        echo json_encode(['error' => 'Route not found']);
    })(),
};
