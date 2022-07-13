package sqlx

import "sync"

// Bindvar types supported by Rebind, BindMap and BindStruct.
const (
	UNKNOWN = iota
	QUESTION
	DOLLAR
	NAMED
	AT
)

var defaultBinds = map[int][]string{
	DOLLAR:   []string{"postgres", "pgx", "pq-timeouts", "cloudsqlpostgres", "ql", "nrpostgres", "cockroach"},
	QUESTION: []string{"mysql", "sqlite3", "nrmysql", "nrsqlite3"},
	NAMED:    []string{"oci8", "ora", "goracle", "godror"},
	AT:       []string{"sqlserver"},
}

var binds sync.Map

func init() {
	for bind, drivers := range defaultBinds {
		for _, driver := range drivers {
			BindDriver(driver, bind)
		}
	}

}

// BindType returns the bindtype for a given database given a drivername.
func BindType(driverName string) int {
	itype, ok := binds.Load(driverName)
	if !ok {
		return UNKNOWN
	}
	return itype.(int)
}

// BindDriver sets the BindType for driverName to bindType.
func BindDriver(driverName string, bindType int) {
	binds.Store(driverName, bindType)
}
