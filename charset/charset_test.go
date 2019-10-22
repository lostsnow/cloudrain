package charset

import (
	"testing"
)

var tests = []struct {
	src      []byte
	dst      []byte
	encoding string
}{
	{
		src:      []byte("测试字符"),
		dst:      []byte{'\xe6', '\xb5', '\x8b', '\xe8', '\xaf', '\x95', '\xe5', '\xad', '\x97', '\xe7', '\xac', '\xa6'},
		encoding: "utf-8",
	},
	{
		src:      []byte("测试字符"),
		dst:      []byte{'\xb2', '\xe2', '\xca', '\xd4', '\xd7', '\xd6', '\xb7', '\xfb'},
		encoding: "gbk",
	},
	{
		src:      []byte("測試字符"),
		dst:      []byte{'\xe6', '\xb8', '\xac', '\xe8', '\xa9', '\xa6', '\xe5', '\xad', '\x97', '\xe7', '\xac', '\xa6'},
		encoding: "utf-8",
	},
	{
		src:      []byte("測試字符"),
		dst:      []byte{'\xb4', '\xfa', '\xb8', '\xd5', '\xa6', '\x72', '\xb2', '\xc5'},
		encoding: "big5",
	},
}

func assert(t *testing.T, src, dst []byte, err error) {
	if err != nil {
		t.Errorf("Failed: %s", err.Error())
	}
	if string(src) != string(dst) {
		t.Errorf("Failed: give: % x, want: % x", src, dst)
		t.Errorf("Failed: give: %s, want: %s", src, dst)
	}
}

func TestEncode(t *testing.T) {
	for _, tt := range tests {
		b, err := Encode(tt.src, tt.encoding)
		assert(t, b, tt.dst, err)
	}
}

func TestDecode(t *testing.T) {
	for _, tt := range tests {
		b, err := Decode(tt.dst, tt.encoding)
		assert(t, b, tt.src, err)
	}
}
