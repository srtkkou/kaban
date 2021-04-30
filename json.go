package kaban

import (
	"strconv"
	"strings"
)

// Convert to JSON bytes.
func (k *Kaban) MarshalJSON() ([]byte, error) {
	jBlob := make([]byte, 0, (2 * len(k.block)))
	jBlob = append(jBlob, '{')
	for key, index := range k.keyMap {
		jBlob = append(jBlob, '"')
		jBlob = append(jBlob, []byte(key)...)
		jBlob = append(jBlob, '"', ':')
		switch k.block[index] {
		case sepString:
			s, err := k.LoadString(key)
			if err != nil {
				return []byte{}, err
			}
			jBlob = append(jBlob, '"')
			jBlob = append(jBlob, []byte(s)...)
			jBlob = append(jBlob, '"', ',')
		case sepInt:
			num, err := k.LoadInt64(key)
			if err != nil {
				return []byte{}, err
			}
			s := strconv.FormatInt(num, 10)
			jBlob = append(jBlob, []byte(s)...)
			jBlob = append(jBlob, ',')
		case sepUint:
			num, err := k.LoadUint64(key)
			if err != nil {
				return []byte{}, err
			}
			s := strconv.FormatUint(num, 10)
			jBlob = append(jBlob, []byte(s)...)
			jBlob = append(jBlob, ',')
		case sepSlice:
			switch k.block[index+1] {
			case sepString:
				strs, err := k.LoadStrings(key)
				if err != nil {
					return []byte{}, err
				}
				jBlob = append(jBlob, '[', '"')
				jBlob = append(jBlob, []byte(strings.Join(strs, `","`))...)
				jBlob = append(jBlob, '"', ']', ',')
				//case sepInt:
				//case sepUint:
			}
		}
	}
	// Replace last comma with brace.
	lastIndex := len(jBlob) - 1
	jBlob[lastIndex] = '}'
	return jBlob, nil
}
