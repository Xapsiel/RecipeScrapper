package domain

import (
	"context"
	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"sync"
	"time"
)

type Recipe struct {
	Name        string            `json:"name"`
	Kitchen     string            `json:"kitchen"`
	Type        string            `json:"type"`
	Description map[int]string    `json:"description"`
	Ingredients map[string]string `json:"ingredients"`
	ID          int               `json:"ID"`
	ShopLink    string            `json:"shopLink"`
}

const KEY = "mongodb+srv://ilyasuseinov3301:mishka_2023@recipebook.xxu8dre.mongodb.net/?retryWrites=true&w=majority&appName=RecipeBook"

type Scrapper struct {
	parseURL  string //Страницу, которую изначально парсим
	domainURL string //Главная страница сайта
}

func NewScrapper(ParseURL string, DomainURL string) *Scrapper {
	return &Scrapper{parseURL: ParseURL, domainURL: DomainURL}
}

func (s *Scrapper) ScrapeURL(URL string) []string {
	var allURL []string
	c := colly.NewCollector()
	c.OnRequest(func(request *colly.Request) {
	})
	c.OnHTML("div.emotion-etyz2y", func(element *colly.HTMLElement) {
		urlToRecipe := element.DOM.Find(".emotion-18hxz5k").Nodes[0].Attr[0].Val
		allURL = append(allURL, s.domainURL+urlToRecipe)
	})
	c.Visit(URL)
	return allURL
}

func (s *Scrapper) ScrapeRecipe(URL string, wg *sync.WaitGroup) {
	var recipe Recipe
	var step = 0
	var name, foodType, kitchen string
	var ingredients = make(map[string]string)
	var description map[int]string = make(map[int]string)

	c := colly.NewCollector()
	c.OnHTML(".emotion-19rdt1j", func(element *colly.HTMLElement) {
		name = element.DOM.Find(".emotion-gl52ge").Nodes[0].FirstChild.Data
		upLine := element.DOM.Find(".emotion-1h6i17m")
		foodType = upLine.Nodes[1].FirstChild.Data
		kitchen = upLine.Nodes[2].FirstChild.Data

		recipe.Name = name
		recipe.Type = foodType
		recipe.Kitchen = kitchen
	})
	c.OnHTML(".emotion-1509vkh", func(element *colly.HTMLElement) {
		ing := element.DOM.Find(".emotion-ydhjlb")
		for _, elem := range ing.Nodes {
			ingredientName := elem.FirstChild.FirstChild.FirstChild.FirstChild.Data
			ingredientCount := elem.LastChild.LastChild.Data
			ingredients[ingredientName] = ingredientCount
		}
		recipe.Ingredients = ingredients

	})
	c.OnHTML(".emotion-wdt5in", func(element *colly.HTMLElement) {
		step++
		description[step] = element.Text
	})

	c.Visit(URL)
	recipe.Description = description
	recipe.ID = rand.Intn(int(time.Now().Unix()))
	s.writeRecipe(recipe, wg)

}

func (s *Scrapper) writeRecipe(recipe Recipe, wg *sync.WaitGroup) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(KEY).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	collection := client.Database("RecipeBook").Collection("recipes")
	_, err = collection.InsertOne(context.TODO(), recipe)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%v was added\n", recipe.Name)
	wg.Done()
}
