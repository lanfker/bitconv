package bitconv

import (
	"fmt"
)

//GMSign represents if an integer we get from CAN bit stream is unsigned or signed
type GMSign int

const (
	//Signed The integer is signed
	Signed GMSign = 0

	//Unsigned The data is unsigned
	Unsigned GMSign = 1
)

//getSignedInt gets a negative integer from the bit pattern represented by val
//Notice that if the value is positive, we do not need to call this function at all..
func getSignedInt(val, msb int64, sign GMSign) int64 {
	var t int64 = 1
	if sign == Signed {
		for i := int64(0); i < msb; i++ {
			t <<= 1
		}

		if val&t == 1 {
			val = (^val + 1)
			val *= -1
		}
		return val
	}
	return val
}

func dataValid(payload []byte, sbyte, sbit, bitlen int) bool {
	ok := false
	if sbit > bitlen && sbyte < len(payload) {
		ok = true
	}
	if bitlen > sbit && sbyte+(bitlen-sbit-1)/8 < len(payload) {
		ok = true
	}
	return ok
}

//GetUnsigned gets an unsigned integer from the bit pattern represented by the payload byte slice
func GetUnsigned(payload []byte, sbyte, sbit, bitlen int) int64 {
	if !dataValid(payload, sbyte, sbit, bitlen) {
		return 0
	}
	var v int64
	return extractBitRange(payload, v, sbyte, sbit, bitlen)
}

//GetSigned gets an signed integer from the bit pattern represented by the payload byte slice
//When the underline integer is unsigned, it behaves the same as GetUnsigned
func GetSigned(payload []byte, sbyte, sbit, bitlen int) int64 {
	if !dataValid(payload, sbyte, sbit, bitlen) {
		return 0
	}

	var i byte = 1
	for k := 0; k < sbit; k++ {
		i <<= 1
	}
	var v int64
	if payload[sbyte]&i == 0 {
		v = extractBitRange(payload, v, sbyte, sbit, bitlen)
		return v
	}
	v = fillMSB(bitlen)
	v = extractBitRange(payload, v, sbyte, sbit, bitlen)
	return getSignedInt(v, int64(bitlen), Signed)
}

//fillMSB fills the most significant bits with ONES if the number is negative
func fillMSB(n int) int64 {
	var v int64
	for i := 0; i < 64-n; i++ {
		v <<= 1
		v |= 1
	}
	return v
}

//extractBitRange Extract a range of bits specified by sbyte sbit bitlen
func extractBitRange(stream []byte, val int64, sbyte, sbit, bitlen int) int64 {
	var t uint8 = 1
	//fmt.Printf("extractBitRange: starting from %b\n", t)
	curByte := sbyte
	curBit := sbit
	curLen := 0

	for curLen < bitlen {
		if curBit < 0 {
			curBit = 7
			curByte++
		}

		t = 1
		for i := 0; i < curBit; i++ {
			t <<= 1
		}

		b := stream[curByte] & t
		if b == 0 {
			val <<= 1
			val |= 0
		} else {
			val <<= 1
			val |= 1
		}
		//fmt.Printf("curBit: %d, t: %b, val: %b, b: %b\n", curBit, t, *val, b)
		curBit--
		curLen++
	}
	return val
}

func printPayload(payload []byte) {
	for k := range payload {
		fmt.Printf("%08b ", payload[k])
	}
	fmt.Println()
}
