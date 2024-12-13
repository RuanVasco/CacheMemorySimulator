package main

import (
	"fmt"
	"math/rand"
	"time"
)

const instructionListSize = 100
const cacheSize = 10

type Instruction struct {
	id     int
	rep    int
	status bool
}

type CacheLine struct {
	instruction *Instruction
}

func roundRobin(cache *[cacheSize]CacheLine, quantum int) {
	completed := 0

	for completed < len(*cache) {
		for j := 0; j < len(*cache); j++ {
			if cache[j].instruction != nil {
				instruction := cache[j].instruction

				if instruction.rep > 0 && instruction.status {
					fmt.Printf("Executing process %d for %d time units, next rp %d\n", instruction.id, quantum, instruction.rep)

					if instruction.rep <= quantum {
						fmt.Printf("Process %d finished execution\n", instruction.id)
						instruction.status = false
						completed++
					} else {
						instruction.rep -= quantum
					}
				}
			}
		}
	}

	time.Sleep(1 * time.Second)
}

func generateProcesses(list *[instructionListSize]Instruction) {
	for i := 0; i < len(list); i++ {
		list[i].id = i
		list[i].rep = rand.Intn(9) + 2
		list[i].status = true
	}
}

func generateCache(cache *[cacheSize]CacheLine, instructions *[instructionListSize]Instruction) {
	for i := 0; i < len(*cache); i++ {
		for {
			instruction := &(*instructions)[rand.Intn(instructionListSize)]

			if instruction.status {
				cache[i].instruction = instruction
				break
			}
		}
	}
}

func main() {
	var instructions [instructionListSize]Instruction
	var cache [cacheSize]CacheLine
	generateProcesses(&instructions)
	generateCache(&cache, &instructions)

	quantum := 3

	roundRobin(&cache, quantum)
}
