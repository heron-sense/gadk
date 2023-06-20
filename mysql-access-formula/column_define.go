package maf

type TableDefinition struct {
	cols []Column
}

type Column struct {
	Name string
	Type byte
}

func DeclAsStr(name string, required bool, maxLen uint16) Column {
	return Column{}
}

func DeclAsU64(name string, required bool) Column {
	return Column{}
}

func DeclAsU16(name string, required bool) Column {
	return Column{}
}
