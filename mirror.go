package main

import (
    "fmt"
    "flag"
    "os"
    "net/http"
    "net/url"
    "io"
    "golang.org/x/net/html"
    "strings"
    "path/filepath"
    //"time"
    "sync"
)

func createPaths (full_path string) *os.File{

  dir, file := filepath.Split(full_path)
  fmt.Println(full_path)
  if ((len(dir) == 0) || (len(file) == 0)){
    file = "index.html"
  }
  if(len(dir)>0){
    err := os.MkdirAll(dir, 0777)
    if(err != nil){
        fmt.Println("Directory Create Error: ",dir, err)
        os.Exit(1)
    }
  }
  fileWriter, err := os.Create(dir+file)
  if(err != nil){
      fmt.Println("File Open Error: ",err)
      os.Exit(1)
  }
  fmt.Println("file created successfully",full_path)
  return fileWriter
}

func generatelinks (resp_reader io.Reader, uri string, ch chan string,wg *sync.WaitGroup){
  z := html.NewTokenizer(resp_reader)
  fmt.Println("finding links in", uri)
  countLinks := 0
  for{
      tt := z.Next();
      switch{
          case tt==html.ErrorToken:
              fmt.Println("Total number of links: ", countLinks)
              return
          case tt==html.StartTagToken:
              t := z.Token()
              if t.Data == "a"{
                  for _,a := range t.Attr{
                      if a.Key == "href"{
                        if(strings.Index(a.Val, uri) == 0){
                          //fmt.Println("Link: ", a.Val)
                          countLinks++
                          fmt.Println("link found", a.Val)
                          ch <- a.Val
                          fmt.Println("link pushed in channel", a.Val)
                          wg.Add(1)
                          go retrieve(ch,wg)
                          }
                          break;
                      }

                  }

              }
      }
  }
  fmt.Println("links finding complete", uri)
}

func retrieve(ch chan string,wg *sync.WaitGroup) {
          uri := <-ch
          parsed_url, err := url.Parse(uri)
          if(err != nil){
            fmt.Println("Url Parsing Error: ",err)
            os.Exit(1)
          }
          fmt.Println("Http request", uri)
          resp, err := http.Get(uri)
          fmt.Println("Http response", uri)

          if(err != nil){
              fmt.Println("Http Transport Error: ",err)
              os.Exit(1)
          }
          full_path := parsed_url.Host+parsed_url.Path
          fileWriter:= createPaths(full_path)
          defer fileWriter.Close()
          resp_reader := io.TeeReader(resp.Body, fileWriter)
          fmt.Println("file saved successfully")
          generatelinks(resp_reader,uri,ch,wg)
          defer resp.Body.Close()
          defer wg.Done()
}


func main(){
    flag.Parse()
    args := flag.Args()

    if(len(args)<1){
        fmt.Println("Specify a start url")
        os.Exit(1)
     }
    ch := make(chan string,1)
    ch <- args[0]
    var wg sync.WaitGroup
    wg.Add(1)
    go retrieve(ch,&wg)
    wg.Wait()
    //time.Sleep(10000 * time.Millisecond)
    fmt.Println("**********operation complete***************")
}
