package ase

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strconv"
)

var (
	ErrInvalidFile      = errors.New("ase: file not an ASE file")
	ErrInvalidVersion   = errors.New("ase: version is not 1.0")
	ErrInvalidBlockType = errors.New("ase: invalid block type")
)

type ASE struct {
	signature [4]uint8
	version   [2]int16
	numBlocks int32
	Colors    []Color
	Groups    []Group
}

//	ASE File Spec http://www.selapa.net/swatches/colors/fileformats.php#adobe_ase

//	Decodes a valid ASE input.
func Decode(r io.Reader) (ase ASE, err error) {
	if err = ase.readSignature(r); err != nil {
		return
	}
	if err = ase.readVersion(r); err != nil {
		return
	}
	if err = ase.readNumBlocks(r); err != nil {
		return
	}

	//	if we encounter groups, store a ref here
	var g Group

	//	itereate based on our block count
	for i := 0; i < int(ase.numBlocks); i++ {
		//	new block
		b := block{}

		//	decode the block container
		if err = b.Read(r); err != nil {
			return
		}

		//	switch on block type
		switch b.Type {
		case colorEntry:
			c := Color{}
			if err = c.read(r); err != nil {
				return
			}

			//	if we have a group, add color to the group
			if g.Name != "" {
				g.Colors = append(g.Colors, c)
			} else {
				//	color is not in a group. add to color slice
				ase.Colors = append(ase.Colors, c)
			}

			break
		case groupStart:
			//	new group
			g = Group{}

			//	read the group
			if err = g.read(r); err != nil {
				return
			}

			break
		case groupEnd:
			//	add the group to our ase struct
			ase.Groups = append(ase.Groups, g)

			//	reset our group struct
			g = Group{}

			break
		default:
			err = ErrInvalidBlockType
			return
		}
	}

	return
}

//	Helper function that decodes a file into an ASE.
func DecodeFile(file string) (ase ASE, err error) {
	//	open the file
	f, err := os.Open(file)
	if err != nil {
		return
	}

	//	decode the file
	return Decode(f)
}

// Encodes an ASE into any `w` that satisfies the io.Writer interface.
func Encode(ase ASE, w io.Writer) (err error) {
	if err = ase.writeSignature(w); err != nil {
		return err
	}
	if err = ase.writeVersion(w); err != nil {
		return err
	}
	if err = ase.writeNumBlocks(w); err != nil {
		return err
	}
	if err = ase.writeColors(w); err != nil {
		return err
	}
	if err = ase.writeGroups(w); err != nil {
		return err
	}

	return nil
}

//	Returns the file signature in a human readable format.
func (ase *ASE) Signature() string {
	return string(ase.signature[0:])
}

// Returns the file version in a human readable format.
func (ase *ASE) Version() string {
	return strconv.Itoa(int(ase.version[0])) + "." + strconv.Itoa(int(ase.version[1]))
}

// Decodes the ASE's signature
func (ase *ASE) readSignature(r io.Reader) (err error) {
	//	Read the signature
	if err = binary.Read(r, binary.BigEndian, &ase.signature); err != nil {
		return
	}

	//	Checks signature is `ASEF`
	if string(ase.signature[0:]) != "ASEF" {
		return ErrInvalidFile
	}

	return
}

//	Decodes the ASE's version
func (ase *ASE) readVersion(r io.Reader) (err error) {
	// Read the version
	if err = binary.Read(r, binary.BigEndian, &ase.version); err != nil {
		return
	}

	// Checks version is 1.0
	if ase.version[0] != 1 && ase.version[1] != 0 {
		return ErrInvalidVersion
	}

	return
}

//	Decodes the ASE's number of blocks
func (ase *ASE) readNumBlocks(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &ase.numBlocks)
}

// Encodes the ASE signature
func (ase *ASE) writeSignature(w io.Writer) (err error) {
	signature := []byte("ASEF")
	return binary.Write(w, binary.BigEndian, signature)
}

// Encodes the ASE version
func (ase *ASE) writeVersion(w io.Writer) (err error) {
	version := [2]int16{1, 0}
	return binary.Write(w, binary.BigEndian, version)
}

// Encodes the ASE's number of blocks
func (ase *ASE) writeNumBlocks(w io.Writer) (err error) {
	numBlocks := ase.calculateNumBlocks()
	return binary.Write(w, binary.BigEndian, numBlocks)
}

// Determines the numBlocks of an ASE on the fly rather than returning its `ase.numBlocks` attribute.
// There is currently no mechanism in place to update numBlocks if a user adds or removes either colors, groups, or colors within groups.
func (ase *ASE) calculateNumBlocks() (numBlocks int32) {
	// A color has only one block.
	colorBlocks := len(ase.Colors)

	// A single group has a start block and an end block.
	groupBlocks := len(ase.Groups) * 2

	// Run a comprehension counting every color inside groups.
	if groupBlocks != 0 {
		for _, group := range ase.Groups {
			colorBlocks += len(group.Colors)
		}
	}

	numBlocks = int32(colorBlocks + groupBlocks)

	return
}

// Encode the data for ase.Colors according to the ASE spec.
func (ase *ASE) writeColors(w io.Writer) (err error) {
	for _, color := range ase.Colors {
		if err = color.write(w); err != nil {
			return err
		}
	}
	return nil
}

// Encode the data for ase.Groups according to the ASE spec.
func (ase *ASE) writeGroups(w io.Writer) (err error) {
	for _, group := range ase.Groups {
		if err = group.write(w); err != nil {
			return err
		}
	}
	return nil
}
