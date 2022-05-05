package symdiff

import (
	"testing"

	"gorgonia.org/gorgonia/execution/engines"
	"gorgonia.org/gorgonia/exprgraph"
	. "gorgonia.org/gorgonia/internal/api"
	"gorgonia.org/tensor"
)

func TestForwardDiffAnalysis(t *testing.T) {
	eng := engines.NewStd()
	g := exprgraph.NewGraph(eng)
	x := exprgraph.NewNode(g, "x", tensor.WithShape(2, 3), tensor.WithBacking([]float64{1, 2, 3, 4, 5, 6}))
	y := exprgraph.NewNode(g, "y", tensor.WithShape(3, 2), tensor.WithBacking([]float64{6, 5, 4, 3, 2, 1}))
	z := exprgraph.NewNode(g, "z", tensor.WithShape(), tensor.Of(tensor.Float64))
	a := exprgraph.NewNode(g, "a", tensor.WithShape(), tensor.Of(tensor.Float64))

	xy, err := MatMul(x, y)
	if err != nil {
		t.Fatalf("MatMul Err %v", err)
	}
	xypz, err := Add(xy, z)
	if err != nil {
		t.Fatalf("Add Err %v", err)
	}

	sorted, err := exprgraph.Sort(g)
	if err != nil {
		t.Fatalf("Sort error %v", err)
	}
	affectsOutput := forwardDiffAnalysis(g, []*exprgraph.Node{xypz.(*exprgraph.Node)}, sorted)
	if affectsOutput.Contains(a.NodeID()) {
		t.Errorf("Affects Output has `a`. This is incorrect.")
	}

	for _, n := range []*exprgraph.Node{x, y, z} {
		if !affectsOutput.Contains(n.NodeID()) {
			t.Errorf("Expected %s to affect output", n)
		}
	}
}
