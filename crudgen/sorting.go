package main

type StructFieldList []StructField

func (l StructFieldList) Len() int {
	return len(l)
}

func (l StructFieldList) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}

func (l StructFieldList) Swap(i, j int) {
	tmp := l[i]
	l[i] = l[j]
	l[j] = tmp
}

type StructTypeList []*StructType

func (l StructTypeList) Len() int {
	return len(l)
}

func (l StructTypeList) Less(i, j int) bool {
	return l[i].Name < l[j].Name
}

func (l StructTypeList) Swap(i, j int) {
	tmp := l[i]
	l[i] = l[j]
	l[j] = tmp
}
