// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

type SessionKeeper interface {
	Get(string) any
	Set(string, any)
	SetKV(KV)
	Del(string)
	Save() error
	Expire(int32)
	SidIsNew() bool // SessionID is new?
	Sid() string    // SessionID
	Destroy()       // 销毁当前Session
}
