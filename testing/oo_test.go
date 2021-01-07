package testing

import (
	db_token "distributed/Models/mytoken"
	"distributed/distributed"
	"fmt"
	"testing"
	"time"
)


func Test_Create(t *testing.T)  {
	cfgList := getConfg()
	p,err:= distributed.NewMyDbDistributed(cfgList)
	if err != nil {
		fmt.Println("error :",err.Error())
	}
	for i:=0;i<10;i++ {
		orderCode := p.Generate(time.Now(),int64(i))
		conn,err := p.AnalysisCode(orderCode)
		if err !=nil {
			fmt.Println(err.Error())
			return
		}
		var mm db_token.MToken
		mm.Token = orderCode
		tableName := fmt.Sprintf("m_token_%v",p.Sup(conn.Cfg.TableNo,2))
		db :=conn.Conn.Table(tableName).Create(&mm)
		if db.Error != nil {
			fmt.Println("插入失败 :","table:",tableName," token:",orderCode," cfg:",conn)
		}else{
			fmt.Println("插入成功 :","table:",tableName," token:",orderCode," cfg:",conn)
		}
	}
}


func Test_Search(t *testing.T) {

	cfgList := getConfg()
	p,err:= distributed.NewMyDbDistributed(cfgList)
	if err != nil {
		fmt.Print(err.Error())
	}

	{
		var ttt db_token.MToken
		token := "2021010618045445860402010020"
		conn,err := p.AnalysisCode(token)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		tableName := fmt.Sprintf("m_token_%v",p.Sup(conn.Cfg.TableNo,2))
		conn.Conn.Table(tableName).Where("token = ?",token).Find(&ttt)
		fmt.Println("查询结果 : ",ttt)
	}
}





func getConfg() []distributed.DbConfig{

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
	{
		var cfg1 distributed.DbConfig
		cfg1.DbNo = 01
		cfg1.TableNo = 02
		cfg1.Host = "127.0.0.1"
		cfg1.Port = "3306"
		cfg1.UserName = "root"
		cfg1.UserPwd = "root"
		cfg1.DbName = "yk_mytoken"
		cfgList = append(cfgList, cfg1)
	}

	{
		var cfg1 distributed.DbConfig
		cfg1.DbNo = 02
		cfg1.TableNo = 01
		cfg1.Host = "127.0.0.1"
		cfg1.Port = "3306"
		cfg1.UserName = "root"
		cfg1.UserPwd = "root"
		cfg1.DbName = "yk_mytoken02"
		cfgList = append(cfgList, cfg1)
	}

	return cfgList
}


func Test_createOneDb(t *testing.T)  {
	cfgList := getOneConfg()
	p,err:= distributed.NewMyDbDistributed(cfgList)
	if err != nil {
		fmt.Println("error :",err.Error())
	}
	for i:=0;i<10;i++ {
		orderCode := p.Generate(time.Now(),int64(i))
		conn,err := p.AnalysisCode(orderCode)
		if err !=nil {
			fmt.Println(err.Error())
			return
		}
		var mm db_token.MToken
		mm.Token = orderCode
		tableName := fmt.Sprintf("m_token_%v",conn.Name)
		db :=conn.Conn.Table(tableName).Create(&mm)
		if db.Error != nil {
			fmt.Println("插入失败 :","table:",tableName," token:",orderCode," x:",conn)
		}else{
			fmt.Println("插入成功 :","table:",tableName," token:",orderCode," x:",conn)
		}
	}
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