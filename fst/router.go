// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 绑定在 RouterGroup 和 RouteItem 上的 不同事件处理函数数组
// RouterGroup 上的事件处理函数 最后需要作用在 RouteItem 上才会有实际的意义
// 事件要尽量少一些，每个路由节点都要分配一个对象
// TODO: 此结构占用空间还是比较大的，可以考虑释放。
type routeEvents struct {
	// 下面的事件类型，按照执行顺序排列
	ePreValidHds  []uint16
	eBeforeHds    []uint16
	eHds          []uint16
	eAfterHds     []uint16
	ePreSendHds   []uint16
	eAfterSendHds []uint16
}

type RouterGroup struct {
	routeEvents              // 直接作用于本节点的事件可能为空
	combEvents   routeEvents // 合并父节点的分组事件，routeEvents可能为空，但是combEvents几乎不会为空
	myApp        *GoFast
	prefix       string
	children     []*RouterGroup
	hdsIdx       int16  // 记录当前分组 对应新事件数组中的起始位置索引
	selfHdsLen   uint16 // 记录当前分组中一共加入的处理函数个数（仅限于本分组加入的事件，不包含合并上级分组的）
	parentHdsLen uint16 // 记录所属上级分组的所有处理函数个数（仅包含上级分组，不含本分组的事件个数）
}

type RouteItem struct {
	group       *RouterGroup // router group
	method      string       // httpMethod
	fullPath    string       // 路由的完整路径
	routeEvents              // all handlers
	routerIdx   uint16       // 此路由在路由数组中的索引值
}

// 每一种事件类型需要占用3个字节(开始索引2字节 + 长度1字节(长度最大255))
// 这里抽象出N种事件类型，应该够用了，这样每个路由节点占用3*N字节空间，64位机器1字长是8字节
// RouterGroup 和 RouteItem 都用这一组数据结构记录事件处理函数
type handlersNode struct {
	hdsIdxChain []uint16 // 执行链的索引数组

	validIdx     uint16
	beforeIdx    uint16
	hdsIdx       uint16
	afterIdx     uint16
	preSendIdx   uint16
	afterSendIdx uint16

	validLen     uint8
	beforeLen    uint8
	afterLen     uint8
	hdsLen       uint8
	preSendLen   uint8
	afterSendLen uint8
}

//// ++++++++++++++++++++++++++++++++++++++++++++++
//// 第二种方案（暂时不用）
//// 将某个路由节点的所有处理函数按顺序全部排序成数组，请求匹配到路由节点之后直接执行这里的队列即可
//// 当节点多的时候这种方式相对第一种占用更多内存。
//type handlersNodePlan2 struct {
//	startIdx uint16 // 2字节
//	hdsLen   uint8  // 1字节
//}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 路由节点支持自定义扩展。可以自定义配置项，和配置项集合。
type RouteConfigs interface {
	Reordering(*GoFast, uint16)
}

type RouteConfig interface {
	AddToList(uint16)
}

type RouteIndex struct {
	Idx uint16 // 此路由在路由数组中的索引值
}
