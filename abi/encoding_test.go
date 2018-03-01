package abi

import (
	"math/big"
	"testing"

	"github.com/filecoin-project/go-filecoin/types"

	"github.com/stretchr/testify/assert"
)

// TODO: tests that check the exact serialization of different inputs.
// defer this until we're reasonably unlikely to change the way we serialize
// things.

func TestBasicEncodingRoundTrip(t *testing.T) {
	cases := map[string][]interface{}{
		"empty":      nil,
		"one-int":    {big.NewInt(579)},
		"one addr":   {types.Address("cats")},
		"two addrs":  {types.Address("cats"), types.Address("dogs")},
		"one []byte": {[]byte("foo")},
		"two []byte": {[]byte("foo"), []byte("bar")},
		"a string":   {"flugzeug"},
		"mixed":      {big.NewInt(17), []byte("beep"), "mr rogers", types.Address("cryptokitties")},
	}

	for tname, tcase := range cases {
		t.Run(tname, func(t *testing.T) {
			assert := assert.New(t)
			vals, err := ToValues(tcase)
			assert.NoError(err)

			data, err := EncodeValues(vals)
			assert.NoError(err)

			var types []Type
			for _, val := range vals {
				types = append(types, val.Type)
			}

			outVals, err := DecodeValues(data, types)
			assert.NoError(err)
			assert.Equal(vals, outVals)

			assert.Equal(tcase, FromValues(outVals))
		})
	}
}

type fooTestStruct struct {
	Bar string
	Baz uint64
}

func TestToValuesFailures(t *testing.T) {
	cases := []struct {
		name   string
		vals   []interface{}
		expErr string
	}{
		{
			name:   "nil value",
			vals:   []interface{}{nil},
			expErr: "unsupported type: <nil>",
		},
		{
			name:   "normal int",
			vals:   []interface{}{17},
			expErr: "unsupported type: int",
		},
		{
			name:   "a struct",
			vals:   []interface{}{&fooTestStruct{"b", 99}},
			expErr: "unsupported type: *abi.fooTestStruct",
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			assert := assert.New(t)
			_, err := ToValues(tcase.vals)
			assert.EqualError(err, tcase.expErr)
		})
	}
}