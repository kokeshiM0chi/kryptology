package bls12381

import (
	crand "crypto/rand"
	"testing"

	"fmt"
	"github.com/stretchr/testify/require"

	"github.com/coinbase/kryptology/pkg/core/curves/native"
)

func TestSinglePairing(t *testing.T) {
	g := new(G1).Generator()
	h := new(G2).Generator()
	fmt.Println("# g:	", g)
	fmt.Println("# h:	", h)

	e := new(Engine)
	fmt.Println("# pairin　before reset:	", e.Result())
	e.AddPair(g, h)
	p := e.Result()
	fmt.Println("# pairing:	", e.Result())

	p.Neg(p)

	e.Reset()
	fmt.Println("# pairin　after reset:	", e.Result())
	e.AddPairInvG2(g, h)
	q := e.Result()
	e.Reset()
	e.AddPairInvG1(g, h)
	r := e.Result()

	require.Equal(t, 1, p.Equal(q))
	require.Equal(t, 1, q.Equal(r))
}

func TestMultiPairing(t *testing.T) {
	const Tests = 10
	e1 := new(Engine)
	e2 := new(Engine)

	g1s := make([]*G1, Tests)
	g2s := make([]*G2, Tests)
	sc := make([]*native.Field, Tests)
	res := make([]*Gt, Tests)
	expected := new(Gt).SetOne()

	for i := 0; i < Tests; i++ {
		var bytes [64]byte
		g1s[i] = new(G1).Generator()
		g2s[i] = new(G2).Generator()
		sc[i] = Bls12381FqNew()
		_, _ = crand.Read(bytes[:])
		sc[i].SetBytesWide(&bytes)
		if i&1 == 0 {
			g1s[i].Mul(g1s[i], sc[i])
		} else {
			g2s[i].Mul(g2s[i], sc[i])
		}
		e1.AddPair(g1s[i], g2s[i])
		e2.AddPair(g1s[i], g2s[i])
		res[i] = e1.Result()
		e1.Reset()
		expected.Add(expected, res[i])
	}

	actual := e2.Result()
	require.Equal(t, 1, expected.Equal(actual))
}

// confirm
// e(a * g1, g2) == e(g1, a * g2)
func TestCoefficients(t *testing.T) {
	e := new(Engine)
	Bls12381FqNew()
	// setting scalar value
	var bytesA [64]byte

	_, _ = crand.Read(bytesA[:])
	var a *native.Field

	a = Bls12381FqNew()
	a.SetBytesWide(&bytesA)

	// generate base point
	g1 := new(G1).Generator()
	g2 := new(G2).Generator()

	g1.Mul(g1, a)

	e.AddPair(g1, g2)
	gt0 := e.Result()

	e.Reset()
	g1 = new(G1).Generator()
	g2 = new(G2).Generator()
	g2.Mul(g2, a)

	e.AddPair(g1, g2)
	gt1 := e.Result()

	require.Equal(t, gt1.Bytes(), gt0.Bytes())
}
