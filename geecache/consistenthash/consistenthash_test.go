package consistenthash

import (
	"log"
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})
	hash.Add("6", "4", "2")
	log.Printf("map.keys:%+v,map.hashMap:%+v", hash.keys, hash.hashMap)
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	hash.Add("8")
	log.Printf("map.keys:%+v,map.hashMap:%+v", hash.keys, hash.hashMap)
	testCases["27"] = "8"
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
