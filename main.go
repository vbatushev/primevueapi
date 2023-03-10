package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	wwwPrefix  = "https://www.primefaces.org/designer/api/primevue"
	appVersion = "primevueapi 1.0"
	usageText  = appVersion + ` — Send HTML to html2docx service
Usage of:
`
)

var (
	apiVer = "3.9.0"
)

// Section — это структура с заголовком и фрагментом элементов раздела.
// @property {string} Title - Название раздела.
// @property {[]SectionItem} Items - Это фрагмент структур SectionItem.
type Section struct {
	Title string        `json:"title"`
	Items []SectionItem `json:"items"`
}

// SectionItem — это структура с тремя полями: «Свойство», «Значение» и «Комментарий».
// @property {string} Property - Имя свойства.
// @property {string} Value - Стоимость имущества.
// @property {string} Comment - Комментарий к разделу.
type SectionItem struct {
	Property string `json:"property"`
	Value    string `json:"value"`
	Comment  string `json:"comment"`
}

func main() {
	version := flag.Bool("v", false, "version")
	flag.StringVar(&apiVer, "ver", "3.9.0", "API version")
	flag.Usage = func() {
		fmt.Printf(usageText)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *version {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	apiPage, err := url.JoinPath(wwwPrefix, apiVer)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.Get(apiPage)
	if err != nil {
		log.Fatal(err)
	}

	result, sections, err := parseContent(res)
	if err != nil {
		log.Fatal(err)
	}

	sections = sortSections(sections)

	os.WriteFile("_variables.scss", []byte(result), 0666)

	if jsonBytes, err := json.MarshalIndent(sections, "", "  "); err == nil {
		os.WriteFile("variables.json", jsonBytes, 0666)
	}
}

// Функция принимает ответ на запрос к веб-сайту MDN, анализирует HTML и возвращает строку CSS и фрагмент
// структур.
// @param resp () - Ответ на http-запрос
//
// @author Vitaly Batushev
func parseContent(resp *http.Response) (result string, sections []Section, err error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return result, sections, err
	}

	doc.Find(".main .main__section").Each(func(i int, sel *goquery.Selection) {
		section := Section{}
		head := sel.Find(".main__heading .container")
		section.Title = head.Text()

		sel.Find(".main__item").Each(func(i int, mainItem *goquery.Selection) {
			item := SectionItem{}
			code := mainItem.Find(".item__code-wrapper pre code")
			codeSplit := strings.Split(code.Text(), ":")
			if len(codeSplit) == 2 {
				item.Property = strings.TrimSpace(codeSplit[0])
				item.Value = strings.TrimSuffix(strings.TrimSpace(codeSplit[1]), ";")
			}
			desc := mainItem.Find(".item__description")
			item.Comment = desc.Text()
			section.Items = append(section.Items, item)
		})

		sections = append(sections, section)
	})

	for _, s := range sections {
		result += fmt.Sprintf("// %s\n", strings.ToUpper(s.Title))
		for _, item := range s.Items {
			result += fmt.Sprintf("\n// %s\n", item.Comment)
			result += fmt.Sprintf("%s: %s;\n", item.Property, item.Value)
		}
		result += "\n\n"
	}
	return result, sections, err
}

// Функция принимает фрагмент разделов и возвращает фрагмент разделов,
// но с разделами, отсортированными в определенном порядке.
// @param sections ([]Section) - Разделы для сортировки.
//
// @author Vitaly Batushev
// @returns Фрагмент структур Section.
func sortSections(sections []Section) []Section {
	sectionNames := []string{"general", "form", "button", "data", "panel", "overlay", "menu", "message", "media", "misc"}
	result := make([]Section, len(sectionNames))
	var unknowns []Section
	for _, section := range sections {
		var found bool
		for i, name := range sectionNames {
			if name == section.Title {
				result[i] = Section{
					Title: name,
					Items: section.Items,
				}
				found = true
				break
			}
		}
		if !found {
			unknowns = append(unknowns, section)
		}
	}

	result = append(result, unknowns...)
	return result
}
