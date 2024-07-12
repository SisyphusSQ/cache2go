package main

import (
	"fmt"
	"time"

	"github.com/SisyphusSQ/cache2go"
)

var fn = func(key interface{}) {
	fmt.Println("About to expire:", key.(string))
}

func main() {
	cache := cache2go.NewCacheTable("test", cache2go.OnCreate, 20*time.Second)
	defer cache.Close()

	cache.SetAddedItemCallback(func(entry *cache2go.CacheItem) {
		fmt.Println("Added Callback 1:", entry.Key(), entry.Data(), entry.CreatedOn())
	})

	cache.SetUpdatedItemCallback(func(entry *cache2go.CacheItem) {
		fmt.Println("Updated Callback 2:", entry.Key(), entry.Data(), entry.CreatedOn())
	})

	cache.SetAboutToDeleteItemCallback(func(entry *cache2go.CacheItem) {
		fmt.Println("Deleting:", entry.Key(), entry.Data(), entry.CreatedOn())
	})

	//cache.Add("someKey", 10*time.Second, "This is a test!")
	cache.AddExpireFunc("someKey", "This is a test!", 10*time.Second, fn)
	res, err := cache.Value("someKey")
	if err == nil {
		fmt.Println("Found value in cache:", res.Data())
	} else {
		fmt.Println("Error retrieving value from cache:", err)
	}

	//res.SetAboutToExpireCallback(func(key interface{}) {
	//	fmt.Println("About to expire:", key.(string))
	//})

	res, err = cache.Update("someKey", "This is a test1111!", 10*time.Second, false)

	time.Sleep(30 * time.Second)

	fmt.Println("end...")
}
