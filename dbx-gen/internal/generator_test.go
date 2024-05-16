package internal

import "testing"

func Test_Generate(t *testing.T) {

	// ACCOUNTX_DB_HOST = mysql
	// ACCOUNTX_DB_USER = root
	// ACCOUNTX_DB_PASSWORD = basebitxdp
	// ACCOUNTX_DB_NAME = enigma2_accountx
	// ACCOUNTX_DB_CONN_POOL_SIZE = 5
	Options.Driver = "mysql"
	Options.DataSourceName = "root:basebitxdp@tcp(172.18.0.210:32600)/enigma2_workflowx?parseTime=True&loc=Local"
	Options.Schema = "enigma2_accountx"
	Options.Output = "module/module.gen.go"
	Options.Package = "module"
	err := generate()
	if err != nil {
		t.Fatal(err)
	}
}
