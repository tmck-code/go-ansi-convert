package convert

import (
	"bytes"
	"encoding/binary"
	"strings"
)

// SAUCE data types
const (
	DataTypeNone       byte = 0
	DataTypeCharacter  byte = 1
	DataTypeBitmap     byte = 2
	DataTypeVector     byte = 3
	DataTypeAudio      byte = 4
	DataTypeBinaryText byte = 5
	DataTypeXBin       byte = 6
	DataTypeArchive    byte = 7
	DataTypeExecutable byte = 8
)

// Character file types
const (
	FileTypeCharacterASCII      byte = 0
	FileTypeCharacterANSI       byte = 1
	FileTypeCharacterANSIMation byte = 2
	FileTypeCharacterRIPScript  byte = 3
	FileTypeCharacterPCBoard    byte = 4
	FileTypeCharacterAvatar     byte = 5
	FileTypeCharacterHTML       byte = 6
	FileTypeCharacterSource     byte = 7
	FileTypeCharacterTundraDraw byte = 8
)

// ANSiFlags bit masks
const (
	ANSiFlagNonBlinkMode  byte = 0x01 // B: iCE Color
	ANSiFlagLetterSpacing byte = 0x06 // LS: 8/9 pixel font (bits 1-2)
	ANSiFlagAspectRatio   byte = 0x18 // AR: Aspect ratio (bits 3-4)
)

// Letter spacing values
const (
	LetterSpacingLegacy byte = 0x00 // 00: Legacy, no preference
	LetterSpacing8Pixel byte = 0x02 // 01: 8 pixel font
	LetterSpacing9Pixel byte = 0x04 // 10: 9 pixel font
)

// Aspect ratio values
const (
	AspectRatioLegacy  byte = 0x00 // 00: Legacy, no preference
	AspectRatioLegacy1 byte = 0x08 // 01: Legacy device (elongated pixels)
	AspectRatioModern  byte = 0x10 // 10: Modern device (square pixels)
)

// TInfoFieldNames holds the semantic names for TInfo1-4 fields
type TInfoFieldNames struct {
	TInfo1Name string
	TInfo2Name string
	TInfo3Name string
	TInfo4Name string
}

// tInfoFieldMap maps DataType and FileType to field names
var tInfoFieldMap = map[byte]map[byte]TInfoFieldNames{
	DataTypeNone: {
		0: {"0", "0", "0", "0"},
	},
	DataTypeCharacter: {
		FileTypeCharacterASCII:      {"Character width", "Number of lines", "0", "0"},
		FileTypeCharacterANSI:       {"Character width", "Number of lines", "0", "0"},
		FileTypeCharacterANSIMation: {"Character width", "Character screen height", "0", "0"},
		FileTypeCharacterRIPScript:  {"Pixel width", "Pixel height", "Number of colors", "0"},
		FileTypeCharacterPCBoard:    {"Character width", "Number of lines", "0", "0"},
		FileTypeCharacterAvatar:     {"Character width", "Number of lines", "0", "0"},
		FileTypeCharacterHTML:       {"0", "0", "0", "0"},
		FileTypeCharacterSource:     {"0", "0", "0", "0"},
		FileTypeCharacterTundraDraw: {"Character width", "Number of lines", "0", "0"},
	},
	DataTypeBitmap: {
		0:  {"Pixel width", "Pixel height", "Pixel depth", "0"}, // GIF
		1:  {"Pixel width", "Pixel height", "Pixel depth", "0"}, // PCX
		2:  {"Pixel width", "Pixel height", "Pixel depth", "0"}, // LBM/IFF
		3:  {"Pixel width", "Pixel height", "Pixel depth", "0"}, // TGA
		4:  {"Pixel width", "Pixel height", "Pixel depth", "0"}, // FLI
		5:  {"Pixel width", "Pixel height", "Pixel depth", "0"}, // FLC
		6:  {"Pixel width", "Pixel height", "Pixel depth", "0"}, // BMP
		7:  {"Pixel width", "Pixel height", "Pixel depth", "0"}, // GL
		8:  {"Pixel width", "Pixel height", "Pixel depth", "0"}, // DL
		9:  {"Pixel width", "Pixel height", "Pixel depth", "0"}, // WPG
		10: {"Pixel width", "Pixel height", "Pixel depth", "0"}, // PNG
		11: {"Pixel width", "Pixel height", "Pixel depth", "0"}, // JPG
		12: {"Pixel width", "Pixel height", "Pixel depth", "0"}, // MPG
		13: {"Pixel width", "Pixel height", "Pixel depth", "0"}, // AVI
	},
	DataTypeVector: {
		0: {"0", "0", "0", "0"}, // DXF
		1: {"0", "0", "0", "0"}, // DWG
		2: {"0", "0", "0", "0"}, // WPG
		3: {"0", "0", "0", "0"}, // 3DS
	},
	DataTypeAudio: {
		0:  {"0", "0", "0", "0"}, // MOD
		1:  {"0", "0", "0", "0"}, // 669
		2:  {"0", "0", "0", "0"}, // STM
		3:  {"0", "0", "0", "0"}, // S3M
		4:  {"0", "0", "0", "0"}, // MTM
		5:  {"0", "0", "0", "0"}, // FAR
		6:  {"0", "0", "0", "0"}, // ULT
		7:  {"0", "0", "0", "0"}, // AMF
		8:  {"0", "0", "0", "0"}, // DMF
		9:  {"0", "0", "0", "0"}, // OKT
		10: {"0", "0", "0", "0"}, // ROL
		11: {"0", "0", "0", "0"}, // CMF
		12: {"0", "0", "0", "0"}, // MID
		13: {"0", "0", "0", "0"}, // SADT
		14: {"0", "0", "0", "0"}, // VOC
		15: {"0", "0", "0", "0"}, // WAV
		16: {"Sample rate", "0", "0", "0"}, // SMP8
		17: {"Sample rate", "0", "0", "0"}, // SMP8S
		18: {"Sample rate", "0", "0", "0"}, // SMP16
		19: {"Sample rate", "0", "0", "0"}, // SMP16S
		20: {"0", "0", "0", "0"}, // PATCH8
		21: {"0", "0", "0", "0"}, // PATCH16
		22: {"0", "0", "0", "0"}, // XM
		23: {"0", "0", "0", "0"}, // HSC
		24: {"0", "0", "0", "0"}, // IT
	},
	DataTypeBinaryText: {
		// FileType is variable for BinaryText (encodes width)
	},
	DataTypeXBin: {
		0: {"Character width", "Number of lines", "0", "0"},
	},
	DataTypeArchive: {
		0: {"0", "0", "0", "0"}, // ZIP
		1: {"0", "0", "0", "0"}, // ARJ
		2: {"0", "0", "0", "0"}, // LZH
		3: {"0", "0", "0", "0"}, // ARC
		4: {"0", "0", "0", "0"}, // TAR
		5: {"0", "0", "0", "0"}, // ZOO
		6: {"0", "0", "0", "0"}, // RAR
		7: {"0", "0", "0", "0"}, // UC2
		8: {"0", "0", "0", "0"}, // PAK
		9: {"0", "0", "0", "0"}, // SQZ
	},
	DataTypeExecutable: {
		0: {"0", "0", "0", "0"},
	},
}

type SAUCE struct {
	ID         string // 5 bytes: Should be "SAUCE"
	Version    string // 2 bytes: Should be "00"
	Title      string // 35 bytes: Title of the file
	Author     string // 20 bytes: Creator (nick)name or handle
	Group      string // 20 bytes: Group or company name
	Date       string // 8 bytes: CCYYMMDD format
	FileSize   uint32 // 4 bytes: Original file size (little-endian)
	DataType   byte   // 1 byte: Type of data
	FileType   byte   // 1 byte: Type of file
	TInfo1Name string // 35 bytes: Name of TInfo1 field
	TInfo2Name string // 35 bytes: Name of TInfo2 field
	TInfo3Name string // 35 bytes: Name of TInfo3 field
	TInfo4Name string // 35 bytes: Name of TInfo4 field
	TInfo1     uint16 // 2 bytes: Type dependent info (little-endian)
	TInfo2     uint16 // 2 bytes: Type dependent info (little-endian)
	TInfo3     uint16 // 2 bytes: Type dependent info (little-endian)
	TInfo4     uint16 // 2 bytes: Type dependent info (little-endian)
	Comments   byte   // 1 byte: Number of comment lines
	TFlags     byte   // 1 byte: Type dependent flags
	TInfoS     string // 22 bytes: Type dependent string (null-terminated)
}

// ParseSAUCE parses SAUCE metadata from the last 128 bytes of data
// Returns nil if no valid SAUCE record is found
func ParseSAUCE(data []byte) *SAUCE {
	// SAUCE record is always 128 bytes at the end of the file
	if len(data) < 128 {
		return nil
	}

	// Read the last 128 bytes
	sauceData := data[len(data)-128:]

	// Create a reader for binary data
	reader := bytes.NewReader(sauceData)

	sauce := &SAUCE{}

	// Read ID (5 bytes)
	idBytes := make([]byte, 5)
	reader.Read(idBytes)
	sauce.ID = string(idBytes)

	// Check if this is a valid SAUCE record
	if sauce.ID != "SAUCE" {
		return nil
	}

	// Read Version (2 bytes)
	versionBytes := make([]byte, 2)
	reader.Read(versionBytes)
	sauce.Version = string(versionBytes)

	// Read Title (35 bytes) - trim spaces
	titleBytes := make([]byte, 35)
	reader.Read(titleBytes)
	sauce.Title = strings.TrimRight(string(titleBytes), " \x00")

	// Read Author (20 bytes) - trim spaces
	authorBytes := make([]byte, 20)
	reader.Read(authorBytes)
	sauce.Author = strings.TrimRight(string(authorBytes), " \x00")

	// Read Group (20 bytes) - trim spaces
	groupBytes := make([]byte, 20)
	reader.Read(groupBytes)
	sauce.Group = strings.TrimRight(string(groupBytes), " \x00")

	// Read Date (8 bytes) - trim spaces
	dateBytes := make([]byte, 8)
	reader.Read(dateBytes)
	sauce.Date = strings.TrimRight(string(dateBytes), " \x00")

	// Read FileSize (4 bytes, little-endian unsigned)
	binary.Read(reader, binary.LittleEndian, &sauce.FileSize)

	// Read DataType (1 byte)
	binary.Read(reader, binary.LittleEndian, &sauce.DataType)

	// Read FileType (1 byte)
	binary.Read(reader, binary.LittleEndian, &sauce.FileType)

	// Read TInfo1 (2 bytes, little-endian)
	binary.Read(reader, binary.LittleEndian, &sauce.TInfo1)

	// Read TInfo2 (2 bytes, little-endian)
	binary.Read(reader, binary.LittleEndian, &sauce.TInfo2)

	// Read TInfo3 (2 bytes, little-endian)
	binary.Read(reader, binary.LittleEndian, &sauce.TInfo3)

	// Read TInfo4 (2 bytes, little-endian)
	binary.Read(reader, binary.LittleEndian, &sauce.TInfo4)

	// Read Comments (1 byte)
	binary.Read(reader, binary.LittleEndian, &sauce.Comments)

	// Read TFlags (1 byte)
	binary.Read(reader, binary.LittleEndian, &sauce.TFlags)

	// Read TInfoS (22 bytes) - null-terminated string
	tinfoSBytes := make([]byte, 22)
	reader.Read(tinfoSBytes)
	// Find the null terminator
	nullIndex := bytes.IndexByte(tinfoSBytes, 0)
	if nullIndex != -1 {
		sauce.TInfoS = string(tinfoSBytes[:nullIndex])
	} else {
		sauce.TInfoS = string(tinfoSBytes)
	}

	// Populate TInfo field names based on DataType and FileType
	if dataTypeMap, exists := tInfoFieldMap[sauce.DataType]; exists {
		// For BinaryText, all field names are "0" regardless of FileType value
		if sauce.DataType == DataTypeBinaryText {
			sauce.TInfo1Name = "0"
			sauce.TInfo2Name = "0"
			sauce.TInfo3Name = "0"
			sauce.TInfo4Name = "0"
		} else if fieldNames, exists := dataTypeMap[sauce.FileType]; exists {
			sauce.TInfo1Name = fieldNames.TInfo1Name
			sauce.TInfo2Name = fieldNames.TInfo2Name
			sauce.TInfo3Name = fieldNames.TInfo3Name
			sauce.TInfo4Name = fieldNames.TInfo4Name
		} else {
			// Unknown FileType - default to "0"
			sauce.TInfo1Name = "0"
			sauce.TInfo2Name = "0"
			sauce.TInfo3Name = "0"
			sauce.TInfo4Name = "0"
		}
	} else {
		// Unknown DataType - default to "0"
		sauce.TInfo1Name = "0"
		sauce.TInfo2Name = "0"
		sauce.TInfo3Name = "0"
		sauce.TInfo4Name = "0"
	}

	return sauce
}

// HasNonBlinkMode returns true if the iCE Color flag is set (ANSi files)
func (s *SAUCE) HasNonBlinkMode() bool {
	return s.TFlags&ANSiFlagNonBlinkMode != 0
}

// GetLetterSpacing returns the letter spacing setting (ANSi files)
func (s *SAUCE) GetLetterSpacing() byte {
	return s.TFlags & ANSiFlagLetterSpacing
}

// GetAspectRatio returns the aspect ratio setting (ANSi files)
func (s *SAUCE) GetAspectRatio() byte {
	return s.TFlags & ANSiFlagAspectRatio
}

// GetFontName returns the font name from TInfoS (ANSi files)
func (s *SAUCE) GetFontName() string {
	return s.TInfoS
}

// IsCharacterFile returns true if this is a character-based file
func (s *SAUCE) IsCharacterFile() bool {
	return s.DataType == DataTypeCharacter
}

// IsANSIFile returns true if this is an ANSI file
func (s *SAUCE) IsANSIFile() bool {
	return s.DataType == DataTypeCharacter && s.FileType == FileTypeCharacterANSI
}
