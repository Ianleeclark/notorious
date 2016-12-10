package sqlStoreImpl

import (
	"testing"
)

var MYSQLSTORE MySQLStore

func TestInitMySQLStore(t *testing.T) {
	MYSQLSTORE = InitMySQLStore()
}

func TestOpenConnection(t *testing.T) {
	_, err := MYSQLSTORE.OpenConnection()
	if err != nil {
		t.Errorf("TestMySQLStore:TestOpenConnection() failed to open a connection with err %v", err)
	}
}
