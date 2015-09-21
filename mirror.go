package main

import (
  "fmt"
  "net/http"    
  "io"
  "golang.org/x/net/html"
  "strings"
  "os"
  "net/url"
) 

func main() {
//  baseurl := "https://www.google.com"
baseurl := "https://library.snu.ac.kr/sites/default/files/library-brochure/leading_the_way_2015.pdf"
  resp, err := http.Get(baseurl)
  //defer resp.Body.Close()
  if (err!=nil){                                                  
    fmt.Println("http transport error is:", err)
  }else
  { 
    u,_ := url.Parse(baseurl)   
    os.MkdirAll(u.Path,0777)
    w,_ := os.Create(baseurl)

   // body, err := ioutil.ReadAll(resp.Body)
    b :=io.TeeReader(resp.Body, w)
    ct := resp.Header["Content-Type"][0]
    fmt.Println(ct)
    if (err == nil){                                              
     // fmt.Println("size", (string(body)))
      /*b := resp.Body
	body, err := ioutil.ReadAll(b)
	if (err!=nil){                                                  
	    fmt.Println("http transport error is:", err)
	  }
	fmt.Println("size", (string(body)))*/
      page := html.NewTokenizer(b)
      fmt.Println("page", page)
      totalCount := 0
      startTag := 0
      links := 0
      inLinks := 0
      
for {
    tokenType := page.Next()
    //fmt.Println("type", tokenType)
    totalCount = totalCount + 1
    if tokenType == html.ErrorToken {
      break;
    }
    if tokenType == html.StartTagToken {
      startTag = startTag + 1
      token := page.Token()
      if(token.Data == "a"){
        links += 1 
	for _, a := range token.Attr {
	    if a.Key == "href" {
		fmt.Println(a.Val)
              if(strings.Index(a.Val, url) == 0){
               inLinks += 1
              // fmt.Println(a.Val);
              }
	      break
	    }
	}



      }
    }
}
fmt.Println("totalCount",totalCount);
fmt.Println("startTag",startTag);
fmt.Println("total links",links);
fmt.Println("total internal links",inLinks);
    }
  }
}

