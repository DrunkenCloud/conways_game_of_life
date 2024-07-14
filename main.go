package main

import (
    "fmt"
    "math/rand"
    "net/http"
    "sync"
    "time"
)

const (
    width  = 200
    height = 100
)

var (
    grid       [height][width]bool
    gridMutex  sync.Mutex
)

func main() {
    rand.Seed(time.Now().UnixNano())
    initializeGrid()

    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)

    http.HandleFunc("/step", stepHandler)
    http.HandleFunc("/grid", gridHandler)

    fmt.Println("Starting server at :8080")
    http.ListenAndServe(":8080", nil)
}

func initializeGrid() {
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            grid[y][x] = rand.Float64() < 0.2
        }
    }
}

func stepHandler(w http.ResponseWriter, r *http.Request) {
    gridMutex.Lock()
    defer gridMutex.Unlock()

    nextGrid := [height][width]bool{}

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            aliveNeighbours := countAliveNeighbors(x, y)
            if grid[y][x] && aliveNeighbours < 2 {
                nextGrid[y][x] = false
            } else if grid[y][x] && (aliveNeighbours == 2 || aliveNeighbours == 3) {
                nextGrid[y][x] = true
            } else if grid[y][x] && aliveNeighbours > 3 {
                nextGrid[y][x] = false
            } else if !grid[y][x] && aliveNeighbours == 3 {
                nextGrid[y][x] = true
            } else {
                nextGrid[y][x] = grid[y][x]
            }
        }
    }

    grid = nextGrid

    renderGrid(w)
}

func gridHandler(w http.ResponseWriter, r *http.Request) {
    gridMutex.Lock()
    defer gridMutex.Unlock()
    renderGrid(w)
}

func renderGrid(w http.ResponseWriter) {
    fmt.Fprintf(w, "<div>")
    for y := 0; y < height; y++ {
        fmt.Fprintf(w, "<div class='grid-row'>")
        for x := 0; x < width; x++ {
            class := ""
            if grid[y][x] {
                class = "alive"
            }
            fmt.Fprintf(w, "<div class='cell %s'></div>", class)
        }
        fmt.Fprintf(w, "</div>")
    }
    fmt.Fprintf(w, "</div>")
}


func countAliveNeighbors(x, y int) int {
    count := 0
    for dy := -1; dy <= 1; dy++ {
        for dx := -1; dx <= 1; dx++ {
            if dx == 0 && dy == 0 {
                continue
            }
            nx, ny := x+dx, y+dy
            if nx >= 0 && nx < width && ny >= 0 && ny < height && grid[ny][nx] {
                count++
            }
        }
    }
    return count
}
