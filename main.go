package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/kelseyhightower/envconfig"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type config struct {
	Seed int64

	Agents int `default:"100"`
	Money  int `default:"10"`
	Steps  int
}

type agent struct {
	ID    string
	Money int
}

type agents []*agent

func (as agents) Sum() (total int) {
	for _, a := range as {
		total += a.Money
	}

	return total
}

func (as agents) Values() plotter.Values {
	vs := make(plotter.Values, len(as))

	for i, a := range as {
		vs[i] = float64(a.Money)
	}

	return vs
}

func (as agents) Select() *agent {
	return as[rand.Intn(len(as))]
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// Step runs one simulation step.
func (as agents) Step() {

	avg := float64(as.Sum()) / float64(len(as))

	//upper := int(10 * avg)
	upper := int(2 * avg)
	//upper := int(1.75 * avg)
	//upper := int(1.5 * avg)
	//upper := int(1.25 * avg)

	//lower := 0
	lower := int(avg) - (upper - int(avg))

	for _, a := range as {
		for {
			p := as.Select()
			if a == p {
				continue
			}

			sign := rand.Intn(2)
			if sign == 0 {
				sign = -1
			}

			delta := rand.Intn(min(a.Money, p.Money) + 1)
			//delta := 1
			//delta := rand.Intn(3) + 1

			am := a.Money + sign*delta
			pm := p.Money - sign*delta

			/*
				spew.Dump(map[string]int{
					"delta": delta,
					"am":    am,
					"pm":    pm,
				})
			*/

			if am >= lower && pm >= lower && am < upper && pm < upper {
				a.Money = am
				p.Money = pm
			}

			break
		}

	}
}

func (as agents) Plot(name string) {
	p := plot.New()

	p.Title.Text = "Histogram"
	p.X.Min = 0
	p.X.Max = 100
	p.Y.Min = 0
	p.Y.Max = 100

	vs := as.Values()

	h, err := plotter.NewHist(vs, 100)
	if err != nil {
		panic(err)
	}

	p.Add(h)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, name); err != nil {
		panic(err)
	}
}

func main() {
	// Parse config.
	c := config{}
	err := envconfig.Process("ephy", &c)
	if err != nil {
		panic(err)
	}

	// Initialize rand with seed.
	if c.Seed == 0 {
		c.Seed = time.Now().UnixNano()
	}

	rand.Seed(c.Seed)

	spew.Dump(c)

	// Create all agents with starting dollars.
	as := agents{}
	for i := 0; i < c.Agents; i++ {
		as = append(as, &agent{
			ID:    strconv.Itoa(i),
			Money: c.Money,
		})
	}

	// Run simulation step.
	digits := int(math.Log10(float64(c.Steps)))
	fmt.Println(digits)
	as.Plot(fmt.Sprintf("hist.%0.*d.svg", digits, 0))

	for i := 1; i < c.Steps; i++ {
		as.Step()
		as.Plot(fmt.Sprintf("hist.%0.*d.svg", digits, i))
	}

	as.Plot("hist.svg")
}
