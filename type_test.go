package gisp

import (
	"reflect"
	"testing"
)

func TestTypeFound(t *testing.T) {
	m := money{9.99, "USD"}
	g := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	g.DefAs("money", reflect.TypeOf(m))
	_, err := g.Parse("(var bill::money)")
	if err != nil {
		t.Fatalf("except define a money var but error: %v", err)
	}
	g.Setvar("bill", m)

	mny, ok := g.Lookup("bill")
	if !ok {
		t.Fatalf("money var bill as %v not found ", m)
	}
	if !reflect.DeepEqual(m, mny) {
		t.Fatalf("except got money var bill as %v but %v", m, mny)
	}
}
