package main

import (
	"fmt"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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

func simulateStep(cache *Cache, instructions []*Instruction, algorithm string) bool {
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

	return instructionsLeft
}

func runAlgorithmStep(algorithm string, cache *Cache, instructions []*Instruction, statsArea, cacheArea, instructionsArea *fyne.Container, manualButton *widget.Button) {
	instructionsLeft := simulateStep(cache, instructions, algorithm)

	statsArea.Objects = []fyne.CanvasObject{
		widget.NewLabel(fmt.Sprintf("Algoritmo: %s", algorithm)),
		widget.NewLabel(fmt.Sprintf("Cache Hits: %d", cache.cacheHits)),
		widget.NewLabel(fmt.Sprintf("Cache Misses: %d", cache.cacheMisses)),
		widget.NewLabel(fmt.Sprintf("Eficiência: %.2f%%", cache.Efficiency())),
	}
	statsArea.Refresh()

	cacheArea.Objects = []fyne.CanvasObject{}
	for _, line := range cache.lines {
		if line != nil {
			cacheArea.Add(widget.NewLabel(fmt.Sprintf("ID: %d, REP: %d, Timestamp: %d", line.id, line.rep, line.time)))
		} else {
			cacheArea.Add(widget.NewLabel("Vazio"))
		}
	}
	cacheArea.Refresh()

	instructionsArea.Objects = []fyne.CanvasObject{}
	for _, instruction := range instructions {
		if instruction != nil {
			instructionsArea.Add(widget.NewLabel(fmt.Sprintf("ID: %d, REP: %d, Timestamp: %d", instruction.id, instruction.rep, instruction.time)))
		} else {
			instructionsArea.Add(widget.NewLabel("Vazio"))
		}
	}
	instructionsArea.Refresh()

	if !instructionsLeft {
		statsArea.Add(widget.NewLabel("Simulação concluída!"))
		statsArea.Refresh()
		manualButton.Disable()
	}
}

func main() {
	cacheSize := 10
	instructionCount := 100
	algorithms := []string{"FIFO", "LRU", "LFU", "Random"}

	app := app.New()
	window := app.NewWindow("Cache Simulation")
	window.Resize(fyne.NewSize(800, 400))

	statsArea := container.NewVBox()
	cacheArea := container.NewVBox()
	instructionsArea := container.NewVBox()

	var currentAlgorithm string
	var cache *Cache
	var instructions []*Instruction
	var manualButton *widget.Button

	buttonsArea := container.NewVBox()
	for _, algorithm := range algorithms {
		alg := algorithm
		btn := widget.NewButton(alg, func() {
			currentAlgorithm = alg
			cache = NewCache(cacheSize)
			manualButton.Enable()
			instructions = generateInstructions(instructionCount)
			statsArea.Objects = nil
			cacheArea.Objects = nil
			instructionsArea.Objects = nil
		})
		buttonsArea.Add(btn)
	}

	manualButton = widget.NewButton("Manual Clock", func() {
		if currentAlgorithm != "" && cache != nil && instructions != nil {
			runAlgorithmStep(currentAlgorithm, cache, instructions, statsArea, cacheArea, instructionsArea, manualButton)
		} else {
			statsArea.Objects = []fyne.CanvasObject{
				widget.NewLabel("Selecione um algoritmo primeiro!"),
			}
			statsArea.Refresh()
		}
	})

	mainContent := container.New(layout.NewHBoxLayout(),
		buttonsArea,
		container.NewVBox(
			manualButton,
			widget.NewLabel("Estatísticas:"),
			statsArea,
			container.NewHBox(
				container.NewVBox(
					widget.NewLabel("Cache:"),
					cacheArea,
				),
				container.NewVScroll(
					container.NewVBox(widget.NewLabel("Instruções:"),
						instructionsArea),
				),
			),
		),
	)

	window.SetContent(mainContent)
	window.ShowAndRun()
}
