// Program to implement matrix multiplication in Go:
// given an input file with lines in the form of "i j k",
// where i is row number, j is the column number and k is the value at row i, col j;
// then:
//     if it is a square matrix, multiply it by its transpose
//     otherwise return an error;
// use arrays for the matrices (and not slices);


// Plan: read file and use a map to store working values: i, j, k
//    key = i_j   value = k
// Keep track of max_i and max_j as we process input file
// After all file is read, check if max_i == max_j for error
// Assume missing elements are valid, and their values are default (0 for int)

// type Matrix struct { row int; col int; data []int }
// Need methods to
//  - get element at row_i, col_j
//  - set element at row_i, col_j
//  - print the matrix
//  - transpose
//  - multiply

package main
import (
  "strconv"
  "strings"
  "os"
  "fmt"
  "errors"
  "bufio"
)

type Matrix struct {
  row int
  col int
  data []int
}

func (m Matrix) elem(i int, j int) int {
  return m.data[i* m.col + j]
}

func (m Matrix) setElem(i int, j int, k int) {
  m.data[i*m.col + j] = k
}

func (m Matrix) print() {
  for i := 0; i < m.row; i++ {
    for j := 0; j < m.col; j++ {
      fmt.Print(" ", m.elem(i,j))
    }
    fmt.Println("")
  }
}

func (m Matrix) transpose() Matrix {
  mm := Matrix{m.col, m.row, make([]int, m.col * m.row)}
  for i := 0; i < m.row; i++ {
    for j := 0; j < m.col; j++ {
      mm.setElem(j,i, m.elem(i,j))
    }
  }
  return mm
}

func (m Matrix) multiply(m2 Matrix) Matrix {
  if m.col != m2.row {
    check(errors.New("Only compatible matrices can multiply!"))
  }
  mm := Matrix{m.col, m2.row, make([]int, m.col * m2.row)}
  for i := 0; i < m.row; i++ {
    for j := 0; j < m2.col; j++ {
      sum := 0
      for k := 0; k < m.col; k++ {
        sum += m.elem(i,k) * m2.elem(k,j)
      }
      mm.setElem(i,j, sum)
    }
  }
  return mm
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func parseLine(line string) (int, int, int) {
  items := strings.Split(line, " ")
  if len(items) != 3 {
    check(errors.New("Each line must have 3 numbers: line= " + line))
  }
  i, err := strconv.Atoi(items[0])
  check(err)
  j, err := strconv.Atoi(items[1])
  check(err)
  k, err := strconv.Atoi(items[2])
  check(err)
  if i < 0 || j < 0 {
    check(errors.New("Coefficient must be non-negative: line= " + line))
  }
  return i, j, k
}

func matrixFromFile(file string) Matrix {
  f, err := os.Open(file)
  check(err)
  defer f.Close()
  reader := bufio.NewReader(f)
  scanner := bufio.NewScanner(reader)
  var tmp = make(map[string]int)
  max_i, max_j := -1, -1

  for scanner.Scan() {
    line := scanner.Text()
    if strings.HasPrefix(line, "//") { continue }
    items := strings.Split(line, " ")
    i, j, k := parseLine(line)
    if (i > max_i) { max_i = i}
    if (j > max_j) { max_j = j}
    tmp[items[0] + "_" + items[1]] = k
  }
  if max_i != max_j {
    check(errors.New("Matrix must be square: row= " + strconv.Itoa(max_i+1) +
      ", col= " + strconv.Itoa(max_j+1)))
  }
  max_i++; max_j++
  mm := Matrix{max_i, max_j, make([]int, max_i * max_j)}
  for x := 0; x < max_i; x++ {
    for y := 0; y < max_j; y++ {
      mm.setElem(x,y, tmp[strconv.Itoa(x) + "_" + strconv.Itoa(y)])
    }
  }
  return mm
}

func main() {
  if len(os.Args) != 2 {
   check(errors.New("Usage: os.Args[0] path_to_data_file"))
  }
  file := os.Args[1]

  mm := matrixFromFile(file)
  fmt.Println("Input matrix")
  mm.print()
  mm1 := mm.transpose()
  fmt.Println("Transpose matrix")
  mm1.print()
  mm2 := mm.multiply(mm1)
  fmt.Println("Product matrix")
  mm2.print()
}
