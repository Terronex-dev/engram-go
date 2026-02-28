# Engram Go SDK

Go SDK for reading and writing Engram memory files.

## Installation

```bash
go get github.com/Terronex-dev/engram-go
```

## Usage

### Reading an Engram File

```go
package main

import (
    "fmt"
    "log"

    engram "github.com/Terronex-dev/engram-go"
)

func main() {
    // Read file
    file, err := engram.ReadFile("memories.engram")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Loaded %d memories\n", len(file.Nodes))

    // Create a tree for navigation and search
    tree := engram.NewMemoryTree(file.Nodes)

    // Get by tag
    important := tree.GetByTag("important")
    fmt.Printf("Found %d important memories\n", len(important))

    // Text search
    matches := tree.SearchByContent("keyword", 10)
    for _, node := range matches {
        fmt.Printf("- %s: %s\n", node.ID, node.Content[:50])
    }
}
```

### Writing an Engram File

```go
package main

import (
    "log"

    engram "github.com/Terronex-dev/engram-go"
)

func main() {
    nodes := []engram.MemoryNode{
        {
            ID:      "memory-1",
            Content: "First memory content",
            Tags:    []string{"important"},
        },
        {
            ID:       "memory-2",
            Content:  "Second memory content",
            ParentID: "memory-1",
        },
    }

    file := &engram.EngramFile{
        Header: engram.EngramHeader{
            Version: "1.0",
            Metadata: engram.FileMeta{
                Title:  "My Memories",
                Author: "Gopher",
            },
        },
        Nodes: nodes,
    }

    if err := engram.WriteFile("output.engram", file); err != nil {
        log.Fatal(err)
    }
}
```

### Semantic Search

```go
// Assuming you have embeddings in your nodes
tree := engram.NewMemoryTree(file.Nodes)

// Search with a query embedding
queryEmbedding := []float32{0.1, 0.2, 0.3, /* ... */}
results := tree.Search(queryEmbedding, 5)

for _, result := range results {
    fmt.Printf("Score: %.3f - %s\n", result.Score, result.Node.Content)
}
```

### Streaming Large Files

```go
import "os"

f, _ := os.Open("large.engram")
defer f.Close()

reader, err := engram.NewStreamReader(f)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("File has %d nodes\n", reader.Header.NodeCount)

for {
    node, err := reader.Next()
    if err != nil {
        log.Fatal(err)
    }
    if node == nil {
        break // End of file
    }
    // Process node...
}
```

## API Reference

### Types

- `MemoryNode` - A single memory with content, embeddings, tags, entities, and links
- `EngramFile` - Complete file with header and nodes
- `EngramHeader` - File metadata including version, schema, and security info
- `MemoryTree` - Navigation and search structure

### Functions

- `ReadFile(path)` - Read an Engram file from disk
- `WriteFile(path, file)` - Write an Engram file to disk
- `Decode(data)` - Decode Engram data from bytes
- `Encode(file)` - Encode an Engram file to bytes
- `VerifyIntegrity(path)` - Check file integrity
- `NewStreamReader(reader)` - Create a streaming reader for large files

### MemoryTree Methods

- `Get(id)` - Get node by ID
- `GetAll()` - Get all nodes
- `Count()` - Total node count
- `GetByTag(tag)` - Get nodes with a specific tag
- `GetTags()` - Get all unique tags
- `GetChildren(parentID)` - Get direct children
- `GetRoots()` - Get nodes without parents
- `Search(embedding, limit)` - Semantic search
- `SearchByContent(query, limit)` - Text search
- `Filter(predicate)` - Filter with custom function

## Cross-SDK Compatibility

This SDK produces files compatible with:
- [@terronex/engram](https://www.npmjs.com/package/@terronex/engram) (TypeScript)
- [engram-py](https://github.com/Terronex-dev/engram-py) (Python)
- [engram-rs](https://github.com/Terronex-dev/engram-rs) (Rust)

## License

MIT License - see [LICENSE](LICENSE)
