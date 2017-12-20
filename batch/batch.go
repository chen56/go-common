package batch

import (
	"github.com/apex/log"
	"time"
)

type Batch struct {
	limit         int
	limitDuration time.Duration
	lastTickTime  time.Time
	total         int
	totalError    int
	totalOk       int
	totalIgnore   int
	items         []interface{}
	callback      Callback

}
type Callback func(items []interface{},batch *Batch)

func (this *Batch) tick() {
	if this.IsReachLimit() {
		if !this.IsItemsEmpty() {
			this.callback(this.items,this)
			this.items = []interface{}{}
		}
	}
	this.lastTickTime = time.Now()
}

func (this *Batch) IsItemsEmpty() bool {
	return len(this.items) == 0
}

func (this *Batch) IsReachLimit() bool {
	if this.limitDuration != maxDuration {
		if time.Now().Sub(this.lastTickTime) > this.limitDuration{
        	return true
		}
	}
	return this.total%this.limit == 0
}


func (this *Batch) TickOk(item interface{}) {
	this.items = append(this.items, item)
	this.total += 1
	this.totalOk += 1
	this.tick()
}

func (this *Batch) TickError() {
	this.total += 1
	this.totalError += 1
	this.tick()
}

func (this *Batch) TickIgnore() {
	this.total += 1
	this.totalIgnore += 1
	this.tick()
}

type Builder struct {
	limit         int
	limitDuration time.Duration
	callback Callback
}

const maxDuration time.Duration = 1<<63 - 1

func NewBuilder() *Builder {
	return &Builder{
		limit:         1,
		limitDuration: maxDuration,
	}
}

func (this *Builder) Limit(limit int) (self *Builder) {
	this.limit = limit
	return this
}

func (self *Builder) LimitDuration(duration time.Duration) *Builder {
	self.limitDuration = duration
	return self
}

func (self *Builder) Callback(callabck Callback) *Builder {
	self.callback = callabck
	return self
}

func (this *Builder) Build() (result *Batch) {
	return &Batch{
		lastTickTime:  time.Now(),
		callback:      this.callback,
		limit:         this.limit,
		limitDuration: this.limitDuration,
	}
}

func (this *Batch) Total() int {
	return this.total
}

func (this *Batch) TotalOk() int {
	return this.totalOk
}

func (this *Batch) TotalError() int {
	return this.totalError
}

func (this *Batch) TotalIgnore() int {
	return this.totalIgnore
}

func (this *Batch) FillLog(log *log.Entry) *log.Entry {
	return log.WithField("totalError", this.TotalError()).
		WithField("totalOk", this.TotalOk()).
		WithField("totalIgnore", this.TotalIgnore()).
		WithField("total", this.Total())
}
