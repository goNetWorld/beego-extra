package main

import (
	//_ "github.com/go-sql-driver/mysql"
)

type mysqlConfig struct {
	name      string
	user      string //用户名
	password  string //用户密码
	location  string //所在位置,地址
	port      string //端口
	character string //默认字符集
}

type UserDataTable struct {
	AutoId int64 `orm:"column(AutoId)" json:"-"`
	RegisterTime int64 `orm:"column(RegisterTime)" valid:"Range(123,322)" cn:"注册时间"`
	UserPhone string `orm:"current(UserPhone)" form:"UserPhone" valid:"Required;" cn:"手机号" en:"phone"`
}

func main()  {
	var user UserDataTable
	user.UserPhone=""
	user.AutoId=12321
	user.RegisterTime=132
	UpdateGlobalFieldTag("cn")
	v:=Validate{
		fieldTag:"en",
	}
	e:=v.Valid(user)
	if e!=nil{
		println(e.Error())
	}
}


//func main1()  {
//	mysqlData := mysqlConfig{}
//	mysqlData.name = "mysql"
//	mysqlData.user = "root"
//	mysqlData.password = "sjk1234"
//	mysqlData.character = "utf-8"
//	mysqlData.port = "3306"
//
//	mysqlData.location = "127.0.0.1"
//	err := orm.RegisterDriver(mysqlData.name, orm.DRMySQL)
//	if err != nil {
//		logs.GetLogger("MysqlDbInit").Println(err)
//	}
//	var dbSource string
//	// root /fwqsjk1122/
//	//grant all privileges  on *.* to "price_monitor_user"@'%';
//	//create user 'price_monitor_user'@'%' identified by 'price_cipher_1122';
//	dbSource=strings.ToLower("price_monitor_user:price_cipher_1122@tcp(127.0.0.1:33445)/PriceMonitor?charset=utf8")
//	if runtime.GOOS=="windows"{
//		dbSource=strings.ToLower("root:sjk1234@tcp(127.0.0.1:3306)/PriceMonitor?charset=utf8")
//	}
//	err = orm.RegisterDataBase("default", mysqlData.name,dbSource , 1)
//	if err != nil {
//		logs.GetLogger("MysqlDbRegister").Println(err)
//		os.Exit(1)
//	}
//	db, _ := orm.GetDB("default")
//	_=db.Ping()
//	db.SetConnMaxLifetime(-1)
//	db.SetMaxIdleConns(1)
//	db.SetMaxOpenConns(1)
//	orm.DefaultTimeLoc = time.UTC
//	o:=orm.NewOrm()
//	var userData []UserDataTable
//	_,err=o.Raw("select * from UserDataTable").QueryRows("column",&userData)
//	if err!=nil{
//		println(err.Error())
//	}
//	var json=jsoniter.Config{
//		TagKey:"form",
//		OnlyTaggedField:true,
//	}.Froze()
//	d1,_:=json.Marshal(userData)
//	println(string(d1))
//}
