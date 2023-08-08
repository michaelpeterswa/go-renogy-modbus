package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"676f.dev/utilities/tools/shortid"
)

func dump(b []byte) error {
	sid := shortid.NewShortID(shortid.Base58CharacterSet)
	shortId, err := sid.Generate(6)
	if err != nil {
		return fmt.Errorf("failed to generate short id: %w", err)
	}

	fileName := fmt.Sprintf("%s-%s.bin", time.Now().Format("2006-01-02"), shortId)

	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func load(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return b, nil
}
