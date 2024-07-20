package lavalinkbot

import (
	"bytes"
	"cmp"
	"embed"
	"fmt"
	"io"
	"path"

	"github.com/adrg/frontmatter"
)

func ReadThings(things embed.FS) (map[string]Thing, error) {
	files, err := things.ReadDir("things")
	if err != nil {
		return nil, err
	}

	thingMap := make(map[string]Thing, len(files))
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		thingFiles, err := things.ReadDir(path.Join("things", file.Name()))
		if err != nil {
			return nil, err
		}

		var thing Thing
		for _, f := range thingFiles {
			if f.Name() == "index.md" {
				data, err := things.ReadFile(path.Join("things", file.Name(), "index.md"))
				if err != nil {
					return nil, fmt.Errorf("failed to read index.md for %s: %w", file.Name(), err)
				}

				var matter thingMatter
				data, err = frontmatter.Parse(bytes.NewBuffer(data), &matter)
				if err != nil {
					return nil, fmt.Errorf("failed to parse frontmatter for %s: %w", file.Name(), err)
				}

				thing.Name = cmp.Or(matter.Name, file.Name())
				thing.FileName = file.Name()
				thing.Content = string(data)
				continue
			}

			data, err := things.ReadFile(path.Join("things", file.Name(), f.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s for %s: %w", f.Name(), file.Name(), err)
			}

			thing.Files = append(thing.Files, ThingFile{
				Name: f.Name(),
				Buf:  data,
			})
		}

		thingMap[file.Name()] = thing
	}

	return thingMap, nil
}

type thingMatter struct {
	Name string `yaml:"name"`
}

type Thing struct {
	Name     string
	FileName string
	Content  string
	Files    []ThingFile
}

type ThingFile struct {
	Name string
	Buf  []byte
}

func (t ThingFile) Reader() io.Reader {
	return bytes.NewReader(t.Buf)
}
