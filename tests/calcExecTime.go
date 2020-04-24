package main

import (
	"io/ioutil"
	"os"
	// "log"
	"regexp"
	"strconv"
	"time"
	"fmt"
)

func main() {
	pathToLogFile := os.Args[1]
	logContent, err := ioutil.ReadFile(pathToLogFile)
	if err != nil {
		panic(err)
	}

	r, _ := regexp.Compile("Testing Timestamp: ([0-9]+)")
	matches := r.FindAllStringSubmatch(string(logContent), -1)
	minT, maxT := time.Now(), *new(time.Time)

	// method 1: loop through every match
	// for _, match := range matches {
	// 	t := unixStampToTimeObj(match[1])
	// 	if t.After(maxT) {
	// 		maxT = t
	// 	}
	// 	if t.Before(minT) {
	// 		minT = t
	// 	}
	// }

	// method 2: just look at first and last
	minT = unixStampToTimeObj(matches[0][1])
	maxT = unixStampToTimeObj(matches[len(matches) - 1][1])

	fmt.Println("max time is: ", maxT)
	fmt.Println("min time is: ", minT)
	elapsedSeconds := maxT.Sub(minT).Seconds()
	appendToFile("execTime.txt", fmt.Sprintf("%f seconds\n\n", elapsedSeconds))
}

func unixStampToTimeObj(stamp string) time.Time {
	i, err := strconv.ParseInt(stamp, 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(i, 0)
}

func appendToFile(filename string, data string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
	    panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(data); err != nil {
	    panic(err)
	}
}