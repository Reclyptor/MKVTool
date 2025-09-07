package main

import (
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func atoi64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

func filter[T any](list []T, predicate func(T) bool) []T {
	var results []T
	for _, item := range list {
		if predicate(item) {
			results = append(results, item)
		}
	}
	return results
}

func sizeToBytes(size string) int64 {
	size = strings.TrimSpace(size)
	if size == "" {
		return 0
	}

	regex := regexp.MustCompile(`^(?P<val>[0-9]+) ?((?P<mult>[Kk]|[Mm]|[Gg]|[Tt]|[Pp])[Bb]?)?$`)
	matches := regex.FindStringSubmatch(size)
	if matches == nil {
		return 0
	}

	value, _ := strconv.ParseInt(matches[1], 10, 64)
	unit := strings.ToUpper(matches[3])
	switch unit {
	case "K":
		return value * 1024
	case "M":
		return value * 1024 * 1024
	case "G":
		return value * 1024 * 1024 * 1024
	case "T":
		return value * 1024 * 1024 * 1024 * 1024
	case "P":
		return value * 1024 * 1024 * 1024 * 1024 * 1024
	default:
		return value
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
