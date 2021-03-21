package main

import (
	"fmt"
	"io"

	"github.com/pippolo84/beetools/internal/torrent"
)

func show(w io.Writer, r io.Reader) error {
	torrent, err := torrent.NewTorrent(r)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, torrent)

	return nil
}
