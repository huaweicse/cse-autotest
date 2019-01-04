package exchanger_test

import (
	"code.huawei.com/cse/pkg/exchanger"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestInterServiceExchger_Init(t *testing.T) {
	t.Log("---remove duplicate elem")
	tmpElem := []string{"a", "b", "b", "c", "c", "d"}
	tmpElemWithNoDuplicat := []string{"a", "b", "c", "d"}
	e := exchanger.InterServiceExchger{
		Elem: tmpElem[:],
	}
	e.Init()
	l := e.Elem[:]
	sort.Strings(l)
	sort.Strings(tmpElemWithNoDuplicat)
	assert.Equal(t, len(tmpElemWithNoDuplicat), len(l))
	for i, _ := range tmpElemWithNoDuplicat {
		assert.Equal(t, tmpElemWithNoDuplicat[i], l[i])
	}
	assert.Equal(t, 4, len(e.Elem))
	assert.Equal(t, 4, len(e.Level2Map))
	assert.Equal(t, 4, len(e.Level3Map))

	t.Log("---length of elem is 0")
	e2 := exchanger.InterServiceExchger{}
	e2.Init()
	assert.Equal(t, 0, len(e2.Level2Map))
	assert.Equal(t, 0, len(e2.Level3Map))

	t.Log("---length of elem is 1")
	e3 := exchanger.InterServiceExchger{
		Elem: []string{"a"},
	}
	e3.Init()
	assert.Equal(t, 0, len(e3.Level2Map))
	assert.Equal(t, 0, len(e3.Level3Map))

	t.Log("---length of elem is 2")
	e4 := exchanger.InterServiceExchger{
		Elem: []string{"a", "b"},
	}
	e4.Init()
	assert.Equal(t, 2, len(e4.Level2Map))
	assert.Equal(t, 0, len(e4.Level3Map))

	t.Log("---length of elem is 3")
	e5 := exchanger.InterServiceExchger{
		Elem: []string{"a", "b", "c"},
	}
	e5.Init()
	assert.Equal(t, 3, len(e5.Level2Map))
	assert.Equal(t, 3, len(e5.Level3Map))
}

func TestInterServiceExchger_GetTargets_Reset(t *testing.T) {
	e := exchanger.InterServiceExchger{
		Elem: []string{"a", "b", "c", "d", "e"},
	}
	e.Init()
	assert.Equal(t, 5, len(e.Level2Map))
	assert.Equal(t, 5, len(e.Level3Map))
	for _, v := range e.Elem {
		assert.Equal(t, 4, len(e.Level2Map[v]))
		m := e.Level3Map[v]
		assert.Equal(t, 4, len(m))
		for _, vv := range m {
			assert.Equal(t, 3, len(vv))
		}
	}

	e.Reset()
	assert.Nil(t, e.Elem)
	assert.Nil(t, e.Level2Map)
	assert.Nil(t, e.Level3Map)
}

func TestNewExchanger(t *testing.T) {
	e := exchanger.NewExchanger("a", "b", "c", "d")
	assert.Equal(t, 4, len(e.Level2Map))
}
