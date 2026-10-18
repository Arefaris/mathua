package graphtheory

import "github.com/chuma-beep/mathua/internal/generator"

// Register adds the graph-theory generators to reg.
func Register(reg *generator.Registry) {
	reg.Register("graph.isomorphism", &isomorphismGen{})
	reg.Register("graph.complement", &complementGen{})
	reg.Register("graph.subgraph", &subgraphGen{})
	reg.Register("graph.homomorphism", &homomorphismGen{})
	reg.Register("graph.product", &productGen{})
	reg.Register("graph.distance", &distanceGen{})
	reg.Register("graph.cut_vertex", &cutVertexGen{})
	reg.Register("graph.bridge", &bridgeGen{})
	reg.Register("graph.connectivity", &connectivityGen{})
	reg.Register("graph.menger", &mengerGen{})
	reg.Register("graph.matching", &matchingGen{})
	reg.Register("graph.perfect_matching", &perfectMatchingGen{})
	reg.Register("graph.tutte", &tutteGen{})
	reg.Register("graph.vertex_cover", &vertexCoverGen{})
	reg.Register("graph.independent_set", &independentSetGen{})
	reg.Register("graph.chromatic_polynomial", &chromaticPolynomialGen{})
	reg.Register("graph.brooks", &brooksGen{})
	reg.Register("graph.edge_coloring", &edgeColoringGen{})
	reg.Register("graph.perfect_graph", &perfectGraphGen{})
	reg.Register("graph.list_coloring", &listColoringGen{})
}
