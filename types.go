package xtypes

import (
	"database/sql"
	"database/sql/driver"
)

type Factory interface {
	New() Object
	Reset(g Object) error
}

type Object interface {
	Comparable
	String() string
}

type Comparable interface {
	Compare(o Object) int
}

type Marshaler interface {
	Marshal([]byte) ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal([]byte) error
}

type SQLValuer interface {
	driver.Valuer
}

type SQLScanner interface {
	sql.Scanner
}
