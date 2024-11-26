package main

import "sync"

// This is a simple sync.Map wrapper for the sample database
type myMap struct {
	sync.Map
}

// since Golang doesn't implement a synchronizedmap by default,
// the sync.Map that is the database is wrapped with these
func (m *myMap) Insert(key string, value *Receipt) {
	m.Store(key, value)
}

// since Golang doesn't implement a synchronizedmap by default,
// the sync.Map that is the database is wrapped with these
func (m *myMap) Select(key string) *Receipt {
	found_receipt, _ := m.Load(key)
	r, _ := found_receipt.(*Receipt)
	if r != nil {
		return r
	} else {
		return nil
	}
}
