package beetea

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
}

// ActionNode represents a leaf node that performs an action.
type ActionNode struct {
    Action func() Status
}

func (a *ActionNode) Tick() Status {
    return a.Action()
}

// ConditionNode represents a leaf node that checks a condition.
type ConditionNode struct {
    Condition func() bool
}

func (c *ConditionNode) Tick() Status {
    if c.Condition() {
        return Success
    }
    return Failure
}

// CompositeNode is a base struct for nodes that have children.
type CompositeNode struct {
    Children []Node
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

type TreeBuilder interface {
    AddTask(taskID string, actionFunc func() Status, dependencies []string)
    Build() Node
}

func NewSelector(children ...Node) *Selector {
    return &Selector{CompositeNode{Children: children}}
}

func NewSequence(children ...Node) *Sequence {
    return &Sequence{CompositeNode{Children: children}}
}

func NewAction(action func() Status) *ActionNode {
    return &ActionNode{Action: action}
}

func NewCondition(condition func() bool) *ConditionNode {
    return &ConditionNode{Condition: condition}
}
