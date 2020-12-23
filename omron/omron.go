// Package omron assists in deciphering omron communication protocol message.
//
// This package is written by inspecting messages from the Chiltrix CX34.
package omron

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// CommandFrame is described here:
// https://paginas.fe.up.pt/~pfs/recursos/plcs/omron/cqm1/sbc_manual/sec4.pdf
// 4-4-1
type CommandFrame struct {
	NodeNumber    uint8
	HeaderCode    HeaderCode
	Text          string
	Check         uint8
	ParsedCommand Command
}

// Command is an interface implemented by all values of ParsedCommand.
type Command interface {
	// String returns a human readable explanation of the command.
	String() string
	parsedCommand()
}

var frameRE = regexp.MustCompile(`^@(\d\d)(..)(.*)(..)\*$`)

const (
	nodeNoGroup        = 1
	headerCodeGroup    = 2
	textGroup          = 3
	frameCheckSeqGroup = 4
)

// ParseCommandFrame returns a CommandFrame parsed from a line of text.
func ParseCommandFrame(line string) (*CommandFrame, error) {
	if !strings.HasPrefix(line, "@") {
		return nil, fmt.Errorf("invallid command frame %q - does not begin with @", line)
	}
	if !strings.HasSuffix(line, "*") {
		return nil, fmt.Errorf("invallid command frame %q - does not end with *\r", line)
	}
	m := frameRE.FindStringSubmatch(line)
	if len(m) == 0 {
		return nil, fmt.Errorf("invalid command frame %q does not match regexp %s", line, frameRE)
	}
	dest, headerCodeText, text, fcs := m[nodeNoGroup], m[headerCodeGroup], m[textGroup], m[frameCheckSeqGroup]
	destInt, err := strconv.ParseUint(dest, 10, 8)
	if err != nil {
		return nil, err
	}
	fcsByte, err := strconv.ParseUint(fcs, 16, 8)
	if err != nil {
		return nil, err
	}
	gotFCS, err := calculateCheckBit(line)
	if err != nil {
		return nil, err
	}
	if gotFCS != uint8(fcsByte) {
		return nil, fmt.Errorf("calcualted frame check %d, got %d: %q", gotFCS, fcsByte, line)
	}

	headerCode := HeaderCode(headerCodeText)
	var parsedCommand Command
	switch headerCode {
	case RD:
		if len(text) == 8 {
			begin, err := strconv.ParseUint(text[0:4], 10, 16)
			if err != nil {
				return nil, fmt.Errorf("could not parse RD begin word: %w", err)
			}
			count, err := strconv.ParseUint(text[4:8], 10, 16)
			if err != nil {
				return nil, fmt.Errorf("could not parse RD word count: %w", err)
			}
			parsedCommand = &DMAreaReadCommand{BeginningWord: uint16(begin), NumberOfWords: uint16(count)}
		} else if len(text) < 2 || ((len(text)-2)%4) != 0 {
			return nil, fmt.Errorf("unsupported RD command %q", text)
		} else {
			endCodeNum, err := strconv.ParseUint(text[0:2], 16, 8)
			if err != nil {
				return nil, err
			}
			endCode := EndCode(endCodeNum)
			var data []uint16
			for i := 2; i < len(text); i += 4 {
				word, err := strconv.ParseUint(text[i:i+4], 16, 16)
				if err != nil {
					return nil, err
				}
				data = append(data, uint16(word))
			}
			parsedCommand = &DMAreaReadResponse{EndCode: endCode, Data: data}
		}
	case WD:
		cmd, resp, err := parseDMAreaWrite(text)
		if err != nil {
			return nil, fmt.Errorf("bad WD command: %w", err)
		}
		if cmd != nil {
			parsedCommand = cmd
		} else {
			parsedCommand = resp
		}
	default:
		return nil, fmt.Errorf("unknown supported header code %q", headerCode)
	}

	return &CommandFrame{
		NodeNumber:    uint8(destInt),
		HeaderCode:    headerCode,
		Text:          text,
		Check:         uint8(fcsByte),
		ParsedCommand: parsedCommand,
	}, nil
}

func calculateCheckBit(line string) (uint8, error) {
	line = strings.TrimSuffix(line, "*")
	if len(line) > 2 {
		line = line[0 : len(line)-2]
	}
	check := uint8(0)
	for _, c := range []byte(line) {
		check ^= c
	}
	return check, nil
}

// HeaderCode is one of the codes from 4-4-3 of https://paginas.fe.up.pt/~pfs/recursos/plcs/omron/cqm1/sbc_manual/sec4.pdf.
// and https://assets.omron.eu/downloads/manual/en/w364_cqm1h_series_programming_manual_en.pdf.
type HeaderCode string

func (hc HeaderCode) Valid() bool {
	return knownCodes[hc]
}

// Valid header codes.
const (
	RD HeaderCode = "RD" // DM AREA READ
	WD HeaderCode = "WD" // DM AREA WRITE
)

var knownCodes = map[HeaderCode]bool{
	WD: true,
	RD: true,
}

type DMAreaReadCommand struct {
	// First word to read 0000-6655
	BeginningWord uint16
	// Number of words to read
	NumberOfWords uint16
}

func (r *DMAreaReadCommand) String() string {
	return fmt.Sprintf("RD command (%d words): %d", r.NumberOfWords, r.BeginningWord)
}

func (*DMAreaReadCommand) parsedCommand() {}

type DMAreaReadResponse struct {
	EndCode EndCode
	// Words of data.
	Data []uint16
}

func (r *DMAreaReadResponse) String() string {
	return fmt.Sprintf("RD response/%d: %v", r.EndCode, r.Data)
}
func (*DMAreaReadResponse) parsedCommand() {}

type EndCode uint8

func (c EndCode) String() string {
	if c == OK {
		return "OK"
	}
	return fmt.Sprintf("UNKNOWN/%d", c)
}

// Valid end codes.
const (
	OK EndCode = 0
)

type DMAreaWriteCommand struct {
	// First word to read 0000-6655
	BeginningWord uint16
	// Words to write starting at BeginningWord.
	Data []uint16
}

func (r *DMAreaWriteCommand) String() string {
	return fmt.Sprintf("WD write %d words from %d: %d", len(r.Data), r.BeginningWord, r.Data)
}

func (*DMAreaWriteCommand) parsedCommand() {}

type DMAreaWriteResponse struct {
	EndCode EndCode
	// Words of data.
	Data []uint16
}

func (r *DMAreaWriteResponse) String() string {
	return fmt.Sprintf("WD response: %s", r.EndCode)
}

func (*DMAreaWriteResponse) parsedCommand() {}

func parseDMAreaWrite(text string) (*DMAreaWriteCommand, *DMAreaWriteResponse, error) {
	if len(text) == 2 {
		endCodeNum, err := strconv.ParseUint(text[0:2], 16, 8)
		if err != nil {
			return nil, nil, err
		}
		endCode := EndCode(endCodeNum)
		return nil, &DMAreaWriteResponse{EndCode: endCode}, nil
	}
	if len(text) < 2 || len(text)%4 != 0 {
		return nil, nil, fmt.Errorf("invalid WD command %q", text)
	}
	begin, err := strconv.ParseUint(text[0:4], 10, 16)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse RD begin word: %w", err)
	}
	var data []uint16
	for i := 4; i < len(text); i += 4 {
		word, err := strconv.ParseUint(text[i:i+4], 16, 16)
		if err != nil {
			return nil, nil, err
		}
		data = append(data, uint16(word))
	}
	return &DMAreaWriteCommand{BeginningWord: uint16(begin), Data: data}, nil, nil

}
