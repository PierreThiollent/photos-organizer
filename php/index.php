<?php

declare(strict_types=1);

$workingDir = __DIR__;
$scriptPath = __FILE__;
$movedCount = 0;
$errorCount = 0;
$folderCount = 0;

echo "🏞️ Photos/videos organizer" . PHP_EOL;
echo "Processing files..." . PHP_EOL . PHP_EOL;

$startedAt = hrtime(true);

if (!chdir($workingDir)) {
    fwrite(STDERR, "Failed to open working directory: {$workingDir}" . PHP_EOL);
    exit(1);
}

$files = glob('*.*');
if ($files === false) {
    fwrite(STDERR, "Failed to read directory: {$workingDir}" . PHP_EOL);
    exit(1);
}

foreach ($files as $file) {
    $path = $workingDir . DIRECTORY_SEPARATOR . $file;

    if ($path === $scriptPath || !is_file($path)) {
        continue;
    }

    $modificationTime = @filemtime($path);
    if ($modificationTime === false) {
        $errorCount++;
        $error = error_get_last();
        $message = $error['message'] ?? 'unknown error';
        fwrite(STDERR, "Failed to get modification date for {$path}: {$message}" . PHP_EOL);
        continue;
    }

    $modificationDate = date('Y-m-d', $modificationTime);
    $newDir = $workingDir . DIRECTORY_SEPARATOR . $modificationDate;

    if (!file_exists($newDir)) {
        if (!@mkdir($newDir)) {
            $errorCount++;
            $error = error_get_last();
            $message = $error['message'] ?? 'unknown error';
            fwrite(STDERR, "Failed to create directory {$newDir}: {$message}" . PHP_EOL);
            continue;
        }

        $folderCount++;
    }

    $newPath = $newDir . DIRECTORY_SEPARATOR . $file;
    if (!@rename($path, $newPath)) {
        $errorCount++;
        $error = error_get_last();
        $message = $error['message'] ?? 'unknown error';
        fwrite(STDERR, "Failed to move file {$path} to {$newPath}: {$message}" . PHP_EOL);
        continue;
    }

    $movedCount++;
}

$elapsedSeconds = (hrtime(true) - $startedAt) / 1_000_000_000;
$filesPerSecond = $elapsedSeconds > 0 ? $movedCount / $elapsedSeconds : 0;
$status = $errorCount === 0 ? 'Success' : 'Completed with errors';
$duration = number_format($elapsedSeconds, 3, '.', '') . ' s';
$throughput = number_format($filesPerSecond, 0, '.', '') . ' files/s';

echo "┌──────────────────────────────────────────┐" . PHP_EOL;
printf("│ %-40s │\n", 'Processing summary');
echo "├──────────────────────────────────────────┤" . PHP_EOL;
printf("│ %-15s : %22s │\n", 'Status', $status);
printf("│ %-15s : %22d │\n", 'Files moved', $movedCount);
printf("│ %-15s : %22d │\n", 'Folders created', $folderCount);
printf("│ %-15s : %22d │\n", 'Errors', $errorCount);
printf("│ %-15s : %22s │\n", 'Duration', $duration);
printf("│ %-15s : %22s │\n", 'Throughput', $throughput);
echo "└──────────────────────────────────────────┘" . PHP_EOL;
