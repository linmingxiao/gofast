package sqlx

import (
	"context"
	msql "database/sql"
	"github.com/qinchende/gofast/logx"
	"github.com/qinchende/gofast/skill/timex"
	"time"
)

// 执行超过500ms的语句需要优化分析，我们先打印出慢日志
const slowThreshold = time.Millisecond * 500

func (conn *MysqlORM) CloneWithCtx(ctx context.Context) *MysqlORM {
	newConn := *conn
	newConn.Ctx = ctx
	return &newConn
}

func (conn *MysqlORM) Exec(sql string, args ...any) msql.Result {
	return conn.ExecCtx(conn.Ctx, sql, args...)
}

func (conn *MysqlORM) ExecCtx(ctx context.Context, sql string, args ...any) msql.Result {
	logx.Debug(sql)

	var result msql.Result
	var err error
	startTime := timex.Now()
	if conn.tx == nil {
		result, err = conn.Writer.ExecContext(ctx, sql, args...)
	} else {
		result, err = conn.tx.ExecContext(ctx, sql, args...)
	}
	dur := timex.Since(startTime)
	if dur > slowThreshold {
		logx.SlowF("[SQL][%dms] exec: slow-call - %s", dur/time.Millisecond, sql)
		//} else {
		//	logx.InfoF("sql exec: %s", sql)
	}
	errPanic(err)
	return result
}

func (conn *MysqlORM) QuerySql(sql string, args ...any) *msql.Rows {
	return conn.QuerySqlCtx(conn.Ctx, sql, args...)
}

func (conn *MysqlORM) QuerySqlCtx(ctx context.Context, sql string, args ...any) *msql.Rows {
	logx.Debug(sql)

	var rows *msql.Rows
	var err error
	startTime := timex.Now()
	if conn.tx == nil {
		rows, err = conn.Reader.QueryContext(ctx, sql, args...)
	} else {
		rows, err = conn.tx.QueryContext(ctx, sql, args...)
	}
	dur := timex.Since(startTime)
	if dur > slowThreshold {
		logx.SlowF("[SQL][%dms] query: slow-call - %s", dur/time.Millisecond, sql)
		//} else {
		//	logx.InfoF("sql query: %s", sql)
	}
	errPanic(err)
	return rows
}

func (conn *MysqlORM) TransBegin() *MysqlORM {
	return conn.TransCtx(conn.Ctx)
}

func (conn *MysqlORM) TransCtx(ctx context.Context) *MysqlORM {
	tx, err := conn.Writer.BeginTx(ctx, nil)
	errPanic(err)
	return &MysqlORM{Attrs: conn.Attrs, Ctx: ctx, tx: tx, rdsNodes: conn.rdsNodes}
}

func (conn *MysqlORM) TransFunc(fn func(newConn *MysqlORM)) {
	conn.TransFuncCtx(conn.Ctx, fn)
}

func (conn *MysqlORM) TransFuncCtx(ctx context.Context, fn func(newConn *MysqlORM)) {
	tx, err := conn.Writer.BeginTx(ctx, nil)
	errPanic(err)

	nConn := MysqlORM{Attrs: conn.Attrs, Ctx: ctx, tx: tx, rdsNodes: conn.rdsNodes}
	defer nConn.TransEnd()
	fn(&nConn)
}

func (conn *MysqlORM) Commit() error {
	return conn.tx.Commit()
}

func (conn *MysqlORM) Rollback() error {
	return conn.tx.Rollback()
}

func (conn *MysqlORM) TransEnd() {
	var err error

	if pic := recover(); pic != nil {
		err = conn.Rollback()
	} else {
		err = conn.Commit()
	}
	if err != nil {
		logx.ErrorF("Terrible mistake. trans end error: %s", err)
	}
}
