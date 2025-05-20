package support

type NilData struct {
	value string
	null  bool
}

func (n *NilData) Value() interface{} {
	if n.null {
		return nil
	}
	return n.value
}

func NewData(x string) NilData {
	return NilData{x, false}
}

func NewNil() NilData {
	return NilData{"", true}
}
