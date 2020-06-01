package main

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)
var emailPattern = regexp.MustCompile(`^[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[a-zA-Z0-9](?:[\w-]*[\w])?$`)
var ipPattern = regexp.MustCompile(`^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`)
var base64Pattern = regexp.MustCompile(`^(?:[A-Za-z0-99+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$`)
var zipCodePattern = regexp.MustCompile(`^[1-9]\d{5}$`)


func Required(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	obj:=fieldValue.Interface()
	success:=true
	defer func() {
		if !success{
			err=errors.New(fieldName+"不能为空")
		}
	}()
	if obj == nil {
		success=false
		return 
	}

	if str, ok := obj.(string); ok {
		success=len(strings.TrimSpace(str)) > 0
		return 
	}
	if _, ok := obj.(bool); ok {
		success=true
		return 
	}
	if i, ok := obj.(int); ok {
		success=i != 0
		return 
	}
	if i, ok := obj.(uint); ok {
		success=i != 0
		return
	}
	if i, ok := obj.(int8); ok {
		success=i != 0
		return 
	}
	if i, ok := obj.(uint8); ok {
		success=i != 0
		return 
	}
	if i, ok := obj.(int16); ok {
		success=i != 0
		return
	}
	if i, ok := obj.(uint16); ok {
		success=i != 0
		return
	}
	if i, ok := obj.(uint32); ok {
		success=i != 0
		return
	}
	if i, ok := obj.(int32); ok {
		success=i != 0
		return
	}
	if i, ok := obj.(int64); ok {
		success=i != 0
		return
	}
	if i, ok := obj.(uint64); ok {
		success=i != 0
		return
	}
	if i, ok := obj.(float32); ok {
		success=i != 0.0
		return
	}
	if i, ok := obj.(float64); ok {
		success=i != 0.0
		return
	}
	if t, ok := obj.(time.Time); ok {
		success=!t.IsZero()
		return 
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		success=v.Len() > 0
		return 
	}
	return 
}
func Max(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	//fieldValue	为该字段的中文名称,而不是简单的字段名
	//obj	是该字段的具体
	intSet,err:=ParseToInt64(args)
	if err!=nil||len(intSet)<1{
		err=errors.New(fieldType.Name+"中Max函数必须传递参数")
		return
	}
	num:=ConverseToInt64(fieldValue)
	if num==math.MaxInt64{
		err=errors.New(fieldType.Name+"只能为Int类型")
		return
	}
	if num>intSet[0]{
		err=errors.New(fieldName+"不得大于"+strconv.FormatInt(intSet[0],10))
		return
	}
	return
}
func Range(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	intSet,err:=ParseToInt64(args)
	if err!=nil||len(intSet)!=2{
		err=errors.New(fieldType.Name+"中Range函数必须传递参数")
		return
	}
	if intSet[0]>intSet[1]{
		//第一个参数值大于第二个参数值
		err=errors.New(fieldType.Name+"的Range函数参数错误")
		return 
	}
	num:=ConverseToInt64(fieldValue)
	if num==math.MaxInt64{
		err=errors.New(fieldType.Name+"只能为Int等类型")
		return
	}
	if intSet[0]>num||num>intSet[1] {
		//小于第一个参数值 或者 大于第二个参数值
		err = errors.New(fieldName + "的范围是" + strconv.FormatInt(intSet[0], 10) + "-" + strconv.FormatInt(intSet[1], 10) + "之间")
	}
	return
}

func Min(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	intSet,err:=ParseToInt64(args)
	if err!=nil||len(intSet)<1{
		err=errors.New(fieldType.Name+"中Min函数必须传递参数")
		return
	}
	num:=ConverseToInt64(fieldValue)
	if num==math.MaxInt64{
		err=errors.New(fieldType.Name+"只能为Int类型")
		return
	}
	if num<intSet[0]{
		err=errors.New(fieldName+"不得小于"+strconv.FormatInt(intSet[0],10))
		return
	}
	return
}
func MinFloat(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	intSet,err:=ParseToFloat64(args)
	if err!=nil||len(intSet)<1{
		err=errors.New(fieldType.Name+"中MinFloat函数必须传递参数")
		return
	}
	num:=ConverseToFloat64(fieldValue)
	if num==math.MaxFloat64{
		err=errors.New(fieldType.Name+"只能为Float类型")
		return
	}
	if num<intSet[0]{
		v:=strconv.FormatFloat(intSet[0],'f',5,10)
		err=errors.New(fieldName+"不得小于"+v)
		return
	}
	return
}
func MaxFloat(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	floatSet,err:=ParseToFloat64(args)
	if err!=nil||len(floatSet)<1 {
		err = errors.New(fieldType.Name + "中MaxFloat函数必须传递参数")
		return
	}
	num:=ConverseToFloat64(fieldValue)
	if num==math.MaxFloat64{
		err=errors.New(fieldType.Name+"只能为Float类型")
		return
	}
	if num>floatSet[0]{
		v:=strconv.FormatFloat(floatSet[0],'f',5,10)
		err=errors.New(fieldName+"不得大于"+v)
		return
	}
	return
}
func RangeFloat(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	floatSet,err:=ParseToFloat64(args)
	if err!=nil||len(floatSet)!=2{
		err=errors.New(fieldType.Name+"中RangeFloat函数必须传递参数")
		return
	}
	if floatSet[0]>floatSet[1]{
		//第一个参数值大于第二个参数值
		err=errors.New(fieldType.Name+"的RangeFloat函数参数设置错误")
		return
	}
	num:=ConverseToFloat64(fieldValue)
	if num==math.MaxFloat64{
		err=errors.New(fieldType.Name+"只能为Float类型")
		return
	}
	if floatSet[0]>num||num>floatSet[1] {
		//小于第一个参数值 或者 大于第二个参数值
		f1:=strconv.FormatFloat(floatSet[0],'f',10,10)
		f2:=strconv.FormatFloat(floatSet[1],'f',10,10)
		err = errors.New(fieldName + "的范围是" + f1 + "-" + f2 + "之间")
	}
	return
}

func MaxSize(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	intSet,err:=ParseToInt64(args)
	if err!=nil||len(intSet)<1{
		err=errors.New(fieldType.Name+"中MaxSize函数必须传递参数")
		return
	}
	str,ok:=fieldValue.Interface().(string)
	if !ok{
		err=errors.New(fieldType.Name+"只能为String类型")
		return
	}
	if int64(utf8.RuneCountInString(str))>intSet[0]{
		err=errors.New(fieldName+"长度不得大于"+strconv.FormatInt(intSet[0],10))
		return
	}
	return
}

func MinSize(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	intSet,err:=ParseToInt64(args)
	if len(intSet)<1||err!=nil{
		err=errors.New(fieldType.Name+"中MinSize函数必须传递参数")
		return
	}
	str,ok:=fieldValue.Interface().(string)
	if !ok{
		err=errors.New(fieldType.Name+"只能为String类型")
		return
	}
	if int64(utf8.RuneCountInString(str))<intSet[0]{
		err=errors.New(fieldName+"长度不得小于"+strconv.FormatInt(intSet[0],10))
		return
	}
	return
}

func FixSize(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	intSet,err:=ParseToInt64(args)
	if err!=nil||len(intSet)<1{
		err=errors.New(fieldType.Name+"中MaxSize函数必须传递参数")
		return
	}
	str,ok:=fieldValue.Interface().(string)
	if !ok{
		err=errors.New(fieldType.Name+"只能为String类型")
		return
	}
	if int64(utf8.RuneCountInString(str))!=intSet[0]{
		err=errors.New(fieldName+"长度必须为"+strconv.FormatInt(intSet[0],10))
		return
	}
	return
}

func Alpha(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	str,ok:=fieldValue.Interface().(string)
	if !ok{
		err=errors.New(fieldType.Name+"只能为String类型")
		return
	}
	for _, v := range str {
		if ('Z' < v || v < 'A') && ('z' < v || v < 'a') {
			err=errors.New(fieldName+"只能为大小写字母")
			return 
		}
	}
	return
}

func Numeric(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	str,ok:=fieldValue.Interface().(string)
	if !ok{
		err=errors.New(fieldType.Name+"只能为String类型")
		return
	}
	for _, v := range str {
		if '9' < v || v < '0' {
			err=errors.New(fieldName+"只能为0-9数字")
			return 
		}
	}
	return
}
func AlphaNumeric(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	str,ok:=fieldValue.Interface().(string)
	if !ok{
		err=errors.New(fieldType.Name+"只能为String类型")
		return
	}
	for _, v := range str {
		if ('Z' < v || v < 'A') && ('z' < v || v < 'a') && ('9' < v || v < '0') {
			err=errors.New(fieldName+"只能为0-9数字")
			return 
		}
	}
	return
}

func Match(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	_,ok:=fieldValue.Interface().(string)
	if !ok{
		err=errors.New(fieldType.Name+"只能为String类型")
		return
	}
	if len(args)!=1{
		println(fieldType.Name+"中Match函数必须传递参数")
		return
	}
	str:=fieldValue.String()
	r:=regexp.MustCompile(args[0])
	b:=r.MatchString(fmt.Sprintf("%v", str))
	if !b{
		err=errors.New(fieldName+"格式错误")
	}
	return
}
func IpAddress(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	str,ok:=fieldValue.Interface().(string)
	if !ok{
		err=errors.New(fieldType.Name+"只能为String类型")
		return
	}
	b:=ipPattern.MatchString(str)
	if !b{
		err=errors.New(fieldName+"格式错误")
	}
	return
}
func ZipCode(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	str,ok:=fieldValue.Interface().(string)
	if !ok{
		err=errors.New(fieldType.Name+"只能为String类型")
		return
	}
	b:=zipCodePattern.MatchString(str)
	if !b{
		err=errors.New(fieldName+"格式错误")
	}
	return
}
func Email(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	str,ok:=fieldValue.Interface().(string)
	if !ok{
		err=errors.New(fieldType.Name+"只能为String类型")
		return
	}
	b:=emailPattern.MatchString(str)
	if !b{
		err=errors.New(fieldName+"格式错误")
	}
	return
}
func Base64(fieldName string,fieldValue reflect.Value,fieldType reflect.StructField,args []string)(err error){
	str,ok:=fieldValue.Interface().(string)
	if !ok{
		err=errors.New(fieldType.Name+"只能为String类型")
		return
	}
	b:=base64Pattern.MatchString(str)
	if !b{
		err=errors.New(fieldName+"格式错误")
	}
	return
}


func ParseToInt64(s []string)(r []int64,err error) {
	var i int64
	defer func() {
		if err!=nil{
			r=[]int64{}
		}
	}()
	for index:=0;index<len(s) ;index++  {
		i,err=strconv.ParseInt(s[index],10,64)
		if err!=nil{
			break
		}
		r = append(r,i)
	}
	return
}
func ParseToFloat64(s []string)(r []float64,err error) {
	var i float64
	defer func() {
		if err!=nil{
			r=[]float64{}
		}
	}()
	for index:=0;index<len(s) ;index++  {
		i,err=strconv.ParseFloat(s[index],64)
		if err!=nil{
			break
		}
		r = append(r,i)
	}
	return
}

func ConverseToInt64(value reflect.Value) (result int64) {
	//由于Uint64中支持的数据值大于int64,故存在数据值丢失的情况,数值超int64的范围,需自主定义函数进行处理
	//为了通用性,未兼容Uint64等极大值的数据类型
	obj:=value.Interface()
	if i, ok := obj.(int); ok {
		result=int64(i)
		return
	}
	if i, ok := obj.(uint); ok {
		result=int64(i)
		return
	}
	if i, ok := obj.(int8); ok {
		result=int64(i)
		return
	}
	if i, ok := obj.(uint8); ok {
		result=int64(i)
		return
	}
	if i, ok := obj.(int16); ok {
		result=int64(i)
		return
	}
	if i, ok := obj.(uint16); ok {
		result=int64(i)
		return
	}
	if i, ok := obj.(uint32); ok {
		result=int64(i)
		return
	}
	if i, ok := obj.(int32); ok {
		result=int64(i)
		return
	}
	if i, ok := obj.(int64); ok {
		result=i
		return
	}
	if i, ok := obj.(uint64); ok {
		result=int64(i)
		return
	}
	result=math.MaxInt64
	return
}

func ConverseToFloat64(value reflect.Value) (result float64) {
	//将数据类型转化为 int64方便进行数据值的比较
	obj:=value.Interface()
	if i, ok := obj.(float32); ok {
		result=float64(i)
		return
	}
	if i, ok := obj.(float64); ok {
		result=i
		return
	}
	result=math.MaxFloat64
	return
}
