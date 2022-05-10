package kite_go

type Type string

const (
	String  Type = "string"
	Number       = "number"
	Version      = "version"
	Float        = "float"
	Boolean      = "bool"
)

type Variable struct {
	Name  string      `bson:"name" json:"name"`
	Type  Type        `bson:"type" json:"type"`
	Value interface{} `bson:"value" json:"value"`
}
