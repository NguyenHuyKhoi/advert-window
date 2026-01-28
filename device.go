package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
)

type DeviceInfo struct {
	DeviceID    string `json:"device_id"`
	DeviceModel string `json:"device_model"`
	DeviceName  string `json:"device_name"`
	OSVersion   string `json:"os_version"`
}

func getDeviceInfo() DeviceInfo {
	hostname, _ := os.Hostname()
	
	return DeviceInfo{
		DeviceID:    getStableDeviceID(),
		DeviceModel: runtime.GOARCH, // Trả về kiến trúc CPU (amd64, arm64...)
		DeviceName:  hostname,       // Tên máy tính
		OSVersion:   runtime.GOOS,   // Hệ điều hành (windows, darwin, linux...)
	}
}

func getStableDeviceID() string {
	dir, _ := os.UserConfigDir()
	path := filepath.Join(dir, "advert", "device_id")

	if b, err := os.ReadFile(path); err == nil {
		return string(b)
	}

	id := uuid.New().String()
	_ = os.MkdirAll(filepath.Dir(path), 0700)
	_ = os.WriteFile(path, []byte(id), 0600)
	return id
}