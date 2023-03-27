package main

import (
	"database/sql"
	"fmt"
	"log"

	lru "github.com/hashicorp/golang-lru/v2"
	_ "github.com/mattn/go-sqlite3"
)

type MDBCache struct {
	Cache *lru.Cache[string, *sql.DB]
}

func LruInstantiate(cassUrl string, cacheSize int) MDBCache {

	cache, _ := lru.NewWithEvict(cacheSize, func(key string, value *sql.DB) {
		log.Println("Evicted key : ", key)
		value.Close()
	})

	return MDBCache{Cache: cache}
}

func (c *MDBCache) GetMDB(key string) *sql.DB {

	if c.Cache.Contains(key) {
		mdb, ok := c.Cache.Get(key)
		if !ok {
			log.Fatal("Error Get and Return MDB from cache")
		}
		return mdb
	}

	destDb, _ := sql.Open("sqlite3", "./data/"+key+"/customer.db")

	row, err := destDb.Query("PRAGMA journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	if row.Next() {
		var res string
		row.Scan(&res)
		fmt.Println("Wal status " + res)
	}

	c.Cache.Add(key, destDb)

	return destDb

}
