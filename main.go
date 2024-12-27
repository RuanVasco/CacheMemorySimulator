package main

import (
	"fmt"
	"math/rand"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

type Instruction struct {
	id   int
	rep  int
	time int64
}

type Cache struct {
	lines       []*Instruction
	cacheHits   int
	cacheMisses int
}

func NewCache(size int) *Cache {
	return &Cache{lines: make([]*Instruction, size)}
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
			id:  i,
			rep: rand.Intn(9) + 2,
		}
	}
	return instructions
}

func getRandomInstruction(instructions []*Instruction) *Instruction {
	i := rand.Intn(len(instructions))
	return instructions[i]
}

func containsInstruction(lines []*Instruction, instruction *Instruction) bool {
	for _, line := range lines {
		if line != nil && line.id == instruction.id {
			return true
		}
	}
	return false
}

func FIFO(c *Cache) int {
	leastId := int(^uint(0) >> 1)
	leastIdIndex := 0

	for i, line := range c.lines {
		if line != nil && line.id < leastId {
			leastId = line.id
			leastIdIndex = i
		}
	}

	return leastIdIndex
}

func LFU(c *Cache) int {
	leastFrequentIndex := -1
	leastFrequency := int64(^uint(0) >> 1)

	for i, line := range c.lines {
		if line != nil && line.time > 0 && line.time < leastFrequency {
			leastFrequency = line.time
			leastFrequentIndex = i
		}
	}

	return leastFrequentIndex
}

func LRU(c *Cache) int {
	index := -1
	oldestTimestamp := int64(^uint64(0) >> 1)

	for i, line := range c.lines {
		if line != nil && line.time < oldestTimestamp {
			oldestTimestamp = line.time
			index = i
		}
	}

	return index
}

func Random(c *Cache) int {
	return rand.Intn(len(c.lines))
}

func simulate(cacheSize int, instructionCount int, algorithm string) {
	instructions := generateInstructions(instructionCount)
	cache := NewCache(cacheSize)

	for {
		instructionsLeft := false

		for i, line := range cache.lines {
			if line == nil {
				cache.lines[i] = getRandomInstruction(instructions)
				instructionsLeft = true
			} else {
				if line.rep > 0 {
					cache.cacheHits++
					line.rep--
					line.time = time.Now().UnixNano()
					instructionsLeft = true
				} else {
					j := 0
					cache.cacheMisses++

					switch algorithm {
					case "FIFO":
						j = FIFO(cache)
					case "LFU":
						j = LFU(cache)
					case "LRU":
						j = LRU(cache)
					case "Random":
						j = Random(cache)
					}

					for _, instruction := range instructions {
						if instruction.rep > 0 && !containsInstruction(cache.lines, instruction) {
							cache.lines[j] = instruction
							instructionsLeft = true
							break
						}
					}

				}
			}
		}

		if !instructionsLeft {
			break
		}
	}

	fmt.Printf("Resultados para o algoritmo %s:\n", algorithm)
	fmt.Printf("Cache Hits: %d\n", cache.cacheHits)
	fmt.Printf("Cache Misses: %d\n", cache.cacheMisses)
	fmt.Printf("Eficiência: %.2f%%\n", cache.Efficiency())
}

func main() {
	cacheSize := 10
	instructionCount := 100
	algorithms := []string{"FIFO", "LFU", "LRU", "Random"}

	for _, algorithm := range algorithms {
		fmt.Println("--------------------------------")
		simulate(cacheSize, instructionCount, algorithm)
	}

	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()
}
