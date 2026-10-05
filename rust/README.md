# tri-photos — Rust

Rust implementation of the [tri-photos](../README.md) organizer. Entry point is `src/main.rs`.

## Requirements

- Rust 1.85+ (the crate uses edition 2024)
- Dependency: [`chrono`](https://crates.io/crates/chrono) 0.4 (local timezone formatting)

## Build

```bash
cargo build --release
```

## Run

The program works on the directory containing the **binary** (`env::current_exe()`), so copy it next to the files you want to sort:

```bash
cp target/release/tri-photos /path/to/photos/
cd /path/to/photos
./tri-photos
```

## Behavior

- Working directory = parent of `env::current_exe()`.
- Iterates the directory with `std::fs::read_dir`.
- Skips dot-files, files without an extension, non-regular files, and the executable itself.
- Folder name is the local modification date formatted as `%Y-%m-%d` via `chrono::Local`.
- Errors are logged to stderr and counted; processing continues.

## Benchmark

Median of 5 runs on 20 000 files (~78 MB), Apple M1:

| Duration  | Throughput      |
| --------- | --------------: |
| 1.70 s    | ≈ 11 700 files/s |

The fastest of the three in this benchmark. See the [root README](../README.md#benchmark) for methodology and the full comparison.

## Notes

- `std::fs::rename` maps to `rename(2)` and is atomic on the same filesystem.
- `DateTime::<Local>::from(...)` converts the system `SystemTime` into local time for folder naming.
