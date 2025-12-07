package server

// TODO: new directory added at runtime? Add to watches, too.
//
// TODO: well probably need IN_MODIFY -> seems all editors (test vscode, nvim, helix) at least perform a MODIFY on the relevant file when saving
// They all emit IN_CLOSE_WRITE on save, too. This is also triggered when creating, although we should probably handle that separately with IN_CREATE

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	// "time"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/flonle/mdbuddy/renderer"
)

type previewServer struct {
	inotifyInstanceFD int            // The inotify instance file descriptor
	watches           map[int]string // watch descriptor : (absolute) path
	previewFile       string         // The last changed markdown file
}

// Start a server that serves a preview of the last changed file
// amongst all the given files. Binds to `addr`.
//
// The server will also *watch* all given file for write events. When
// detected, the server will (re)render the affected file, and show that instead.
func ServePreview(addr string, paths []string) error {
	// Initialize inotify instance
	inotifyfd, err := unix.InotifyInit1(unix.IN_CLOEXEC)
	if err != nil {
		return fmt.Errorf("Failed to initialize inotify: %v", err)
	}
	defer unix.Close(inotifyfd)

	// Create previewServer
	server := &previewServer{
		inotifyInstanceFD: inotifyfd,
		watches:           map[int]string{},
		previewFile:       "",
	}

	// Add watches to inotify instance
	absPaths, err := normalizePaths(paths)
	if err != nil {
		return err
	}
	for _, absPath := range absPaths {
		server.addWatchRecursively(absPath)
	}

	// Start HTTP server
	http.HandleFunc("/", server.servePreview)
	http.HandleFunc("/sse-refresh", server.handleSSERefresh)
	addr_, port, _ := strings.Cut(addr, ":")
	if addr_ == "" {
		addr_ = "localhost"
	}
	log.Printf("Preview server running on http://%s:%s\n", addr_, port)
	return http.ListenAndServe(addr, nil)
}

func (s *previewServer) servePreview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	input, err := os.ReadFile(s.previewFile)
	if err != nil {
		input = []byte("# Live Preview\n\nPlease write to a watched file to see its preview.")
	}
	renderer.RenderBareNote(input, w)
}

func (s *previewServer) handleSSERefresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	buf := make([]byte, unix.SizeofInotifyEvent*4096)
	for {
		n, err := unix.Read(s.inotifyInstanceFD, buf)
		if err != nil {
			log.Fatalf("Failed to read events: %v", err)
		}

		offset := 0
		for offset < n {
			event := (*unix.InotifyEvent)(unsafe.Pointer(&buf[offset]))
			path := s.watches[int(event.Wd)]

			// Extract filename if present
			nameLen := int(event.Len)
			var filename string
			if nameLen > 0 {
				nameBytes := buf[offset+unix.SizeofInotifyEvent : offset+unix.SizeofInotifyEvent+nameLen]
				filename = string(nameBytes[:clen(nameBytes)])
			}

			if strings.HasSuffix(filename, ".md") {
				log.Println("sending refresh signal")
				s.previewFile = filepath.Join(path, filename)
				w.Write([]byte("data: refresh\n\n"))
				flusher.Flush()
			}

			// Print event description
			fmt.Printf("[%s] ", path)
			if filename != "" {
				fmt.Printf("%s: ", filename)
			}
			fmt.Println(describeEvent(event.Mask))

			offset += unix.SizeofInotifyEvent + nameLen
		}
	}
}

func clen(b []byte) int {
	for i, c := range b {
		if c == 0 {
			return i
		}
	}
	return len(b)
}

func describeEvent(mask uint32) string {
	events := []string{}

	if mask&unix.IN_ACCESS != 0 {
		events = append(events, "ACCESS")
	}
	if mask&unix.IN_MODIFY != 0 {
		events = append(events, "MODIFY")
	}
	if mask&unix.IN_ATTRIB != 0 {
		events = append(events, "ATTRIB")
	}
	if mask&unix.IN_CLOSE_WRITE != 0 {
		events = append(events, "CLOSE_WRITE")
	}
	if mask&unix.IN_CLOSE_NOWRITE != 0 {
		events = append(events, "CLOSE_NOWRITE")
	}
	if mask&unix.IN_OPEN != 0 {
		events = append(events, "OPEN")
	}
	if mask&unix.IN_MOVED_FROM != 0 {
		events = append(events, "MOVED_FROM")
	}
	if mask&unix.IN_MOVED_TO != 0 {
		events = append(events, "MOVED_TO")
	}
	if mask&unix.IN_CREATE != 0 {
		events = append(events, "CREATE")
	}
	if mask&unix.IN_DELETE != 0 {
		events = append(events, "DELETE")
	}
	if mask&unix.IN_DELETE_SELF != 0 {
		events = append(events, "DELETE_SELF")
	}
	if mask&unix.IN_MOVE_SELF != 0 {
		events = append(events, "MOVE_SELF")
	}
	if mask&unix.IN_ISDIR != 0 {
		events = append(events, "ISDIR")
	}
	if mask&unix.IN_IGNORED != 0 {
		events = append(events, "IGNORED")
	}

	if len(events) == 0 {
		return fmt.Sprintf("UNKNOWN(0x%x)", mask)
	}

	result := events[0]
	for i := 1; i < len(events); i++ {
		result += " | " + events[i]
	}
	return result
}

// Add a new inotify watch for the given path. If that path is a directory,
// traverse it recursively to add a watch for all subdirectories, too.
//
// Ignores directories that start with '.' !
func (s *previewServer) addWatchRecursively(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("Failed to stat %s: %v", path, err)
	}

	if info.IsDir() {
		err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				return nil
			}

			name := d.Name()
			if len(name) > 0 && name[0] == '.' {
				return fs.SkipDir // Skip directories starting with .
			}
			s.addWatch(path)
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		s.addWatch(path)
	}

	return nil
}

// Add a new inotify watch for a single given path.
func (s *previewServer) addWatch(path string) error {
	wd, err := unix.InotifyAddWatch(s.inotifyInstanceFD, path, unix.IN_CLOSE_WRITE)
	if err != nil {
		return err
	}
	s.watches[wd] = path
	return nil
}

// Return a []string with the absolute version of all paths in `paths`.
func normalizePaths(paths []string) ([]string, error) {
	absPaths := make([]string, len(paths))
	for i, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return absPaths, fmt.Errorf("failed to get absolute path for %s: %v", path, err)
		}
		absPaths[i] = absPath
	}
	return absPaths, nil
}
