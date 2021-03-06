package distributed

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/w3liu/go-common/constant/timeformat"
	"os"
	"strconv"
	"time"
)

const DB_MaxConns = 10
const DB_MaxOpenConns = 15

//配置文件
type DbConfig struct {
	Host           string // 主机名称
	Port           string // 端口
	UserName       string // 用户名
	UserPwd        string // 密码
	DbName         string // 数据库名称
	DbNo           int64  "" //2位   - 库名
	TableNo        int64  "" //2位   - 表名
	DbMaxConns     int
	DBMaxOpenConns int
}

//分布式服务
type DbDistributed struct {
	Key        	string   //2位(库)2位(表名) = key
	TableName  	string   //解析采用的到 - table的名称
	UniqueCode  string 	 //解析采用的到 - 解析的code
	Conn      *gorm.DB //本数据的CONN
	Cfg       DbConfig //配置文件
}

type MyDbDistributed struct {
	cfgList   []DbConfig
	connlist  map[string]DbDistributed
	connlist2 []DbDistributed
}

func NewMyDbDistributed(configList []DbConfig) (*MyDbDistributed, error) {
	p := MyDbDistributed{}
	p.cfgList = configList
	p.connlist = make(map[string]DbDistributed)
	//循环打开
	for _, cfg := range p.cfgList {
		var item DbDistributed
		item.Cfg = cfg
		conn, err := p.OpenDb(cfg)
		if err != nil {
			return nil, err
		}
		item.Key = fmt.Sprintf("%v%v", p.Sup(cfg.DbNo, 2), p.Sup(cfg.TableNo, 2))
		item.Conn = conn
		p.connlist[item.Key] = item
		p.connlist2 = append(p.connlist2, item)
	}
	return &p, nil
}

func (p *MyDbDistributed) OpenDb(cfg DbConfig) (*gorm.DB, error) {
	CloudOrderDb, err := gorm.Open("mysql", cfg.UserName+":"+cfg.UserPwd+"@tcp("+cfg.Host+")/"+cfg.DbName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}
	CloudOrderDb.SingularTable(true)
	CloudOrderDb.DB().SetMaxIdleConns(cfg.DbMaxConns)
	CloudOrderDb.DB().SetMaxOpenConns(cfg.DBMaxOpenConns)
	return CloudOrderDb, nil
}

func (p *MyDbDistributed) CloseDb(cfg DbConfig) {
	keyName := fmt.Sprintf("%v%v", p.Sup(cfg.DbNo, 2), p.Sup(cfg.TableNo, 2))
	if _, ok := p.connlist[keyName]; ok {
		if p.connlist[keyName].Conn != nil {
			p.connlist[keyName].Conn.Close()
		}
	}
}

func (p *MyDbDistributed) CloseDbAll() {
	for _, conn := range p.connlist {
		p.CloseDb(conn.Cfg)
	}
}

//生成24位订单号
// num 为序号
//前面17位代表时间精确到毫秒，中间3位代表进程id，    (22位)2位表名     序号最后4位代表序号
// 20201231161128|753|164|02|01|0002
func (p *MyDbDistributed) Generate(t time.Time, num int64) string {
	//time
	s := t.Format(timeformat.Continuity)
	m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3
	ms := p.Sup(m, 3)
	//pid
	pid := os.Getpid() % 1000
	ps := p.Sup(p.pidSub(int64(pid)), 3)
	//i :=  	atomic.AddInt64(&num, 1)
	qm := num % int64(len(p.connlist))
	db := p.connlist2[qm]

	//fmt.Println(num,"%",len(p.DbList)," = ",p.sup(qm,2))
	//index
	r := num % 10000
	rs := p.Sup(r, 4)
	n := fmt.Sprintf("%s%s%s%s%s%s", s, ms, ps, p.Sup(db.Cfg.DbNo, 2), p.Sup(db.Cfg.TableNo, 2), rs)
	return n
}

//分析出表明  链接
func (p *MyDbDistributed) AnalysisCode(UniqueCode string, tableName string) (DbDistributed, error) {
	var dinfo DbDistributed
	dinfo.Cfg.DbNo = -1
	dinfo.Cfg.TableNo = -1
	dinfo.Key = ""
	dinfo.TableName = ""
	dinfo.UniqueCode = ""
	if len(UniqueCode) == 28 {
		byteCode := []byte(UniqueCode)
		dbNo, _ := strconv.ParseInt(string(byteCode[20:22]), 10, 64)
		dinfo.Cfg.DbNo = dbNo
		tableNo, _ := strconv.ParseInt(string(byteCode[22:24]), 10, 64)
		dinfo.Cfg.TableNo = tableNo
		dinfo.Key = fmt.Sprintf("%v%v", p.Sup(dbNo, 2), p.Sup(tableNo, 2))

		if _, ok := p.connlist[dinfo.Key]; ok {
			rConn := p.connlist[dinfo.Key]
			rConn.TableName = fmt.Sprintf("%v_%v", tableName, p.Sup(p.connlist[dinfo.Key].Cfg.TableNo, 2))
			rConn.UniqueCode = UniqueCode
			return rConn, nil
		}
		return dinfo, errors.New("不存在库和表")
	}
	return dinfo, errors.New("长度必须是28位")
}

func (p *MyDbDistributed) AnalysisListCode(listCode []string, tableName string) (map[string][]DbDistributed, error) {
	var list map[string][]DbDistributed
	list = make(map[string][]DbDistributed)
	for _, item := range listCode {
		d, err := p.AnalysisCode(item, tableName)
		if err != nil {
			fmt.Println(err.Error())
			return list, err
		}
		list[d.Key] = append(list[d.Key], d)
	}
	return list, nil
}

func (p *MyDbDistributed) pidSub(i int64) int64 {
	if i > 999 {
		str := fmt.Sprintf("%d", i)
		iStr := string([]byte(str)[0:3])
		i, _ = strconv.ParseInt(iStr, 10, 64)
	}
	return i
}

//对长度不足n的数字前面补0
func (p *MyDbDistributed) Sup(i int64, n int) string {
	m := fmt.Sprintf("%d", i)
	for len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}
