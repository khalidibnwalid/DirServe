package templviews

import (
	"fmt"
	"path"
	"strings"
)

type FileItem struct {
	Name    string
	IsDir   bool
	Size    int64
	Path    string
	IsImage bool
	IsVideo bool
	IsAudio bool
}

func isImage(name string) bool {
	ext := strings.ToLower(path.Ext(name))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp" || ext == ".svg"
}

func isVideo(name string) bool {
	ext := strings.ToLower(path.Ext(name))
	return ext == ".mp4" || ext == ".webm" || ext == ".ogv" || ext == ".mov" || ext == ".mkv" || ext == ".avi"
}

func isAudio(name string) bool {
	ext := strings.ToLower(path.Ext(name))
	return ext == ".mp3" || ext == ".ogg" || ext == ".wav" || ext == ".flac"
}

func formatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
	}
	return fmt.Sprintf("%.1f GB", float64(size)/(1024*1024*1024))
}

type breadcrumb struct {
	Name string
	Path string
}
