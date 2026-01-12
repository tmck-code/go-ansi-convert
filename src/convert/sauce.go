package convert

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
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

// TInfo field names
const (
	TInfoNameNone                  = "0"
	TInfoNameCharacterWidth        = "Character width"
	TInfoNameNumberOfLines         = "Number of lines"
	TInfoNameCharacterScreenHeight = "Character screen height"
	TInfoNamePixelWidth            = "Pixel width"
	TInfoNamePixelHeight           = "Pixel height"
	TInfoNameNumberOfColors        = "Number of colors"
	TInfoNamePixelDepth            = "Pixel depth"
	TInfoNameSampleRate            = "Sample rate"
)

// TInfoFieldNames holds the semantic names for TInfo1-4 fields
type TInfoFieldNames struct {
	TInfo1Name string
	TInfo2Name string
	TInfo3Name string
	TInfo4Name string
}

// TInfoField represents a SAUCE TInfo field with its semantic name and value
type TInfoField struct {
	Name  string
	Value uint16
}

var dataTypeNames = map[byte]string{
	DataTypeNone:       "None",
	DataTypeCharacter:  "Character",
	DataTypeBitmap:     "Bitmap",
	DataTypeVector:     "Vector",
	DataTypeAudio:      "Audio",
	DataTypeBinaryText: "BinaryText",
	DataTypeXBin:       "XBin",
	DataTypeArchive:    "Archive",
	DataTypeExecutable: "Executable",
}

var fileTypeNames = map[byte]string{
	FileTypeCharacterASCII:      "ASCII",
	FileTypeCharacterANSI:       "ANSi",
	FileTypeCharacterANSIMation: "ANSiMation",
	FileTypeCharacterRIPScript:  "RIP Script",
	FileTypeCharacterPCBoard:    "PCBoard",
	FileTypeCharacterAvatar:     "Avatar",
	FileTypeCharacterHTML:       "HTML",
	FileTypeCharacterSource:     "Source",
	FileTypeCharacterTundraDraw: "TundraDraw",
}

// Common TInfo field name patterns
var (
	tInfoNone                 = TInfoFieldNames{TInfoNameNone, TInfoNameNone, TInfoNameNone, TInfoNameNone}
	tInfoCharacterWidthLines  = TInfoFieldNames{TInfoNameCharacterWidth, TInfoNameNumberOfLines, TInfoNameNone, TInfoNameNone}
	tInfoCharacterWidthHeight = TInfoFieldNames{TInfoNameCharacterWidth, TInfoNameCharacterScreenHeight, TInfoNameNone, TInfoNameNone}
	tInfoPixelDimColors       = TInfoFieldNames{TInfoNamePixelWidth, TInfoNamePixelHeight, TInfoNameNumberOfColors, TInfoNameNone}
	tInfoPixelDimDepth        = TInfoFieldNames{TInfoNamePixelWidth, TInfoNamePixelHeight, TInfoNamePixelDepth, TInfoNameNone}
	tInfoSampleRate           = TInfoFieldNames{TInfoNameSampleRate, TInfoNameNone, TInfoNameNone, TInfoNameNone}
)

// tInfoFieldMap maps DataType and FileType to field names
var tInfoFieldMap = map[byte]map[byte]TInfoFieldNames{
	DataTypeNone: {
		0: tInfoNone,
	},
	DataTypeCharacter: {
		FileTypeCharacterASCII:      tInfoCharacterWidthLines,
		FileTypeCharacterANSI:       tInfoCharacterWidthLines,
		FileTypeCharacterANSIMation: tInfoCharacterWidthHeight,
		FileTypeCharacterRIPScript:  tInfoPixelDimColors,
		FileTypeCharacterPCBoard:    tInfoCharacterWidthLines,
		FileTypeCharacterAvatar:     tInfoCharacterWidthLines,
		FileTypeCharacterHTML:       tInfoNone,
		FileTypeCharacterSource:     tInfoNone,
		FileTypeCharacterTundraDraw: tInfoCharacterWidthLines,
	},
	DataTypeBitmap: {
		0:  tInfoPixelDimDepth, // GIF
		1:  tInfoPixelDimDepth, // PCX
		2:  tInfoPixelDimDepth, // LBM/IFF
		3:  tInfoPixelDimDepth, // TGA
		4:  tInfoPixelDimDepth, // FLI
		5:  tInfoPixelDimDepth, // FLC
		6:  tInfoPixelDimDepth, // BMP
		7:  tInfoPixelDimDepth, // GL
		8:  tInfoPixelDimDepth, // DL
		9:  tInfoPixelDimDepth, // WPG
		10: tInfoPixelDimDepth, // PNG
		11: tInfoPixelDimDepth, // JPG
		12: tInfoPixelDimDepth, // MPG
		13: tInfoPixelDimDepth, // AVI
	},
	DataTypeVector: {
		0: tInfoNone, // DXF
		1: tInfoNone, // DWG
		2: tInfoNone, // WPG
		3: tInfoNone, // 3DS
	},
	DataTypeAudio: {
		0:  tInfoNone,       // MOD
		1:  tInfoNone,       // 669
		2:  tInfoNone,       // STM
		3:  tInfoNone,       // S3M
		4:  tInfoNone,       // MTM
		5:  tInfoNone,       // FAR
		6:  tInfoNone,       // ULT
		7:  tInfoNone,       // AMF
		8:  tInfoNone,       // DMF
		9:  tInfoNone,       // OKT
		10: tInfoNone,       // ROL
		11: tInfoNone,       // CMF
		12: tInfoNone,       // MID
		13: tInfoNone,       // SADT
		14: tInfoNone,       // VOC
		15: tInfoNone,       // WAV
		16: tInfoSampleRate, // SMP8
		17: tInfoSampleRate, // SMP8S
		18: tInfoSampleRate, // SMP16
		19: tInfoSampleRate, // SMP16S
		20: tInfoNone,       // PATCH8
		21: tInfoNone,       // PATCH16
		22: tInfoNone,       // XM
		23: tInfoNone,       // HSC
		24: tInfoNone,       // IT
	},
	DataTypeBinaryText: {
		// FileType is variable for BinaryText (encodes width)
	},
	DataTypeXBin: {
		0: tInfoCharacterWidthLines,
	},
	DataTypeArchive: {
		0: tInfoNone, // ZIP
		1: tInfoNone, // ARJ
		2: tInfoNone, // LZH
		3: tInfoNone, // ARC
		4: tInfoNone, // TAR
		5: tInfoNone, // ZOO
		6: tInfoNone, // RAR
		7: tInfoNone, // UC2
		8: tInfoNone, // PAK
		9: tInfoNone, // SQZ
	},
	DataTypeExecutable: {
		0: tInfoNone,
	},
}

type SAUCE struct {
	ID       string     // 5 bytes: Should be "SAUCE"
	Version  string     // 2 bytes: Should be "00"
	Title    string     // 35 bytes: Title of the file
	Author   string     // 20 bytes: Creator (nick)name or handle
	Group    string     // 20 bytes: Group or company name
	Date     string     // 8 bytes: CCYYMMDD format
	FileSize uint32     // 4 bytes: Original file size (little-endian)
	DataType byte       // 1 byte: Type of data
	FileType byte       // 1 byte: Type of file
	TInfo1   TInfoField // 2 bytes: Type dependent info (little-endian)
	TInfo2   TInfoField // 2 bytes: Type dependent info (little-endian)
	TInfo3   TInfoField // 2 bytes: Type dependent info (little-endian)
	TInfo4   TInfoField // 2 bytes: Type dependent info (little-endian)
	Comments byte       // 1 byte: Number of comment lines
	TFlags   byte       // 1 byte: Type dependent flags
	TInfoS   string     // 22 bytes: Type dependent string (null-terminated)
}

// DetectEncoding attempts to determine if file data is CP437 or ISO-8859-1
// Uses chardet library for detection, with fallback heuristics for ANSI art
func DetectEncoding(data []byte) string {
	pointsForCP437, pointsForISO := 0, 0

	// first, try CP437
	decoder := charmap.CodePage437.NewDecoder()
	CP437Translated, err := decoder.Bytes(data)
	if err != nil {
		pointsForISO += 1
	} else {
		// now try encoding the cp437 data as ISO-8859-1
		encoder := charmap.ISO8859_1.NewEncoder()
		_, err = encoder.Bytes(CP437Translated)
		// if that failed, then it's definitely cp437
		if err != nil {
			fmt.Fprintln(os.Stderr, "  - translation from CP437 to ISO-8859-1 failed!!")
			pointsForCP437 += 1
		}
	}

	// second, try ISO-8859-1, then re-encode as CP437
	decoder = charmap.ISO8859_1.NewDecoder()
	ISOTranslated, err := decoder.Bytes(data)
	if err != nil {
		pointsForCP437 += 1
	} else {
		encoder := charmap.CodePage437.NewEncoder()
		_, err = encoder.Bytes(ISOTranslated)
		// if that failed, then it's definitely iso-8859-1
		if err != nil {
			fmt.Fprintln(os.Stderr, "  - translation from ISO-8859-1 to CP437 failed!!")
			pointsForISO += 1
		}
	}

	// having more of these chars is usually bad
	for _, ch := range [][]byte{[]byte("»"), []byte("Ü"), []byte("╖")} {
		origCount := bytes.Count(data, ch)
		cp437Count := bytes.Count(CP437Translated, ch)
		isoCount := bytes.Count(ISOTranslated, ch)

		if cp437Count > isoCount {
			fmt.Fprintf(os.Stderr, "  - char %q: orig=%d, \x1b[91mcp437=%d\x1b[0m, iso=%d\n", ch, origCount, cp437Count, isoCount)
			pointsForISO += 1
		} else if isoCount > cp437Count {
			fmt.Fprintf(os.Stderr, "  - char %q: orig=%d, cp437=%d, \x1b[91miso=%d\x1b[0m\n", ch, origCount, cp437Count, isoCount)
			pointsForCP437 += 1
		} else {
			fmt.Fprintf(os.Stderr, "  - char %q: orig=%d, cp437=%d, iso=%d\n", ch, origCount, cp437Count, isoCount)
		}
	}
	// having more of these chars is usually good
	// other candidates: []byte("²")
	for _, ch := range [][]byte{[]byte("█"), []byte("¯"), []byte("░"), []byte("┌")} {
		origCount := bytes.Count(data, ch)
		cp437Count := bytes.Count(CP437Translated, ch)
		isoCount := bytes.Count(ISOTranslated, ch)

		if cp437Count > isoCount {
			fmt.Fprintf(os.Stderr, "  - char %q: orig=%d, \x1b[92mcp437=%d\x1b[0m, iso=%d\n", ch, origCount, cp437Count, isoCount)
			pointsForCP437 += 1
		} else if isoCount > cp437Count {
			fmt.Fprintf(os.Stderr, "  - char %q: orig=%d, cp437=%d, \x1b[92miso=%d\x1b[0m\n", ch, origCount, cp437Count, isoCount)
			pointsForISO += 1
		} else {
			fmt.Fprintf(os.Stderr, "  - char %q: orig=%d, cp437=%d, iso=%d\n", ch, origCount, cp437Count, isoCount)
		}
	}
	fmt.Fprintf(os.Stderr, "Points:\n- CP437:      %d\n- ISO-8859-1: %d\n", pointsForCP437, pointsForISO)

	if pointsForCP437 > pointsForISO {
		return "cp437"
	} else if pointsForISO > pointsForCP437 {
		return "iso-8859-1"
	}
	fmt.Fprintln(os.Stderr, "Unable to detect encoding, assuming ASCII")
	return "ascii"
}

func ParseSAUCEFromFile(path string) (*SAUCE, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("error reading file %s: %v", path, err)
	}
	return ParseSAUCE(data)
}

// ParseSAUCE parses SAUCE metadata from the last 128 bytes of data
// Returns nil if no valid SAUCE record is found
func ParseSAUCE(data []byte) (*SAUCE, string, error) {
	// SAUCE record is always 128 bytes at the end of the file
	if len(data) < 128 {
		return nil, "", fmt.Errorf("data too short to contain SAUCE record")
	}

	// Find the "\x1a" character that indicates the start of the SAUCE record
	sauceIdx := bytes.Index(data, []byte{0x1a})
	sauceData := data[sauceIdx+1:]
	// sauceData := data[len(data)-128:]

	// Detect encoding and decode file data to UTF-8
	fileDataRaw := data[:len(data)-128]
	encoding := DetectEncoding(fileDataRaw)

	var decoder *charmap.Charmap
	decoder = charmap.ISO8859_1

	isASCII := false
	switch encoding {
	case "cp437":
		fmt.Fprintf(os.Stderr, "Detected encoding: CP437\n")
		decoder = charmap.CodePage437
	case "iso-8859-1":
		fmt.Fprintf(os.Stderr, "Detected encoding: ISO-8859-1\n")
		decoder = charmap.ISO8859_1
	case "ascii":
		fmt.Fprintf(os.Stderr, "Detected encoding: ASCII\n")
		isASCII = true
	default:
		fmt.Fprintf(os.Stderr, "Unknown encoding detected, defaulting to cp437\n")
		decoder = charmap.CodePage437
	}

	var fileData []byte
	if isASCII {
		fileData = fileDataRaw
	} else {
		decodedFileData, err := decoder.NewDecoder().Bytes(fileDataRaw)
		if err != nil {
			log.Fatalf("Error decoding file data: %v", err)
		}
		fileData = decodedFileData
	}

	// strip the \x1a EOF character if present
	fileData = bytes.TrimRight(fileData, "\x1a")

	// Create a reader for binary data
	reader := bytes.NewReader(sauceData)

	sauce := &SAUCE{}

	// Read ID (5 bytes)
	idBytes := make([]byte, 5)
	reader.Read(idBytes)
	sauce.ID = string(idBytes)

	// Check if this is a valid SAUCE record
	if sauce.ID != "SAUCE" {
		return nil, "", fmt.Errorf("no valid SAUCE record found")
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

	binary.Read(reader, binary.LittleEndian, &sauce.FileSize)     // Read FileSize (4 bytes, little-endian unsigned)
	binary.Read(reader, binary.LittleEndian, &sauce.DataType)     // Read DataType (1 byte)
	binary.Read(reader, binary.LittleEndian, &sauce.FileType)     // Read FileType (1 byte)
	binary.Read(reader, binary.LittleEndian, &sauce.TInfo1.Value) // Read TInfo1 (2 bytes, little-endian)
	binary.Read(reader, binary.LittleEndian, &sauce.TInfo2.Value) // Read TInfo2 (2 bytes, little-endian)
	binary.Read(reader, binary.LittleEndian, &sauce.TInfo3.Value) // Read TInfo3 (2 bytes, little-endian)
	binary.Read(reader, binary.LittleEndian, &sauce.TInfo4.Value) // Read TInfo4 (2 bytes, little-endian)
	binary.Read(reader, binary.LittleEndian, &sauce.Comments)     // Read Comments (1 byte)
	binary.Read(reader, binary.LittleEndian, &sauce.TFlags)       // Read TFlags (1 byte)
	tinfoSBytes := make([]byte, 22)                               // Read TInfoS (22 bytes) - null-terminated string
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
			sauce.TInfo1.Name = TInfoNameNone
			sauce.TInfo2.Name = TInfoNameNone
			sauce.TInfo3.Name = TInfoNameNone
			sauce.TInfo4.Name = TInfoNameNone
		} else if fieldNames, exists := dataTypeMap[sauce.FileType]; exists {
			sauce.TInfo1.Name = fieldNames.TInfo1Name
			sauce.TInfo2.Name = fieldNames.TInfo2Name
			sauce.TInfo3.Name = fieldNames.TInfo3Name
			sauce.TInfo4.Name = fieldNames.TInfo4Name
		} else {
			// Unknown FileType - default to "0"
			sauce.TInfo1.Name = TInfoNameNone
			sauce.TInfo2.Name = TInfoNameNone
			sauce.TInfo3.Name = TInfoNameNone
			sauce.TInfo4.Name = TInfoNameNone
		}
	} else {
		// Unknown DataType - default to "0"
		sauce.TInfo1.Name = TInfoNameNone
		sauce.TInfo2.Name = TInfoNameNone
		sauce.TInfo3.Name = TInfoNameNone
		sauce.TInfo4.Name = TInfoNameNone
	}

	return sauce, string(fileData), nil
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

func (s *SAUCE) ToJSON() (string, error) {
	jsonBytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func (s *SAUCE) ToString() string {
	cYellow := "\033[33m"
	cBold := "\033[1m"
	cReset := "\033[0m"
	cItalic := "\033[3m"

	fmtStr := fmt.Sprintf("%%s%%-%ds%%s %%v\n", 10)

	var sb strings.Builder
	sb.WriteString(cYellow + cBold + "SAUCE Metadata:\n" + cReset)
	sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "ID:", cReset, fmt.Sprintf("%s%s%s", cItalic, s.ID, cReset)))
	sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "Version:", cReset, fmt.Sprintf("%s%s%s", cItalic, s.Version, cReset)))
	sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "Title:", cReset, s.Title))
	sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "Author:", cReset, s.Author))
	sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "Group:", cReset, s.Group))
	sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "Date:", cReset, s.Date))
	sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "File Size:", cReset, s.FileSize))
	sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "Data Type:", cReset, fmt.Sprintf("(%v) %s", s.DataType, dataTypeNames[s.DataType])))
	sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "File Type:", cReset, fmt.Sprintf("(%v) %s", s.FileType, fileTypeNames[s.FileType])))
	if s.TInfo1.Name != TInfoNameNone {
		sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "TInfo1:", cReset, ""))
		sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "  Name:", cReset, s.TInfo1.Name))
		sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "  Value:", cReset, s.TInfo1.Value))
	}
	if s.TInfo2.Name != TInfoNameNone {
		sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "TInfo2:", cReset, ""))
		sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "  Name:", cReset, s.TInfo2.Name))
		sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "  Value:", cReset, s.TInfo2.Value))
	}
	if s.TInfo3.Name != TInfoNameNone {
		sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "TInfo3:", cReset, ""))
		sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "  Name:", cReset, s.TInfo3.Name))
		sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "  Value:", cReset, s.TInfo3.Value))
	}
	if s.TInfo4.Name != TInfoNameNone {
		sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "TInfo4:", cReset, ""))
		sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "  Name:", cReset, s.TInfo4.Name))
		sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "  Value:", cReset, s.TInfo4.Value))
	}
	sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "Comments:", cReset, string(s.Comments)))
	sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "TFlags:", cReset, string(s.TFlags)))
	sb.WriteString(fmt.Sprintf(fmtStr, cYellow, "TInfoS:", cReset, s.TInfoS))
	return sb.String()
}
