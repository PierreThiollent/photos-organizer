# tri-photos — Go

Go implementation of the [tri-photos](../README.md) organizer. A single file, `main.go`, using only the standard library.

## Requirements

- Go 1.26+ (see `go.mod`)

## Build

```bash
go build -o tri-photos .
```

## Run

The program works on the directory containing the **binary**, so copy it next to the files you want to sort:

```bash
cp tri-photos /path/to/photos/
cd /path/to/photos
./tri-photos
```

## Behavior

- Working directory = `filepath.Dir(os.Executable())`.
- Enumerates the directory once with `os.ReadDir` before moving anything.
- Skips dot-files, files without an extension (e.g. `README`), subdirectories, and the executable itself.
- Folder name is the local modification date formatted as `YYYY-MM-DD`.
- Errors are logged to stderr and counted; processing continues.

## Benchmark

Median of 5 runs on 20 000 files (~78 MB), Apple M1:

| Duration  | Throughput     |
| --------- | -------------: |
| 2.49 s    | ≈ 8 000 files/s |

See the [root README](../README.md#benchmark) for methodology and the full comparison.

## Notes

- `os.Rename` is atomic on the same filesystem, which is what makes the move cheap.
- The date format is Go's reference layout `2006-01-02`.
