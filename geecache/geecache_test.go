package geecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}
/**
匿名函数转换为GetterFunc类型
GetterFunc实现了Getter接口
f通过调用Get方法调用自身
 */
func TestGetter(t *testing.T) {
	var f GetterFunc = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	except := []byte("key")
	if v,_ := f.Get("key"); !reflect.DeepEqual(except, v) {
		t.Errorf("callback failed")
	}
}

func TestGroup(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key ", key)
		if v,ok := db[key]; ok {
			if _,b := loadCounts[key]; !b {
				loadCounts[key] = 0
			}
			loadCounts[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s key not exists", key)
	}))

	for k,v := range db {
		view,err := gee.Get(k)
		if err != nil || view.String() != v {
			t.Fatalf("failed to get value of %s. error: %s\n", k, err)
		}
		if _,err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss\n", k)
		}
	}

	if view,err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}