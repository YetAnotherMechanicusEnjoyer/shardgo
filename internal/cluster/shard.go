package cluster

import (
	"hash/fnv"
	"slices"
	"sort"
)

type ConsistentHash struct {
	nodes    []uint32
	nodeMap  map[uint32]string
	replicas int
}

func NewConsistentHash(replicas int) *ConsistentHash {
	return &ConsistentHash{
		nodeMap:  make(map[uint32]string),
		replicas: replicas,
	}
}

func (c *ConsistentHash) AddNode(node string) {
	for i := 0; i < c.replicas; i++ {
		virtualNode := node + "#" + string(rune(i))
		hash := c.hash(virtualNode)
		c.nodes = append(c.nodes, hash)
		c.nodeMap[hash] = node
	}
	slices.Sort(c.nodes)
}

func (c *ConsistentHash) RemoveNode(node string) {
	for i := 0; i < c.replicas; i++ {
		virtualNode := node + "#" + string(rune(i))
		hash := c.hash(virtualNode)
		idx := sort.Search(len(c.nodes), func(i int) bool {
			return c.nodes[i] >= hash
		})
		if idx < len(c.nodes) && c.nodes[idx] == hash {
			c.nodes = append(c.nodes[:idx], c.nodes[idx+1:]...)
			delete(c.nodeMap, hash)
		}
	}
}

func (c *ConsistentHash) GetNode(key string) (string, bool) {
	if len(c.nodes) == 0 {
		return "", false
	}
	hash := c.hash(key)
	idx := sort.Search(len(c.nodes), func(i int) bool { return c.nodes[i] >= hash })
	if idx == len(c.nodes) {
		idx = 0
	}
	node, ok := c.nodeMap[c.nodes[idx]]
	return node, ok
}

func (c *ConsistentHash) hash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func (c *ConsistentHash) GetAllNodes() []string {
	nodeSet := make(map[string]bool)
	for _, node := range c.nodeMap {
		nodeSet[node] = true
	}
	nodes := make([]string, 0, len(nodeSet))
	for node := range nodeSet {
		nodes = append(nodes, node)
	}
	return nodes
}
