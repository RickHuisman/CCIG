package parser

import "CCIG/ast"

type OffsetGenerator struct {
	currentOffset int
	offsetWidth   int
	knownVars     []*ast.VarStatement
}

func assignOffsets(nodes []ast.Statement) int {
	g := OffsetGenerator{
		currentOffset: 0,
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
	case *ast.ReturnStatement:
		r := node.(*ast.ReturnStatement)
		g.assignOffset(r.Value)
	case *ast.IfElseStatement:
		i := node.(*ast.IfElseStatement)
		g.assignOffset(i.Condition)
		g.assignOffset(i.Then)
		g.assignOffset(i.Else)
	case *ast.BlockStatement:
		b := node.(*ast.BlockStatement)
		for _, stmt := range b.Statements {
			g.assignOffset(stmt)
		}
	case *ast.FunStatement:
		f := node.(*ast.FunStatement)
		g.assignOffset(f.Body)
	case *ast.ExprStatement:
		e := node.(*ast.ExprStatement)
		g.assignOffset(e.Value)
	case *ast.PrefixExpr:
		p := node.(*ast.PrefixExpr)
		g.assignOffset(p.Right)
	case *ast.InfixExpr:
		i := node.(*ast.InfixExpr)
		g.assignOffset(i.Left)
		g.assignOffset(i.Right)
	case *ast.CallExpr:
		return
	case *ast.IdentifierExpr:
		s := node.(*ast.IdentifierExpr)
		s.Offset = findVar(g.knownVars, s.Value).Offset
	case *ast.NumberExpr:
		return
	default:
		panic("TODO") // TODO
	}
}

func (g *OffsetGenerator) addVar(node *ast.VarStatement) {
	g.currentOffset += 1
	node.Offset = g.currentOffset * g.offsetWidth

	g.knownVars = append(g.knownVars, node)
}

func findVar(knownVars []*ast.VarStatement, ident string) *ast.VarStatement {
	for _, knownVar := range knownVars {
		if knownVar.Name == ident {
			return knownVar
		}
	}
	panic("TODO")
}
