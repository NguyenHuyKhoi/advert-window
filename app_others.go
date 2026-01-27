//go:build !windows

package main

import "fmt"

// Hàm này sẽ được gọi khi chạy trên macOS hoặc Linux
func (a *App) enableAutoStart() {
    fmt.Println("Auto-start is only supported on Windows. Skipping...")
}