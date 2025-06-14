package templviews

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/a-h/templ"
)

func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path

		rootDir := "."
		for i, arg := range os.Args {
			if arg == "-dir" && i+1 < len(os.Args) {
				rootDir = os.Args[i+1]
				break
			}
		}

		absPath, err := filepath.Abs(rootDir)
		if err != nil {
			http.Error(w, "Error getting absolute path", http.StatusInternalServerError)
			return
		}

		// Get the physical file path
		physicalPath := filepath.Join(absPath, filepath.FromSlash(requestPath))

		// Check if file exists
		fileInfo, err := os.Stat(physicalPath)
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Error accessing file", http.StatusInternalServerError)
			return
		}

		breadcrumbs := generateBreadcrumbs(requestPath)

		// If directory, show gallery
		if fileInfo.IsDir() {
			// Debug info
			log.Printf("Serving directory: %s (Physical path: %s)", requestPath, physicalPath)

			files, err := os.ReadDir(physicalPath)
			if err != nil {
				http.Error(w, "Error reading directory", http.StatusInternalServerError)
				return
			}

			// Convert files to FileItems
			fileItems := []FileItem{}
			for _, file := range files {
				// Skip hidden files
				if strings.HasPrefix(file.Name(), ".") {
					continue
				}

				info, err := file.Info()
				if err != nil {
					continue
				}

				fileItem := FileItem{
					Name:    file.Name(),
					Path:    path.Join(requestPath, file.Name()),
					IsDir:   file.IsDir(),
					Size:    info.Size(),
					IsImage: isImage(file.Name()),
					IsVideo: isVideo(file.Name()),
					IsAudio: isAudio(file.Name()),
				}
				fileItems = append(fileItems, fileItem)
			}

			// Sort files: directories first, then by name
			sort.Slice(fileItems, func(i, j int) bool {
				if fileItems[i].IsDir != fileItems[j].IsDir {
					return fileItems[i].IsDir
				}
				return fileItems[i].Name < fileItems[j].Name
			})
			w.Header().Set("Content-Type", "text/html; charset=utf-8")

			log.Printf("Rendering gallery with %d files, path: %s", len(fileItems), requestPath)

			component := Gallery(requestPath, breadcrumbs, fileItems)
			templ.Handler(component).ServeHTTP(w, r)
		} else {
			// It's a file, show file viewer
			fileItem := FileItem{
				Name:    fileInfo.Name(),
				Path:   requestPath,
				IsDir:   fileInfo.IsDir(),
				Size:    fileInfo.Size(),
				IsImage: isImage(fileInfo.Name()),
				IsVideo: isVideo(fileInfo.Name()),
				IsAudio: isAudio(fileInfo.Name()),
			} // Render file viewer
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			component := FileViewer(fileItem, breadcrumbs)
			templ.Handler(component).ServeHTTP(w, r)
		}
	})
}

func generateBreadcrumbs(urlPath string) []breadcrumb {
	var breadcrumbs []breadcrumb

	// Root
	breadcrumbs = append(breadcrumbs, breadcrumb{
		Name: "Home",
		Path: "/",
	})

	segments := strings.Split(strings.Trim(urlPath, "/"), "/")
	if len(segments) > 0 && segments[0] != "" {
		currentPath := ""
		for _, segment := range segments {
			currentPath = path.Join(currentPath, segment)
			breadcrumbs = append(breadcrumbs, breadcrumb{
				Name: segment,
				Path: currentPath,
			})
		}
	}

	return breadcrumbs
}
