package domain

type (
	Package struct {
		FilePath   string
		Name       string
		Variables  Fields
		Constants  Fields
		Interfaces Interfaces
		Classes    Classes
		Functions  Functions
	}
	Packages []Package
)

func (packages Packages) AllClasses() Classes {
	var classes Classes
	for _, domainPackage := range packages {
		classes = append(classes, domainPackage.Classes...)
	}
	return classes
}

func (packages Packages) AllInterfaces() Interfaces {
	var interfaces Interfaces
	for _, domainPackage := range packages {
		interfaces = append(interfaces, domainPackage.Interfaces...)
	}
	return interfaces
}
