package trackdecode

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"
	"github.com/disgoorg/lavasrc-plugin"
)

const trackInfoVersioned int32 = 1

func decodeLavaSrcExtendedFields(track *lavalink.Track, r io.Reader) error {
	var info lavasrc.TrackInfo
	defer func() {
		raw, _ := json.Marshal(info)
		track.PluginInfo = raw
	}()

	albumName, err := readNullableString(r)
	if err != nil {
		return fmt.Errorf("failed to read track album name: %w", err)
	}
	if albumName != nil {
		info.AlbumName = *albumName
	}

	albumURL, err := readNullableString(r)
	if err != nil {
		return fmt.Errorf("failed to read track album url: %w", err)
	}
	if albumURL != nil {
		info.ArtistURL = *albumURL
	}

	artistURL, err := readNullableString(r)
	if err != nil {
		return fmt.Errorf("failed to read track artist url: %w", err)
	}
	if artistURL != nil {
		info.ArtistURL = *artistURL
	}

	artistArtworkURL, err := readNullableString(r)
	if err != nil {
		return fmt.Errorf("failed to read track artist artwork url: %w", err)
	}
	if artistArtworkURL != nil {
		info.ArtistArtworkURL = *artistArtworkURL
	}

	previewURL, err := readNullableString(r)
	if err != nil {
		return fmt.Errorf("failed to read track preview url: %w", err)
	}
	if previewURL != nil {
		info.PreviewURL = *previewURL
	}

	info.IsPreview, err = readBool(r)
	if err != nil {
		return fmt.Errorf("failed to read track is preview: %w", err)
	}

	return nil
}

type probeInfo struct {
	ProbeInfo string `json:"probeInfo"`
}

func decodeProbeInfo(track *lavalink.Track, r io.Reader) error {
	info, err := readString(r)
	if err != nil {
		return fmt.Errorf("failed to read track probe info: %w", err)
	}

	raw, _ := json.Marshal(probeInfo{
		ProbeInfo: info,
	})
	track.PluginInfo = raw
	return nil
}

var customTrackDecoders = map[string]func(track *lavalink.Track, r io.Reader) error{
	"deezer": func(track *lavalink.Track, r io.Reader) error {
		return decodeLavaSrcExtendedFields(track, r)
	},
	"spotify": func(track *lavalink.Track, r io.Reader) error {
		return decodeLavaSrcExtendedFields(track, r)
	},
	"applemusic": func(track *lavalink.Track, r io.Reader) error {
		return decodeLavaSrcExtendedFields(track, r)
	},
	"http": func(track *lavalink.Track, r io.Reader) error {
		return decodeProbeInfo(track, r)
	},
	"local": func(track *lavalink.Track, r io.Reader) error {
		return decodeProbeInfo(track, r)
	},
}

func DecodeString(encoded string) (*lavalink.Track, int, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid base64: %w", err)
	}

	var (
		r     = bytes.NewReader(data)
		track = &lavalink.Track{
			Encoded:    encoded,
			PluginInfo: []byte("{}"),
			UserData:   []byte("{}"),
		}
		value int32
	)

	if value, err = readInt32(r); err != nil {
		return nil, 0, fmt.Errorf("failed to read track flags: %w", err)
	}

	flags := int32(int64(value) & 0xC0000000 >> 30)
	messageSize := int64(value) & 0x3FFFFFFF
	if messageSize == 0 {
		return nil, 0, errors.New("message size: 0")
	}

	var version int
	if flags&trackInfoVersioned == 0 {
		version = 1
	} else {
		var v uint8
		if v, err = readUInt8(r); err != nil {
			return nil, 0, fmt.Errorf("failed to read track version: %w", err)
		}
		version = int(v)
	}

	if track.Info.Title, err = readString(r); err != nil {
		return track, version, fmt.Errorf("failed to read track title: %w", err)
	}
	if track.Info.Author, err = readString(r); err != nil {
		return track, version, fmt.Errorf("failed to read track author: %w", err)
	}

	var length int64
	if length, err = readInt64(r); err != nil {
		return track, version, fmt.Errorf("failed to read track length: %w", err)
	}
	track.Info.Length = lavalink.Duration(length)

	if track.Info.Identifier, err = readString(r); err != nil {
		return track, version, fmt.Errorf("failed to read track identifier: %w", err)
	}
	if track.Info.IsStream, err = readBool(r); err != nil {
		return track, version, fmt.Errorf("failed to read track is stream: %w", err)
	}
	if version >= 2 {
		if track.Info.URI, err = readNullableString(r); err != nil {
			return track, version, fmt.Errorf("failed to read track uri: %w", err)
		}
	}
	if version >= 3 {
		if track.Info.ArtworkURL, err = readNullableString(r); err != nil {
			return track, version, fmt.Errorf("failed to read track artwork url: %w", err)
		}
		if track.Info.ISRC, err = readNullableString(r); err != nil {
			return track, version, fmt.Errorf("failed to read track isrc: %w", err)
		}
	}
	if track.Info.SourceName, err = readString(r); err != nil {
		return track, version, fmt.Errorf("failed to read track source name: %w", err)
	}

	if decoder, ok := customTrackDecoders[track.Info.SourceName]; ok {
		if err = decoder(track, r); err != nil {
			return track, version, fmt.Errorf("failed to read track source fields: %w", err)
		}
	}

	var position int64
	if position, err = readInt64(r); err != nil {
		return track, version, fmt.Errorf("failed to read track position: %w", err)
	}
	track.Info.Position = lavalink.Duration(position)

	return track, version, nil
}

func readInt64(r io.Reader) (i int64, err error) {
	return i, binary.Read(r, binary.BigEndian, &i)
}

func readInt32(r io.Reader) (i int32, err error) {
	return i, binary.Read(r, binary.BigEndian, &i)
}

func readInt16(r io.Reader) (i int16, err error) {
	return i, binary.Read(r, binary.BigEndian, &i)
}

func readUInt8(r io.Reader) (i uint8, err error) {
	return i, binary.Read(r, binary.BigEndian, &i)
}

func readBool(r io.Reader) (b bool, err error) {
	return b, binary.Read(r, binary.BigEndian, &b)
}

func readString(r io.Reader) (string, error) {
	size, err := readInt16(r)
	if err != nil {
		return "", err
	}
	b := make([]byte, size)
	if err = binary.Read(r, binary.BigEndian, &b); err != nil {
		return "", err
	}
	return string(b), nil
}

func readNullableString(r io.Reader) (*string, error) {
	b, err := readBool(r)
	if err != nil || !b {
		return nil, err
	}

	s, err := readString(r)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
