package kaban

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

type (
	// Dictionary JSON辞書
	Dictionary struct {
		blobMap sync.Map // JSONバイト列マップ
	}
)

// NewDictionary 辞書の新規作成
func NewDictionary() *Dictionary {
	return &Dictionary{}
}

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
