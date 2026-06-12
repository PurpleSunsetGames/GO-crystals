package main

import (
	"fmt"
	"os"
	"log"
	"os/exec"
	"encoding/csv"
	"sync"
	"math"
	"math/rand"
	"strconv"
)

const (
	numWorkers = 8
	numParticles = 625
	particleRadius = 3.
	tStep = .0001
	screenW = 50
	screenH = 50
	stepsPerFrame = 100
)
var tempRangeMin = 270.
var tempRangeMax = 400.
var epsilon = 3.

type group struct {
	groupParticles []particle
}

type particle struct {
	pos point
	vel point
	F point
	temp float64
	size float64
}

type point struct {
	X float64
	Y float64
}

// Methods
func startGroup() []particle {
	gParticles := make([]particle, numParticles)

	for i:=0; i<int(math.Sqrt(numParticles)); i++ {
		for j:=0; j<int(math.Sqrt(numParticles)); j++ {
			newPoint := new(point)
			newPoint.X = float64(i)*screenW/math.Sqrt(numParticles) - screenW/2
			newPoint.Y = float64(j)*screenH/math.Sqrt(numParticles) - screenH/2
			gParticles[j+i*int(math.Sqrt(numParticles))].pos = newPoint.addmult(randPoint(1,1),.5)
			gParticles[j+i*int(math.Sqrt(numParticles))].temp = rand.Float64()*(tempRangeMax - tempRangeMin) + tempRangeMin
			gParticles[j+i*int(math.Sqrt(numParticles))].size = particleRadius
		}
	}
	return gParticles
}

func (p1 point) addmult(p2 point, a float64) point{
	newPoint := new(point)
	newPoint.X = p1.X + a*p2.X
	newPoint.Y = p1.Y + a*p2.Y
	return *newPoint
}

func (p point) mag() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (this *particle) applyForce() {
	if math.Abs(this.pos.X) > screenW/2 {
		this.F.X = this.F.X - 1000*(this.pos.X - (screenW/2)*this.pos.X/math.Abs(this.pos.X))
	}
	if math.Abs(this.pos.Y) > screenH/2 {
		this.F.Y = this.F.Y - 1000*(this.pos.Y - (screenH/2)*this.pos.Y/math.Abs(this.pos.Y))
	}
	this.vel = this.vel.addmult(this.F.addmult(randPoint(1,1), this.temp), tStep)
	this.pos = this.pos.addmult(this.vel, tStep)
	newPoint := new(point)
	this.F = *newPoint
}

// Leonard-Jones as a place holder
func (this *particle) interact(other particle) {
	distTo := math.Sqrt(math.Pow(this.pos.X - other.pos.X, 2) + math.Pow(this.pos.Y - other.pos.Y, 2))
	vLJ := 2*epsilon*(math.Pow(particleRadius/distTo, 12) - math.Pow(particleRadius/distTo, 6))
	this.F = this.F.addmult(this.pos.addmult(other.pos, -1), vLJ)
	other.F = other.F.addmult(other.pos.addmult(this.pos, -1), vLJ)
}

// One worker's computation
func computeForces(id int, nPerWorker int, pList *group) {
	start := id*nPerWorker
	end := min(start+nPerWorker, numParticles)
	for i := start; i < end; i++ {
		for j := 0; j < numParticles; j++ {
			if (j != i) {
				pList.groupParticles[i].interact(pList.groupParticles[j])
			}
		}
		pList.groupParticles[i].applyForce()
	}
}
// Helper func
func randPoint(sizeX float64, sizeY float64) point {
	newPoint := new(point)
	newPoint.X = (rand.Float64() - .5) * sizeX
	newPoint.Y = (rand.Float64() - .5) * sizeY
	return *newPoint
}

// One cycle of the sim
func step(ps *group) {
	n := numParticles / numWorkers
	var wg sync.WaitGroup
	for i := 0; i<numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			computeForces(i, n, ps)
		}()
	} 
	wg.Wait()
}

func renderVid() {
	// FFmpeg command to turn sequential PNGs into a video
	// -i: input pattern (e.g., image-1.png, image-2.png)
	// -c:v libx264: use H.264 video codec
	// -pix_fmt yuv420p: ensures browser and media player compatibility
	cmd := exec.Command("ffmpeg", 
		"-framerate", "30", 
		"-i", "image-%d.png", 
		"-c:v", "libx264", 
		"-pix_fmt", "yuv420p", 
		"output.mp4",
	)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error converting PNGs to video:", err)
		return
	}
	fmt.Println("Video successfully created: output.mp4")
}

func main() {
	g := new(group)
	g.groupParticles = startGroup()
	numSteps := 100
	pdata := make([][]string, numSteps)
	for i:=0; i<numSteps; i++ {
		pdata[i] = make([]string, 2*numParticles)
		for j:=0; j<numParticles; j++ {
			pdata[i][j*2] = strconv.FormatFloat(g.groupParticles[j].pos.X, 'f', 6, 64)
			pdata[i][j*2 + 1] = strconv.FormatFloat(g.groupParticles[j].pos.Y, 'f', 6, 64)
		}
		for k:=0; k<stepsPerFrame; k++{
			step(g)
		}
		if tempRangeMax > 40 {
			tempRangeMax -= 2
		}
		if tempRangeMin > 20 {
			tempRangeMin -= 2
		}	
	}
	file, err := os.Create("output.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.WriteAll(pdata)
}

