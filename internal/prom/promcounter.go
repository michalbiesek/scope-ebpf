package prom

import (
	"fmt"
)

// PromMetricCounter describes the Prometheus metrics counter
type PromMetricCounter struct {
	Name    string
	Valobj  string
	Counter int
	Unit    string
}

// Returns string value of prometheus metrics
func (p *PromMetricCounter) Add() {
	p.Counter += 1
}

// Returns string value of prometheus metrics
func (p *PromMetricCounter) String() string {
	return fmt.Sprintf("# TYPE %s counter\n%s{%s,unit=\"%s\"} %d", p.Name, p.Name, p.Valobj, p.Unit, p.Counter)
}
