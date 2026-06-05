package main

import (
	"fmt"
	"os/exec"
	"sync"
	"math"
	"math/rand"
)

const (
	numWorkers = 2
	numParticles = 8
	particleRadius = 1.
	influenceRadius = 1.5
	tStep = .02
	screenW = 200
	screenH = 200
)

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
	for i:=0; i<numParticles; i++ {
		gParticles[i].pos = randPoint(screenW, screenH)
		gParticles[i].temp = rand.Float64() * 400
		gParticles[i].size = particleRadius
	}
	return gParticles
}

func (p1 point) addmult(p2 point, a float64) point{
	newPoint := new(point)
	newPoint.X = p1.X + a*p2.X
	newPoint.Y = p1.Y + a*p2.Y
	return *newPoint
}

func (this particle) applyForce() {
	this.vel.addmult(this.F.addmult(randPoint(this.temp, this.temp), .001), tStep)
	this.pos.addmult(this.vel, tStep)
}

func (this particle) interact(other particle) particle {
	distTo := math.Sqrt(math.Pow(this.pos.X - other.pos.X, 2) + math.Pow(this.pos.Y - other.pos.Y, 2))
	if distTo <= this.size {
		this.F.addmult(this.pos.addmult(other.pos, -1./(math.Pow(distTo,2.))), 1.)
	} else if distTo <= influenceRadius {
		this.F.addmult(this.pos.addmult(other.pos, -1./(math.Pow(distTo,4.))), -1.)
	}
	return this
}

// One worker's computation
func computeForces(id int, nPerWorker int, pList *group) {
	start := id*nPerWorker
	end := min(start+nPerWorker, numParticles)
	for i := start; i < end; i++ {
		for j := 0; j < numParticles; j++ {
			pList.groupParticles[i] = pList.groupParticles[i].interact(pList.groupParticles[j])
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
	fmt.Println(g.groupParticles)
	step(g)
	step(g)
	step(g)
	fmt.Println(g.groupParticles)
}

