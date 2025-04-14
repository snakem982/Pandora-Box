package cache

import (
	"fmt"
	"github.com/metacubex/bbolt"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/pandora/pkg/constant"
	"os"
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

func Dump(dstDBPath string) error {
	_, err := os.Stat(dstDBPath)
	if !os.IsNotExist(err) {
		_ = os.Remove(dstDBPath)
	}

	// 创建目标数据库
	dstDB, err := bbolt.Open(dstDBPath, 0600, nil)
	if err != nil {
		return err
	}
	defer func(dstDB *bbolt.DB) {
		err := dstDB.Close()
		if err != nil {

		}
	}(dstDB)

	return dstDB.Batch(func(tx *bbolt.Tx) error {
		newBucket, err := tx.CreateBucketIfNotExists(BName)
		if err != nil {
			log.Warnln("[DumpFile] can't create bucket: %s", err.Error())
			return fmt.Errorf("create bucket: %v", err)
		}

		return BDb.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket(BName)
			return b.ForEach(func(k, v []byte) error {
				key := string(k)
				if key == constant.SecretKey {
					return nil
				}

				if key == constant.QuitSignal {
					return nil
				}

				if !strings.HasPrefix(key, constant.RealIpHeader) {
					return newBucket.Put(k, v)
				}

				return nil
			})
		})
	})
}

func Recovery(srcDBPath string) error {
	_, err := os.Stat(srcDBPath)
	if os.IsNotExist(err) {
		return err
	}

	// 打开源数据库
	srcDB, err := bbolt.Open(srcDBPath, 0600, nil)
	if err != nil {
		return err
	}
	defer func(srcDB *bbolt.DB) {
		err := srcDB.Close()
		if err != nil {

		}
	}(srcDB)

	return srcDB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BName)
		return b.ForEach(func(k, v []byte) error {
			return Put(string(k), v)
		})
	})
}
