// web server that implements mini clone of redis: a service with REST APIs
//   and 4 methods: PUT, GET, DELETE, COUNT

// Test the server by using curl
// curl -X PUT -d total_records=100 localhost:8082
// curl -X PUT -d total_bytes=10000 localhost:8082
// curl -X PUT -d something_else="hello, world" localhost:8082
// curl -X PUT -d key1=value1 localhost:8082
// curl -X GET -d "key1" localhost:8082
// curl -X COUNT localhost:8082
// curl -X COUNT -d "total" localhost:8082
// curl -X DELETE -d "key1" localhost:8082

package main
import (
  "net/http"
  "io/ioutil"
  "strconv"
  "strings"
  "regexp"
)

var storage = make(map[string]string)

func serve(w http.ResponseWriter, req *http.Request) {
  body, err := ioutil.ReadAll(req.Body)
  if err != nil {
      panic("Error in ioutil.ReadAll(req.Body)")
  }

  method := req.Method
  switch method {
  case "PUT":
    data := strings.Split(string(body), "=")
    storage[data[0]] = data[1]
    w.Write([]byte("OK"))
  case "GET":
    w.Write([]byte(storage[string(body)]))
  case "DELETE":
    delete(storage,string(body))
    w.Write([]byte("OK"))
  case "COUNT":
    if len(body) == 0 {
      w.Write([]byte(strconv.Itoa(len(storage))))
    } else {
      pattern := string(body) + ".*"
      r, _ := regexp.Compile(pattern)
      count := 0
      for key := range storage {
        if r.MatchString(key) {
          count++
        }
      }
      w.Write([]byte(strconv.Itoa(count)))
    }
  }
}

func main() {
  http.HandleFunc("/", serve)
  http.ListenAndServe(":8082", nil)
}
