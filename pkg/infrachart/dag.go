package infrachart

// findCycle detects a dependency cycle in the chart's resource graph (nodes
// are document indexes, edges point from a resource to the resource it
// references). It returns the nodes of one cycle in reference order, or nil
// when the graph is a DAG. One cycle is enough for the gate — the author
// fixes it and re-runs.
func findCycle(edges [][]int) []int {
	const (
		unvisited = 0
		inStack   = 1
		done      = 2
	)
	state := make([]int, len(edges))
	var stack []int

	var visit func(node int) []int
	visit = func(node int) []int {
		state[node] = inStack
		stack = append(stack, node)
		for _, next := range edges[node] {
			switch state[next] {
			case inStack:
				// Found a back edge — slice the current stack from the first
				// occurrence of next to get the cycle members in order.
				for i, n := range stack {
					if n == next {
						cycle := make([]int, len(stack)-i)
						copy(cycle, stack[i:])
						return cycle
					}
				}
			case unvisited:
				if cycle := visit(next); cycle != nil {
					return cycle
				}
			}
		}
		stack = stack[:len(stack)-1]
		state[node] = done
		return nil
	}

	for node := range edges {
		if state[node] == unvisited {
			if cycle := visit(node); cycle != nil {
				return cycle
			}
		}
	}
	return nil
}
