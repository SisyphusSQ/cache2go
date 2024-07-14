/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2012, Radu Ioan Fericean
 *                   2013-2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package cache2go

import (
	"fmt"
	"sync"
	"time"
)

var (
	cache = make(map[string]*CacheTable)
	mutex sync.RWMutex
)

// Cache returns the existing cache table with given name or creates a new one
// if the table does not exist yet.
func Cache(table string) *CacheTable {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		mutex.Lock()
		t, ok = cache[table]
		// Double check whether the table exists or not.
		if !ok {
			cleanup := 30 * time.Second

			t = &CacheTable{
				name:            table,
				items:           make(map[any]*CacheItem),
				cleanupInterval: cleanup,
				cleanupTicker:   time.NewTicker(cleanup),
				doneChan:        make(chan struct{}),
			}
			cache[table] = t
		}

		mutex.Unlock()
		go func() {
			defer t.cleanupTicker.Stop()
			for {
				select {
				case <-t.cleanupTicker.C:
					t.expirationCheck()
				case <-t.doneChan:
					return
				}
			}
		}()

	}

	return t
}

// NewCacheTable returns the existing cache table with given name or creates a new one
// if the table does not exist yet.
// add expireType when new
func NewCacheTable(table string, e expireType, cleanup time.Duration) *CacheTable {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		mutex.Lock()
		t, ok = cache[table]
		// Double check whether the table exists or not.
		if !ok {
			// bind expireType
			if e != OnCreate && e != OnLastAccess {
				panic(fmt.Sprintf("cache2go: invalid expireType: %s", e))
			}

			if cleanup == time.Duration(0) {
				cleanup = 30 * time.Second
			}

			t = &CacheTable{
				name:            table,
				items:           make(map[any]*CacheItem),
				expireType:      e,
				cleanupInterval: cleanup,
				cleanupTicker:   time.NewTicker(cleanup),
				doneChan:        make(chan struct{}),
			}
			cache[table] = t
		}

		mutex.Unlock()
		go func() {
			defer t.cleanupTicker.Stop()
			for {
				select {
				case <-t.cleanupTicker.C:
					t.expirationCheck()
				case <-t.doneChan:
					return
				}
			}
		}()
	}

	return t
}
