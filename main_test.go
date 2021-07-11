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

func getStatusCode() int {
    cmd := exec.Command("./temp")

    if err := cmd.Start(); err != nil {
        // t.Errorf("Abs(-1) = %d; want 1", err)
        // log.Fatalf("cmd.Start: %v", err)
    }

    if err := cmd.Wait(); err != nil {
        if exiterr, ok := err.(*exec.ExitError); ok {
            // The program has exited with an exit code != 0

            if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
                return status.ExitStatus()
                // t.Errorf("Exit Status: %d", status.ExitStatus())
            }
        } else {
            // t.Errorf("cmd.Wait: %v", err)
        }
    }
    // t.Errorf("cmd.Wait: %v", err)
    return 0 // TODO
}

func runTestCases(t *testing.T,tests []TestCase) {
    for _, test := range tests {
        run(test.expr)
        got := getStatusCode()

        if got != test.want {
            t.Errorf("got = %d; want %d, input: %s", got, test.want, test.expr)
        }
    }
}

func TestAddition(t *testing.T) {
    var tests = []TestCase {
        {"return 100 + 50;", 150},
        {"return -5 + 10;", 5},
        {"return 5 + 0;", 5},
        {"return 5 + 2 + 3 - 5;", 5},
    }

    runTestCases(t, tests)
}

func TestSubtraction(t *testing.T) {
    var tests = []TestCase {
        {"return 100 - 50;", 50},
        {"return 5 - 0;", 5},
        {"return 5 - 2 - 3 + 5;", 5},
    }

    runTestCases(t, tests)
}

func TestMultiplication(t *testing.T) {
    var tests = []TestCase {
        {"return 10 * 5;", 50},
        {"return -5 * -5;", 25},
        {"return 5 * 0;", 0},
        {"return 5 + 2 * 3;", 11},
    }

    runTestCases(t, tests)
}

func TestLocal(t *testing.T) {
    var tests = []TestCase {
        {"var foo = 5; foo + 10;", 15},
        {"var foo = 5; var bar = 10; foo + bar;", 15},
    }

    runTestCases(t, tests)
}

func TestReturn(t *testing.T) {
    var tests = []TestCase {
        {"return 10;", 10},
    }

    runTestCases(t, tests)
}

//func TestIfElse(t *testing.T) {
//    var tests = []struct {
//        expr string
//        want int
//    }{
//        {"if (10 == 10) { return }", 15},
//    }
//
//    for _, test := range tests {
//        run(test.expr)
//        got := getStatusCode()
//
//        if got != test.want {
//            t.Errorf("got = %d; want %d", got, test.want)
//        }
//    }
//}