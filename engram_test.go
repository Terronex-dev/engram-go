package engram

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRoundtrip(t *testing.T) {
	// Create test nodes
	nodes := []MemoryNode{
		{
			ID:      "node-1",
			Content: "First memory",
			Tags:    []string{"test", "first"},
		},
		{
			ID:       "node-2",
			Content:  "Second memory",
			Tags:     []string{"test", "second"},
			ParentID: "node-1",
		},
	}

	file := &EngramFile{
		Header: EngramHeader{
			Version: "1.0",
			Metadata: FileMeta{
				Title: "Test File",
			},
		},
		Nodes: nodes,
	}

	// Encode
	data, err := Encode(file)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Verify magic bytes
	if len(data) < 6 {
		t.Fatal("Data too short")
	}
	for i, b := range MagicBytes {
		if data[i] != b {
			t.Fatalf("Magic byte %d mismatch: got %x, want %x", i, data[i], b)
		}
	}

	// Decode
	decoded, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Verify
	if decoded.Header.NodeCount != 2 {
		t.Errorf("NodeCount: got %d, want 2", decoded.Header.NodeCount)
	}
	if len(decoded.Nodes) != 2 {
		t.Errorf("Nodes length: got %d, want 2", len(decoded.Nodes))
	}
	if decoded.Nodes[0].ID != "node-1" {
		t.Errorf("First node ID: got %s, want node-1", decoded.Nodes[0].ID)
	}
	if decoded.Nodes[1].ParentID != "node-1" {
		t.Errorf("Second node ParentID: got %s, want node-1", decoded.Nodes[1].ParentID)
	}
}

func TestIntegrity(t *testing.T) {
	nodes := []MemoryNode{
		{ID: "test", Content: "Test content"},
	}

	file := &EngramFile{
		Header: EngramHeader{Version: "1.0"},
		Nodes:  nodes,
	}

	data, err := Encode(file)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Verify integrity hash was set
	decoded, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded.Header.Security.Integrity == "" {
		t.Error("Integrity hash not set")
	}

	// Corrupt the data and verify detection
	corrupted := make([]byte, len(data))
	copy(corrupted, data)
	corrupted[len(corrupted)-1] ^= 0xFF // Flip last byte

	_, err = Decode(corrupted)
	if err != ErrIntegrityFailed {
		t.Errorf("Expected ErrIntegrityFailed, got: %v", err)
	}
}

func TestFileIO(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.engram")

	nodes := []MemoryNode{
		{ID: "file-test", Content: "File I/O test"},
	}

	file := &EngramFile{
		Header: EngramHeader{Version: "1.0"},
		Nodes:  nodes,
	}

	// Write
	if err := WriteFile(path, file); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("File was not created")
	}

	// Read
	read, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if read.Nodes[0].ID != "file-test" {
		t.Errorf("Read node ID: got %s, want file-test", read.Nodes[0].ID)
	}

	// Verify integrity
	valid, err := VerifyIntegrity(path)
	if err != nil {
		t.Fatalf("VerifyIntegrity failed: %v", err)
	}
	if !valid {
		t.Error("Integrity check failed")
	}
}

func TestMemoryTree(t *testing.T) {
	nodes := []MemoryNode{
		{ID: "root", Content: "Root node", Tags: []string{"important"}},
		{ID: "child-1", Content: "Child one", Tags: []string{"child"}, ParentID: "root"},
		{ID: "child-2", Content: "Child two", Tags: []string{"child", "important"}, ParentID: "root"},
	}

	tree := NewMemoryTree(nodes)

	// Test Count
	if tree.Count() != 3 {
		t.Errorf("Count: got %d, want 3", tree.Count())
	}

	// Test Get
	node := tree.Get("child-1")
	if node == nil || node.Content != "Child one" {
		t.Error("Get failed")
	}

	// Test GetByTag
	important := tree.GetByTag("important")
	if len(important) != 2 {
		t.Errorf("GetByTag: got %d nodes, want 2", len(important))
	}

	// Test GetTags
	tags := tree.GetTags()
	if len(tags) != 2 {
		t.Errorf("GetTags: got %d tags, want 2", len(tags))
	}

	// Test GetChildren
	children := tree.GetChildren("root")
	if len(children) != 2 {
		t.Errorf("GetChildren: got %d, want 2", len(children))
	}

	// Test GetRoots
	roots := tree.GetRoots()
	if len(roots) != 1 || roots[0].ID != "root" {
		t.Error("GetRoots failed")
	}

	// Test SearchByContent
	matches := tree.SearchByContent("child", 10)
	if len(matches) != 2 {
		t.Errorf("SearchByContent: got %d, want 2", len(matches))
	}

	// Test Filter
	filtered := tree.Filter(func(n *MemoryNode) bool {
		return n.ParentID != ""
	})
	if len(filtered) != 2 {
		t.Errorf("Filter: got %d, want 2", len(filtered))
	}
}

func TestCosineSimilarity(t *testing.T) {
	// Identical vectors
	a := []float32{1, 0, 0}
	b := []float32{1, 0, 0}
	sim := cosineSimilarity(a, b)
	if sim < 0.999 {
		t.Errorf("Identical vectors: got %f, want ~1.0", sim)
	}

	// Orthogonal vectors
	c := []float32{0, 1, 0}
	sim = cosineSimilarity(a, c)
	if sim > 0.001 {
		t.Errorf("Orthogonal vectors: got %f, want ~0.0", sim)
	}

	// Opposite vectors
	d := []float32{-1, 0, 0}
	sim = cosineSimilarity(a, d)
	if sim > -0.999 {
		t.Errorf("Opposite vectors: got %f, want ~-1.0", sim)
	}
}

func TestSemanticSearch(t *testing.T) {
	nodes := []MemoryNode{
		{ID: "a", Content: "Apple", Embedding: []float32{1, 0, 0}},
		{ID: "b", Content: "Banana", Embedding: []float32{0, 1, 0}},
		{ID: "c", Content: "Cherry", Embedding: []float32{0.9, 0.1, 0}},
	}

	tree := NewMemoryTree(nodes)

	// Search for something similar to Apple
	query := []float32{1, 0, 0}
	results := tree.Search(query, 2)

	if len(results) != 2 {
		t.Fatalf("Search: got %d results, want 2", len(results))
	}

	// Apple should be first (exact match)
	if results[0].Node.ID != "a" {
		t.Errorf("First result: got %s, want a", results[0].Node.ID)
	}

	// Cherry should be second (similar)
	if results[1].Node.ID != "c" {
		t.Errorf("Second result: got %s, want c", results[1].Node.ID)
	}
}
