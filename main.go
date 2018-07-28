package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
)

const MaxNumber = 45
const AmountNumbers = 6
const SupplementaryNumbers = 2

type ISSResponse struct {
	Message     string `json:"message"`
	Timestamp   int64  `json:"timestamp"`
	ISSPosition struct {
		Lat  string `json:"latitude"`
		Long string `json:"longitude"`
	} `json:"iss_position"`
}

func main() {

	numbers := make([]int, AmountNumbers)
	supps := make([]int, SupplementaryNumbers)

	i := 0
	for i < len(numbers) {
		num := contactISSForNumber()
		if contains(numbers, num) {
			continue
		}
		numbers[i] = num
		i++
	}

	i = 0
	for i < len(supps) {
		num := contactISSForNumber()
		if contains(supps, num) || contains(numbers, num) {
			continue
		}
		supps[i] = num
		i++
	}

	fmt.Println(numbers, supps)

}

func contains(sl []int, num int) bool {
	for _, a := range sl {
		if a == num {
			return true
		}
	}
	return false
}

func contactISSForNumber() int {
	resp, err := http.Get("http://api.open-notify.org/iss-now.json")
	must(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	must(err)

	var issresult ISSResponse
	err = json.Unmarshal(body, &issresult)
	must(err)

	rand.Seed(issresult.Timestamp)

	allNumbers := stripCharacters(issresult.ISSPosition.Lat + issresult.ISSPosition.Lat)

	return getNumber(allNumbers)
}

func stripNullCharacter(b []byte) string {
	return string(bytes.Replace(b, []byte("\x00"), []byte{}, -1))
}

func stripCharacters(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	must(err)

	return reg.ReplaceAllString(s, "")
}

func getNumber(s string) int {
	to := 2
	startsat := 1 + rand.Intn(len(s)-to)

	runes := []rune(s)
	r := runes[startsat : startsat+to]

	ns := string(r)

	if string(r[0:1]) == "0" || string(r) == "" {
		return getNumber(s)
	}

	n, err := strconv.Atoi(stripNullCharacter([]byte(ns)))
	must(err)

	if n > MaxNumber {
		n, err = strconv.Atoi(string(r[0:1]))
		must(err)
	}

	return n

}

func must(err error) {
	if err != nil {
		panic(err) // if one thing doesn't work, this whole thing won't work
	}
}
