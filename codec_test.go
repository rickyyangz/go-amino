package amino_test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-amino"
)

type SimpleStruct struct {
	String string
	Bytes  []byte
	Time   time.Time
}

func newSimpleStruct() SimpleStruct {
	s := SimpleStruct{
		String: "hello",
		Bytes:  []byte("goodbye"),
		Time:   time.Now().UTC().Truncate(time.Millisecond), // strip monotonic and timezone.
	}
	return s
}

func TestMarshalUnmarshalBinaryPointer0(t *testing.T) {

	var s = newSimpleStruct()
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinaryLengthPrefixed(s) // no indirection
	t.Log(b)
	assert.Nil(t, err)

	var s2 SimpleStruct
	err = cdc.UnmarshalBinaryLengthPrefixed(b, &s2) // no indirection
	assert.Nil(t, err)
	assert.Equal(t, s, s2)

}

func TestMarshalUnmarshalBinaryPointer1(t *testing.T) {

	var s = newSimpleStruct()
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinaryLengthPrefixed(&s) // extra indirection
	assert.Nil(t, err)

	var s2 SimpleStruct
	err = cdc.UnmarshalBinaryLengthPrefixed(b, &s2) // no indirection
	assert.Nil(t, err)
	assert.Equal(t, s, s2)

}

func TestMarshalUnmarshalBinaryPointer2(t *testing.T) {

	var s = newSimpleStruct()
	var ptr = &s
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinaryLengthPrefixed(&ptr) // double extra indirection
	assert.Nil(t, err)

	var s2 SimpleStruct
	err = cdc.UnmarshalBinaryLengthPrefixed(b, &s2) // no indirection
	assert.Nil(t, err)
	assert.Equal(t, s, s2)

}

func TestMarshalUnmarshalBinaryPointer3(t *testing.T) {

	var s = newSimpleStruct()
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinaryLengthPrefixed(s) // no indirection
	assert.Nil(t, err)

	var s2 *SimpleStruct
	err = cdc.UnmarshalBinaryLengthPrefixed(b, &s2) // extra indirection
	assert.Nil(t, err)
	assert.Equal(t, s, *s2)
}

func TestMarshalUnmarshalBinaryPointer4(t *testing.T) {

	var s = newSimpleStruct()
	var ptr = &s
	cdc := amino.NewCodec()
	b, err := cdc.MarshalBinaryLengthPrefixed(&ptr) // extra indirection
	assert.Nil(t, err)

	var s2 *SimpleStruct
	err = cdc.UnmarshalBinaryLengthPrefixed(b, &s2) // extra indirection
	assert.Nil(t, err)
	assert.Equal(t, s, *s2)

}

func TestDecodeInt8(t *testing.T) {
	// DecodeInt8 uses binary.Varint so we need to make
	// sure that all the values out of the range of [-128, 127]
	// return an error.
	tests := []struct {
		in      int8
		wantErr string
		want    int8
	}{
		{in: 0x7F, want: 0x7F},
		{in: -0x7F, want: -0x7F},
		{in: -0x80, want: -0x80},
		{in: 0x10, want: 0x10},
		{in: math.MinInt8, want: math.MinInt8},
		{in: math.MaxInt8, want: math.MaxInt8},
	}

	for i, tt := range tests {
		var w bytes.Buffer
		//w := bufio.NewWriter(&buf)
		if err := amino.EncodeInt8(&w, tt.in); err != nil {
			panic(err)
		}
		//w.Flush()
		gotI8, gotN, err := amino.DecodeInt8(w.Bytes())
		if tt.wantErr != "" {
			if err == nil {
				t.Errorf("#%d expected error=%q", i, tt.wantErr)
			} else if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("#%d\ngotErr=%q\nwantSegment=%q", i, err, tt.wantErr)
			}
			continue
		}
		if err != nil {
			t.Errorf("#%d unexpected error: %v", i, err)
			continue
		}
		if wantI8 := tt.want; gotI8 != wantI8 {
			t.Errorf("#%d gotI8=%d wantI8=%d", i, gotI8, wantI8)
		}
		if wantN := 1; gotN != wantN {
			t.Errorf("#%d gotN=%d wantN=%d", i, gotN, wantN)
		}
	}
}

func TestDecodeInt16(t *testing.T) {
	// DecodeInt16 uses binary.Varint so we need to make
	// sure that all the values out of the range of [-32768, 32767]
	// return an error.
	tests := []struct {
		in      int16
		wantErr string
		want    int16
	}{
		{in: -0x8000, want: -0x8000},
		{in: -0x7FFF, want: -0x7FFF},
		{in: -0x7F, want: -0x7F},
		{in: -0x80, want: -0x80},
		{in: 0x10, want: 0x10},
		{in: math.MinInt16, want: math.MinInt16},
		{in: math.MaxInt16, want: math.MaxInt16},
	}

	for i, tt := range tests {
		var buf bytes.Buffer
		//w := bufio.NewWriter(&buf)
		if err := amino.EncodeInt16(&buf, tt.in); err != nil {
			panic(err)
		}
		//w.Flush()
		gotI16, gotN, err := amino.DecodeInt16(buf.Bytes())
		if tt.wantErr != "" {
			if err == nil {
				t.Errorf("#%d in=(%X) expected error=%q", i, tt.in, tt.wantErr)
			} else if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("#%d\ngotErr=%q\nwantSegment=%q", i, err, tt.wantErr)
			}
			continue
		}
		if err != nil {
			t.Errorf("#%d unexpected error: %v", i, err)
			continue
		}
		if wantI16 := tt.want; gotI16 != wantI16 {
			t.Errorf("#%d gotI16=%d wantI16=%d", i, gotI16, wantI16)
		}
		if wantN := 2; gotN != wantN {
			t.Errorf("#%d gotN=%d wantN=%d", i, gotN, wantN)
		}
	}
}

func TestEncodeDecodeString(t *testing.T) {
	s := "🔌🎉⛵︎♠️⎍"
	bs := []byte(s)
	di := len(bs) * 3 / 4
	b1 := bs[:di]
	b2 := bs[di:]

	// Encoding phase
	buf1 := new(bytes.Buffer)
	if err := amino.EncodeByteSlice(buf1, b1); err != nil {
		t.Fatalf("EncodeByteSlice(b1) = %v", err)
	}
	buf2 := new(bytes.Buffer)
	if err := amino.EncodeByteSlice(buf2, b2); err != nil {
		t.Fatalf("EncodeByteSlice(b2) = %v", err)
	}

	// Decoding phase
	e1 := buf1.Bytes()
	dec1, n, err := amino.DecodeByteSlice(e1)
	if err != nil {
		t.Errorf("DecodeByteSlice(e1) = %v", err)
	}
	if g, w := n, len(e1); g != w {
		t.Errorf("e1: length:: got = %d want = %d", g, w)
	}
	e2 := buf2.Bytes()
	dec2, n, err := amino.DecodeByteSlice(e2)
	if err != nil {
		t.Errorf("DecodeByteSlice(e2) = %v", err)
	}
	if g, w := n, len(e2); g != w {
		t.Errorf("e2: length:: got = %d want = %d", g, w)
	}
	joined := bytes.Join([][]byte{dec1, dec2}, []byte(""))
	if !bytes.Equal(joined, bs) {
		t.Errorf("got joined=(% X) want=(% X)", joined, bs)
	}
	js := string(joined)
	if js != s {
		t.Errorf("got string=%q want=%q", js, s)
	}
}

func TestCodecSeal(t *testing.T) {

	type Foo interface{}
	type Bar interface{}

	cdc := amino.NewCodec()
	cdc.RegisterInterface((*Foo)(nil), nil)
	cdc.Seal()

	assert.Panics(t, func() { cdc.RegisterInterface((*Bar)(nil), nil) })
	assert.Panics(t, func() { cdc.RegisterConcrete(int(0), "int", nil) })
}


type Outer struct {
	A []byte
	I *Inner
	S string
}

type Inner struct {
	In int64
}


func TestCodec_Bytes(t *testing.T) {
	cdc := amino.NewCodec()
	// cdc.RegisterConcrete(&TT{}, "tt", nil)
	aa := make([]byte, 10)
	o := Outer {
		A: aa,
		I: &Inner{},
		S: "1",
	}

	bz, err := cdc.MarshalBinaryBare(o)
	require.NoError(t, err)
	t.Log(bz)

	var o2 Outer
	err = cdc.UnmarshalBinaryBare(bz, &o2)
	require.NoError(t, err)
	require.Equal(t, o, o2)
}

func TestCodec_Bytes2(t *testing.T) {

	ss := SimpleStruct{
		Time: time.Time{},
	}
	fmt.Println(ss)
	fmt.Println(ss.Time.Unix())

	cdc := amino.NewCodec()
	bz := cdc.MustMarshalBinaryBare(ss)
	fmt.Println(bz)

	var ss2 SimpleStruct
	cdc.MustUnmarshalBinaryBare(bz, &ss2)
	fmt.Println(ss2)
	require.Equal(t, ss, ss2)
}
