package main

import (
	"database/sql"
	"fmt"
	"github.com/gocolly/colly"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	Title 		string
	Describe	string
	Cover		string
	Content		string
	About		string
	Category	string
	Author		string
	Label		string
	Score		string
	Download 	string
}


func main() {
	// Instantiate default collector
	mysql := fmt.Sprintf("%s:%s@(%s:%d)/%s","root", "2014gaokao","127.0.0.1",3306,"dkread")
	db, err := sql.Open("mysql", mysql)
	if err != nil {
		fmt.Println(err)
	}

	var data = make([]Item,12)

	c := colly.NewCollector(colly.MaxDepth(1),)
	var page = 1
	var max = 2
	// On every a element which has href attribute call callback
	c.OnHTML(".archive-scroll .post .clearfix", func(e *colly.HTMLElement) {
		href := e.ChildAttr(".post-img","href")
		title := strings.TrimSpace(e.ChildAttr(".post-img","title"))
		src := e.ChildAttr(".post-img img","src")
		img_id := e.ChildAttr(".post-img img","alt")
		img := "images/" + img_id + ".jpg"
		DownloadFileProgress(src,"../images/" + img_id + ".jpg")


		// Print link
		//fmt.Printf("href: %s\n", href)
		//fmt.Printf("title found: %s\n",title)
		//fmt.Printf("src found: %s\n", src)
		item := Item{}
		tip := e.ChildTexts(".post-preview p")
		//fmt.Printf("tip found: %s\n",tip )

		item.Title = title
		item.Describe = tip[1]
		item.Cover = img

		c1 := colly.NewCollector()

		c1.OnHTML(".left-part", func(e *colly.HTMLElement) {
			detailsTitle := strings.TrimSpace(e.ChildText(".post-intro h1"))
			detailsTip := e.ChildTexts(".post-intro p")
			detailsContent := e.ChildTexts(".describe p")

			content := ""
			about := ""
			flag := false
			for k,v := range detailsContent {
				if v == "🌹 恩+京-的-书+房+ ww w +E nJin g - C o m +" {
					detailsContent = append(detailsContent[:k], detailsContent[k+1:]...)
					break
				}
				if !flag {
					content += v
				}
				if v == "作者简介" {
					flag = true
				}
				if flag{
					about += v
				}
			}

			item.Content = content
			item.About = about

			var category string
			for _,v := range detailsTip {
				s := strings.Split(v, "：")

				if s[0] == "分类"{
					category = strings.TrimSpace(s[1])
					item.Category = strings.TrimSpace(s[1])
				}

				if s[0] == "作者"{
					item.Author = strings.TrimSpace(s[1])
				}

				if s[0] == "标签"{
					item.Label = strings.TrimSpace(s[1])
				}

				if s[0] == "豆瓣评分"{
					item.Score = strings.TrimSpace(s[1])
				}

				//[作者：孙力科 分类：人物传记 标签：传记、华为、管理 格式：epub/mobi/azw3 豆瓣评分：5.5]
			}

			//fmt.Println(category)

			// Print link
			//fmt.Printf("details title: %s\n", detailsTitle)
			//fmt.Printf("details tip: %s\n",detailsTip)
			//fmt.Printf("details content: %s\n", detailsContent)

			downUrl := e.ChildAttr("#paydown .downbtn","href")

			//fmt.Printf("downUrl: %s\n", downUrl)
			/*go func(href string) {
				detailsUrl <- href
			}(href)*/


			c2 := colly.NewCollector()

			c2.OnHTML(".download-text", func(e *colly.HTMLElement) {
				//detailsTitle := e.ChildText(".post-intro h1")
				down := e.ChildAttrs("span a","href")


				// Print link
				//fmt.Printf("down: %s\n", down)
				//for _, v := range down{
					u := strings.Split(down[0], ".")

					CreateDir("../books/" + category + "/")
					path := "books/" + category + "/" + detailsTitle + "." + u[len(u) -1]
					DownloadFileProgress(down[0],"../books/" + category + "/" + detailsTitle + "." + u[len(u) -1])
					item.Download = path
				//}

			})
			c2.OnRequest(func(r *colly.Request) {
				//fmt.Println("Visiting", r.URL.String())
			})

			c2.Visit(downUrl)

		})
		c1.OnRequest(func(r *colly.Request) {
			//fmt.Println("Visiting", r.URL.String())
		})

		c1.Visit(href)

		if page <= max {
			page += 1
			c.Visit("https://www.enjing.com/wenxue/page/" + strconv.Itoa(page) + "/")
			//DownloadFileProgress(downUrl,"../books/"+category + "/" + detailsTitle + "")
		}


		//fmt.Println(item)
		data = append(data,item)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://www.enjing.com/wenxue/")

	fmt.Println(data)

	insert(db,data)


}

type Reader struct {
	io.Reader
	Total int64
	Current int64
}

func (r *Reader) Read(p []byte) (n int, err error){
	n, err = r.Reader.Read(p)

	r.Current += int64(n)
	fmt.Printf("\r进度 %.2f%%", float64(r.Current * 10000/ r.Total)/100)

	return
}

func DownloadFileProgress(url, filename string) {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer func() {_ = r.Body.Close()}()

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer func() {_ = f.Close()}()

	reader := &Reader{
		Reader: r.Body,
		Total: r.ContentLength,
	}

	_, _ = io.Copy(f, reader)

}

func CreateDir(path string){
	if _, err := os.Stat(path); err == nil {
		fmt.Println("path exists 1", path)
	} else {
		fmt.Println("path not exists ", path)
		err := os.MkdirAll(path, 0711)

		if err != nil {
			log.Println("Error creating directory")
			log.Println(err)
			return
		}
	}

	// check again
	if _, err := os.Stat(path); err == nil {
		fmt.Println("path exists 2", path)
	}
}

var fields = []string{
	"title" ,
	"describe",
	"cover",
	"content",
	"about",
	"category",
	"author",
	"label",
	"score",
	"download",
	"created",
	"status",
}
func insert(db *sql.DB,d []Item)  {

	flows_replace := make([]string, len(fields))
	for i := range fields {
		flows_replace[i] = fmt.Sprintf("$%v", i+1)
	}
	query := fmt.Sprintf("INSERT INTO flow_statistics (%v) VALUES (%v)", strings.Join(fields, ", "), strings.Join(flows_replace, ", "))


	_, err := db.Exec(query, )
	if err != nil {
		fmt.Println(err)
	}
}