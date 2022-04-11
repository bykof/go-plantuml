package domain

type (
	Interface struct {
		Name      string
		Functions Functions
	}

	Interfaces []Interface
)

func (domainInterface Interface) IsImplementedByClass(class Class) bool {
	for _, interfaceFunction := range domainInterface.Functions {
		var interfaceFunctionIsImplemented = false
		for _, classFunction := range class.Functions {
			if interfaceFunction.EqualImplementation(classFunction) {
				interfaceFunctionIsImplemented = true
				break
			}
		}
		if !interfaceFunctionIsImplemented {
			return false
		}
	}
	return true
}
