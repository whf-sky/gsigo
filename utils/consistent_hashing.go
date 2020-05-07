package utils

import (
	"hash/crc32"
	"sort"
	"strconv"
)

func NewConsistentHashing(virtualNodeNum int) *ConsistentHashing {
	return &ConsistentHashing{
		locationsMap: map[string]string{},
		locations: []string{},
		virtualNodeNum: virtualNodeNum,
		nodes: map[string][]string{},
	}
}

//一致性hash
type ConsistentHashing struct {
	//落点集合映射
	locationsMap map[string]string
	//落点集合
	locations []string
	//虚拟节点数
	virtualNodeNum int
	//节点信息
	nodes map[string][]string
}

//设置虚拟节点数量
//num 节点数
func (c *ConsistentHashing) SetVirtualNodeNum (num int) {
	c.virtualNodeNum = num
}

//字符串转为数字字符串
//str 字符串
func (c *ConsistentHashing) crc32ToString (str string) string{
	return strconv.FormatInt(int64(crc32.ChecksumIEEE([]byte(str))), 10)
}

//添加节点
//node 节点
func (c *ConsistentHashing) AddNode (node string) {
	for i := 0;i < c.virtualNodeNum;i++ {
		tmp := c.crc32ToString(node+strconv.Itoa(i))
		c.locationsMap[tmp] = node
		if c.nodes[node] == nil {
			c.nodes[node] = []string{}
		}
		c.nodes[node] = append(c.nodes[node], tmp)
		sort.Strings(c.nodes[node])
	}
	c.setLocations()
}

//寻找字符串所在位置
//str 字符串
func (c *ConsistentHashing) GetLocation (str string) string {
	if c.locations == nil || len(c.locations) == 0 {
		return ""
	}
	position := c.crc32ToString(str)
	node := ""
	for index, val := range c.locations {
		if index == 0 {
			node = c.locationsMap[val]
		}
		index++
		if position <= val {
			node = c.locationsMap[val]
			break
		}
	}
	return node
}

//设置落点集合
func (c *ConsistentHashing) setLocations () {
	locations := []string{}
	for key, _ := range c.locationsMap  {
		locations = append(locations, key)
	}
	sort.Strings(locations)

	c.locations = locations
}

//删除一个节点
//node 节点
func (c *ConsistentHashing) DeleteNode (node string) {
	for _, n := range c.nodes[node] {
		delete(c.locationsMap, n)
	}
	delete(c.nodes, node)
	c.setLocations()
}
