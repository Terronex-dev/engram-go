// Package engram provides a Go SDK for reading and writing Engram memory files.
package engram

import (
	"time"
)

// MemoryNode represents a single memory in an Engram file.
type MemoryNode struct {
	ID        string     `msgpack:"id"`
	Content   string     `msgpack:"content"`
	Embedding []float32  `msgpack:"embedding,omitempty"`
	Tags      []string   `msgpack:"tags,omitempty"`
	Entities  []Entity   `msgpack:"entities,omitempty"`
	Links     []Link     `msgpack:"links,omitempty"`
	Metadata  NodeMeta   `msgpack:"metadata,omitempty"`
	Children  []string   `msgpack:"children,omitempty"`
	ParentID  string     `msgpack:"parentId,omitempty"`
}

// Entity represents a named entity extracted from content.
type Entity struct {
	Name       string  `msgpack:"name"`
	Type       string  `msgpack:"type"`
	Confidence float32 `msgpack:"confidence,omitempty"`
	Start      int     `msgpack:"start,omitempty"`
	End        int     `msgpack:"end,omitempty"`
}

// Link represents a relationship between memory nodes.
type Link struct {
	TargetID string  `msgpack:"targetId"`
	Type     string  `msgpack:"type"`
	Weight   float32 `msgpack:"weight,omitempty"`
	Metadata map[string]interface{} `msgpack:"metadata,omitempty"`
}

// NodeMeta contains metadata for a memory node.
type NodeMeta struct {
	Source      string                 `msgpack:"source,omitempty"`
	CreatedAt   *time.Time             `msgpack:"createdAt,omitempty"`
	UpdatedAt   *time.Time             `msgpack:"updatedAt,omitempty"`
	Confidence  float32                `msgpack:"confidence,omitempty"`
	Importance  float32                `msgpack:"importance,omitempty"`
	AccessCount int                    `msgpack:"accessCount,omitempty"`
	Custom      map[string]interface{} `msgpack:"custom,omitempty"`
}

// EngramHeader contains file-level metadata.
type EngramHeader struct {
	Version   string       `msgpack:"version"`
	Created   string       `msgpack:"created"`
	Modified  string       `msgpack:"modified"`
	NodeCount int          `msgpack:"nodeCount"`
	Schema    SchemaInfo   `msgpack:"schema,omitempty"`
	Security  SecurityInfo `msgpack:"security,omitempty"`
	Metadata  FileMeta     `msgpack:"metadata,omitempty"`
}

// SchemaInfo describes the schema version and features.
type SchemaInfo struct {
	Version        string   `msgpack:"version,omitempty"`
	EmbeddingModel string   `msgpack:"embeddingModel,omitempty"`
	EmbeddingDim   int      `msgpack:"embeddingDim,omitempty"`
	Features       []string `msgpack:"features,omitempty"`
}

// SecurityInfo contains integrity and encryption information.
type SecurityInfo struct {
	Integrity  string `msgpack:"integrity,omitempty"`
	Encryption string `msgpack:"encryption,omitempty"`
	KeyID      string `msgpack:"keyId,omitempty"`
}

// FileMeta contains file-level custom metadata.
type FileMeta struct {
	Title       string                 `msgpack:"title,omitempty"`
	Description string                 `msgpack:"description,omitempty"`
	Author      string                 `msgpack:"author,omitempty"`
	License     string                 `msgpack:"license,omitempty"`
	Tags        []string               `msgpack:"tags,omitempty"`
	Custom      map[string]interface{} `msgpack:"custom,omitempty"`
}

// EngramFile represents a complete Engram file with header and nodes.
type EngramFile struct {
	Header EngramHeader
	Nodes  []MemoryNode
}
