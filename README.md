# tri-photos

> The same tiny file organizer, written three times — in **Go**, **PHP** and **Rust** — to compare how each language handles a real, I/O-bound task.

**tri-photos** scans the directory it lives in, groups every file by its **last-modified date**, and moves it into a subfolder named `YYYY-MM-DD`. It then prints a summary: files moved, folders created, errors, duration and throughput.

It is a deliberately small program. The point is not the feature, it is the comparison.

## The task

Turn a flat folder full of photos/videos:

```
photos/
├── photo_00001.bmp   # modified 2024-01-01
├── photo_00002.bmp   # modified 2024-09-11
└── photo_00003.bmp   # modified 2024-10-16
```

into a folder tree sorted by date:

```
photos/
├── 2024-01-01/
│   └── photo_00001.bmp
├── 2024-09-11/
│   └── photo_00002.bmp
└── 2024-10-16/
    └── photo_00003.bmp
```

## Benchmark

Each implementation runs the same algorithm on the same dataset. Numbers below are the **median of 5 interleaved runs** of the program's own reported duration.

| Implementation | Directory | Median duration |       Throughput |  Relative |
| -------------- | --------- | --------------: | ---------------: | --------: |
| 🦀 Rust        | `rust/`   |      **1.70 s** | ≈ 11 700 files/s | **1.00×** |
| 🐘 PHP         | `php/`    |          2.35 s |  ≈ 8 500 files/s |     1.38× |
| 🐹 Go          | `go/`     |          2.49 s |  ≈ 8 000 files/s |     1.46× |

**Setup**

- **Dataset:** 20 000 BMP files (~78 MB), spread over 365 distinct modification dates.
- **Machine:** Apple M1 (8 cores, 16 GB RAM), macOS 15.7.1, APFS filesystem.
- **Toolchains:** Go 1.27.1, PHP 8.5.9, Rust 1.97.1 (`--release`).
- **Method:** 5 runs per language, interleaved to smooth out cache/thermal drift; median reported. Each run starts from a fresh copy of the dataset with timestamps preserved.

**How to read this.** The workload is 20 000 `rename()` calls — it is **I/O-bound, not CPU-bound**. The differences come mostly from syscall patterns, per-file allocations and runtime overhead, not from raw language speed. A different filesystem, file count or directory layout can easily change the ranking. Treat this as "on this machine, for this workload", not as a general verdict on the languages.

Example output (Rust run):

```
🏞️ Photos/videos organizer
Processing files...

┌──────────────────────────────────────────┐
│ Processing summary                       │
├──────────────────────────────────────────┤
│ Status          :                Success │
│ Files moved     :                  20000 │
│ Folders created :                    365 │
│ Errors          :                      0 │
│ Duration        :                1.718 s │
│ Throughput      :          11640 files/s │
└──────────────────────────────────────────┘
```

### Reproducing the benchmark

```bash
# 1. Build (PHP needs no build step)
(cd go   && go build -o tri-photos .)
(cd rust && cargo build --release)

# 2. For each language, copy the dataset + program into a scratch folder and run
cp -a benchmark-photos /tmp/run-go
cp go/tri-photos /tmp/run-go/
(cd /tmp/run-go && ./tri-photos)
```

## Repository layout

```
tri-photos/
├── go/                # Go implementation
├── php/               # PHP implementation
├── rust/              # Rust implementation
├── benchmark-photos/  # 20 000-file dataset (git-ignored, not committed)
└── README.md
```

Each subdirectory has its own README with build/run instructions and implementation notes:

- [Go → `go/README.md`](go/README.md)
- [PHP → `php/README.md`](php/README.md)
- [Rust → `rust/README.md`](rust/README.md)

## Requirements

| Language | Version | Notes                           |
| -------- | ------- | ------------------------------- |
| Go       | 1.26+   | standard library only           |
| PHP      | 8.0+    | standard library only           |
| Rust     | 1.85+   | edition 2024, depends on chrono |

## Quick start

Each program must be **copied into the folder you want to organize** and run from there — that folder is its working directory.

```bash
# Go
(cd go && go build -o tri-photos .)
cp go/tri-photos /path/to/photos/
(cd /path/to/photos && ./tri-photos)

# Rust
(cd rust && cargo build --release)
cp rust/target/release/tri-photos /path/to/photos/
(cd /path/to/photos && ./tri-photos)

# PHP (no build step)
cp php/index.php /path/to/photos/
(cd /path/to/photos && php index.php)
```

> ⚠️ **This tool moves files.** Always try it on a copy first — there is no undo.

Skipped entries, in all implementations: the program itself, dot-files, files without an extension, and subdirectories.

## Recreating the dataset

The 20 000 `photo_*.bmp` files (~78 MB) are intentionally **not committed** (see `.gitignore`). To rebuild an equivalent dataset:

```bash
mkdir -p benchmark-photos && cd benchmark-photos

for i in $(seq -w 1 20000); do
  f="photo_${i}.bmp"
  printf 'BM' > "$f"                                   # minimal BMP header
  m=$(printf '%02d' $((1 + RANDOM % 12)))              # random month in 2024
  d=$(printf '%02d' $((1 + RANDOM % 28)))              # random day
  touch -t "2024${m}${d}120000" "$f"                   # spread over the year
done
```
