package parser

import "CCIG/ast"

type OffsetGenerator struct {
	currentOffset int
	offsetWidth   int
	knownVars     []*ast.VarStatement
}

func assignOffsets(nodes []ast.Node) int {
	g := OffsetGenerator{
		currentOffset: 1,
		offsetWidth:   8,
	}

	for _, node := range nodes {
		g.assignOffset(node)
	}

	return g.currentOffset * g.offsetWidth
}

func (g *OffsetGenerator) assignOffset(node ast.Node) {
	switch node.(type) {
	case *ast.VarStatement:
		s := node.(*ast.VarStatement)
		g.addVar(s)
	case *ast.ExprStatement:
		e := node.(*ast.ExprStatement)
		g.assignOffset(e.Value)
	case *ast.NumberExpr:
		// TODO
		return
	case *ast.PrefixExpr:
		p := node.(*ast.PrefixExpr)
		g.assignOffset(p.Right)
	case *ast.InfixExpr:
		i := node.(*ast.InfixExpr)
		g.assignOffset(i.Left)
		g.assignOffset(i.Right)
	case *ast.IdentifierExpr:
		s := node.(*ast.IdentifierExpr)
		s.Offset = findVar(g.knownVars, s.Value).Offset
	default:
		panic("TODO")
	}
}

func (g *OffsetGenerator) addVar(node *ast.VarStatement) {
	node.Offset = g.currentOffset * g.offsetWidth

	g.knownVars = append(g.knownVars, node)
	g.currentOffset += 1
}

func findVar(knownVars []*ast.VarStatement, ident string) *ast.VarStatement {
	for _, knownVar := range knownVars {
		if knownVar.Name == ident {
			return knownVar
		}
	}
	panic("TODO")
}
