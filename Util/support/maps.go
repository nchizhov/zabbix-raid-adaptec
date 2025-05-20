package support

func InMap(data map[string]any, field string) bool {
	_, ok := data[field]
	return ok
}

func addNested(m map[string]any, name string, val any) {
	if _, ok := m[name]; !ok {
		m[name] = val
	}
}

func getNestedVal(last bool, val any) any {
	if last {
		return val
	}
	return make(map[string]any)
}

func AddDeepNested(m map[string]any, nests []string, val any) {
	totalNests := len(nests) - 1
	if totalNests == -1 {
		return
	}
	if totalNests == 0 {
		addNested(m, nests[0], val)
		return
	}
	nestTmp := m
	for idx, nest := range nests {
		isLast := idx == totalNests
		nestVal := getNestedVal(isLast, val)

		addNested(nestTmp, nest, nestVal)
		if !isLast {
			nestTmp = nestTmp[nest].(map[string]any)
		}
	}
}
