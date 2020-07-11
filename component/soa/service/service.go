package service

import (
	"math/rand"
	"net"
	"strconv"
	"time"
)

// Node 节点
type Node struct {
	// 服务ID
	ID string
	// 服务名称
	Name string
	// 服务地址
	Address string
	// 服务端口
	Port int
	// 服务标签
	Tags []string
}

// GetHost 获取服务完整的地址
func (Node Node) GetHost() string {
	return net.JoinHostPort(Node.Address, strconv.Itoa(Node.Port))
}

// Group
type Group struct {
	items  []Node
	cursor int
}

// Count return node count
func (group *Group) Count() int {
	return len(group.items)
}

// First return fist node
func (group *Group) First() Node {
	return group.items[0]
}

// Rand return a rand node
func (group *Group) Rand() Node {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator

	return group.items[rand.Intn(group.Count())]
}

// Next return the next node
func (group *Group) Next() Node {
	if group.cursor == group.Count() {
		group.cursor = 0
	}

	Node := group.items[group.cursor]
	group.cursor++
	return Node
}

// Add add node
func (group *Group) Add(node Node) {
	group.items = append(group.items, node)
}
