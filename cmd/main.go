package main

import "scrapper/app"

func main() {
	app.ScrapeRecipe("https://eda.ru/recepty?page=", "https://eda.ru", 250)
	//repo.UpdateDB(domain.Recipe{Name: "d"})
}
