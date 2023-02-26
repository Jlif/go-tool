package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
)

func main() {

	flomoPath := "/Users/jiangchen/Downloads/flomo@江城子-20230130"
	pwd, _ := os.Getwd()

	file, err := os.OpenFile(flomoPath+"/index.html", os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}

	doc, err := goquery.NewDocumentFromReader(file)
	doc.Find("div").Find(".memo").Each(func(i int, selection *goquery.Selection) {

		time := selection.Find(".time").Text()[0:10]
		err := os.Mkdir(pwd+"/flomo", 0777)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}

		file, err = os.Create(pwd + "/flomo/" + time + ".md")
		if os.IsExist(err) {
			file, _ = os.OpenFile(pwd+"/flomo/"+time+".md", os.O_APPEND, 0777)
		}
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}

		count := selection.Find(".content").Find("p").Length()

		file.WriteString("- " + selection.Find(".time").Text()[11:16] + " ")
		p := selection.Find(".content").Find("p")
		file.WriteString(p.First().Text() + "<br>")

		p.First().NextAll().Each(func(i int, selection *goquery.Selection) {
			if selection.Text() != "" {
				if count-1 == selection.Index() {
					file.WriteString(selection.Text())
				} else {
					file.WriteString(selection.Text() + "<br><br>")
				}
			}
		})
		file.WriteString("\n")

	})

}
