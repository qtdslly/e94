package fms

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Storage struct {
	root   string
	module string
	dir    string
}

var smap map[string]*Storage

func NewStorage(root, module string, infix ...string) *Storage {
	dir := filepath.Join(root, module)
	for _, v := range infix {
		dir = filepath.Join(dir, v)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil
	}
	s := new(Storage)
	s.root = filepath.Clean(root)
	s.module = module
	s.dir = dir
	if smap == nil {
		smap = make(map[string]*Storage)
	}
	smap[module] = s
	return s
}

func FindStorage(module string) *Storage {
	return smap[module]
}

func (s *Storage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ss := strings.Split(filepath.Clean(r.URL.Path), string(os.PathSeparator))

	l := len(ss)
	if l < 2 {
		http.NotFound(w, r)
		return
	}
	found := false
	index := l
	for ; index > 0; index-- {
		if ss[index-1] == s.module {
			found = true
			break
		}
	}
	if !found {
		fmt.Println(ss, s.module)
		http.NotFound(w, r)
		return
	}

	path := s.root
	for ; index <= l; index++ {
		path = path + "/" + ss[index-1]
	}

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if isThumb(path) && createThumb(path) == nil {
			} else {
				http.NotFound(w, r)
				return
			}
		} else {
			http.NotFound(w, r)
			return
		}
	}

	if filepath.Ext(path) == "apk" {
		w.Header().Set("Content-Disposition", "attachment;filename="+path)
		w.Header().Set("ContentType", "application/vnd.android")
	}

	http.ServeFile(w, r, path)
}

// httpRange specifies the byte range to be sent to the client.
type httpRange struct {
	start, length int64
}

func (r httpRange) contentRange(size int64) string {
	return fmt.Sprintf("bytes %d-%d/%d", r.start, r.start+r.length-1, size)
}

func (r httpRange) mimeHeader(contentType string, size int64) textproto.MIMEHeader {
	return textproto.MIMEHeader{
		"Content-Range": {r.contentRange(size)},
		"Content-Type":  {contentType},
	}
}

// parseRange parses a Range header string as per RFC 2616.
func parseRange(s string, size int64) ([]httpRange, error) {
	if s == "" {
		return nil, nil // header not present
	}
	const b = "bytes="
	if !strings.HasPrefix(s, b) {
		return nil, errors.New("invalid range")
	}
	var ranges []httpRange
	for _, ra := range strings.Split(s[len(b):], ",") {
		ra = strings.TrimSpace(ra)
		if ra == "" {
			continue
		}
		i := strings.Index(ra, "-")
		if i < 0 {
			return nil, errors.New("invalid range")
		}
		start, end := strings.TrimSpace(ra[:i]), strings.TrimSpace(ra[i+1:])
		var r httpRange
		if start == "" {
			// If no start is specified, end specifies the
			// range start relative to the end of the file.
			i, err := strconv.ParseInt(end, 10, 64)
			if err != nil {
				return nil, errors.New("invalid range")
			}
			if i > size {
				i = size
			}
			r.start = size - i
			r.length = size - r.start
		} else {
			i, err := strconv.ParseInt(start, 10, 64)
			if err != nil || i >= size || i < 0 {
				return nil, errors.New("invalid range")
			}
			r.start = i
			if end == "" {
				// If no end is specified, range extends to end of the file.
				r.length = size - r.start
			} else {
				i, err := strconv.ParseInt(end, 10, 64)
				if err != nil || r.start > i {
					return nil, errors.New("invalid range")
				}
				if i >= size {
					i = size - 1
				}
				r.length = i - r.start + 1
			}
		}
		ranges = append(ranges, r)
	}
	return ranges, nil
}

func (s *Storage) ServeHTTPWithPartial(w http.ResponseWriter, r *http.Request, psize int64) {
	ss := strings.Split(filepath.Clean(r.URL.Path), string(os.PathSeparator))

	l := len(ss)
	if l < 2 {
		http.NotFound(w, r)
		return
	}
	found := false
	index := l
	for ; index > 0; index-- {
		if ss[index-1] == s.module {
			found = true
			break
		}
	}
	if !found {
		fmt.Println(ss, s.module)
		http.NotFound(w, r)
		return
	}

	path := s.root
	for ; index <= l; index++ {
		path = path + "/" + ss[index-1]
	}

	fInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if isThumb(path) && createThumb(path) == nil {
			} else {
				http.NotFound(w, r)
				return
			}
		} else {
			http.NotFound(w, r)
			return
		}
	}

	if fInfo.Size() > psize*2 {
		rangeHeader := r.Header.Get("Range")
		if rangeHeader != "" {
			r.Header.Del("Range")
			rs, _ := parseRange(rangeHeader, fInfo.Size())
			if len(rs) == 1 {
				if rs[0].length > psize {
					rs[0].length = psize
				}
				r.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", rs[0].start, rs[0].start+rs[0].length-1))
			}
		} else {
			r.Header.Add("Range", fmt.Sprintf("bytes=0-%d", psize-1))
		}

	}

	http.ServeFile(w, r, path)
}

// TODO
func (s *Storage) ServeFile(w http.ResponseWriter, r *http.Request, filename string) {
}

func (s *Storage) BuildStoragePath(dst, ext string, createdAt *time.Time) string {
	if filepath.Ext(dst) == "" {
		dst = dst + strings.ToLower(ext)
	}
	folder := s.dir
	if createdAt != nil && !createdAt.IsZero() {
		folder = filepath.Join(folder, createdAt.Format("20060102"))
	}
	folder = filepath.Join(folder, filepath.Dir(dst))
	_, err := os.Stat(folder)
	if err != nil {
		os.MkdirAll(folder, 0755)
	}
	return filepath.Join(folder, filepath.Base(dst))
}

func (s *Storage) SaveFile(dst, path string, createdAt *time.Time) (string, error) {
	dstPath := s.BuildStoragePath(dst, filepath.Ext(path), createdAt)
	if dstPath == "" {
		return "", nil
	}

	src, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer src.Close()

	os.MkdirAll(filepath.Dir(dstPath), 0755)
	dstFile, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, src)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(dstPath, s.root), nil
}

func (s *Storage) CreateThumb(path string, width, height int, createdAt time.Time) string {
	return buildThumbPath(path, width, height)
}
