package cache

import (
	"github.com/metacubex/bbolt"
	"strings"
)

var BName = []byte("Pandora-Box")
var BDb *bbolt.DB

// Put 将给定的键值对存储到数据库中
func Put(key string, value []byte) error {
	return BDb.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BName)            // 获取数据库中的指定桶
		err := b.Put([]byte(key), value) // 将键值对存储到桶中
		return err                       // 返回可能发生的错误
	})
}

// Get 函数接受一个字符串类型的键值(key)作为参数，并返回一个字节数组
func Get(key string) (value []byte) {
	_ = BDb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BName)      // 获取数据库对应的桶
		value = b.Get([]byte(key)) // 根据提供的键获取与之对应的值
		return nil
	})

	return value // 返回获取到的值
}

func GetList(key string) (values [][]byte) {
	values = make([][]byte, 0)
	_ = BDb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BName)
		return b.ForEach(func(k, v []byte) error {
			if strings.HasPrefix(string(k), key) {
				values = append(values, v)
			}
			return nil
		})
	})

	return values
}

func Delete(key string) error {
	return BDb.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BName)
		err := b.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	})
}

func DeleteList(m map[string]any) error {
	return BDb.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BName)
		return b.ForEach(func(k, v []byte) error {
			if _, find := m[string(k)]; find {
				_ = b.Delete(k)
			}
			return nil
		})
	})
}
