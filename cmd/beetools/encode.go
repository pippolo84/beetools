package main

import (
	"encoding/json"
	"io"

	"github.com/pippolo84/beetools/internal/torrent"
	"github.com/pippolo84/beetools/pkg/bencode"
)

func encode(w io.Writer, r io.Reader) error {
	var t torrent.Torrent
	if err := json.NewDecoder(r).Decode(&t); err != nil {
		return err
	}

	d := t.ToDict()

	if err := bencode.NewEncoder(w).Encode(d); err != nil {
		return err
	}

	return nil
}
