package main

import (
	"fmt"
	"strings"
)

type builder struct {
	sql *strings.Builder
	q   *queue
}

func createBuilder() *builder {
	return &builder{sql: &strings.Builder{}, q: createQueue()}
}

func (b *builder) buildSql() {
	b.q.buildSql(b.sql)
}

type queue struct {
	head *node
	size int
}

type node struct {
	value string
	next  *node
}

func createQueue() *queue {
	return &queue{}
}

func (q *queue) add(v string) {
	n := &node{value: v}

	if q.head == nil {
		q.head = n
		q.size++
		return
	}

	tail := getTail(q.head, n.value)
	if tail == nil {
		return
	}
	tail.next = n
	q.size++
}

func (q *queue) get() string {
	if q.head == nil {
		return ""
	}

	n := q.head
	q.head = q.head.next
	q.size--
	return n.value
}

func getTail(n *node, v string) *node {
	if n.value == v {
		return nil
	}
	if n.next != nil {
		return getTail(n.next, v)
	}
	return n
}

func (q *queue) buildSql(sql *strings.Builder) {
	stringNode(sql, q)
}

func stringNode(sql *strings.Builder, q *queue) {
	if q.head != nil {
		sql.WriteString(q.head.value)
		q.head = q.head.next
		q.size--
		stringNode(sql, q)
		return
	}

	sql.WriteString(";")
}

func main() {
	builder := createBuilder()

	for v := range 10 {
		builder.q.add(fmt.Sprint(v))
		builder.q.add(fmt.Sprint(v))
	}

	fmt.Println(builder.q)
	builder.buildSql()
	fmt.Println(builder.sql.String())
	fmt.Println(builder.q)
}
