package torrent

import (
	"encoding/json"
	"io"
	"time"

	"github.com/pippolo84/beetools/pkg/bencode"
)

// Info holds the "info" part of a .torrent file.
type Info struct {
	Length      int64  `json:"length"`
	Name        string `json:"name"`
	PieceLength int64  `json:"piece length"`
	Pieces      []byte `json:"pieces"`
}

// Torrent represents all the information in a .torrent file.
type Torrent struct {
	Announce     string    `json:"announce"`
	Comment      string    `json:"comment"`
	CreationDate time.Time `json:"creation date"`
	HTTPSeeds    []string  `json:"httpseeds"`
	Info         Info      `json:"info"`
}

// NewTorrent returns a new Torrent initialized with bencode-data from
// read the r io.Reader.
func NewTorrent(r io.Reader) (*Torrent, error) {
	d := bencode.Dict{}
	dec := bencode.NewDecoder(r)
	if err := dec.Decode(&d); err != nil {
		return nil, err
	}

	mapValue := d.Value()

	infoValue := mapValue["info"].(map[string]interface{})
	info := Info{
		Length:      infoValue["length"].(int64),
		Name:        infoValue["name"].(string),
		PieceLength: infoValue["piece length"].(int64),
		Pieces:      []byte(infoValue["pieces"].(string)),
	}

	seedsValue := mapValue["httpseeds"].([]interface{})
	seeds := make([]string, 0, len(seedsValue))
	for _, v := range seedsValue {
		seeds = append(seeds, v.(string))
	}
	return &Torrent{
		Announce:     mapValue["announce"].(string),
		Comment:      mapValue["comment"].(string),
		CreationDate: time.Unix(mapValue["creation date"].(int64), 0),
		HTTPSeeds:    seeds,
		Info:         info,
	}, nil
}

// ToDict returns a bencode package Dict representation of the torrent.
func (t *Torrent) ToDict() bencode.Dict {
	seedsValues := make([]interface{}, 0, len(t.HTTPSeeds))
	for _, v := range t.HTTPSeeds {
		seedsValues = append(seedsValues, bencode.NewByteString(v))
	}

	return bencode.NewDict(map[bencode.ByteString]interface{}{
		bencode.NewByteString("announce"):      bencode.NewByteString(t.Announce),
		bencode.NewByteString("comment"):       bencode.NewByteString(t.Comment),
		bencode.NewByteString("creation date"): bencode.NewInteger(t.CreationDate.Unix()),
		bencode.NewByteString("httpseeds"):     bencode.NewList(seedsValues),
		bencode.NewByteString("info"): bencode.NewDict(map[bencode.ByteString]interface{}{
			bencode.NewByteString("length"):       bencode.NewInteger(t.Info.Length),
			bencode.NewByteString("name"):         bencode.NewByteString(t.Info.Name),
			bencode.NewByteString("piece length"): bencode.NewInteger(t.Info.PieceLength),
			bencode.NewByteString("pieces"):       bencode.NewByteString(string(t.Info.Pieces)),
		}),
	})
}

// String satisfies the fmt.Stringer interface.
func (t Torrent) String() string {
	// filter pieces data away for stringification
	t.Info.Pieces = []byte("")

	buf, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return ""
	}
	return string(buf)
}
