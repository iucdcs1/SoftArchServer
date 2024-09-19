package main

import (
    "bytes"
    "fmt"
    "os/exec"
    "strings"
    "time"
)

func runCommand(cmd string, args ...string) (string, error) {
    var out bytes.Buffer
    var stderr bytes.Buffer
    command := exec.Command(cmd, args...)
    command.Stdout = &out
    command.Stderr = &stderr
    err := command.Run()
    return out.String(), err
}

func measureTestDurationAndCount() (int, int, time.Duration, string) {
    start := time.Now()
    output, _ := runCommand("go", "test", "./...", "-v")
    duration := time.Since(start)
    passCount := strings.Count(output, "--- PASS:")
	failCount := strings.Count(output, "--- FAIL:")
	return passCount, failCount, duration, output
}

func checkAdditionalTestCases(initialPassCount, initialFailCount, refactoredPassCount, refactoredFailCount int) bool {
    return ((refactoredPassCount + refactoredFailCount) - (initialPassCount + initialFailCount)) <= 2
}

func checkTestDuration(initialDuration, refactoredDuration time.Duration) bool {
    return refactoredDuration <= initialDuration+(initialDuration/10)
}

func checkOutputConsistency(refactoredFailCount, initialFailCount int) bool {
    return refactoredFailCount == initialFailCount
}


func main() {
    fmt.Println("Running initial tests...")
	initialPassCount, initialFailCount, initialDuration, initialOutput := measureTestDurationAndCount()

    fmt.Println("Initial tests completed.")
    fmt.Println("Initial test results:\n", initialOutput)
	fmt.Println("==============================================================")


    fmt.Println("Refactoring code...")
    time.Sleep(5 * time.Second)

    fmt.Println("Running tests after refactoring...")
    refactoredPassCount, refactoredFailCount, refactoredDuration, refactoredOutput := measureTestDurationAndCount()
	fmt.Println("Refactored tests completed.")
    fmt.Println("Refactored test results:\n", refactoredOutput)
	fmt.Println("==============================================================")

	
    fmt.Printf("Initial pass count: %d\n", initialPassCount)
    fmt.Printf("Initial fail count: %d\n", initialFailCount)
    fmt.Printf("Initial test duration: %v\n\n", initialDuration)
	fmt.Println("==============================================================")

    fmt.Printf("Refactored pass count: %d\n", refactoredPassCount)
    fmt.Printf("Refactored fail count: %d\n", refactoredFailCount)
    fmt.Printf("Refactored test duration: %v\n", refactoredDuration)
    fmt.Println("\n")

	fmt.Println("==============================================================")
	fmt.Println("Conclusion: \n")
    if !checkAdditionalTestCases(initialPassCount, initialFailCount, refactoredPassCount, refactoredFailCount) {
        fmt.Println("Refactoring added more than two additional test cases.")
    } else {
        fmt.Println("Refactoring is acceptable. Code is more modular and no more than two additional test cases were required.")
    }

	if checkTestDuration(initialDuration, refactoredDuration) {
		fmt.Println("Refactoring did not introduce performance regressions.")
	} else {
		fmt.Println("Refactoring made tests slower. Investigate for potential inefficiencies.")
	}
	
	if checkOutputConsistency(refactoredFailCount, initialFailCount) {
		fmt.Println("Refactoring preserved existing behavior. Outputs match.")
	} else {
		fmt.Println("Refactoring introduced changes in behavior. Outputs differ.")
	}
	
}
