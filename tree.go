package engram

import (
	"math"
	"sort"
	"strings"
)

// MemoryTree provides navigation and search over memory nodes.
type MemoryTree struct {
	nodes    []MemoryNode
	byID     map[string]*MemoryNode
	byTag    map[string][]*MemoryNode
	children map[string][]*MemoryNode
}

// NewMemoryTree creates a new tree from a slice of nodes.
func NewMemoryTree(nodes []MemoryNode) *MemoryTree {
	tree := &MemoryTree{
		nodes:    nodes,
		byID:     make(map[string]*MemoryNode),
		byTag:    make(map[string][]*MemoryNode),
		children: make(map[string][]*MemoryNode),
	}

	for i := range nodes {
		node := &nodes[i]
		tree.byID[node.ID] = node

		// Index by tags
		for _, tag := range node.Tags {
			tree.byTag[tag] = append(tree.byTag[tag], node)
		}

		// Index children
		if node.ParentID != "" {
			tree.children[node.ParentID] = append(tree.children[node.ParentID], node)
		}
	}

	return tree
}

// Get returns a node by ID.
func (t *MemoryTree) Get(id string) *MemoryNode {
	return t.byID[id]
}

// GetAll returns all nodes.
func (t *MemoryTree) GetAll() []MemoryNode {
	return t.nodes
}

// Count returns the total number of nodes.
func (t *MemoryTree) Count() int {
	return len(t.nodes)
}

// GetByTag returns all nodes with a specific tag.
func (t *MemoryTree) GetByTag(tag string) []*MemoryNode {
	return t.byTag[tag]
}

// GetTags returns all unique tags.
func (t *MemoryTree) GetTags() []string {
	tags := make([]string, 0, len(t.byTag))
	for tag := range t.byTag {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// GetChildren returns direct children of a node.
func (t *MemoryTree) GetChildren(parentID string) []*MemoryNode {
	return t.children[parentID]
}

// GetRoots returns all nodes without parents.
func (t *MemoryTree) GetRoots() []*MemoryNode {
	var roots []*MemoryNode
	for i := range t.nodes {
		if t.nodes[i].ParentID == "" {
			roots = append(roots, &t.nodes[i])
		}
	}
	return roots
}

// SearchResult represents a search match with score.
type SearchResult struct {
	Node  *MemoryNode
	Score float32
}

// Search performs semantic search using cosine similarity.
func (t *MemoryTree) Search(queryEmbedding []float32, limit int) []SearchResult {
	if len(queryEmbedding) == 0 {
		return nil
	}

	var results []SearchResult
	for i := range t.nodes {
		node := &t.nodes[i]
		if len(node.Embedding) == 0 {
			continue
		}

		score := cosineSimilarity(queryEmbedding, node.Embedding)
		results = append(results, SearchResult{
			Node:  node,
			Score: score,
		})
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Limit results
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results
}

// SearchByContent performs text search on content.
func (t *MemoryTree) SearchByContent(query string, limit int) []*MemoryNode {
	query = strings.ToLower(query)
	var matches []*MemoryNode

	for i := range t.nodes {
		if strings.Contains(strings.ToLower(t.nodes[i].Content), query) {
			matches = append(matches, &t.nodes[i])
			if limit > 0 && len(matches) >= limit {
				break
			}
		}
	}

	return matches
}

// Filter returns nodes matching a predicate.
func (t *MemoryTree) Filter(predicate func(*MemoryNode) bool) []*MemoryNode {
	var matches []*MemoryNode
	for i := range t.nodes {
		if predicate(&t.nodes[i]) {
			matches = append(matches, &t.nodes[i])
		}
	}
	return matches
}

// cosineSimilarity computes the cosine similarity between two vectors.
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return float32(dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)))
}
