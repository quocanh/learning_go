package main

// Program to calculate md5 of a file
// The filename is provided in command line
// If no filename is provided, it calculate the md5 of its executable

// In this implementation it reads file one chunk at a time, so it consume less memory

import (
    "fmt"
    "os"
)

// rotate a 32-bit word to the left.
func left_rotate(num, cnt uint32) uint32 {
  return (num << cnt) | (num >> (32 - cnt))
}

/*
 * These functions implement the four basic operations.
 */
func md5_action(m, a, b , x , s, t uint32) uint32 {
  return left_rotate(a + m + x + t, s) + b
}
func md5_ff(a, b, c, d, x, s, t uint32) uint32 {
  return md5_action((b & c) | ((^b) & d), a, b, x, s, t)
}
func md5_gg(a, b, c, d, x, s, t uint32) uint32 {
  return md5_action((b & d) | (c & (^d)), a, b, x, s, t)
}
func md5_hh(a, b, c, d, x, s, t uint32) uint32 {
  return md5_action(b ^ c ^ d, a, b, x, s, t)
}
func md5_ii(a, b, c, d, x, s, t uint32) uint32 {
  return md5_action(c ^ (b | (^d)), a, b, x, s, t)
}

// convert string to 32-bit little-endian words
//   and perform padding. The string is the full message to be hashed
func byte2words(input []byte) []uint32 {
  var byte_len uint32 = uint32(len(input)) // input length in bytes
  var word_len uint32 = (byte_len + 3 ) >> 2  // input size in words
  var total_len uint32 = (((word_len+2) >> 4) << 4) + 14  // rounding to 16n+14 with spec_14_15

  // copy string content to output
  var output = make([]uint32, total_len + 2) // two extra words to keep 64-bit message length
  var i uint32
  for i = 0; i < byte_len ; i++ {
    output[i>>2] |= uint32(input[i] & 0xFF) << ((i<<3)%32)
  }

  // padding
  // output[(byte_len << 3) >> 5] |= 0x80 << ((byte_len << 3) % 32)
  output[byte_len >> 2] |= 0x80 << ((byte_len << 3) % 32)

  // appending message length
  output[total_len] = (byte_len << 3) & 0xffffffff
  output[total_len+1] = (byte_len << 3) >> 32

  return output;
}

// convert string to 32-bit little-endian words. Its length must be 64*n
// The string is a part of the message to be hashed, no padding for it
func byte2words_no_padding(input []byte) []uint32 {
  var byte_len uint32 = uint32(len(input)) // input length in bytes
  if byte_len % 64 != 0 {
    panic("This function only works on string with length of 64*n characters!")
  }
  var word_len uint32 = byte_len >> 2  // input size in words

  // copy string content to output
  var output = make([]uint32, word_len)
  var i uint32
  for i = 0; i < byte_len ; i++ {
    output[i>>2] |= uint32(input[i] & 0xFF) << ((i<<3)%32)
  }

  return output;
}

// convert string to 32-bit little-endian words. Its length must be 64*n
// The string should be the last part of the message to be hashed,
//   the length (in bytes) of previous part is required
func byte2words_with_padding(input []byte, prev_len uint32) []uint32 {
  var byte_len uint32 = uint32(len(input)) // input length in bytes
  var word_len uint32 = byte_len >> 2  // input size in words
  var buffer_len uint32 = (((word_len+2) >> 4) << 4) + 14  // rounding to 16n+14 with spec_14_15

  // copy string content to output
  var output = make([]uint32, buffer_len + 2) // two extra words to keep 64-bit message length
  var i uint32
  for i = 0; i < byte_len ; i++ {
    output[i>>2] |= uint32(input[i] & 0xFF) << ((i<<3)%32)
  }

  // padding
  output[byte_len >> 2] |= 0x80 << ((byte_len << 3) % 32)

  // appending message length
  msg_len := byte_len + prev_len
  output[buffer_len] = (msg_len << 3) & 0xffffffff
  output[buffer_len+1] = (msg_len << 3) >> 32

  return output;
}

func words2str(x *[4]uint32) string {
  var out string = ""
  for i := 0; i<4; i++ {
    s := fmt.Sprintf("%.2x%.2x%.2x%.2x", byte(x[i]&0xff), byte((x[i]>>8)&0xff),
      byte((x[i]>>16)&0xff), byte((x[i]>>24)&0xff))
    out = out + s
  }
  return out
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

// Calculate md5 hash for a file, reading one chunk a time until EOF
func file_md5(file string) *[4]uint32 {
  f, err := os.Open(file)
  check(err)
  defer f.Close()
  fi, err := f.Stat()
  check(err)
  file_size := fi.Size()

  var a uint32 = 0x67452301
  var b uint32 = 0xefcdab89
  var c uint32 = 0x98badcfe
  var d uint32 = 0x10325476
  regs := [...]uint32{a,b,c,d}

  var chunk_size uint32 = 10240
  buf := make([]byte, chunk_size)
  n1, err := f.Read(buf)
  check(err)

  var wb []uint32
  var cumulative_len uint32 = 0
  for uint32(n1) == chunk_size {
    cumulative_len += chunk_size
    wb = byte2words_no_padding(buf)
    md5_cycle_with_registers(wb, uint32(n1 * 8), &regs)
    if int64(cumulative_len) == file_size {
      n1 = 0
      break
    }
    n1, err = f.Read(buf)
    check(err)
  }
  last_buf := make([]byte, n1)
  for i:=0; i<n1; i++ {
    last_buf[i] = buf[i]
  }
  fmt.Printf("Last chunk size = %d, cumulative_len = %d\n", n1, cumulative_len)
  wb = byte2words_with_padding(last_buf, cumulative_len)
  md5_cycle_with_registers(wb, uint32(n1 * 8), &regs)
  return &regs
}

/*
 * Calculate the MD5 of an array of little-endian words on provided registers
 */
func md5_cycle_with_registers(x []uint32, leng uint32, registers *[4]uint32) {
  var a uint32 = registers[0]
  var b uint32 = registers[1]
  var c uint32 = registers[2]
  var d uint32 = registers[3]

  for i := 0; i < len(x); i += 16 {
    var olda uint32 = a
    var oldb uint32 = b
    var oldc uint32 = c
    var oldd uint32 = d

    a = md5_ff(a, b, c, d, x[i+ 0], 7 , 0xd76aa478)
    d = md5_ff(d, a, b, c, x[i+ 1], 12, 0xe8c7b756)
    c = md5_ff(c, d, a, b, x[i+ 2], 17, 0x242070db)
    b = md5_ff(b, c, d, a, x[i+ 3], 22, 0xc1bdceee)
    a = md5_ff(a, b, c, d, x[i+ 4], 7 , 0xf57c0faf)
    d = md5_ff(d, a, b, c, x[i+ 5], 12, 0x4787c62a)
    c = md5_ff(c, d, a, b, x[i+ 6], 17, 0xa8304613)
    b = md5_ff(b, c, d, a, x[i+ 7], 22, 0xfd469501)
    a = md5_ff(a, b, c, d, x[i+ 8], 7 , 0x698098d8)
    d = md5_ff(d, a, b, c, x[i+ 9], 12, 0x8b44f7af)
    c = md5_ff(c, d, a, b, x[i+10], 17, 0xffff5bb1)
    b = md5_ff(b, c, d, a, x[i+11], 22, 0x895cd7be)
    a = md5_ff(a, b, c, d, x[i+12], 7 , 0x6b901122)
    d = md5_ff(d, a, b, c, x[i+13], 12, 0xfd987193)
    c = md5_ff(c, d, a, b, x[i+14], 17, 0xa679438e)
    b = md5_ff(b, c, d, a, x[i+15], 22, 0x49b40821)

    a = md5_gg(a, b, c, d, x[i+ 1], 5 , 0xf61e2562)
    d = md5_gg(d, a, b, c, x[i+ 6], 9 , 0xc040b340)
    c = md5_gg(c, d, a, b, x[i+11], 14, 0x265e5a51)
    b = md5_gg(b, c, d, a, x[i+ 0], 20, 0xe9b6c7aa)
    a = md5_gg(a, b, c, d, x[i+ 5], 5 , 0xd62f105d)
    d = md5_gg(d, a, b, c, x[i+10], 9 , 0x02441453)
    c = md5_gg(c, d, a, b, x[i+15], 14, 0xd8a1e681)
    b = md5_gg(b, c, d, a, x[i+ 4], 20, 0xe7d3fbc8)
    a = md5_gg(a, b, c, d, x[i+ 9], 5 , 0x21e1cde6)
    d = md5_gg(d, a, b, c, x[i+14], 9 , 0xc33707d6)
    c = md5_gg(c, d, a, b, x[i+ 3], 14, 0xf4d50d87)
    b = md5_gg(b, c, d, a, x[i+ 8], 20, 0x455a14ed)
    a = md5_gg(a, b, c, d, x[i+13], 5 , 0xa9e3e905)
    d = md5_gg(d, a, b, c, x[i+ 2], 9 , 0xfcefa3f8)
    c = md5_gg(c, d, a, b, x[i+ 7], 14, 0x676f02d9)
    b = md5_gg(b, c, d, a, x[i+12], 20, 0x8d2a4c8a)

    a = md5_hh(a, b, c, d, x[i+ 5], 4 , 0xfffa3942)
    d = md5_hh(d, a, b, c, x[i+ 8], 11, 0x8771f681)
    c = md5_hh(c, d, a, b, x[i+11], 16, 0x6d9d6122)
    b = md5_hh(b, c, d, a, x[i+14], 23, 0xfde5380c)
    a = md5_hh(a, b, c, d, x[i+ 1], 4 , 0xa4beea44)
    d = md5_hh(d, a, b, c, x[i+ 4], 11, 0x4bdecfa9)
    c = md5_hh(c, d, a, b, x[i+ 7], 16, 0xf6bb4b60)
    b = md5_hh(b, c, d, a, x[i+10], 23, 0xbebfbc70)
    a = md5_hh(a, b, c, d, x[i+13], 4 , 0x289b7ec6)
    d = md5_hh(d, a, b, c, x[i+ 0], 11, 0xeaa127fa)
    c = md5_hh(c, d, a, b, x[i+ 3], 16, 0xd4ef3085)
    b = md5_hh(b, c, d, a, x[i+ 6], 23, 0x04881d05)
    a = md5_hh(a, b, c, d, x[i+ 9], 4 , 0xd9d4d039)
    d = md5_hh(d, a, b, c, x[i+12], 11, 0xe6db99e5)
    c = md5_hh(c, d, a, b, x[i+15], 16, 0x1fa27cf8)
    b = md5_hh(b, c, d, a, x[i+ 2], 23, 0xc4ac5665)

    a = md5_ii(a, b, c, d, x[i+ 0], 6 , 0xf4292244)
    d = md5_ii(d, a, b, c, x[i+ 7], 10, 0x432aff97)
    c = md5_ii(c, d, a, b, x[i+14], 15, 0xab9423a7)
    b = md5_ii(b, c, d, a, x[i+ 5], 21, 0xfc93a039)
    a = md5_ii(a, b, c, d, x[i+12], 6 , 0x655b59c3)
    d = md5_ii(d, a, b, c, x[i+ 3], 10, 0x8f0ccc92)
    c = md5_ii(c, d, a, b, x[i+10], 15, 0xffeff47d)
    b = md5_ii(b, c, d, a, x[i+ 1], 21, 0x85845dd1)
    a = md5_ii(a, b, c, d, x[i+ 8], 6 , 0x6fa87e4f)
    d = md5_ii(d, a, b, c, x[i+15], 10, 0xfe2ce6e0)
    c = md5_ii(c, d, a, b, x[i+ 6], 15, 0xa3014314)
    b = md5_ii(b, c, d, a, x[i+13], 21, 0x4e0811a1)
    a = md5_ii(a, b, c, d, x[i+ 4], 6 , 0xf7537e82)
    d = md5_ii(d, a, b, c, x[i+11], 10, 0xbd3af235)
    c = md5_ii(c, d, a, b, x[i+ 2], 15, 0x2ad7d2bb)
    b = md5_ii(b, c, d, a, x[i+ 9], 21, 0xeb86d391)

    a = a + olda
    b = b + oldb
    c = c + oldc
    d = d + oldd
  }
  registers[0] = a
  registers[1] = b
  registers[2] = c
  registers[3] = d
}

func main() {
  var file string
  if len(os.Args) == 1 {
    file = os.Args[0]
  } else {
    file = os.Args[1]
  }

  a := file_md5(file)
  fmt.Printf("md5 for the file's content = %s \n", words2str(a))

}
