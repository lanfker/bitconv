## What this tiny library does:

Suppose we have a `[]byte` slice with length `N`. 

This program views the `N` bytes as consecutive bits with possible values `0` or `1`. 

In terms of bit ordering, this is how it should be: 

`76543210` `76543210` `76543210` ... 

`BYTE 1`     `BYTE2`        `BYTE3`      ...

The wrapper function `GetSigned` and `GetUnsigned` are all you needed. They both take the same set of arguments and have prototypes of the following: 

```go
func GetUnsigned (payload []byte, sbyte, sbit, bitlen int) int64
```

```go
func GetSigned (payload []byte, sbyte, sbit, bitlen int) int64
```

* __payload__ the []byte slice

* __sbyte__ (index starts with 0)

* __sbit__ within the start byte (index starting with 0, and this can only be [0,7] since a byte has only 8 bits)

* __bitlen__ represents the number of bits the resulting integer should have. 

  

  __Example__ For example, to parse a 13-bit unsigned integer from byte 0, starting from the 3rd bit, we would have something like this: 

  ```go
  GetUnsigned/GetSigned (0,3,13)
  ```

  This function call will take the following bold bits as values

  0000**0000 00000000 0**0000000 

__Special case__: if the specified bit range goes beyond the []byte slice, you always get __ZERO__. If zero is a valid value for your application, you are on your own. 