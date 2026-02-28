package engram

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

// Magic bytes for Engram files
var MagicBytes = []byte{0x45, 0x4E, 0x47, 0x52, 0x41, 0x4D} // "ENGRAM"

// ErrInvalidMagic is returned when the file doesn't have valid magic bytes.
var ErrInvalidMagic = errors.New("invalid magic bytes: not an Engram file")

// ErrIntegrityFailed is returned when the integrity check fails.
var ErrIntegrityFailed = errors.New("integrity check failed: file may be corrupted")

// ReadFile reads an Engram file from disk.
func ReadFile(path string) (*EngramFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return Decode(data)
}

// Decode decodes Engram data from bytes.
func Decode(data []byte) (*EngramFile, error) {
	// Check magic bytes
	if len(data) < 6 || !bytes.Equal(data[:6], MagicBytes) {
		return nil, ErrInvalidMagic
	}

	// Skip magic bytes
	reader := bytes.NewReader(data[6:])
	decoder := msgpack.NewDecoder(reader)

	// Decode header
	var header EngramHeader
	if err := decoder.Decode(&header); err != nil {
		return nil, fmt.Errorf("failed to decode header: %w", err)
	}

	// Read remaining bytes as payload
	payloadStart := 6 + (len(data) - 6 - reader.Len())
	payloadBytes := data[payloadStart:]

	// Verify integrity if present
	if header.Security.Integrity != "" {
		hash := sha256.Sum256(payloadBytes)
		computed := hex.EncodeToString(hash[:])
		if computed != header.Security.Integrity {
			return nil, ErrIntegrityFailed
		}
	}

	// Decode nodes from payload
	var nodes []MemoryNode
	payloadDecoder := msgpack.NewDecoder(bytes.NewReader(payloadBytes))
	if err := payloadDecoder.Decode(&nodes); err != nil {
		return nil, fmt.Errorf("failed to decode nodes: %w", err)
	}

	return &EngramFile{
		Header: header,
		Nodes:  nodes,
	}, nil
}

// WriteFile writes an Engram file to disk.
func WriteFile(path string, file *EngramFile) error {
	data, err := Encode(file)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Encode encodes an Engram file to bytes.
func Encode(file *EngramFile) ([]byte, error) {
	// Encode payload first to compute integrity hash
	payloadBytes, err := msgpack.Marshal(file.Nodes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode nodes: %w", err)
	}

	// Compute integrity hash
	hash := sha256.Sum256(payloadBytes)
	integrity := hex.EncodeToString(hash[:])

	// Update header
	header := file.Header
	header.NodeCount = len(file.Nodes)
	header.Modified = time.Now().UTC().Format(time.RFC3339)
	if header.Created == "" {
		header.Created = header.Modified
	}
	if header.Version == "" {
		header.Version = "1.0"
	}
	header.Security.Integrity = integrity

	// Encode header
	headerBytes, err := msgpack.Marshal(header)
	if err != nil {
		return nil, fmt.Errorf("failed to encode header: %w", err)
	}

	// Combine: magic + header + payload
	var buf bytes.Buffer
	buf.Write(MagicBytes)
	buf.Write(headerBytes)
	buf.Write(payloadBytes)

	return buf.Bytes(), nil
}

// VerifyIntegrity checks the integrity of an Engram file.
func VerifyIntegrity(path string) (bool, error) {
	_, err := ReadFile(path)
	if err != nil {
		if errors.Is(err, ErrIntegrityFailed) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// StreamReader provides streaming read access to large Engram files.
type StreamReader struct {
	reader  io.Reader
	decoder *msgpack.Decoder
	Header  EngramHeader
	count   int
	index   int
}

// NewStreamReader creates a streaming reader from an io.Reader.
func NewStreamReader(r io.Reader) (*StreamReader, error) {
	// Read and verify magic bytes
	magic := make([]byte, 6)
	if _, err := io.ReadFull(r, magic); err != nil {
		return nil, fmt.Errorf("failed to read magic bytes: %w", err)
	}
	if !bytes.Equal(magic, MagicBytes) {
		return nil, ErrInvalidMagic
	}

	decoder := msgpack.NewDecoder(r)

	// Decode header
	var header EngramHeader
	if err := decoder.Decode(&header); err != nil {
		return nil, fmt.Errorf("failed to decode header: %w", err)
	}

	return &StreamReader{
		reader:  r,
		decoder: decoder,
		Header:  header,
		count:   header.NodeCount,
		index:   0,
	}, nil
}

// Next returns the next memory node, or nil if done.
func (sr *StreamReader) Next() (*MemoryNode, error) {
	if sr.index >= sr.count {
		return nil, nil
	}

	var node MemoryNode
	if err := sr.decoder.Decode(&node); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to decode node: %w", err)
	}

	sr.index++
	return &node, nil
}
