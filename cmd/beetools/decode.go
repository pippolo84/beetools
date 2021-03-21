package main

import (
	"encoding/json"
	"io"

	"github.com/pippolo84/beetools/internal/torrent"
)

func decode(w io.Writer, r io.Reader) error {
	torrent, err := torrent.NewTorrent(r)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(torrent); err != nil {
		return err
	}

	return nil
}
