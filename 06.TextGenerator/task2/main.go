package main

import (
	"06.TextGenerator/generator"
	"fmt"
	"os"
	"runtime/trace"
	"sync"
	"sync/atomic"
	"time"
)

const TextSize = 200_000
const MinWordLength = 7_000
const MaxWordLength = 8_000

const GoroutinesCount = 4
const Range = TextSize / GoroutinesCount

type Counter struct {
	counts map[string]*int32
	mutex  sync.RWMutex
}

func (c *Counter) Count(word string) {
	var count *int32

	if count = c.GetCount(word); count == nil {
		count = c.InitCount(word)
	}

	atomic.AddInt32(count, 1)
}

func (c *Counter) GetCount(word string) *int32 {
	defer c.mutex.RUnlock()
	c.mutex.RLock()
	return c.counts[word]
}

func (c *Counter) InitCount(word string) *int32 {
	count := c.GetCount(word)

	if count == nil {
		defer c.mutex.Unlock()
		c.mutex.Lock()
		value := int32(0)
		count = &value
		c.counts[word] = count
	}

	return count
}

func RunProcess(c *Counter, text []string, from, to int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := from; i <= to; i++ {
		c.Count(text[i])
	}
}

func (c *Counter) Print() {
	for char, count := range c.counts {
		fmt.Printf("%c -> %d \n", char, *count)
	}
}

func main() {
	f, err := os.Create("trace2.out")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	var wg sync.WaitGroup
	wg.Add(GoroutinesCount)

	words := generator.RandomStringArray(TextSize, MinWordLength, MaxWordLength)

	c := &Counter{counts: make(map[string]*int32)}

	start := time.Now()

	for i := 0; i < TextSize; i += Range {
		go RunProcess(c, words, i, i+(Range-1), &wg)
	}

	wg.Wait()

	fmt.Println(time.Now().Sub(start).Seconds())
	//c.Print()
}
