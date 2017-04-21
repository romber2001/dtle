package api

import (
	"sort"
)

// Nodes is used to query node-related API endpoints
type Nodes struct {
	client *Client
}

// Nodes returns a handle on the node endpoints.
func (c *Client) Nodes() *Nodes {
	return &Nodes{client: c}
}

// List is used to list out all of the nodes
func (n *Nodes) List(q *QueryOptions) ([]*NodeListStub, *QueryMeta, error) {
	var resp NodeIndexSort
	qm, err := n.client.query("/v1/nodes", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(NodeIndexSort(resp))
	return resp, qm, nil
}

func (n *Nodes) PrefixList(prefix string) ([]*NodeListStub, *QueryMeta, error) {
	return n.List(&QueryOptions{Prefix: prefix})
}

// Info is used to query a specific node by its ID.
func (n *Nodes) Info(nodeID string, q *QueryOptions) (*Node, *QueryMeta, error) {
	var resp Node
	qm, err := n.client.query("/v1/node/"+nodeID, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// Allocations is used to return the allocations associated with a node.
func (n *Nodes) Allocations(nodeID string, q *QueryOptions) ([]*Allocation, *QueryMeta, error) {
	var resp []*Allocation
	qm, err := n.client.query("/v1/node/"+nodeID+"/allocations", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(AllocationSort(resp))
	return resp, qm, nil
}

// ForceEvaluate is used to force-evaluate an existing node.
func (n *Nodes) ForceEvaluate(nodeID string, q *WriteOptions) (string, *WriteMeta, error) {
	var resp nodeEvalResponse
	wm, err := n.client.write("/v1/node/"+nodeID+"/evaluate", nil, &resp, q)
	if err != nil {
		return "", nil, err
	}
	return resp.EvalID, wm, nil
}

// Node is used to deserialize a node entry.
type Node struct {
	ID                string
	Datacenter        string
	Name              string
	HTTPAddr          string
	Attributes        map[string]string
	Meta              map[string]string
	Status            string
	StatusDescription string
	StatusUpdatedAt   int64
	CreateIndex       uint64
	ModifyIndex       uint64
}

// NodeListStub is a subset of information returned during
// node list operations.
type NodeListStub struct {
	ID                string
	Datacenter        string
	Name              string
	Status            string
	StatusDescription string
	CreateIndex       uint64
	ModifyIndex       uint64
}

// NodeIndexSort reverse sorts nodes by CreateIndex
type NodeIndexSort []*NodeListStub

func (n NodeIndexSort) Len() int {
	return len(n)
}

func (n NodeIndexSort) Less(i, j int) bool {
	return n[i].CreateIndex > n[j].CreateIndex
}

func (n NodeIndexSort) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// nodeEvalResponse is used to decode a force-eval.
type nodeEvalResponse struct {
	EvalID string
}

// AllocationSort reverse sorts allocs by CreateIndex.
type AllocationSort []*Allocation

func (a AllocationSort) Len() int {
	return len(a)
}

func (a AllocationSort) Less(i, j int) bool {
	return a[i].CreateIndex > a[j].CreateIndex
}

func (a AllocationSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
