package beetea

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Status represents the status of a node's execution.
type Status int

const (
	Running Status = iota
	Success
	Warning
	Failure
)

// Node is the basic interface for all nodes in the behavior tree.
type Node interface {
	Tick() Status
	CalculateHash() string
	GetVersion() int
	UpdateVersion()
}

// BaseNode includes fields and methods common to all nodes, such as hash and version.
type BaseNode struct {
	hash    string
	version int
}

func (b *BaseNode) CalculateHash() string {
	// Implementation will depend on the specific node type
	return ""
}

func (b *BaseNode) GetVersion() int {
	return b.version
}

// When the version is updated, the hash should be recalculated.
func (b *BaseNode) UpdateVersion() {
	b.version++
	b.hash = b.CalculateHash()
}

// ActionNode represents a leaf node that performs an action.
type ActionNode struct {
	BaseNode
	Action func() Status
	ID     string
}

func (a *ActionNode) Tick() Status {
	return a.Action()
}

func (a *ActionNode) CalculateHash() string {
	data := fmt.Sprintf("ActionNode:%s:Version:%d", a.ID, a.version)
	hash := sha256.Sum256([]byte(data))
	a.hash = hex.EncodeToString(hash[:])
	return a.hash
}

// ConditionNode represents a leaf node that checks a condition
type ConditionNode struct {
	BaseNode
	ID        string
	Condition func() bool
}

func (c *ConditionNode) Tick() Status {
	if c.Condition() {
		return Success
	}
	return Failure
}

func (a *ConditionNode) CalculateHash() string {
	data := fmt.Sprintf("ConditionNode:%s:Version:%d", a.ID, a.version)
	hash := sha256.Sum256([]byte(data))
	a.hash = hex.EncodeToString(hash[:])
	return a.hash
}

// CompositeNode is a base struct for nodes that have children.
type CompositeNode struct {
	BaseNode
	Children []Node
	ID       string
}

// Selector executes its children until one of them succeeds.
type Selector struct {
	CompositeNode
}

func (s *Selector) Tick() Status {
	for _, child := range s.Children {
		status := child.Tick()
		if status != Failure {
			return status
		}
	}
	return Failure
}

func (a *Selector) CalculateHash() string {
	data := fmt.Sprintf("SelectorNode:%s:Version:%d", a.ID, a.version)
	hash := sha256.Sum256([]byte(data))
	a.hash = hex.EncodeToString(hash[:])
	return a.hash
}

// Sequence executes its children in order, succeeding if all succeed.
type Sequence struct {
	CompositeNode
}

func (seq *Sequence) Tick() Status {
	for _, child := range seq.Children {
		status := child.Tick()
		if status != Success {
			return status
		}
	}
	return Success
}

func (seq *Sequence) CalculateHash() string {
	data := fmt.Sprintf("SequenceNode:%s:Version:%d", seq.ID, seq.version)

	for _, child := range seq.Children {
		data += ":" + child.CalculateHash()
	}

	hash := sha256.Sum256([]byte(data))
	seq.hash = hex.EncodeToString(hash[:])
	return seq.hash
}

type TreeBuilder interface {
	AddTask(taskID string, actionFunc func() Status, dependencies []string)
	Build() Node
}

func NewSelector(id string, children ...Node) *Selector {
	node := &Selector{CompositeNode{Children: children, ID: id}}
	return node
}

func NewSequence(id string, children ...Node) *Sequence {
	node := &Sequence{CompositeNode{Children: children, ID: id}}
	return node
}

func NewAction(id string, action func() Status) *ActionNode {
	node := &ActionNode{
		Action: action,
		ID:     id,
	}
	node.UpdateVersion()
	return node
}

func NewCondition(id string, condition func() bool) *ConditionNode {
	node := &ConditionNode{
		Condition: condition,
		ID:        id,
	}
	node.UpdateVersion()
	return node
}
