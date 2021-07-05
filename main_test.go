package main

import (
    "os/exec"
    "syscall"
    "testing"
)

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

func TestAddition(t *testing.T) {
    var tests = []struct {
        expr string
        want int
    }{
        {"100 + 50;", 150},
        {"-5 + 10;", 5},
        {"5 + 0;", 5},
        {"5 + 2 + 3 - 5;", 5},
    }

    for _, test := range tests {
        run(test.expr)
        got := getStatusCode()

        if got != test.want {
            t.Errorf("got = %d; want %d", got, test.want)
        }
    }
}

func TestSubtraction(t *testing.T) {
    var tests = []struct {
        expr string
        want int
    }{
        {"100 - 50;", 50},
        {"5 - 0;", 5},
        {"5 - 2 - 3 + 5;", 5},
    }

    for _, test := range tests {
        run(test.expr)
        got := getStatusCode()

        if got != test.want {
            t.Errorf("got = %d; want %d", got, test.want)
        }
    }
}

func TestMultiplication(t *testing.T) {
    var tests = []struct {
        expr string
        want int
    }{
        {"100 * 50;", 5000},
        {"-5 * -5;", 25},
        {"5 * 0;", 0},
        {"5 + 2 * 3;", 11},
    }

    for _, test := range tests {
        run(test.expr)
        got := getStatusCode()

        if got != test.want {
            t.Errorf("got = %d; want %d", got, test.want)
        }
    }
}
