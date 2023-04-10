package main

import (
	"06.TextGenerator/generator"
	"fmt"
	"time"
)

const TextSize = 200_000
const MinWordLength = 7_000
const MaxWordLength = 8_000

type Counter struct {
	counts map[string]int
}

func (c *Counter) Count(word string) {
	c.counts[word]++
}

func (c *Counter) Print() {
	for word, count := range c.counts {
		fmt.Printf("%c -> %d \n", word, count)
	}
}

func main() {
	words := generator.RandomStringArray(TextSize, MinWordLength, MaxWordLength)

	c := &Counter{counts: make(map[string]int)}

	start := time.Now()
	for _, word := range words {
		c.Count(word)
	}
	fmt.Println(time.Now().Sub(start).Seconds())
	//c.Print()
}
