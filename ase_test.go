package ase

import (
	"bytes"
	"testing"
)

var testColors = []Color{
	Color{
		Name:   "RGB",
		Model:  "RGB",
		Values: []float32{1, 1, 1},
		Type:   "Normal",
	},
	Color{
		Name:   "Grayscale",
		Model:  "CMYK",
		Values: []float32{0, 0, 0, 0.47},
		Type:   "Spot",
	},
	Color{
		Name:   "cmyk",
		Model:  "CMYK",
		Values: []float32{0, 1, 0, 0},
		Type:   "Spot",
	},
	Color{
		Name:   "LAB",
		Model:  "RGB",
		Values: []float32{0, 0.6063648, 0.524658},
		Type:   "Global",
	},
	Color{
		Name:   "PANTONE P 1-8 C",
		Model:  "LAB",
		Values: []float32{0.9137255, -5, 94},
		Type:   "Spot",
	},
}

var testGroup = Group{
	Name: "A Color Group",
	Colors: []Color{
		Color{
			Name:   "Red",
			Model:  "RGB",
			Values: []float32{1, 0, 0},
			Type:   "Global",
		},
		Color{
			Name:   "Green",
			Model:  "RGB",
			Values: []float32{0, 1, 0},
			Type:   "Global",
		},
		Color{
			Name:   "Blue",
			Model:  "RGB",
			Values: []float32{0, 0, 1},
			Type:   "Global",
		},
	},
}

func TestSignature(t *testing.T) {
	testFile := "samples/test.ase"

	ase, err := DecodeFile(testFile)
	if err != nil {
		t.Error(err)
	}

	expectedSignature := "ASEF"
	if ase.Signature() != expectedSignature {
		t.Error("expected signature of", expectedSignature, ", got:", ase.Signature())
	}
}

func TestVersion(t *testing.T) {
	testFile := "samples/test.ase"

	ase, err := DecodeFile(testFile)
	if err != nil {
		t.Error(err)
	}

	expectedVersion := "1.0"
	if ase.Version() != expectedVersion {
		t.Error("expected version of", expectedVersion, ", got:", ase.Version())
	}
}

func TestDecode(t *testing.T) {
	testFile := "samples/test.ase"

	ase, err := DecodeFile(testFile)
	if err != nil {
		t.Error(err)
	}

	// Check ASE's metadata (Signature and Version tested in separate functions)
	expectedNumBlocks := int32(10)
	if ase.numBlocks != expectedNumBlocks {
		t.Error("expected ", expectedNumBlocks, " numBlocks, got ", ase.numBlocks)
	}

	// Check the ASE's Colors
	expectedColors := testColors
	expectedNumColors := len(expectedColors)
	actualNumColors := len(ase.Colors)

	if actualNumColors != expectedNumColors {
		t.Error("expected number of colors to be", expectedNumColors,
			"got ", actualNumColors)
	}

	for i, color := range ase.Colors {
		expectedColor := expectedColors[i]

		if color.Name != expectedColor.Name {
			t.Error("expected initial color with name ", expectedColor.Name,
				"got ", color.Name)
		}

		if color.Model != expectedColor.Model {
			t.Error("expected initial color of Model ", expectedColor.Model,
				"got ", color.Model)
		}

		for j, _ := range expectedColor.Values {
			if color.Values[j] != expectedColor.Values[j] {
				t.Error("expected color value ", expectedColor.Values[j],
					"got ", color.Values[j])
			}
		}

		if color.Type != expectedColor.Type {
			t.Error("expected color type ", expectedColor.Type,
				"got ", color.Type)
		}
	}

	// Check the ASE's Groups' data
	expectedGroupsLen := 1
	actualGroupLen := len(ase.Groups)

	if actualGroupLen != expectedGroupsLen {
		t.Error("expected group length of ", expectedGroupsLen,
			"got ", actualGroupLen)
	}

	group := ase.Groups[0]

	if group.Name != testGroup.Name {
		t.Error("expected group name to be ", testGroup.Name,
			", got: ", group.Name)
	}

	// Check the ASE's Groups' Color Data
	for i, color := range group.Colors {
		expectedColor := testGroup.Colors[i]

		if color.Name != expectedColor.Name {
			t.Error("expected initial color with name ", expectedColor.Name,
				"got ", color.Name)
		}

		if color.Model != expectedColor.Model {
			t.Error("expected initial color of Model ", expectedColor.Model,
				"got ", color.Model)
		}

		for j, _ := range expectedColor.Values {
			if color.Values[j] != expectedColor.Values[j] {
				t.Error("expected color value ", expectedColor.Values[j],
					"got ", color.Values[j])
			}
		}

		if color.Type != expectedColor.Type {
			t.Error("expected color type ", expectedColor.Type,
				"got ", color.Type)
		}
	}

}

func TestEncode(t *testing.T) {

	// Initialize a sample ASE
	sampleAse := ASE{}
	sampleAse.Colors = testColors
	sampleAse.Groups = append(sampleAse.Groups, testGroup)

	// Encode the sampleAse into the buffer and immediately decode it.
	b := new(bytes.Buffer)
	Encode(sampleAse, b)
	ase, _ := Decode(b)

	// Check the ASE's decoded values.
	if string(ase.signature[0:]) != "ASEF" {
		t.Error("ase: file not an ASE file")
	}

	if ase.version[0] != 1 && ase.version[1] != 0 {
		t.Error("ase: version is not 1.0")
	}

	expectedNumBlocks := int32(10)
	actualNumBlocks := ase.numBlocks
	if actualNumBlocks != expectedNumBlocks {
		t.Error("ase: expected", expectedNumBlocks,
			" blocks to be present, got: ", actualNumBlocks)
	}

	expectedAmountOfColors := 5
	if len(ase.Colors) != expectedAmountOfColors {
		t.Error("ase: expected", expectedAmountOfColors, " colors to be present")
	}

	for i, color := range ase.Colors {
		expectedColor := testColors[i]

		if color.Name != expectedColor.Name {
			t.Error("expected initial color with name ", expectedColor.Name,
				"got ", color.Name)
		}

		if color.Model != expectedColor.Model {
			t.Error("expected initial color of Model ", expectedColor.Model,
				"got ", color.Model)
		}

		for j, _ := range expectedColor.Values {
			if color.Values[j] != expectedColor.Values[j] {
				t.Error("expected color value ", expectedColor.Values[j],
					"got ", color.Values[j])
			}
		}

		if color.Type != expectedColor.Type {
			t.Error("expected color type ", expectedColor.Type,
				"got ", color.Type)
		}
	}

	expectedAmountOfGroups := 1
	actualAmountOfGroups := len(ase.Groups)
	if actualAmountOfGroups != expectedAmountOfGroups {
		t.Error("expected ", expectedAmountOfGroups,
			"amount of groups, got: ", actualAmountOfGroups)
	}

	group := ase.Groups[0]

	if group.Name != testGroup.Name {
		t.Error("expected group name to be ", testGroup.Name,
			", got: ", group.Name)
	}

	for i, color := range group.Colors {
		expectedColor := testGroup.Colors[i]

		if color.Name != expectedColor.Name {
			t.Error("expected initial color with name ", expectedColor.Name,
				"got ", color.Name)
		}

		if color.Model != expectedColor.Model {
			t.Error("expected initial color of Model ", expectedColor.Model,
				"got ", color.Model)
		}

		for j, _ := range expectedColor.Values {
			if color.Values[j] != expectedColor.Values[j] {
				t.Error("expected color value ", expectedColor.Values[j],
					"got ", color.Values[j])
			}
		}

		if color.Type != expectedColor.Type {
			t.Error("expected color type ", expectedColor.Type,
				"got ", color.Type)
		}
	}

}
