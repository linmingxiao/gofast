// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"github.com/qinchende/gofast/fst/render"
	"net/http"
	"net/url"
	"sync"
)

// Context is the most important part of GoFast. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	*GFResponse               // response (请求前置拦截器 要用到的上下文)
	ReqRaw      *http.Request // request
	matchRst    matchResult   // 路由匹配结果

	Pms        map[string]string // 所有Request参数的map（[Params] ? + queryCache + formCache）
	Params     Params            // : 或 * 对应的参数
	queryCache url.Values        // param query result from c.ReqRaw.URL.Query()
	formCache  url.Values        // the parsed form data from POST, PATCH, or PUT body parameters.

	// Session数据，这里不规定Session的载体，可以自定义
	Sess *CtxSession
	// 设置成 true ，将中断后面的所有handlers
	aborted bool
	// render.Render 对象
	PRender *render.Render // render 对象
	PCode   *int           // status code

	// -----------------------------
	// This mutex protect Keys map
	mu sync.RWMutex
	// Accepted defines a list of manually accepted formats for content negotiation.
	Accepted []string
	// SameSite allows a server to define a cookie attribute making it impossible for
	// the browser to send this cookie along with cross-site requests.
	sameSite http.SameSite
	// Keys is a key/value pair exclusively for the context of each request.
	// 上下文传值
	Keys map[string]interface{}
}

/************************************/
/********** context creation ********/
/************************************/

func (c *Context) reset() {
	c.Keys = nil
	c.Sess = nil
	c.Accepted = nil

	c.Pms = nil
	c.Params = c.Params[0:0]
	c.queryCache = nil
	c.formCache = nil
	c.aborted = false

	// add by sdx 2021.01.06
	c.matchRst.ptrNode = nil
	c.matchRst.params = &c.Params
	c.matchRst.rts = false
	c.matchRst.allowRTS = c.gftApp.RedirectTrailingSlash
}

// 如果在当前请求上下文中需要新建goroutine，那么新的 goroutine 中必须要用 copy 后的 Context
// Copy returns a copy of the current context that can be safely used outside the request's scope.
// This has to be used when the context has to be passed to a goroutine.
func (c *Context) Copy() *Context {
	cp := Context{
		GFResponse: c.GFResponse,
		ReqRaw:     c.ReqRaw,
		Params:     c.Params,
		matchRst:   c.matchRst,
		Pms:        c.Pms,
		Sess:       c.Sess,
		aborted:    c.aborted,
	}
	cp.ResWrap.ResponseWriter = nil

	cp.Keys = map[string]interface{}{}
	for k, v := range c.Keys {
		cp.Keys[k] = v
	}
	paramCopy := make([]Param, len(cp.Params))
	copy(paramCopy, cp.Params)
	cp.Params = paramCopy
	return &cp
}
