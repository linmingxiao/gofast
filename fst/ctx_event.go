// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a BSD-style license
package fst

// TODO： 底层大致有两种设计思路，目前采用第一种方案， 没有优化
// 方案1. 依次执行分组和节点自己的事件中间件函数
func (c *Context) execHandlers(ptrMini *radixMiniNode) {
	it := fstMem.hdsMiniNodes[ptrMini.hdsItemIdx]
	gp := fstMem.hdsMiniNodes[ptrMini.hdsGroupIdx]

	// 1.valid
	for gp.validLen > 0 {
		fstMem.hdsList[gp.validIdx](c)
		gp.validLen--
		gp.validIdx++
	}
	for it.validLen > 0 {
		fstMem.hdsList[it.validIdx](c)
		it.validLen--
		it.validIdx++
	}

	// 2.before
	for gp.beforeLen > 0 {
		fstMem.hdsList[gp.beforeIdx](c)
		gp.beforeLen--
		gp.beforeIdx++
	}
	for it.beforeLen > 0 {
		fstMem.hdsList[it.beforeIdx](c)
		it.beforeLen--
		it.beforeIdx++
	}

	// 3.handler
	for it.hdsLen > 0 {
		fstMem.hdsList[it.hdsIdx](c)
		it.hdsLen--
		it.hdsIdx++
	}

	// 4.after
	for it.afterLen > 0 {
		fstMem.hdsList[it.afterIdx](c)
		it.afterLen--
		it.afterIdx++
	}
	for gp.afterLen > 0 {
		fstMem.hdsList[gp.afterIdx](c)
		gp.afterLen--
		gp.afterIdx++
	}

	// 5.send
	for it.sendLen > 0 {
		fstMem.hdsList[it.sendIdx](c)
		it.sendLen--
		it.sendIdx++
	}
	for gp.sendLen > 0 {
		fstMem.hdsList[gp.sendIdx](c)
		gp.sendLen--
		gp.sendIdx++
	}

	// 6.afterSend
	for it.afterSendLen > 0 {
		fstMem.hdsList[it.afterSendIdx](c)
		it.afterSendLen--
		it.afterSendIdx++
	}
	for gp.afterSendLen > 0 {
		fstMem.hdsList[gp.afterSendIdx](c)
		gp.afterSendLen--
		gp.afterSendIdx++
	}
}

// 方案2. 基于已经将所有的事件函数组织成了一个有序的索引数组。只需要一次循环就执行所有的中间件函数
// 这种实现其实不现实，不同类型的事件是在框架封装过程中分开执行的
func (c *Context) execHandlersMini(ptrMini *radixMiniNode) {
	it := fstMem.hdsNodesPlan2[ptrMini.hdsItemIdx]

	for it.hdsLen > 0 {
		fstMem.hdsList[it.startIdx](c)
		it.hdsLen--
		it.startIdx++
	}
}
