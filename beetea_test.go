package beetea

import (
	"fmt"
	"testing"
	"time"
)

func TestActionNode(t *testing.T) {
	actionNode := NewAction("action1", func() Status {
		return Success
	})

	if actionNode.Tick() != Success {
		t.Errorf("Expected ActionNode to return Success")
	}

	initialVersion := actionNode.GetVersion()
	actionNode.UpdateVersion()
	if actionNode.GetVersion() == initialVersion {
		t.Errorf("Expected version to be incremented")
	}

	if actionNode.CalculateHash() == "" {
		t.Errorf("Expected non-empty hash for ActionNode")
	}
}

func TestConditionNode(t *testing.T) {
	conditionNode := NewCondition("condition1", func() bool {
		return true
	})

	if conditionNode.Tick() != Success {
		t.Errorf("Expected ConditionNode to return Success on true condition")
	}

	if conditionNode.CalculateHash() == "" {
		t.Errorf("Expected non-empty hash for ConditionNode")
	}
}

func TestSelectorNode(t *testing.T) {
	successNode := NewAction("action1", func() Status {
		return Success
	})
	failureNode := NewAction("action1", func() Status {
		return Failure
	})
	selector := NewSelector("selctor1", failureNode, successNode)

	if selector.Tick() != Success {
		t.Errorf("Expected Selector to succeed when one child succeeds")
	}

	if selector.CalculateHash() == "" {
		t.Errorf("Expected non-empty hash for Selector")
	}
}

func TestSequenceNode(t *testing.T) {
	successNode := NewAction("action1", func() Status {
		return Success
	})
	failureNode := NewAction("action1", func() Status {
		return Failure
	})
	sequence := NewSequence("sequence1", successNode, failureNode)

	if sequence.Tick() != Failure {
		t.Errorf("Expected Sequence to fail when one child fails")
	}

	if sequence.CalculateHash() == "" {
		t.Errorf("Expected non-empty hash for Sequence")
	}
}

func TestHashRecalculationOnModification(t *testing.T) {
	actionNode := NewAction("action2", func() Status {
		return Running
	})
	initialHash := actionNode.CalculateHash()

	// Simulate modification by changing the action
	actionNode.Action = func() Status {
		return Success
	}
	actionNode.UpdateVersion()

	if actionNode.CalculateHash() == initialHash {
		t.Errorf("Expected hash to change after modification")
	}
}

func TestCompositeNodeWithVariousChildStatuses(t *testing.T) {
	alwaysRunning := NewAction("alwaysRunning", func() Status {
		return Running
	})
	alwaysSuccess := NewAction("alwaysSuccess", func() Status {
		return Success
	})
	alwaysFailure := NewAction("alwaysFailure", func() Status {
		return Failure
	})

	// Selector should succeed if any child succeeds
	selector := NewSelector("selector2", alwaysFailure, alwaysRunning, alwaysSuccess)
	if selector.Tick() != Running {
		t.Errorf("Expected Selector to return Running when a child is running")
	}

	// Sequence should fail if any child fails
	sequence := NewSequence("sequence3", alwaysSuccess, alwaysRunning, alwaysFailure)
	if sequence.Tick() != Running {
		t.Errorf("Expected Sequence to return Running when a child is running before encountering a failure")
	}
}

// This test will print the time iff the current minute is even
// Run it a few times. If the minute is even, the time will be printed.
func TestMain(m *testing.T) {
	// 1. Define the action function for even minutes
	logTimeAction := func() Status {
		fmt.Println("Current Time:", time.Now())
		return Success
	}

	// 2. Define the action function for odd minutes
	oddAction := func() Status {
		fmt.Println("Current time was not even")
		return Success
	}

	// 3. Define a condition function
	isEvenMinuteCondition := func() bool {
		return time.Now().Minute()%2 == 0
	}

	// 4. Create action nodes
	evenActionNode := NewAction("logTimeEven", logTimeAction)
	oddActionNode := NewAction("logTimeOdd", oddAction)

	// 5. Create a condition node
	conditionNode := NewCondition("isEvenMinute", isEvenMinuteCondition)

	// 6. Adjust the tree structure
	// This is where we need to be creative with the existing API
	// For the sake of simplicity in this example, let's directly use the Selector for demonstration
	selector := NewSelector("timeBasedActionSelector", conditionNode, oddActionNode, evenActionNode)

	// 7. Execute the tree
	status := selector.Tick()
	fmt.Println("Tree execution status:", status)

	// Optionally, inspect the tree's structure or the status of individual nodes
	fmt.Println("Even Action Node Hash:", evenActionNode.CalculateHash())
	fmt.Println("Odd Action Node Hash:", oddActionNode.CalculateHash())
	fmt.Println("Condition Node Hash:", conditionNode.CalculateHash())
	fmt.Println("Selector Node Hash:", selector.CalculateHash())
}
