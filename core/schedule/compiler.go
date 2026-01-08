package schedule

import "slices"

type compiler struct {
	idToIdx map[uint64]int

	adj   [][]int
	indeg []int

	ready []int
	order []uint64
}

func newCompiler() *compiler {
	return &compiler{
		idToIdx: make(map[uint64]int),
	}
}

func (c *compiler) compileDirty(store *storage, li *labelIndex, log *logger) {
	stages := make([]stage, 0)
	for st := range store.dirty.Drain() {
		stages = append(stages, st)
	}
	slices.Sort(stages)

	for _, st := range stages {
		c.compileStage(store, li, log, st)
	}
}

func (c *compiler) compileStage(store *storage, li *labelIndex, log *logger, stage stage) {
	ids := store.stages[stage]
	if len(ids) <= 1 {
		c.setPlan(store, stage, ids)
		return
	}

	c.resetStage(ids)
	c.indexStage(ids)
	c.buildGraph(store, li, log, ids)

	order, ok := c.topo(ids)
	if !ok {
		log.add("cycle detected in stage; using insertion order", 0, "", stage)
		c.setPlan(store, stage, ids)
		return
	}

	c.setPlan(store, stage, order)
}

func (c *compiler) resetStage(ids []uint64) {
	n := len(ids)

	for k := range c.idToIdx {
		delete(c.idToIdx, k)
	}

	if cap(c.adj) < n {
		c.adj = make([][]int, n)
	} else {
		c.adj = c.adj[:n]
		for i := range c.adj {
			c.adj[i] = c.adj[i][:0]
		}
	}

	if cap(c.indeg) < n {
		c.indeg = make([]int, n)
	} else {
		c.indeg = c.indeg[:n]
		for i := range c.indeg {
			c.indeg[i] = 0
		}
	}

	c.ready = c.ready[:0]
	c.order = c.order[:0]
}

func (c *compiler) indexStage(ids []uint64) {
	for i, id := range ids {
		c.idToIdx[id] = i
	}
}

func (c *compiler) addEdge(from, to int) {
	if from == to {
		return
	}
	if slices.Contains(c.adj[from], to) {
		return
	}

	c.adj[from] = append(c.adj[from], to)
	c.indeg[to]++
}

func (c *compiler) buildGraph(store *storage, li *labelIndex, log *logger, ids []uint64) {
	for currIdx, id := range ids {
		n, ok := store.get(id)
		if !ok || n == nil {
			continue
		}
		c.applyAfter(li, log, currIdx, n)
		c.applyBefore(li, log, currIdx, n)
	}
}

func (c *compiler) applyAfter(li *labelIndex, log *logger, currIdx int, curr *node) {
	for _, lbl := range curr.after {
		ref, ok := c.resolve(li, log, curr, lbl)
		if !ok {
			continue
		}
		j, ok := c.idToIdx[ref.id]
		if !ok {
			continue
		}
		c.addEdge(j, currIdx)
	}
}

func (c *compiler) applyBefore(li *labelIndex, log *logger, currIdx int, curr *node) {
	for _, lbl := range curr.before {
		ref, ok := c.resolve(li, log, curr, lbl)
		if !ok {
			continue
		}
		j, ok := c.idToIdx[ref.id]
		if !ok {
			continue
		}
		c.addEdge(currIdx, j)
	}
}

func (c *compiler) resolve(li *labelIndex, log *logger, curr *node, lbl string) (label, bool) {
	ref, ok := li.get(lbl)
	if !ok {
		log.add("unknown dependency label: "+lbl, curr.id, curr.label, curr.stage)
		return label{}, false
	}
	if ref.stage != curr.stage {
		log.add("cross-stage dependency ignored: "+lbl, curr.id, curr.label, curr.stage)
		return label{}, false
	}
	return ref, true
}

func (c *compiler) topo(ids []uint64) ([]uint64, bool) {
	n := len(ids)

	for i := range n {
		if c.indeg[i] == 0 {
			c.ready = append(c.ready, i)
		}
	}

	for head := 0; head < len(c.ready); head++ {
		v := c.ready[head]
		c.order = append(c.order, ids[v])

		for _, w := range c.adj[v] {
			c.indeg[w]--
			if c.indeg[w] == 0 {
				c.ready = append(c.ready, w)
			}
		}
	}

	if len(c.order) != n {
		return nil, false
	}
	return c.order, true
}

func (c *compiler) setPlan(store *storage, stage stage, ids []uint64) {
	dst := store.plan[stage]
	dst = append(dst[:0], ids...)
	store.plan[stage] = dst
}
