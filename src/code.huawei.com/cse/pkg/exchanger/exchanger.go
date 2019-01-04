package exchanger

// use it after initialize
type InterServiceExchger struct {
	Elem      []string
	Level2Map map[string][]string
	Level3Map map[string]map[string][]string
}

func (e *InterServiceExchger) Level2Targets(name string) []string {
	if _, ok := e.Level2Map[name]; !ok {
		return nil
	}
	return e.Level2Map[name][:]
}

func (e *InterServiceExchger) Level3Targets(name string) map[string][]string {
	if _, ok := e.Level3Map[name]; !ok {
		return nil
	}
	r := make(map[string][]string)
	for k, v := range e.Level3Map[name] {
		r[k] = v[:]
	}
	return r
}

func (e *InterServiceExchger) Init() {
	if len(e.Elem) <= 1 {
		return
	}

	// remove duplicate elem
	m := make(map[string]bool)
	for _, k := range e.Elem {
		m[k] = true
	}
	tmpElem := make([]string, 0)
	for k, _ := range m {
		tmpElem = append(tmpElem, k)
	}
	e.Elem = tmpElem[:]

	e.Level2Map = make(map[string][]string)
	for _, i1 := range e.Elem {
		e.Level2Map[i1] = make([]string, 0)
		for _, i2 := range e.Elem {
			if i1 != i2 {
				e.Level2Map[i1] = append(e.Level2Map[i1], i2)
			}
		}
	}

	if len(e.Elem) <= 2 {
		return
	}
	e.Level3Map = make(map[string]map[string][]string)
	for _, i1 := range e.Elem {
		e.Level3Map[i1] = make(map[string][]string)
		for _, i2 := range e.Elem {
			if i2 == i1 {
				continue
			}
			e.Level3Map[i1][i2] = make([]string, 0)
			for _, i3 := range e.Elem {
				if i1 != i2 && i2 != i3 && i1 != i3 {
					e.Level3Map[i1][i2] = append(e.Level3Map[i1][i2], i3)
				}
			}
		}
	}
}

func (e *InterServiceExchger) Reset() {
	e.Elem = nil
	e.Level2Map = nil
	e.Level3Map = nil
}

func NewExchanger(names ...string) *InterServiceExchger {
	e := &InterServiceExchger{}
	e.Elem = names
	e.Init()
	return e
}
