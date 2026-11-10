package atlas

import (
	"hash/fnv"
	"math"
	"sort"
)

// Cosmos is a precomputed 3D layout of people and places. The vertical axis is
// time; the horizontal plane is a force-directed layout of co-occurrence.
type Cosmos struct {
	IDs    []string  `json:"ids"`
	Names  []string  `json:"names"`
	Type   []int     `json:"t"` // 0 person, 1 place
	X      []float32 `json:"x"`
	Y      []float32 `json:"y"`
	Z      []float32 `json:"z"`
	W      []float32 `json:"w"`
	Group  []int     `json:"g"`
	Year   []int     `json:"yr"`
	Edges  []int32   `json:"e"` // a,b,kind (0 cooccur,1 family,2 person-place)
	Groups []string  `json:"groups"`
}

func hash01(s string, salt byte) float32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	h.Write([]byte{salt})
	return float32(h.Sum32()%100000) / 100000
}

func (a *Atlas) Cosmos() *Cosmos {
	a.cosmosOnce.Do(func() { a.cosmos = a.buildCosmos() })
	return a.cosmos
}

func (a *Atlas) buildCosmos() *Cosmos {
	c := &Cosmos{}
	idx := map[string]int{}
	add := func(id, name string, t int, w float32, g, yr int) {
		idx[id] = len(c.IDs)
		c.IDs = append(c.IDs, id)
		c.Names = append(c.Names, name)
		c.Type = append(c.Type, t)
		c.W = append(c.W, w)
		c.Group = append(c.Group, g)
		c.Year = append(c.Year, yr)
	}
	for _, p := range a.People {
		if p.Verses < 2 && len(p.Children) == 0 && p.Father == "" {
			continue
		}
		yr := (p.Y0 + p.Y1) / 2
		add(p.ID, p.Name, 0, float32(math.Log1p(float64(p.Verses))), p.Group, yr)
	}
	for _, p := range a.Places {
		if !p.HasGeo || p.Verses < 2 {
			continue
		}
		add(p.ID, p.Name, 1, float32(math.Log1p(float64(p.Verses))), 0, (p.Y0+p.Y1)/2)
	}
	for _, g := range a.Groups {
		c.Groups = append(c.Groups, g.Name)
	}

	pair := map[[2]int32]int{}
	key := func(i, j int) [2]int32 {
		if i > j {
			i, j = j, i
		}
		return [2]int32{int32(i), int32(j)}
	}
	for _, v := range a.Verses {
		var ps []int
		if len(v.People) <= 12 {
			for _, id := range v.People {
				if i, ok := idx[id]; ok {
					ps = append(ps, i)
				}
			}
		}
		for x := 0; x < len(ps); x++ {
			for y := x + 1; y < len(ps); y++ {
				if ps[x] != ps[y] {
					pair[key(ps[x], ps[y])] += 2
				}
			}
		}
		if len(v.Places) <= 6 {
			for _, pid := range v.Places {
				pi, ok := idx[pid]
				if !ok {
					continue
				}
				for _, i := range ps {
					pair[key(i, pi)] += 1
				}
			}
		}
	}
	type edge struct {
		a, b int32
		w    int
		kind int32
	}
	var edges []edge
	for k, w := range pair {
		kind := int32(0)
		if c.Type[k[0]] != c.Type[k[1]] {
			kind = 2
			if w < 3 {
				continue
			}
		} else if w < 6 {
			continue
		}
		edges = append(edges, edge{k[0], k[1], w, kind})
	}
	for _, p := range a.People {
		pi, ok := idx[p.ID]
		if !ok {
			continue
		}
		for _, ch := range p.Children {
			if ci, ok := idx[ch]; ok {
				edges = append(edges, edge{int32(pi), int32(ci), 12, 1})
			}
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].w != edges[j].w {
			return edges[i].w > edges[j].w
		}
		if edges[i].a != edges[j].a {
			return edges[i].a < edges[j].a
		}
		return edges[i].b < edges[j].b
	})
	if len(edges) > 14000 {
		edges = edges[:14000]
	}
	for _, e := range edges {
		c.Edges = append(c.Edges, e.a, e.b, e.kind)
	}

	// layout
	n := len(c.IDs)
	px := make([]float32, n)
	py := make([]float32, n)
	pz := make([]float32, n)
	ty := make([]float32, n)
	for i := 0; i < n; i++ {
		ang := hash01(c.IDs[i], 1) * 2 * math.Pi
		r := 20 + 40*hash01(c.IDs[i], 2)
		ty[i] = float32(c.Year[i]) / 4100 * 90
		px[i] = float32(math.Cos(float64(ang))) * r
		pz[i] = float32(math.Sin(float64(ang))) * r
		py[i] = ty[i]
	}
	vx := make([]float32, n)
	vy := make([]float32, n)
	vz := make([]float32, n)
	const iters = 90
	for it := 0; it < iters; it++ {
		cool := float32(1) - float32(it)/float32(iters)*0.9
		// repulsion (naive but n is a few thousand)
		for i := 0; i < n; i++ {
			xi, yi, zi := px[i], py[i], pz[i]
			var fx, fy, fz float32
			for j := i + 1; j < n; j++ {
				dx, dy, dz := xi-px[j], yi-py[j], zi-pz[j]
				d2 := dx*dx + dy*dy + dz*dz + 0.5
				if d2 > 900 {
					continue
				}
				f := 3.0 / d2
				fx += dx * f
				fy += dy * f
				fz += dz * f
				vx[j] -= dx * f
				vy[j] -= dy * f
				vz[j] -= dz * f
			}
			vx[i] += fx
			vy[i] += fy
			vz[i] += fz
		}
		for k := 0; k < len(c.Edges); k += 3 {
			a1, b1 := c.Edges[k], c.Edges[k+1]
			dx, dy, dz := px[b1]-px[a1], py[b1]-py[a1], pz[b1]-pz[a1]
			d := float32(math.Sqrt(float64(dx*dx+dy*dy+dz*dz))) + 0.01
			f := (d - 6) * 0.012
			dx, dy, dz = dx/d*f, dy/d*f, dz/d*f
			vx[a1] += dx
			vy[a1] += dy
			vz[a1] += dz
			vx[b1] -= dx
			vy[b1] -= dy
			vz[b1] -= dz
		}
		for i := 0; i < n; i++ {
			// gravity to axis, spring to time
			vx[i] -= px[i] * 0.004
			vz[i] -= pz[i] * 0.004
			vy[i] += (ty[i] - py[i]) * 0.06
			px[i] += clamp(vx[i]*cool, 3)
			py[i] += clamp(vy[i]*cool, 3)
			pz[i] += clamp(vz[i]*cool, 3)
			vx[i], vy[i], vz[i] = 0, 0, 0
		}
	}
	c.X, c.Y, c.Z = px, py, pz
	return c
}

func clamp(v, m float32) float32 {
	if v > m {
		return m
	}
	if v < -m {
		return -m
	}
	return v
}
