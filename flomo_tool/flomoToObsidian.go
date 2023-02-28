package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	flomoPath := "/Users/jiangchen/Downloads/flomo@江城子-20230130"

	file, err := os.OpenFile(flomoPath+"/index.html", os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(file)
	memoArr := buildMemoArr(doc)
	writeToFile(memoArr)
}

func writeToFile(memoArr []*memo) {
	pwd, _ := os.Getwd()
	err := os.Mkdir(pwd+"/flomo", 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	for _, memo := range memoArr {
		file, _ := os.OpenFile(pwd+"/flomo/"+memo.time[0:10]+".md", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)

		//写入文件时，使用带缓存的 *Writer
		write := bufio.NewWriter(file)
		write.WriteString("- " + memo.time[11:16] + " ")

		for i, v := range memo.content {
			v = strings.Trim(v, " ")
			if i == len(memo.content)-1 {
				write.WriteString(v)
			} else {
				if strings.HasPrefix(v, "#") {
					write.WriteString(v + "<br>")
				} else {
					write.WriteString(v + "<br><br>")
				}
			}
		}

		for _, v := range memo.files {
			write.WriteString("![[" + v + "]]")
		}
		write.WriteString("\n")
		//Flush将缓存的文件真正写入到文件中
		write.Flush()

	}

}

func buildMemoArr(doc *goquery.Document) []*memo {
	var memoArr []*memo
	doc.Find("div").Find(".memo").Each(func(i int, selection *goquery.Selection) {

		// if strings.Contains(selection.Find(".time").Text(), "2022-12-18") {
		curMemo := memo{}
		time := selection.Find(".time").Text()
		curMemo.time = time

		var content []string
		selection.Find(".content").Find("p").Each(func(i int, selection *goquery.Selection) {
			if selection.Text() != "" && selection.Text() != "null" {
				content = append(content, selection.Text())
			}
		})
		curMemo.content = content

		var files []string
		selection.Find(".files").Find("img").Each(func(i int, selection *goquery.Selection) {
			src, _ := selection.Attr("src")
			arr := strings.Split(src, "/")
			files = append(files, arr[3])
		})
		curMemo.files = files

		memoArr = append(memoArr, &curMemo)
		// }
	})
	return reverse(memoArr)
}

func reverse(slice []*memo) []*memo {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

type memo struct {
	time    string
	content []string
	files   []string
}
