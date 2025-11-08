# Go Image Resizer

A concurrent image resizing tool written in Go that processes multiple images in parallel.

## Features

- Concurrent image processing with configurable worker pool
- Supports PNG, JPG, and JPEG formats
- Maintains aspect ratio during resize
- Recursive directory scanning
- Performance metrics (processing time, average per image)

## Usage

```bash
go run main.go [flags]
```

### Flags

- `-dir` - Input directory to scan for images (default: current directory)
- `-width` - Target width for resized images (default: 400)
- `-workers` - Number of concurrent workers (default: 5)

### Examples

Resize all images in current directory to 800px width:
```bash
go run main.go -width 800
```

Process images from a specific folder with 10 workers:
```bash
go run main.go -dir ./test-images -workers 10
```

## Output

Resized images are saved to the `output/` directory with the same filename as the original.

## Build

```bash
go build -o image-resizer main.go
./image-resizer -width 600
```

