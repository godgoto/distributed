package main

import "distributed/distributed"

func main()  {
	getOneConfg()
}
func getOneConfg() []distributed.DbConfig {

	var cfgList []distributed.DbConfig
	{
		var cfg1 distributed.DbConfig
		cfg1.DbNo = 01
		cfg1.TableNo = 01
		cfg1.Host = "localhost"
		cfg1.Port = "3306"
		cfg1.UserName = "root"
		cfg1.UserPwd = "root"
		cfg1.DbName = "yk_mytoken"
		cfgList = append(cfgList, cfg1)
	}
	return cfgList
}