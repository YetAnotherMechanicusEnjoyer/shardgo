package cluster

import (
	"sync"
)

type Manager struct {
	shard *ConsistentHash
	mutex sync.RWMutex
}

func NewManager(replicas int) *Manager {
	return &Manager{
		shard: NewConsistentHash(replicas),
	}
}

func (m *Manager) AddNode(node string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.shard.AddNode(node)
}

func (m *Manager) RemoveNode(node string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.shard.RemoveNode(node)
}

func (m *Manager) GetNodeForKey(key string) (string, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.shard.GetNode(key)
}

func (m *Manager) GetAllNodes() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.shard.GetAllNodes()
}
