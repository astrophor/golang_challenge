package drum

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
// TODO: implement
func DecodeFile(path string) (*Pattern, error) {
	p := &Pattern{}

	fd, err := os.Open(path)
	if err != nil {
		return p, err
	}
	defer fd.Close()

	if err := p.Parse(fd); err != nil {
		fmt.Println("parse failed: ", err)
	} else {
		fmt.Println(fmt.Sprint(p))
	}

	return p, nil
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
// TODO: implement
type Pattern struct {
	Version [6]byte
	Length  uint64
	Name    [32]byte
	Tempo   float32
	Data    []Track
}

type Track struct {
	Id     uint8
	Length uint32
	Name   string
	Step   [16]uint8
}

func (p *Pattern) Parse(r io.Reader) error {
	if _, err := r.Read(p.Version[:]); err != nil {
		return err
	}

	binary.Read(r, binary.BigEndian, &p.Length)

	if _, err := r.Read(p.Name[:]); err != nil {
		return err
	}

	binary.Read(r, binary.LittleEndian, &p.Tempo)

	buffer := make([]byte, p.Length-36)
	if _, err := r.Read(buffer); err != nil {
		return err
	}
	reader := bytes.NewReader(buffer[:])

	var name_buffer [1024]byte

	for reader.Len() > 0 {
		var track Track

		binary.Read(reader, binary.BigEndian, &track.Id)
		binary.Read(reader, binary.BigEndian, &track.Length)
		if _, err := reader.Read(name_buffer[:track.Length]); err != nil {
			return err
		}
		track.Name = string(name_buffer[:track.Length])
		for idx := 0; idx < 16; idx++ {
			binary.Read(reader, binary.BigEndian, &track.Step[idx])
		}
		p.Data = append(p.Data, track)
	}

	return nil
}

func (p *Pattern) String() string {
	var res string

	idx := 0
	for ; idx < 32; idx++ {
		if p.Name[idx] == '\x00' {
			break
		}
	}
	res = fmt.Sprintf("Saved with HW Version: %v\nTempo: %v\n", string(p.Name[:idx]), p.Tempo)

	for _, track := range p.Data {
		res += fmt.Sprintf("(%v) %v\t|%v|%v|%v|%v|\n", track.Id, track.Name,
			trans(track.Step[:4]),
			trans(track.Step[4:8]),
			trans(track.Step[8:12]),
			trans(track.Step[12:16]))
	}

	return res
}

func trans(a []uint8) string {
	res := ""
	for idx, _ := range a {
		if a[idx] == 0 {
			res += "-"
		} else {
			res += "x"
		}
	}

	return res
}
