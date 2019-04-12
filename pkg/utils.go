package pkg

import (
	"github.com/vetcher/go-astra/types"
)

func FindInterface(f *types.File, name string) *types.Interface {
	for _, iface := range f.Interfaces {
		if iface.Name == name {
			return &iface
		}
	}

	return nil
}

func GetPlainParam(params interface{}, paramName string) interface{} {
	p, ok := params.([]interface{})
	if !ok {
		return nil
	}

	if len(p) == 0 {
		return nil
	}

	v, ok := p[0].(map[interface{}]interface{})

	return v[paramName]
}
