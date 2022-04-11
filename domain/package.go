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

func (domainPackage Package) Add(otherPackage Package) Package {
	return Package{
		Name:       domainPackage.Name,
		FilePath:   domainPackage.FilePath,
		Variables:  append(domainPackage.Variables, otherPackage.Variables...),
		Constants:  append(domainPackage.Constants, otherPackage.Constants...),
		Interfaces: append(domainPackage.Interfaces, otherPackage.Interfaces...),
		Classes:    append(domainPackage.Classes, otherPackage.Classes...),
		Functions:  append(domainPackage.Functions, otherPackage.Functions...),
	}
}

func (domainPackage Package) IsEmpty() bool {
	return len(
		domainPackage.Variables,
	) == 0 && len(
		domainPackage.Constants,
	) == 0 && len(
		domainPackage.Interfaces,
	) == 0 && len(
		domainPackage.Classes,
	) == 0 && len(
		domainPackage.Functions,
	) == 0
}

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
