package cache

import (
	"encoding/json"
	"fmt"
	"github.com/metacubex/bbolt"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/pkg/constant"
	"github.com/snakem982/pandora-box/pkg/utils"
	"os"
	"reflect"
	"strings"
)

var BName = []byte("Pandora-Box")
var BDb *bbolt.DB

func GetDBInstance() *bbolt.DB {
	path := utils.GetUserHomeDir(constant.DefaultServerDB)

	var err error
	BDb, err = bbolt.Open(path, os.ModePerm, nil)
	if err != nil {
		panic(err)
	}
	_ = utils.SetPermissions(path)
	_ = BDb.Batch(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(BName)
		if err != nil {
			log.Warnln("[CacheFile] can't create bucket: %s", err.Error())
			return fmt.Errorf("create bucket: %v", err)
		}
		return nil
	})

	return BDb
}

func GetMetaDB() {
	path := utils.GetUserHomeDir("cache.db")

	metaDB, err := bbolt.Open(path, os.ModePerm, nil)
	if err != nil {
		_ = utils.SetPermissions(path)
		return
	}
	_ = metaDB.Close()
	_ = utils.SetPermissions(path)
}

// Put 将任意类型的键值对存储到数据库中
func Put(key string, value interface{}) error {
	return BDb.Update(func(tx *bbolt.Tx) error {
		// 将任意类型的值编码为 JSON 格式的 []byte
		encodedValue, err := json.Marshal(value)
		if err != nil {
			return err // 如果编码失败，则返回错误
		}
		b := tx.Bucket(BName)                  // 获取数据库中指定的桶
		err = b.Put([]byte(key), encodedValue) // 将键值对存储到桶中
		return err                             // 返回可能的错误
	})
}

// Get 从数据库中获取键对应的值并解码为指定类型
func Get(key string, value interface{}) error {
	return BDb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BName)              // 获取指定的桶
		encodedValue := b.Get([]byte(key)) // 根据键获取值
		if encodedValue == nil {
			return fmt.Errorf("key not found") // 若键不存在，返回错误
		}
		// 将 JSON 格式的 []byte 解码为原始类型
		err := json.Unmarshal(encodedValue, value)
		return err // 返回可能的解码错误
	})
}

// GetList 获取与指定前缀键相关的所有值，并解码为指定类型的切片
func GetList(key string, values interface{}) error {
	// 检查传入的 values 是否为切片的指针
	v := reflect.ValueOf(values)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("values must be a pointer to a slice")
	}

	// 临时存储 JSON 格式的字节切片
	tempValues := make([][]byte, 0)

	err := BDb.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BName) // 获取桶
		return b.ForEach(func(k, v []byte) error {
			if strings.HasPrefix(string(k), key) {
				tempValues = append(tempValues, v)
			}
			return nil
		})
	})
	if err != nil {
		return err // 返回可能的错误
	}

	// 解码 JSON 数据到目标切片
	sliceValue := v.Elem()
	for _, encodedValue := range tempValues {
		elem := reflect.New(sliceValue.Type().Elem()).Interface()
		if err := json.Unmarshal(encodedValue, elem); err != nil {
			return err // 解码错误
		}
		sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(elem).Elem()))
	}

	return nil
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

	return BDb.Batch(func(tx *bbolt.Tx) error {
		newBucket, err := tx.CreateBucketIfNotExists(BName)
		if err != nil {
			log.Warnln("[DumpFile] can't create bucket: %s", err.Error())
			return fmt.Errorf("create bucket: %v", err)
		}

		return srcDB.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket(BName)
			return b.ForEach(func(k, v []byte) error {
				return newBucket.Put(k, v)
			})
		})
	})
}

func Close() {
	if BDb != nil {
		_ = BDb.Close()
	}
}
