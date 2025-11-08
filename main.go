package main

import (
    "flag"
    "fmt"
    "image"
    "image/jpeg"
    "image/png"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Job struct {
    Path string
}

type Result struct {
    Path     string
    Duration time.Duration
    Err      error
}

func main() {
    dir := flag.String("dir", ".", "input folder")
    width := flag.Int("width", 400, "resize width")
    workers := flag.Int("workers", 5, "worker count")
    flag.Parse()

    start := time.Now()

    images := findImages(*dir)
    jobs := make(chan Job, len(images))
    results := make(chan Result, len(images))

    var wg sync.WaitGroup

    for i := 0; i < *workers; i++ {
        wg.Add(1)
        go resizeWorker(i, jobs, results, &wg, *width)
    }

    for _, p := range images {
        jobs <- Job{Path: p}
    }
    close(jobs)

    go func() {
        wg.Wait()
        close(results)
    }()

    var count int
    var total time.Duration

    for r := range results {
        count++
        total += r.Duration
        if r.Err != nil {
            fmt.Println("error:", r.Path, r.Err)
        }
    }

    fmt.Println("processed:", count)
    if count > 0 {
        fmt.Println("avg per image:", total/time.Duration(count))
    }
    fmt.Println("total time:", time.Since(start))
}

func findImages(dir string) []string {
    list := []string{}
    filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }
        if info.IsDir() {
            return nil
        }
        ext := filepath.Ext(path)
        switch ext {
        case ".png", ".jpg", ".jpeg":
            list = append(list, path)
        }
        return nil
    })
    return list
}

func resizeWorker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, width int) {
    defer wg.Done()

    os.MkdirAll("output", 0755)

    for job := range jobs {
        t0 := time.Now()

        f, err := os.Open(job.Path)
        if err != nil {
            results <- Result{Path: job.Path, Duration: 0, Err: err}
            continue
        }

        img, format, err := image.Decode(f)
        f.Close()
        if err != nil {
            results <- Result{Path: job.Path, Duration: 0, Err: err}
            continue
        }

        resized := resize(img, width)

        outPath := filepath.Join("output", filepath.Base(job.Path))
        out, err := os.Create(outPath)
        if err != nil {
            results <- Result{Path: job.Path, Duration: 0, Err: err}
            continue
        }

        switch format {
        case "png":
            err = png.Encode(out, resized)
        default:
            err = jpeg.Encode(out, resized, &jpeg.Options{Quality: 90})
        }
        out.Close()

        results <- Result{Path: job.Path, Duration: time.Since(t0), Err: err}
    }
}

func resize(src image.Image, targetW int) image.Image {
    b := src.Bounds()
    w := b.Dx()
    h := b.Dy()

    if w <= targetW {
        return src
    }

    scale := float64(targetW) / float64(w)
    targetH := int(float64(h) * scale)

    dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))

    for y := 0; y < targetH; y++ {
        sy := int(float64(y) / scale)
        for x := 0; x < targetW; x++ {
            sx := int(float64(x) / scale)
            dst.Set(x, y, src.At(sx, sy))
        }
    }

    return dst
}
