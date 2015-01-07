package main
import (
  "fmt"
  "io/ioutil"
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
//   and perform padding
func byte2words(input []byte) []uint32 {
  var byte_len uint32 = uint32(len(input)) // input length in bytes
  var word_len uint32 = (byte_len + 3 ) >> 2  // input size in words
  var total_len uint32
  if word_len % 16 == 14 {
    total_len = word_len + 16
  } else {
    total_len = (((byte_len * 8 + 64) >> 9) << 4) + 14
  }

  // copy string content to output
  var output = make([]uint32, total_len + 2) // two extra words to keep 64-bit message length
  var i uint32
  for i = 0; i < byte_len * 8; i += 8 {
    output[i>>5] |= uint32(input[i>>3] & 0xFF) << (i%32)
  }

  // padding
  output[(byte_len * 8) >> 5] |= 0x80 << ((byte_len * 8) % 32)

  // appending message length
  output[total_len] = (byte_len * 8) & 0xffffffff
  output[total_len+1] = (byte_len * 8) >> 32

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

func calc_md5(s []byte) *[4]uint32 {
  a1 := byte2words(s)  // convert byte_array to array of little-endian words
  return md5_cycle(a1, uint32(len(s)) * 8)  // apply md5 algorithm on little-endian words
}

/*
 * Calculate the MD5 of an array of little-endian words
 */
func md5_cycle(x []uint32, leng uint32) *[4]uint32 {
  var a uint32 = 0x67452301
  var b uint32 = 0xefcdab89
  var c uint32 = 0x98badcfe
  var d uint32 = 0x10325476

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

  o := new([4]uint32)
  o[0] = a; o[1] = b; o[2] = c; o[3] = d
  return o
}

func main() {
  var file string
  if len(os.Args) == 1 {
    file = os.Args[0]
  } else {
    file = os.Args[1]
  }
  tmp, err := ioutil.ReadFile(file)
  if err != nil {
    panic(err)
  }
  a := calc_md5(tmp)
  fmt.Printf("md5 for the file's content = %s \n", words2str(a))
}
