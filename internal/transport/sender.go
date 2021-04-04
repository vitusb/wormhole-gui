package transport

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"github.com/psanford/wormhole-william/wormhole"
)

// NewFileSend takes the chosen file and sends it using wormhole-william.
func (c *Client) NewFileSend(ctx context.Context, file fyne.URIReadCloser, progress wormhole.SendOption) (string, chan wormhole.SendResult, error) {
	if fyne.CurrentDevice().IsMobile() {
		reader, err := newSizeSeekCloser(file)
		if err != nil {
			return "", nil, err
		}

		return c.SendFile(ctx, file.URI().Name(), reader, progress)
	}

	return c.SendFile(ctx, file.URI().Name(), file.(io.ReadSeeker), progress)
}

// NewDirSend takes a listable URI and sends it using wormhole-william.
func (c *Client) NewDirSend(ctx context.Context, dir fyne.ListableURI, progress wormhole.SendOption) (string, chan wormhole.SendResult, error) {
	prefixStr, _ := filepath.Split(dir.Path())
	prefix := len(prefixStr) // Where the prefix ends. Doing it this way is faster and works when paths don't use same separator (\ or /).

	var files []wormhole.DirectoryEntry
	if err := filepath.Walk(dir.Path(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() || !info.Mode().IsRegular() {
			return nil
		}

		files = append(files, wormhole.DirectoryEntry{
			Path: path[prefix:], // Instead of strings.TrimPrefix. Paths don't need match (e.g. "C:/home/dir" == "C:\home\dir").
			Mode: info.Mode(),
			Reader: func() (io.ReadCloser, error) {
				return os.Open(path) // #nosec - path is already cleaned by filepath.Walk
			},
		})

		return nil
	}); err != nil {
		fyne.LogError("Error on walking directory", err)
		return "", nil, err
	}

	return c.SendDirectory(ctx, dir.Name(), files, progress)
}

// NewTextSend takes a text input and sends the text using wormhole-william.
func (c *Client) NewTextSend(ctx context.Context, text string, progress wormhole.SendOption) (string, chan wormhole.SendResult, error) {
	return c.SendText(ctx, text, progress)
}

func newSizeSeekCloser(reader io.ReadCloser) (*sizeSeekCloser, error) {
	buffer := &bytes.Buffer{}
	size, err := io.Copy(buffer, reader)
	if err != nil {
		fyne.LogError("Could not buffer file contents", err)
		return nil, err
	}

	return &sizeSeekCloser{buffer, size}, nil
}

type sizeSeekCloser struct {
	reader *bytes.Buffer
	size   int64
}

// Seek fakes the size check done by wormhole-william. It does not use anything else.
func (s *sizeSeekCloser) Seek(offset int64, whence int) (int64, error) {
	return s.size, nil
}

func (s *sizeSeekCloser) Read(p []byte) (n int, err error) {
	return s.reader.Read(p)
}

func (s *sizeSeekCloser) Close() error {
	s.reader.Reset()
	return nil
}
