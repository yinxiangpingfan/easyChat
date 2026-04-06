package gorm_plugin

import (
	"easyChat/pkg/log"
	"time"

	"gorm.io/gorm"
)

const startTimeKey = "slow_query:start_time"

type SlowQueryPlugin struct {
	Threshold time.Duration // 超过此阈值记录慢查询日志
}

func NewSlowQueryPlugin(threshold time.Duration) *SlowQueryPlugin {
	return &SlowQueryPlugin{Threshold: threshold}
}

func (p *SlowQueryPlugin) Name() string { return "slow_query" }

func (p *SlowQueryPlugin) Initialize(db *gorm.DB) error {
	before := func(db *gorm.DB) {
		db.Set(startTimeKey, time.Now())
	}
	after := func(db *gorm.DB) {
		start, ok := db.Get(startTimeKey)
		if !ok {
			return
		}
		elapsed := time.Since(start.(time.Time))
		if elapsed >= p.Threshold {
			logger := log.FromContext(db.Statement.Context)
			logger.Warnf("慢查询: sql=%s, 耗时=%v", db.Statement.SQL.String(), elapsed)
		}
	}

	db.Callback().Create().Before("gorm:create").Register("slow_query:before_create", before)
	db.Callback().Create().After("gorm:create").Register("slow_query:after_create", after)
	db.Callback().Query().Before("gorm:query").Register("slow_query:before_query", before)
	db.Callback().Query().After("gorm:query").Register("slow_query:after_query", after)
	db.Callback().Update().Before("gorm:update").Register("slow_query:before_update", before)
	db.Callback().Update().After("gorm:update").Register("slow_query:after_update", after)
	db.Callback().Delete().Before("gorm:delete").Register("slow_query:before_delete", before)
	db.Callback().Delete().After("gorm:delete").Register("slow_query:after_delete", after)
	return nil
}
