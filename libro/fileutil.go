package libro

import (
	"io"
	"os"
	"path/filepath"
)

// copyFile copies src to dst. Directories hosting dst are created as needed.
// if dst exists, copy does not happen and an error is returned.
// copyFile forces write to disk (Sync() method of os.File), and value
// certainty that write operation happens correctly over performance.
func copyFile(dst string, src string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0777); err != nil {
		return err
	}

	w, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err = io.Copy(w, r); err != nil {
		return err
	}

	return w.Sync()
}

// samePath checks if two path strings are representing the same path.
// Limitation: this version is only comparing the path strings and do not
// consider situations where path string are different but are pointing to the
// same location on file-system.
func samePath(path1, path2 string) (bool, error) {
	abspath1, err := filepath.Abs(path1)
	if err != nil {
		return false, err
	}

	abspath2, err := filepath.Abs(path2)
	if err != nil {
		return false, err
	}

	return abspath1 == abspath2, nil
}
