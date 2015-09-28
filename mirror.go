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
    "sync"
)

func createPaths (parsed_url *url.URL) *os.File{
  var dir,file string
  if(strings.Index(parsed_url.Path, ".") >= 0){
   dir, file = filepath.Split(parsed_url.Path)
  } else {
   dir  = parsed_url.Path + "/"
   file = "index.html"
  }

  if(len(dir)>0){
    err := os.MkdirAll(parsed_url.Host + dir, 0777)
    if(err != nil){
        fmt.Println("Directory Create Error: ",dir, err)
        os.Exit(1)
    }
  }
  fileWriter, err := os.Create(parsed_url.Host + dir + file)
  if(err != nil){
      fmt.Println("File Open Error: ",err)
      os.Exit(1)
  }
  //fmt.Println("file created successfully",parsed_url.Host + dir + file)
  return fileWriter
}

func isLocal(link, host string) bool{
  if(strings.LastIndex(link, "/") > strings.LastIndex(link, ".")){
    return false
  }
  if(strings.Index(link, "#") > 0){
   return false
  }
  if(strings.Index(link, host) == 7){
   return true
  }
  return false
}

func shouldSearch() bool{
  return true
}

func generatelinks (resp_reader io.Reader, host string, ch chan string,wg *sync.WaitGroup,set map[string]bool){
  z := html.NewTokenizer(resp_reader)
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

                        if(isLocal(a.Val, host)){
                          //fmt.Println("Link: ", a.Val)
                          //fmt.Println("link found", a.Val)
                          if(!set[a.Val]){
                            countLinks++
                            ch <- a.Val
                          //  fmt.Println("link pushed in channel", a.Val)
                            wg.Add(1)
                            go retrieve(ch,wg,set)}
                          }
                          break;
                      }

                  }

              }
      }
  }
}

func retrieve(ch chan string,wg *sync.WaitGroup,set map[string]bool) {
          uri := <-ch
          if(!set[uri]){
          set[uri] = true;
          fmt.Println(uri)
          parsed_url, err := url.Parse(uri)
          if(err != nil){
            fmt.Println("Url Parsing Error: ",err)
            os.Exit(1)
          }
          //fmt.Println("Http request", uri)
          resp, err := http.Get(uri)
          //fmt.Println("Http response", uri)

          if(err != nil){
              fmt.Println("Http Transport Error: ",err)
              os.Exit(1)
          }
          fileWriter:= createPaths(parsed_url)
          defer fileWriter.Close()
          resp_reader := io.TeeReader(resp.Body, fileWriter)
        //  fileWriter.Close()
        //  fmt.Println("file saved successfully")
          if(shouldSearch()){
           generatelinks(resp_reader,parsed_url.Host,ch,wg,set)
          }
          defer resp.Body.Close()
          }
          defer wg.Done()
}


func main(){
    flag.Parse()
    args := flag.Args()
    set := make(map[string]bool)

    if(len(args)<1){
        fmt.Println("Specify a start url")
        os.Exit(1)
     }
    ch := make(chan string,1)
    ch <- args[0]

    parsed_url, err := url.Parse(args[0])
    if(err != nil){
      fmt.Println("Url Parsing Error: ",err)
      os.Exit(1)
    }
    err = os.MkdirAll(parsed_url.Host, 0777)
    if(err != nil){
        fmt.Println("Directory Create Error: ",parsed_url.Host, err)
        os.Exit(1)
    }

    var wg sync.WaitGroup
    wg.Add(1)
    go retrieve(ch,&wg,set)
    wg.Wait()
    fmt.Println("**********operation complete***************")
}
