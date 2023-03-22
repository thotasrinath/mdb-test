package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/gocql/gocql"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/mattn/go-sqlite3"
)

type MDBCache struct {
	Cache      *lru.Cache[string, *sql.DB]
	CqlSession *gocql.Session
}

func LruInstantiate(cassUrl string, cacheSize int) MDBCache {

	cluster := gocql.NewCluster(cassUrl)
	cluster.Keyspace = "test_db"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	//defer session.Close()

	cache, _ := lru.NewWithEvict(cacheSize, func(key string, value *sql.DB) {
		SerializeToCassandra(key, value, session)
	})

	return MDBCache{Cache: cache, CqlSession: session}
}
func SerializeToCassandra(key string, value *sql.DB, session *gocql.Session) {

	log.Println("Evicted key : ", key)
	srcConn, err := value.Conn(context.Background())
	if err != nil {
		log.Fatal("Failed to get connection to source database:", err)
	}
	defer srcConn.Close()

	var serialized []byte
	if err := srcConn.Raw(func(raw interface{}) error {
		var err error
		serialized, err = raw.(*sqlite3.SQLiteConn).Serialize("")
		return err
	}); err != nil {
		log.Fatal("Failed to serialize source database:", err)
	}
	srcConn.Close()

	ctx := context.Background()

	if err := session.Query(`INSERT INTO teacher (id, details) VALUES (?, ?)`,
		key, serialized).WithContext(ctx).Exec(); err != nil {
		log.Fatal(err)
	}

	log.Println("Stored MDB to cassandra for key : ", key)

	value.Close()
}

func (c *MDBCache) GetMDB(key string) *sql.DB {

	if c.Cache.Contains(key) {
		mdb, ok := c.Cache.Get(key)
		if !ok {
			log.Fatal("Error Get and Return MDB from cache")
		}
		return mdb
	}

	var deSerialized []byte

	err := c.CqlSession.Query("SELECT details FROM teacher WHERE id = '" + key + "'").Scan(&deSerialized)
	if err != nil {
		log.Printf("select from teacher failed for key '%v', error '%v'", key, err)
	}

	//fmt.Println(deSerialized)

	destDb, _ := sql.Open("sqlite3", ":memory:")
	//defer destDb.Close()

	if len(deSerialized) > 0 {
		destConn, err := destDb.Conn(context.Background())
		if err != nil {
			log.Fatal("Failed to get connection to destination database:", err)
		}
		defer destConn.Close()

		if err := destConn.Raw(func(raw interface{}) error {
			return raw.(*sqlite3.SQLiteConn).Deserialize(deSerialized, "")
		}); err != nil {
			log.Fatal("Failed to deserialize source database:", err)
		}
		destConn.Close()
	}

	c.Cache.Add(key, destDb)

	return destDb

}
