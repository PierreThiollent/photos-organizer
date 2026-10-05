# tri-photos — PHP

PHP implementation of the [tri-photos](../README.md) organizer. A single script, `index.php`, with no dependencies and no build step.

## Requirements

- PHP 8.0+ (uses `declare(strict_types=1)` and `hrtime`)

## Run

The script works on its own directory (`__DIR__`), so copy it into the folder you want to sort:

```bash
cp index.php /path/to/photos/
cd /path/to/photos
php index.php
```

## Behavior

- Working directory = `__DIR__`.
- Enumerates `*.*` with `glob()`, so **files without a dot are ignored** (this differs slightly from the Go and Rust versions, which skip only extensionless files).
- Skips the script itself and anything that is not a regular file.
- Folder name is the local modification date formatted as `Y-m-d`.
- Errors are logged to stderr via `error_get_last()` and counted; processing continues.

## Benchmark

Median of 5 runs on 20 000 files (~78 MB), Apple M1:

| Duration  | Throughput     |
| --------- | -------------: |
| 2.35 s    | ≈ 8 500 files/s |

PHP lands between Rust and Go here — a reminder that on an I/O-bound workload the runtime's per-file overhead matters more than the language's reputation. See the [root README](../README.md#benchmark) for methodology and the full comparison.

## Notes

- `hrtime(true)` is used for a monotonic nanosecond timer.
- `rename()` is atomic on the same filesystem.
