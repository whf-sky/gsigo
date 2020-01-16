package config

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func NewIni() *Ini {
	return &Ini{}
}

//ini配置文件解析
type Ini struct {
	configs map[string]map[string]interface{}
}

//解析ini文件
func (r *Ini) resolving (filename string) error {
	// 打开文件
	file, err := os.Open(filename)
	// 文件找不到，返回空
	if err != nil {
		return err
	}
	// 在函数结束时，关闭文件
	defer file.Close()
	// 使用读取器读取文件
	reader := bufio.NewReader(file)

	// 当前读取的段的名字
	var section string

	for {
		// 读取文件的一行
		linestr, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		// 切掉行的左右两边的空白字符
		linestr = strings.TrimSpace(linestr)
		// 忽略空行
		if linestr == "" {
			continue
		}
		// 忽略注释
		if linestr[0] == ';' {
			continue
		}
		// 行首和尾巴分别是方括号的，说明是段标记的起止符
		if linestr[0] == '[' && linestr[len(linestr)-1] == ']' {
			// 将段名取出
			section = linestr[1 : len(linestr)-1]
		}

		// 切开等号分割的键值对
		pair := strings.Split(linestr, "=")
		key := strings.TrimSpace(pair[0])
		val := strings.Join(pair[1:], "=")
		val = strings.TrimSpace(val)
		//值是否使用了引号
		if val != "" &&  ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
			val = val[1 : len(val)-1]
		}
		//当存储的值为空初始化时
		if r.configs == nil {
			r.configs = map[string]map[string]interface{}{}
		}
		//如果section没有初始化值进行初始化
		sectionVal, ok := (r.configs)[section];
		if !ok {
			sectionVal = map[string]interface{}{}
		}
		//处理非数组数据
		if key[len(key)-2:] != "[]" {
			sectionVal[key] = val
			(r.configs)[section] = sectionVal
			continue
		}
		//处理数组数据
		key = key[0:len(key)-2]
		if (r.configs)[section][key] == nil {
			sectionVal[key] = []string{val}
		} else {
			sectionVal[key] = append(sectionVal[key].([]string), val)
		}
		(r.configs)[section] = sectionVal
	}
	return nil
}

//合并数据到map建为空的元素上
func (r *Ini) merge (sections ...string) map[string]interface{} {
	common, ok := r.configs[""]
	if !ok {
		common = map[string]interface{}{}
	}
	if len(sections) == 0 {
		return common
	}
	section, ok :=r.configs[sections[0]]
	if !ok {
		return common
	}
	for key, val := range section  {
		common[key] = val
	}
	return common
}

//给int类型赋值
func (r *Ini) toInt (val interface{}, field reflect.Value) {
	valstr := val.(string)
	vint, _ := strconv.ParseInt(valstr, 10, 0)
	field.SetInt(vint)
}

//给正整数类型赋值
func (r *Ini) toUint (val interface{}, field reflect.Value) {
	valstr := val.(string)
	vint, _ := strconv.ParseUint(valstr, 10, 0)
	field.SetUint(vint)
}

//给浮点型赋值
func (r *Ini) toFloat (val interface{}, field reflect.Value) {
	valstr := val.(string)
	vint, _ := strconv.ParseFloat(valstr, 0)
	field.SetFloat(vint)
}

//给布尔类型赋值
func (r *Ini) toBool (field reflect.Value, val interface{}){
	valstr := val.(string)
	valstr = strings.ToLower(valstr)
	switch valstr{
	case "true":
		field.SetBool(true)
	case "false":
		field.SetBool(false)
	default:
		field.SetBool(false)
	}
}

//给结构体类型赋值
func (r *Ini) toStruct (keys []string, field reflect.Value, val interface{}){
	if len(keys) == 0 {
		return
	}
	fields := r.fields(reflect.New(field.Type()).Type())
	field = field.Addr().Elem().FieldByName(fields[keys[0]])
	r.toType(keys[1:], field, val)
}

//给map类型赋值
func (r *Ini) toMap (keys []string, field reflect.Value, val interface{}){
	if len(keys) == 0 {
		return
	}
	//当map为空初始化map
	ftype := field.Type()
	if field.IsNil() {
		field.Set(reflect.MakeMap(ftype))
	}
	//获取map 键的类型
	ktype := ftype.Key()
	//获取map 元素类型
	etype := ftype.Elem()
	//实例map键并赋值
	k := reflect.New(ktype).Elem()
	r.toType(keys, k, keys[0])
	//实例map元素并赋值
	e := reflect.New(etype).Elem()
	r.toType(keys, e, val)
	//给map赋值
	field.SetMapIndex(k, e)
}

//给切片类型赋值
func (r *Ini) toSlice (keys []string, field reflect.Value, slice interface{}){
	ftype := field.Type()
	//当切片为空时初始化切片
	sliceStr := slice.([]string)
	if field.IsNil() {
		len := len(sliceStr)
		field.Set(reflect.MakeSlice(ftype, len, len))
	}
	//切片元素类型
	etype := ftype.Elem()
	for index, elem := range sliceStr {
		//给切片元素赋值
		e := reflect.New(etype).Elem()
		r.toType(keys, e, elem)
		//给切片赋值
		field.Index(index).Set(e)
	}
}

//给对应的类型赋值
func (r *Ini) toType (keys []string, field reflect.Value, val interface{})  {
	switch field.Kind() {
	case reflect.Struct:
		r.toStruct(keys, field, val)
	case reflect.Map:
		r.toMap(keys, field, val)
	case reflect.Slice, reflect.Array:
		r.toSlice(keys, field, val)
	case reflect.Bool:
		r.toBool(field, val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r.toInt(val, field)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r.toUint(val, field)
	case reflect.Float32, reflect.Float64:
		r.toFloat(val, field)
	case reflect.String:
		field.SetString(val.(string))
	default:
	}
}

//获取结构体字段
func (r *Ini) fields(rtype reflect.Type) map[string]string {
	elem := rtype.Elem()
	fields := map[string]string{}
	for i := 0; i < elem.NumField(); i++ {
		//设置标签与字段的映射
		field := elem.Field(i)
		fields[field.Tag.Get("ini")] = field.Name
	}
	return fields
}

//解析map
func (r *Ini) parseMap(out reflect.Value, ftype reflect.Type) {
	//map 键类型
	keyType := ftype.Key()
	//map 元素类型
	elemType := ftype.Elem()
	//根据元素类型进行分别处理
	switch elemType.Kind() {
	case reflect.Struct:
		//获取结构体的标签与字段对应关系
		fields := r.fields(reflect.New(elemType).Type())
		for name, section := range r.configs  {
			//分别实例键和元素
			keyVal := reflect.New(keyType).Elem()
			elemVal := reflect.New(elemType).Elem()
			//对键进行赋值
			r.toType([]string{name}, keyVal, name)
			for keystr, val := range section {
				//获取多个键
				keys := strings.Split(keystr, ".")
				//获取结构体字段
				field := elemVal.FieldByName(fields[keys[0]])
				if field.IsValid() {
					//对map 元素进行赋值
					r.toType(keys[1:], field, val)
				}
			}
			//对map进行赋值
			out.Elem().SetMapIndex(keyVal, elemVal)
		}
	case reflect.Map:
		//存在已赋值的键
		existKey := map[string]string{}
		var keyVal ,elemVal reflect.Value
		for keystr, val := range r.configs[""] {
			//获取多个键
			keys := strings.Split(keystr, ".")
			//键与元素没有实例时进行实例，并且只实例一次
			if _, ok := existKey[keys[0]];!ok {
				keyVal = reflect.New(keyType).Elem()
				elemVal = reflect.New(elemType).Elem()
				existKey[keys[0]] = ""
			}
			//分别对map的键和元素进行赋值
			r.toType(nil, keyVal, keys[0])
			r.toType(keys[1:], elemVal, val)
			//对map进行赋值
			out.Elem().SetMapIndex(keyVal, elemVal)
		}
	default:
		//标量或者切片数组不存在section都基于此处理
		for key, val := range r.configs[""] {
			//对map键赋值
			keyVal := reflect.New(keyType).Elem()
			r.toType(nil, keyVal, key)
			//对map元素赋值
			elemVal := reflect.New(elemType).Elem()
			r.toType(nil, elemVal, val)
			//对map赋值
			out.Elem().SetMapIndex(keyVal, elemVal)
		}

	}
}

func(r *Ini) Read(out interface{}, filename string) (err error) {
	//解析ini文件
	err = r.resolving(filename)
	if err != nil {
		panic(err)
	}
	//反射out，类型必须为指针类型
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Ptr {
		panic("Config is pointer to interface")
	}
	//根据数据类型分别处理
	switch outVal.Type().Elem().Kind() {
	case reflect.Struct:
		//获取结构体标签字段映射
		fields := r.fields(outVal.Type())
		//对结构体进行处理
		elem := outVal.Elem()
		for skey, section := range r.configs  {
			field := elem.FieldByName(fields[skey])
			//处理keys
			for keystr, val := range section {
				keys := strings.Split(keystr, ".")
				r.toType(keys, field, val)
			}
		}
	case reflect.Map:
		//如果map为空并对其初始化值
		ftype := outVal.Type().Elem()
		if outVal.Elem().IsNil() {
			outVal.Elem().Set(reflect.MakeMap(ftype))
		}
		r.parseMap(outVal, ftype)
	}
	return
}

// 根据文件名，段名，键名获取ini的值
func (r *Ini) ReadMerge(out interface{}, filename string, sections ...string) error {
	//解析ini文件
	err := r.resolving(filename)
	if err != nil {
		return err
	}
	//反射out，类型必须为指针类型
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Ptr &&  outVal.Type().Elem().Kind() == reflect.Struct{
		return errors.New("Config is struct pointer to interface ")
	}
	//合并数据
	section := r.merge(sections...)
	//获取结构体标签字段映射
	fields := r.fields(outVal.Type())
	//对结构体进行处理
	elem := outVal.Elem()
	for keystr, val := range section  {
		keys := strings.Split(keystr, ".")
		field := elem.FieldByName(fields[keys[0]])
		r.toType(keys[1:], field, val)
	}
	return nil
}
