package app

import (
	"fmt"
	"scrapper/domain"
	"strconv"
	"sync"
	"time"
)

func ScrapeRecipe(ParseURL, DomainURL string, pages int) {
	wg := sync.WaitGroup{}
	scrapper := domain.NewScrapper(ParseURL, DomainURL)
	st := time.Now().UnixMicro()
	for page := 1; page <= pages; page++ {
		allURL := scrapper.ScrapeURL(ParseURL + strconv.Itoa(page))
		for index := range allURL {
			wg.Add(1)
			go scrapper.ScrapeRecipe(allURL[index], &wg)
		}
		time.Sleep(time.Second)
		fmt.Println(page, "страница")
	}
	end := time.Now().UnixMicro()
	fmt.Printf("%d страниц было спаршено за %f секунд\n", pages, float64(end-st)/1000000.0)
	wg.Wait()

}
