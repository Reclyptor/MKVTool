package main

import (
	"bufio"
	"regexp"
	"strings"
)

var regex = regexp.MustCompile(`^(CINFO|TINFO|SINFO):([\d]+)(?:,([\d]+))?(?:,([\d]+))?,\d+,"(.*)"$`)

type Record string

const (
	CINFO Record = "CINFO"
	TINFO Record = "TINFO"
	SINFO Record = "SINFO"
)

type CInfo struct {
	Field int
	Value string
}

type TInfo struct {
	Track int
	Field int
	Value string
}

type SInfo struct {
	Track  int
	Stream int
	Field  int
	Value  string
}

func tokenize(line string) interface{} {
	matches := regex.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	switch Record(matches[1]) {
	case CINFO:
		return CInfo{
			Field: atoi(matches[2]),
			Value: matches[5],
		}
	case TINFO:
		return TInfo{
			Track: atoi(matches[2]),
			Field: atoi(matches[3]),
			Value: matches[5],
		}
	case SINFO:
		return SInfo{
			Track:  atoi(matches[2]),
			Stream: atoi(matches[3]),
			Field:  atoi(matches[4]),
			Value:  matches[5],
		}
	default:
		return nil
	}
}

func Tokenize(input string) ([]CInfo, []TInfo, []SInfo) {
	var cinfos []CInfo
	var tinfos []TInfo
	var sinfos []SInfo
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		if record := tokenize(line); record != nil {
			switch record.(type) {
			case CInfo:
				cinfos = append(cinfos, record.(CInfo))
			case TInfo:
				tinfos = append(tinfos, record.(TInfo))
			case SInfo:
				sinfos = append(sinfos, record.(SInfo))
			}
		}
	}
	return cinfos, tinfos, sinfos
}
