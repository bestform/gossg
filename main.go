package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bestform/gossg/config"
	"github.com/bestform/gossg/page"
)

func main() {
	initBuildDir()

	config, err := config.ReadConfig("config.json")
	if err != nil {
		fmt.Println("Error reading config file")
		panic(err)
	}

	var toplevelPages []page.RenderedPage
	for _, pageConf := range config.Pages {
		page := page.CreatePage(pageConf, config.HeaderFile, config.FooterFile)
		toplevelPages = append(toplevelPages, page)
	}

	renderFooterLinks(toplevelPages)
	for _, renderedPage := range toplevelPages {
		if len(renderedPage.SubPages) > 0 {
			renderSubPageLinks(&renderedPage)
		}
		os.WriteFile("build/"+renderedPage.Filename, renderedPage.Content, os.ModePerm)
		for _, subPage := range renderedPage.SubPages {
			os.WriteFile("build/"+subPage.Filename, subPage.Content, os.ModePerm)
		}
	}
}

func renderSubPageLinks(renderedPage *page.RenderedPage) {
	linkTemplate, err := os.ReadFile("templates/subpage_link.html")
	if err != nil {
		fmt.Println("Error reading subpage link template")
		panic(err)
	}
	var subPageLinks []byte
	for _, subPage := range renderedPage.SubPages {
		subPageLink := strings.ReplaceAll(string(linkTemplate), "%%TITLE%%", subPage.Title)
		subPageLink = strings.ReplaceAll(subPageLink, "%%LINK%%", subPage.Filename)
		subPageLinks = append(subPageLinks, []byte(subPageLink)...)
	}

	content := strings.ReplaceAll(string(renderedPage.Content), "%%BLOGLINKS%%", string(subPageLinks))
	renderedPage.Content = []byte(content)
}

func renderFooterLinks(toplevelPages []page.RenderedPage) {
	footerLinkTemplate, err := os.ReadFile("templates/footer_link.html")
	if err != nil {
		fmt.Println("Error reading footer link template")
		panic(err)
	}

	var footerLinks []byte
	for _, renderedPage := range toplevelPages {
		footerLink := strings.ReplaceAll(string(footerLinkTemplate), "%%TITLE%%", renderedPage.Title)
		footerLink = strings.ReplaceAll(footerLink, "%%LINK%%", renderedPage.Filename)
		footerLinks = append(footerLinks, []byte(footerLink)...)
	}
	for i := 0; i < len(toplevelPages); i++ {
		content := strings.ReplaceAll(string(toplevelPages[i].Content), "%%FOOTER_LINKS%%", string(footerLinks))
		toplevelPages[i].Content = []byte(content)
		for j := 0; j < len(toplevelPages[i].SubPages); j++ {
			if len(toplevelPages[i].SubPages[j].SubPages) != 0 {
				panic("Only one level deep subpages are supported")
			}
			content := strings.ReplaceAll(string(toplevelPages[i].SubPages[j].Content), "%%FOOTER_LINKS%%", string(footerLinks))
			toplevelPages[i].SubPages[j].Content = []byte(content)
		}
	}
}

func initBuildDir() {
	err := os.RemoveAll("build")
	if err != nil {
		fmt.Println("Error removing build directory")
		panic(err)
	}
	err = os.Mkdir("build", os.ModeDir|os.ModeType|os.ModePerm)
	if err != nil {
		fmt.Println("Error creating build directory")
		panic(err)
	}

	rootFs := os.DirFS("rootfiles")
	err = os.CopyFS("build/", rootFs)
	if err != nil {
		fmt.Println("Error copying root files")
		panic(err)
	}

	assetsFs := os.DirFS("assets")
	err = os.CopyFS("build/assets", assetsFs)
	if err != nil {
		fmt.Println("Error copying assets")
		panic(err)
	}
}
