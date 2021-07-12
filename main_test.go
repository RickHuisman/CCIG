package main

import (
	"os/exec"
	"syscall"
	"testing"
)

type TestCase struct {
	expr string
	want int
}

func getStatusCode(t *testing.T) int {
	cmd := exec.Command("./temp")

	if err := cmd.Start(); err != nil {
		t.Errorf("cmd.Start: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		} else {
			t.Errorf("cmd.Wait: %v", err)
		}
	}
	panic("TODO")
}

func runTestCases(t *testing.T, tests []TestCase) {
	for _, test := range tests {
		run(test.expr)
		got := getStatusCode(t)

		if got != test.want {
			t.Errorf("got = %d; want %d, input: %s", got, test.want, test.expr)
		}
	}
}

func TestAddition(t *testing.T) {
	var tests = []TestCase{
		{"return 100 + 50;", 150},
		{"return -5 + 10;", 5},
		{"return 5 + 0;", 5},
		{"return 5 + 2 + 3 - 5;", 5},
	}

	runTestCases(t, tests)
}

func TestSubtraction(t *testing.T) {
	var tests = []TestCase{
		{"return 100 - 50;", 50},
		{"return 5 - 0;", 5},
		{"return 5 - 2 - 3 + 5;", 5},
	}

	runTestCases(t, tests)
}

func TestMultiplication(t *testing.T) {
	var tests = []TestCase{
		{"return 10 * 5;", 50},
		{"return -5 * -5;", 25},
		{"return 5 * 0;", 0},
		{"return 5 + 2 * 3;", 11},
	}

	runTestCases(t, tests)
}

// TODO Test division

func TestLogicalOperators(t *testing.T) {
	var tests = []TestCase{
		{"return 5 == 5;", 1},
		{"return 5 != 5;", 0},
		{"return 10 > 5;", 1},
		{"return 10 < 5;", 0},
		{"return 10 <= 10;", 1},
		{"return 5 >= 5;", 1},
	}

	runTestCases(t, tests)
}

func TestLocal(t *testing.T) {
	var tests = []TestCase{
		{"var foo = 5; return foo + 10;", 15},
		{"var foo = 5; var bar = 10; return foo + bar;", 15},
	}

	runTestCases(t, tests)
}

func TestReturn(t *testing.T) {
	var tests = []TestCase{
		{"return 10;", 10},
	}

	runTestCases(t, tests)
}

func TestIfElse(t *testing.T) {
	var tests = []TestCase{
		{"if (10 == 10) { return 5; }", 5},
		{"if (10 != 10) { return 5; } else { return 3; } ", 3},
		{"var x = 5; if (5 < 10) { return 5; }", 5},
		// Test for multiple if statements
		{"if (5 < 10) { } if (10 > 5) { return 3; }", 3},
	}

	runTestCases(t, tests)
}
