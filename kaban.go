package kaban

import (
	"bytes"
	//	"encoding/json"
	"fmt"
	//	"sort"
	"strconv"
	"sync"
	//	"sync/atomic"
)

type (
	Kaban struct {
		keyMap map[string]int
		chunk  []byte
	}
)

const (
	chunkSize  = 1024 * 1024
	intBase    = 36 // 36進数
	intBitSize = 64 // 64ビット整数
)

const (
	sepDead   = 0xFF // Dead value
	sepNull   = 0xFE // JSON null
	sepString = 0xFD // JSON string
	sepInt    = 0xFC // JSON int
	sepUint   = 0xFB // JSON int
	sepFloat  = 0xFA // JSON float
	sepBool   = 0xF9 // JSON bool
	// JSON null
	// JSON array
	// JSON object
	//sepAny = sepNull + sepString + sepInt + sepFloat +
	//	sepBool
)

var (
	kabanMtx = new(sync.RWMutex)
)

func New() *Kaban {
	k := &Kaban{
		keyMap: make(map[string]int),
		chunk:  make([]byte, 0, chunkSize),
	}
	return k
}

func (k *Kaban) Store(key string, value interface{}) error {
	// キー長の確認
	if len(key) == 0 {
		return fmt.Errorf("len() empty key")
	}
	// インデックスの格納
	func() {
		kabanMtx.Lock()
		defer kabanMtx.Unlock()
		k.keyMap[key] = len(k.chunk)
	}()
	// バイト列化
	var blob []byte
	switch v := value.(type) {
	case string:
		blob = make([]byte, 0, len(v)+1)
		blob = append(blob, sepString)
		blob = append(blob, []byte(v)...)
	case int:
		s := strconv.FormatInt(int64(v), intBase)
		blob = make([]byte, 0, len(s)+1)
		blob = append(blob, sepInt)
		blob = append(blob, []byte(s)...)
	case int8:
		s := strconv.FormatInt(int64(v), intBase)
		blob = make([]byte, 0, len(s)+1)
		blob = append(blob, sepInt)
		blob = append(blob, []byte(s)...)
	case int16:
		s := strconv.FormatInt(int64(v), intBase)
		blob = make([]byte, 0, len(s)+1)
		blob = append(blob, sepInt)
		blob = append(blob, []byte(s)...)
	case int32:
		s := strconv.FormatInt(int64(v), intBase)
		blob = make([]byte, 0, len(s)+1)
		blob = append(blob, sepInt)
		blob = append(blob, []byte(s)...)
	case int64:
		s := strconv.FormatInt(v, intBase)
		blob = make([]byte, 0, len(s)+1)
		blob = append(blob, sepInt)
		blob = append(blob, []byte(s)...)
	case uint:
		s := strconv.FormatUint(uint64(v), intBase)
		blob = make([]byte, 0, len(s)+1)
		blob = append(blob, sepUint)
		blob = append(blob, []byte(s)...)
	case uint8:
		s := strconv.FormatUint(uint64(v), intBase)
		blob = make([]byte, 0, len(s)+1)
		blob = append(blob, sepUint)
		blob = append(blob, []byte(s)...)
	case uint16:
		s := strconv.FormatUint(uint64(v), intBase)
		blob = make([]byte, 0, len(s)+1)
		blob = append(blob, sepUint)
		blob = append(blob, []byte(s)...)
	case uint32:
		s := strconv.FormatUint(uint64(v), intBase)
		blob = make([]byte, 0, len(s)+1)
		blob = append(blob, sepUint)
		blob = append(blob, []byte(s)...)
	case uint64:
		s := strconv.FormatUint(v, intBase)
		blob = make([]byte, 0, len(s)+1)
		blob = append(blob, sepUint)
		blob = append(blob, []byte(s)...)
	default:
		return fmt.Errorf("v=%v %t\n", v, v)
	}
	blob = append(blob, sepDead)
	func() {
		kabanMtx.Lock()
		defer kabanMtx.Unlock()
		k.chunk = append(k.chunk, blob...)
	}()
	return nil
}

func (k *Kaban) Load(key string, ptr interface{}) error {
	// キー長の確認
	if len(key) == 0 {
		return fmt.Errorf("len() empty key")
	}
	// インデックスの取得
	var index int
	var ok bool
	func() {
		kabanMtx.RLock()
		defer kabanMtx.RUnlock()
		index, ok = k.keyMap[key]
	}()
	if !ok {
		return fmt.Errorf("key %s not found", key)
	}
	if k.chunk[index] == sepDead {
		return fmt.Errorf("key %s not found", key)
	}
	// null値のチェック
	if k.chunk[index] == sepNull {
		ptr = nil
		return nil
	}
	// 他の型のチェック
	tailIndex := bytes.IndexByte(k.chunk[index:], sepDead)
	var blob []byte
	func() {
		kabanMtx.RLock()
		defer kabanMtx.RUnlock()
		blob = k.chunk[(index + 1):tailIndex]
	}()
	str := string(bytes.Runes(blob))
	switch k.chunk[index] {
	case sepString:
		p, ok := ptr.(*string)
		if !ok {
			return fmt.Errorf("cast() *string error")
		}
		*p = str
	case sepInt:
		num, err := strconv.ParseInt(str, intBase, intBitSize)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt() %s", err.Error())
		}
		switch p := ptr.(type) {
		case *int:
			*p = (int(num))
		case *int8:
			*p = (int8(num))
		case *int16:
			*p = (int16(num))
		case *int32:
			*p = (int32(num))
		case *int64:
			*p = num
		default:
			return fmt.Errorf("invalid pointer type")
		}
	case sepUint:
		num, err := strconv.ParseUint(str, intBase, intBitSize)
		if err != nil {
			return fmt.Errorf("strconv.ParseUint() %s", err.Error())
		}
		switch p := ptr.(type) {
		case *uint:
			*p = uint(num)
		case *uint8:
			*p = uint8(num)
		case *uint16:
			*p = uint16(num)
		case *uint32:
			*p = uint32(num)
		case *uint64:
			*p = num
		default:
			return fmt.Errorf("invalid pointer type")
		}
	}
	return nil
}

// NewDictionary 辞書の新規作成
//func NewDictionary() *Dictionary {
//	return &Dictionary{}
//}

/*
// Store 属性の格納
func (d *Dictionary) Store(key string, value interface{}) error {
	// キー長の確認
	if len(key) == 0 {
		return fmt.Errorf("len() empty key")
	}
	// JSONエンコード
	blob, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json.Marshal() %s", err.Error())
	}
	d.blobMap.Store(key, blob)
	return nil
}

// Load 属性の取得
func (d *Dictionary) Load(key string, ptr interface{}) error {
	// キー長の確認
	if len(key) == 0 {
		return fmt.Errorf("len() empty key")
	}
	// JSONバイト列の取得
	v, ok := d.blobMap.Load(key)
	if !ok {
		return fmt.Errorf("sync.Map.Load() key %s not found", key)
	}
	jsonBlob, ok := v.([]byte)
	if !ok {
		return fmt.Errorf("cast []byte error")
	}
	// JSONデコード
	if err := json.Unmarshal(jsonBlob, ptr); err != nil {
		return fmt.Errorf("json.Unmarshal() %s", err.Error())
	}
	return nil
}

// Delete 属性の削除
func (d *Dictionary) Delete(key string) error {
	// キー長の確認
	if len(key) == 0 {
		return fmt.Errorf("len() empty key")
	}
	// キーの削除
	d.blobMap.Delete(key)
	return nil
}

// StoreAll 複数属性の格納
func (d *Dictionary) StoreAll(keyValuePairs ...interface{}) error {
	// 引数長の確認
	if len(keyValuePairs)%2 != 0 {
		return fmt.Errorf("len() invalid argument length")
	}
	// キーと値の格納
	for i := 0; i < len(keyValuePairs); i += 2 {
		key, ok := keyValuePairs[i].(string)
		if !ok {
			return fmt.Errorf("cast string invalid key")
		}
		value := keyValuePairs[i+1]
		if err := d.Store(key, value); err != nil {
			return err
		}
	}
	return nil
}

// Keys ソート済みのキー群の取得
func (d *Dictionary) Keys() []string {
	keys := []string{}
	d.blobMap.Range(func(k, v interface{}) bool {
		if key, ok := k.(string); ok {
			keys = append(keys, key)
		}
		return true
	})
	sort.Strings(keys)
	return keys
}

// Merge 他のディクショナリで上書き合成する。
func (d *Dictionary) Merge(dict *Dictionary) {
	dict.blobMap.Range(func(k, v interface{}) bool {
		d.blobMap.Store(k, v)
		return true
	})
}

// String 文字列化
func (d *Dictionary) String() string {
	jsonBlob, err := d.MarshalJSON()
	if err != nil {
		return err.Error()
	}
	return string(jsonBlob)
}

// MarshalJSON JSONバイト列に変換する。
func (d *Dictionary) MarshalJSON() ([]byte, error) {
	keys := d.Keys()
	// バイト列の組み立て
	buf := new(bytes.Buffer)
	fmt.Fprint(buf, `{`)
	for i, key := range keys {
		// 区切り文字の挿入
		if i > 0 {
			fmt.Fprint(buf, ",")
		}
		// JSONバイト列の取得
		v, ok := d.blobMap.Load(key)
		if !ok {
			continue
		}
		jsonBlob, ok := v.([]byte)
		if !ok {
			continue
		}
		// キーと値の書き込み
		fmt.Fprint(buf, `"`, key, `":`)
		if _, err := buf.Write(jsonBlob); err != nil {
			return nil, fmt.Errorf("bytes.Buffer.Write() %s", err.Error())
		}
	}
	fmt.Fprint(buf, "}")
	return buf.Bytes(), nil
}

// UnmarshalJSON JSONバイト列からデータを復元する。
func (d *Dictionary) UnmarshalJSON(blob []byte) error {
	// JSONバイト列をマップに変換
	var values map[string](interface{})
	if err := json.Unmarshal(blob, &values); err != nil {
		return fmt.Errorf("json.Unmarshal() %s", err.Error())
	}
	// マップの中身をJSONエンコードして格納する。
	for key, value := range values {
		blob, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("json.Marshal() %s", err.Error())
		}
		d.blobMap.Store(key, blob)
	}
	return nil
}
*/
