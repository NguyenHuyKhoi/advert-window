//go:build windows
package main

import (
    "os"
    "golang.org/x/sys/windows/registry"
)

func (a *App) enableAutoStart() {
    exe, err := os.Executable()
    if err != nil {
        return
    }
    
    // Mở registry để đăng ký ứng dụng khởi chạy cùng Windows
    key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
    if err != nil {
        return
    }
    defer key.Close()
    
    // Đặt giá trị thực thi
    _ = key.SetStringValue("ForlifeMediaPlayer", exe)
}