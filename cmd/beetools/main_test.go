package main

import (
	"bytes"
	"crypto/md5"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	in, err := os.Open(
		filepath.Join(
			"testdata",
			"debian-10.8.0-amd64-netinst.iso.torrent",
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		in.Close()
	})

	testDir := t.TempDir()

	testJSON := filepath.Join(testDir, "test.json")
	out, err := os.Create(testJSON)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		out.Close()
	})

	if err := decode(out, in); err != nil {
		t.Fatal(err)
	}
	if err := out.Sync(); err != nil {
		t.Fatal(err)
	}

	inJSON, err := os.Open(testJSON)
	if err != nil {
		t.Fatal(err)
	}

	testTorrent := filepath.Join(testDir, "test.torrent")
	outTorrent, err := os.Create(testTorrent)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		outTorrent.Close()
	})

	if err := encode(outTorrent, inJSON); err != nil {
		t.Fatal(err)
	}
	if err := out.Sync(); err != nil {
		t.Fatal(err)
	}

	golden, err := os.Open(filepath.Join(
		"testdata",
		"debian-10.8.0-amd64-netinst.iso.torrent",
	))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		golden.Close()
	})

	goldenHash := md5.New()
	if _, err := io.Copy(goldenHash, golden); err != nil {
		t.Fatal(err)
	}

	generated, err := os.Open(testTorrent)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		generated.Close()
	})

	generatedHash := md5.New()
	if _, err := io.Copy(generatedHash, generated); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(goldenHash.Sum(nil), generatedHash.Sum(nil)) {
		t.Fatal("md5sum of generated torrent differs")
	}
}
