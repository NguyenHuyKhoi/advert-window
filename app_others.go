//go:build !windows

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type APIResponse struct {
	Success bool    `json:"success"`
	Data    AppData `json:"data"`
}

type AppData struct {
	Version int    `json:"version"`
	URL     string `json:"url"`
}

func (a *App) enableAutoStart() {
	fmt.Println("[Mac/Debug] Auto-start logic skipped")
}

func (a *App) silentUpdate() {
	fmt.Printf("[CheckUpdate] Đang gọi API: %s\n", CHECK_UPDATE_URL)
	
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(CHECK_UPDATE_URL)
	if err != nil {
		fmt.Printf("[CheckUpdate] Lỗi kết nối: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var apiRes APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiRes); err != nil {
		fmt.Printf("[CheckUpdate] Lỗi giải mã JSON: %v\n", err)
		return
	}

	// Log toàn bộ dữ liệu nhận được để kiểm tra
	fmt.Printf("[CheckUpdate] Kết quả từ Server -> Success: %v, Version: %d, URL: %s\n", 
		apiRes.Success, apiRes.Data.Version, apiRes.Data.URL)

	if apiRes.Success && apiRes.Data.Version > AppVersionInt {
		if apiRes.Data.URL == "" {
			fmt.Println("[CheckUpdate] CẢNH BÁO: Có version mới nhưng URL trống!")
		} else {
			fmt.Printf("[CheckUpdate] Tìm thấy bản cập nhật mới! Server: %d, Local: %d\n", 
				apiRes.Data.Version, AppVersionInt)
		}
	} else {
		fmt.Println("[CheckUpdate] Hiện tại chưa có bản cập nhật mới.")
	}
}