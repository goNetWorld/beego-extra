package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type ValidFunc struct {
	Name   string
	Params []interface{}
}

var (
	ValidTag="valid"
	LabelTag="label"
)

func getRegFuncs(tag, key string) (vfs []ValidFunc, str string, err error) {
	tag = strings.TrimSpace(tag)
	index := strings.Index(tag, "Match(/")
	if index == -1 {
		str = tag
		return
	}
	end := strings.LastIndex(tag, "/)")
	if end < index {
		err = fmt.Errorf("invalid Match function")
		return
	}
	reg, err := regexp.Compile(tag[index+len("Match(/") : end])
	if err != nil {
		return
	}
	vfs = []ValidFunc{{"Match", []interface{}{reg, key + ".Match"}}}
	str = strings.TrimSpace(tag[:index]) + strings.TrimSpace(tag[end+len("/)"):])
	return
}

func isStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

type CustomFunc func(cnName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(error)

/*
cnName	用来指示字段的名称,通过Tag获得到的字符串信息
field	通过 field.int(),field.string(),来获取不同的数据值
args	为Tag中用户传递过来的参数信息,用于在内部进行数据验证
*/

type ErrorCustomFunc func(err string,cnName string,field reflect.Value)interface{}
//中文错误信息和字段信息全部进行传递

type validatorMap struct {
	funcs map[string]CustomFunc
	errorFuncs map[string]ErrorCustomFunc
}

var validator=validatorMap{
	funcs: map[string]CustomFunc{
		"Max":Max,
		"Min":Min,
		"Range":Range,
		"MaxSize":MaxSize,
		"MinSize":MinSize,
		"Alpha":Alpha,
		"AlphaNumeric":AlphaNumeric,
		"Numeric":Numeric,
		"ZipCode":ZipCode,
		"Base64":Base64,
		"IpAddress":IpAddress,
		"FixSize":FixSize,
		"Email":Email,
		"Match":Match,
		"Required":Required,
		"MaxFloat":MaxFloat,
		"MinFloat":MinFloat,
		"RangeFloat":RangeFloat,
	},
	errorFuncs:map[string]ErrorCustomFunc{
		"DefaultErrorFunc":DefaultErrorFunc,
	},
}

func DefaultErrorFunc (err string,cnName string,field reflect.Value)interface{}  {
	return err
}

func AddCustomValidatorFunc(name string, f CustomFunc,override ...bool) {
	_,ok:=validator.funcs[name]
	if ok&&len(override)>0&&override[0]{
		//传递了覆盖的操作,同时值为true
		//如果存在,同时未传递覆盖
		println(name+"已存在,未覆盖原函数,funcs.")
		return
	}
	validator.funcs[name]=f
}

func AddCustomErrorFunc(name string,f ErrorCustomFunc,override ...bool) {
	//当出现错误信息的时候,会调用指定的函数进行处理.
	//可以返回给前端,提示多个信息错误
	_,ok:=validator.errorFuncs[name]
	if ok&&len(override)>0&&override[0]{
		//传递了覆盖的操作,同时值为true
		println(name+"已存在,未覆盖原函数,errorFuncs.")
		return
	}
	validator.errorFuncs[name]=f
}

type Validate struct {
	fieldTag string	//该字段在用户看来的名称,即Form中展示的字段名,如手机号,密码
	multiError bool		//是否返回多个错误信息
	errSet []interface{}	//多个错误信息可以在这里进行获取
	errFuncName string		//错误函数的处理名称
}

var defaultErrorFuncName="DefaultErrorFunc"
var defaultFieldTag="nameTag"
//默认的错误处理函数

func UpdateGlobalErrorFuncName(newName string)  {
	//注册修改全局型的错误处理函数名称
	defaultErrorFuncName=newName
}

func UpdateGlobalFieldTag(newName string)  {
	//注册修改全局型的字段名称信息
	defaultFieldTag=newName
}

func (v *Validate)Valid(obj interface{}) (err error) {
	objV := reflect.ValueOf(obj)
	objT := objV.Type()
	if v.errFuncName==""{
		//默认使用自定义的函数进行处理
		v.errFuncName=defaultErrorFuncName
	}
	if v.fieldTag==""{
		//默认使用自定义的函数进行处理
		v.fieldTag=defaultFieldTag
	}
	switch {
	case isStruct(objT):
	case isStructPtr(objT):
		objT = objT.Elem()
		objV = objV.Elem()
	default:
		err = fmt.Errorf("%v 必须是结构体类型或者结构体指针", obj)
		return
	}
	for i := 0; i < objT.NumField(); i++ {
		fieldValue:=objV.Field(i)
		fieldType:=objT.Field(i)
		tag:=strings.TrimSpace(fieldType.Tag.Get(ValidTag))
		cnName:=strings.TrimSpace(fieldType.Tag.Get(v.fieldTag))
		if tag==""{continue}
		funcs:=strings.Split(tag,";")
		for index:=0;index<len(funcs);index++{
			f:=strings.Split(funcs[index],"(")
			funName:=funcs[index]
			if len(f)>0{
				//Tag中可以加括号也可以不加括号
				funName=f[0]
			}
			var params []string
			if len(f)==2{
				params=strings.Split(strings.TrimRight(f[1],")"),",")
			}
			funName=strings.Trim(funName,"")
			if funName==""{continue}
			fun,ok:=validator.funcs[funName]
			if ok{
				err=fun(cnName,fieldValue,fieldType,params)
				if err!=nil&&!v.multiError{
					//只收集一个错误即返回
					return
				}
				if err!=nil{
					fun,ok:=validator.errorFuncs[v.errFuncName]
					if ok{
						//收集多个错误信息,将用户数据信息全部返回给用户
						v.errSet = append(v.errSet, fun(err.Error(),cnName,objV.Field(i)))
					}else{
						println(v.errFuncName+"函数不存在,无法进行错误处理.")
					}
				}
			}else{
				println(funName+"函数未注册,请先注册再使用.")
			}
		}
	}
	return
}