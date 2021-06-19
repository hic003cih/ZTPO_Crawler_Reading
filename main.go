package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	// Instantiate default collector
	//c := colly.NewCollector(
	// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
	//	colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
	//)

	c := colly.NewCollector()

	detail := colly.NewCollector()

	c.Limit(&colly.LimitRule{Parallelism: 5, RandomDelay: 15 * time.Second})

	detail.Limit(&colly.LimitRule{Parallelism: 5, RandomDelay: 30 * time.Second})

	// topic := model.Topic{
	// 	Title:     "",
	// 	Content:   "",
	// 	Questions: model.Question{},
	// 	URL:       "",
	// }

	type Question struct {
		Topic  string
		Item   []string
		Answer string
	}

	type Topic struct {
		Title     string
		Content   string
		Questions []Question
		URL       string
	}

	data := &Topic{}

	// Before making a request print "Visiting ..."
	// c.OnRequest(func(r *colly.Request) {
	// 	log.Println("visiting", r.URL.String())
	// })

	// 抓類別Class 名稱
	// 當Visit訪問網頁後，在網頁響應(Response)之後、發現這是HTML格式 執行的事情

	//取出每題的標題
	/*c.OnHTML(".q-c-title", func(e *colly.HTMLElement) {
		// 每找到一個符合 goquerySelector字樣的結果，便會進這個OnHTML一次
		//e.Request.Visit(e.Attr("href"))
		fmt.Println(e.Text)
	})*/

	//取出題目
	c.OnHTML("span[id='top-title-strong']", func(e *colly.HTMLElement) {
		// 每找到一個符合 goquerySelector字樣的結果，便會進這個OnHTML一次
		//e.Request.Visit(e.Attr("href"))
		//fmt.Println("find")
		//fmt.Println("題目")
		//fmt.Println(e.Text)
		data.Title = strings.TrimSpace(e.Text)
	})

	//取出主內容
	//取不出來題目,因為題目是用“”包起來
	c.OnHTML("p[id='artic-en']", func(e *colly.HTMLElement) {
		// 每找到一個符合 goquerySelector字樣的結果，便會進這個OnHTML一次
		//e.Request.Visit(e.Attr("href"))
		//fmt.Println("find")
		//fmt.Println("主內容")
		//fmt.Println(e.Text)
		data.Content = e.Text
	})

	//取出第一個每題的標題
	c.OnHTML("div[id='question-title']", func(e *colly.HTMLElement) {
		// 每找到一個符合 goquerySelector字樣的結果，便會進這個OnHTML一次
		//e.Request.Visit(e.Attr("href"))
		//fmt.Println("問題")
		//fmt.Println(e.Text)

		//先加第一題Question的位置
		//fmt.Println(e.Text)
		data.Questions = append(data.Questions, Question{})
		data.Questions[0].Topic = strings.TrimSpace(e.Text)

	})

	questionsItemCount := 0
	//取出第一個每題的每個選項
	c.OnHTML("div[class='ans-line']", func(e *colly.HTMLElement) {
		// 每找到一個符合 goquerySelector字樣的結果，便會進這個OnHTML一次
		//e.Request.Visit(e.Attr("href"))
		//fmt.Println("find")
		//fmt.Println("選項")
		//fmt.Println(e.Text)

		//先加第一題Question的選項的位置
		data.Questions[0].Item = append(data.Questions[0].Item, "")
		data.Questions[0].Item[questionsItemCount] = strings.TrimSpace(e.Text)
		questionsItemCount = questionsItemCount + 1
	})

	//取出答案
	c.OnHTML("span[class='true-answer-content']", func(e *colly.HTMLElement) {
		// 每找到一個符合 goquerySelector字樣的結果，便會進這個OnHTML一次
		//e.Request.Visit(e.Attr("href"))
		//fmt.Println("find")
		//fmt.Println("答案")
		//fmt.Println(e.Text)

		//data.Questions = append(data.Questions, Question{})
		data.Questions[0].Answer = strings.TrimSpace(e.Text)
	})
	//取出音擋
	// c.OnHTML("audio[id='jp_audio_0']", func(e *colly.HTMLElement) {
	// 	// 每找到一個符合 goquerySelector字樣的結果，便會進這個OnHTML一次
	// 	//e.Request.Visit(e.Attr("href"))
	// 	//fmt.Println("find")
	// 	//fmt.Println("音檔")
	// 	fmt.Println(e.Attr("src"))
	// })

	// 除了第一個選項以外的連結
	c.OnHTML("a[class='que-index']", func(e *colly.HTMLElement) {

		//fmt.Println("連結")
		//fmt.Println(e.Attr("ahref"))
		url := e.Attr("href")
		//進去連結取選項

		// start scaping the page under the link found
		//e.Request.Visit(url)

		//用另一個collector去訪問各選項
		detail.Visit("https://t.weixue100.com"+url)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	detailCount := 0
	detail.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)

		//訪問新網頁前歸零計算
		questionsItemCount = 0
		detailCount = detailCount + 1
		//每訪問一題,就加一個Question的位置
		data.Questions = append(data.Questions, Question{})
	})

	//訪問各選項,取出每題的標題
	detail.OnHTML("div[id='question-title']", func(e *colly.HTMLElement) {
		// 每找到一個符合 goquerySelector字樣的結果，便會進這個OnHTML一次
		//e.Request.Visit(e.Attr("href"))
		//fmt.Println("選項題目")
		//fmt.Println(e.Text)

		data.Questions[detailCount].Topic = strings.TrimSpace(e.Text)
	})

	//訪問各選項,取出每題的每個選項的標題
	detail.OnHTML("div[class='ans-line']", func(e *colly.HTMLElement) {
		// 每找到一個符合 goquerySelector字樣的結果，便會進這個OnHTML一次
		//e.Request.Visit(e.Attr("href"))
		//fmt.Println("find")
		//fmt.Println("選項")
		//fmt.Println(e.Text)

		//增加每一個Question選項的位置
		data.Questions[detailCount].Item = append(data.Questions[detailCount].Item, "")

		data.Questions[detailCount].Item[questionsItemCount] = strings.TrimSpace(e.Text)
		questionsItemCount = questionsItemCount + 1
	})
	//取出答案
	detail.OnHTML("div[id='true-answer-box']", func(e *colly.HTMLElement) {
		// 每找到一個符合 goquerySelector字樣的結果，便會進這個OnHTML一次
		//e.Request.Visit(e.Attr("href"))
		//fmt.Println("find")
		//fmt.Println("答案")'

		//sspanClass := e.ChildAttrs("span", text)
		//spanClass :=e.ChildText("div")
		spanClass :=e.DOM.Find("div").Find("span").Eq(1).Text()
		fmt.Println(strings.TrimSpace(spanClass))
		data.Questions[detailCount].Answer = strings.TrimSpace(e.Text)


		//strings.Split(e.ChildAttr("a", "href"), "/")[4],
			//strings.TrimSpace(e.DOM.Find("span.title").Eq(0).Text())
	})
	detail.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	//c.Visit("http://go-colly.org/")
	c.Visit("https://t.weixue100.com/toefl/read/284.html")
	c.Wait()

	file, err := os.Create("myfile.txt")
	if err != nil {
		log.Fatal(err)
	}
	//f, err := os.Open("myfile.txt")

	// fmt.Println("This goes to standard output.")
	// fmt.Fprintln(f, "And this goes to the file")
	// fmt.Fprintf(f, "Also to file, with some formatting. Time: %v, line: %d\n",
	// 	time.Now(), 2)

	//矩陣還沒完成
	ItemCombine := []string{}

	for i := 0; i < len(data.Questions); i++ {
		//動態append一個空的array到ItemCombine
		ItemCombine = append(ItemCombine, "")
		for y := 0; y < len(data.Questions[i].Item); y++ {
			//將問題的選項一個一個加到ItemCombine[i]裡面
			ItemCombine[i] += data.Questions[i].Item[y] + "\n"
			//ItemCombine = append(ItemCombine, data.Questions[i].Item[y]+"\n")
		}
	}
	// for i := 0; i < len(data.Questions[1].Item); i++ {
	// 	ItemCombine += data.Questions[1].Item[i] + "\n"
	// }

	//先將title和第一題的所有內容拼起來
	allCombine := fmt.Sprintf("%v\n%v\n%v\nAnswer:%v\n", data.Title, data.Questions[0].Topic, ItemCombine[0], data.Questions[0].Answer)

	//將所有題目的所有內容拼起來
	for i := 1; i < len(data.Questions); i++ {
		//fmt.Printf("%v\n%v", i, allCombine)
		allCombine = fmt.Sprintf("%s\n%v\n%v\nAnswer:%v\n", allCombine, data.Questions[i].Topic, ItemCombine[i], data.Questions[i].Answer)
	}

	//先將title和第一題的所有內容拼起來
	allCombine = fmt.Sprintf("%v\nContent:\n%v\n", allCombine, data.Content)

	//fmt.Println(allCombine)
	fmt.Fprintf(file, allCombine)

	//mw := io.MultiWriter("sss", file)
	//fmt.Fprintln(mw, "This line will be written to stdout and also to a file")
}