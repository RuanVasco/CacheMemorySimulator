package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Instruction struct {
	id       int
	rep      int
	lastTime time.Time
}

type Cache struct {
	lines       []*Instruction
	cacheHits   int
	cacheMisses int
}

func NewCache(size int) *Cache {
	return &Cache{lines: make([]*Instruction, size)}
}

func (c *Cache) Access(instruction *Instruction, algorithm string) {
	for _, line := range c.lines {
		if line != nil && line.id == instruction.id {
			c.cacheHits++
			instruction.rep--
			instruction.lastTime = time.Now()
			return
		}
	}

	c.cacheMisses++
	emptyIndex := -1
	for i, line := range c.lines {
		if line == nil {
			emptyIndex = i
			break
		}
	}

	if emptyIndex != -1 {
		c.lines[emptyIndex] = instruction
		instruction.rep--
		instruction.lastTime = time.Now()
	} else {
		switch algorithm {
		case "FIFO":
			c.FIFO(instruction)
		case "LRU":
			c.LRU(instruction)
		case "Random":
			c.Random(instruction)
		case "LFU":
			c.LFU(instruction)
		default:
			fmt.Println("Algoritmo desconhecido:", algorithm)
		}
	}
}

func (c *Cache) FIFO(instruction *Instruction) {
	c.lines = append(c.lines[1:], instruction)
	instruction.rep--
	instruction.lastTime = time.Now()
}

func (c *Cache) LRU(instruction *Instruction) {
	oldestIndex := 0
	oldestTime := c.lines[0].lastTime

	for i, line := range c.lines {
		if line.lastTime.Before(oldestTime) {
			oldestIndex = i
			oldestTime = line.lastTime
		}
	}

	c.lines[oldestIndex] = instruction
	instruction.rep--
	instruction.lastTime = time.Now()
}

func (c *Cache) Random(instruction *Instruction) {
	randomIndex := rand.Intn(len(c.lines))
	c.lines[randomIndex] = instruction
	instruction.rep--
	instruction.lastTime = time.Now()
}

func (c *Cache) LFU(instruction *Instruction) {
	leastFrequentIndex := 0
	leastRep := c.lines[0].rep

	for i, line := range c.lines {
		if line.rep < leastRep {
			leastFrequentIndex = i
			leastRep = line.rep
		}
	}

	c.lines[leastFrequentIndex] = instruction
	instruction.rep--
	instruction.lastTime = time.Now()
}

func (c *Cache) Efficiency() float64 {
	totalAccesses := c.cacheHits + c.cacheMisses
	if totalAccesses == 0 {
		return 0
	}
	return float64(c.cacheHits) / float64(totalAccesses) * 100
}

func generateInstructions(size int) []*Instruction {
	instructions := make([]*Instruction, size)
	for i := 0; i < size; i++ {
		instructions[i] = &Instruction{
			id:       i,
			rep:      rand.Intn(9) + 2,
			lastTime: time.Now(),
		}
	}
	return instructions
}

func simulate(cacheSize, instructionCount int, algorithm string) {
	rand.Seed(time.Now().UnixNano())

	instructions := generateInstructions(instructionCount)
	cache := NewCache(cacheSize)

	for _, instruction := range instructions {
		cache.Access(instruction, algorithm)
	}

	fmt.Printf("Resultados para o algoritmo %s:\n", algorithm)
	fmt.Printf("Cache Hits: %d\n", cache.cacheHits)
	fmt.Printf("Cache Misses: %d\n", cache.cacheMisses)
	fmt.Printf("Eficiência: %.2f%%\n", cache.Efficiency())
}

func main() {
	cacheSize := 10
	instructionCount := 100
	algorithms := []string{"FIFO", "LRU", "Random", "LFU"}

	for _, algorithm := range algorithms {
		fmt.Println("--------------------------------")
		simulate(cacheSize, instructionCount, algorithm)
	}
}
