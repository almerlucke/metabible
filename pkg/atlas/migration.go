package atlas

import "sort"

const (
	migFrom = -4100
	migStep = 100
	migN    = 42 // buckets of 100 years, 4100 BC .. AD 100
)

type MigPlace struct {
	ID   string  `json:"id"`
	Name string  `json:"n"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}

type MigTrailPoint struct {
	Bucket int      `json:"b"`
	Lat    float64  `json:"lat"`
	Lng    float64  `json:"lng"`
	Weight int      `json:"w"`
	Top    []string `json:"top"`
}

type MigGroup struct {
	Index int             `json:"g"`
	Name  string          `json:"n"`
	Trail []MigTrailPoint `json:"trail"`
}

type Migration struct {
	From   int        `json:"from"`
	Step   int        `json:"step"`
	N      int        `json:"nb"`
	Places []MigPlace `json:"places"`
	Counts [][3]int   `json:"counts"` // place index, bucket, mentions
	Groups []MigGroup `json:"groups"`
}

func bucketOf(year int) int {
	b := (year - migFrom) / migStep
	if b < 0 {
		return 0
	}
	if b >= migN {
		return migN - 1
	}
	return b
}

// Migration aggregates place mentions per century overall and, for each tribe or
// group, a weighted centroid per century: an approximate trail of where the group
// appears in the text over time.
func (a *Atlas) Migration() *Migration {
	m := &Migration{From: migFrom, Step: migStep, N: migN}
	idx := map[string]int{}
	for _, p := range a.Places {
		if p.HasGeo && p.Verses > 0 {
			idx[p.ID] = len(m.Places)
			m.Places = append(m.Places, MigPlace{p.ID, p.Name, p.Lat, p.Lng})
		}
	}
	counts := map[[2]int]int{}

	type acc struct {
		lat, lng float64
		w        int
		names    map[string]int
	}
	groups := make([][]*acc, len(a.Groups)+1)
	for g := range groups {
		groups[g] = make([]*acc, migN)
	}

	for _, v := range a.Verses {
		if len(v.Places) == 0 || len(v.Places) > 5 {
			continue
		}
		b := bucketOf(v.Year)
		var gs map[int]bool
		for _, pid := range v.People {
			if pe := a.PersonByID[pid]; pe != nil && pe.Group > 0 {
				if gs == nil {
					gs = map[int]bool{}
				}
				gs[pe.Group] = true
			}
		}
		for _, plid := range v.Places {
			pi, ok := idx[plid]
			if !ok {
				continue
			}
			counts[[2]int{pi, b}]++
			for g := range gs {
				x := groups[g][b]
				if x == nil {
					x = &acc{names: map[string]int{}}
					groups[g][b] = x
				}
				pl := m.Places[pi]
				x.lat += pl.Lat
				x.lng += pl.Lng
				x.w++
				x.names[pl.Name]++
			}
		}
	}

	for k, n := range counts {
		m.Counts = append(m.Counts, [3]int{k[0], k[1], n})
	}
	sort.Slice(m.Counts, func(i, j int) bool {
		if m.Counts[i][0] != m.Counts[j][0] {
			return m.Counts[i][0] < m.Counts[j][0]
		}
		return m.Counts[i][1] < m.Counts[j][1]
	})

	for gi, g := range a.Groups {
		mg := MigGroup{Index: gi + 1, Name: g.Name}
		for b, x := range groups[gi+1] {
			if x == nil || x.w < 3 {
				continue
			}
			type kv struct {
				k string
				v int
			}
			var top []kv
			for k, v := range x.names {
				top = append(top, kv{k, v})
			}
			sort.Slice(top, func(i, j int) bool {
				if top[i].v != top[j].v {
					return top[i].v > top[j].v
				}
				return top[i].k < top[j].k
			})
			pt := MigTrailPoint{Bucket: b, Lat: x.lat / float64(x.w), Lng: x.lng / float64(x.w), Weight: x.w}
			for i := 0; i < len(top) && i < 3; i++ {
				pt.Top = append(pt.Top, top[i].k)
			}
			mg.Trail = append(mg.Trail, pt)
		}
		if len(mg.Trail) >= 2 {
			m.Groups = append(m.Groups, mg)
		}
	}
	return m
}
