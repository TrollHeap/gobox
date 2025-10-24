package sysfs

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func readFile(basePath, filename string) (string, error) {
	data, err := os.ReadFile(filepath.Join(basePath, filename))
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", filename, err)
	}
	return strings.TrimSpace(string(data)), nil
}

func ReadFileOptional(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// Ne PAS retourner l'erreur si fichier absent
		if os.IsNotExist(err) {
			return "", nil // Pas d'erreur, juste vide
		}
		return "", err // Autre erreur (permissions, etc.)
	}
	return strings.TrimSpace(string(data)), nil
}

func ReadInt(dir, filename string) (int, error) {
	s, err := readFile(dir, filename)
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parsing %s as int: %w", filename, err)
	}
	return val, nil
}

func ReadFloat(dir, filename string) (float64, error) {
	s, err := readFile(dir, filename)
	if err != nil {
		return 0, err
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing %s as float: %w", filename, err)
	}
	return val, nil
}

func ReadString(dir string, filename string) (string, error) {
	return readFile(dir, filename)
}
