package utils

import (
	"archive/tar"
	"bytes"
)

func TarFile(filename, content string) ([]byte, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	hdr := &tar.Header{
		Name: filename,
		Mode: 0644,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return nil, err
	}

	if _, err := tw.Write([]byte(content)); err != nil {
		return nil, err
	}

	tw.Close()
	return buf.Bytes(), nil
}
