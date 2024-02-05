package main

import (
	model "Go-Country/models"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var content1 *fyne.Container

func showFirstPage(window fyne.Window) {
	window.SetContent(content1)
}

var countries []model.Countries

func main() {
	var tabs *container.AppTabs
	var entryButton *fyne.Container

	myApp := app.New()
	myApp.Settings().SetTheme(theme.DarkTheme())
	w := myApp.NewWindow("Countries")
	w.Resize(fyne.NewSize(900, 400))
	r, _ := LoadResourceFromURLString("https://dbdzm869oupei.cloudfront.net/img/sticker/preview/1243.png")

	w.SetIcon(r)

	imageURL := "https://upload.wikimedia.org/wikipedia/commons/thumb/7/74/Mercator-projection.jpg/1280px-Mercator-projection.jpg"
	imageResource := fyne.NewStaticResource("background.jpg", loadImageFromURL(imageURL))

	background := canvas.NewImageFromResource(imageResource)
	background.FillMode = canvas.ImageFillStretch
	countryLabel := widget.NewLabel("")

	background.Translucency = 0.4
	backgroundContainer := container.New(layout.NewMaxLayout(), background)

	input := widget.NewEntry()
	input.SetPlaceHolder("name")

	entry := container.NewWithoutLayout(
		container.NewVBox(
			container.NewCenter(&fyne.Container{
				Objects: []fyne.CanvasObject{
					fyne.NewContainerWithLayout(
						&centerLayout{},
						input,
					),
				},
			}),
		),
	)

	response, err := http.Get("https://restcountries.com/v2/all")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &countries)
	if err != nil {
		panic(err)
	}
	var image *canvas.Image
	var country model.Countries

	var countryName []string
	grid := container.NewGridWithColumns(5)
	for j, _ := range countries {
		grid.Add(widget.NewLabel(countries[j].Name))
		countryName = append(countryName, countries[j].Name)
	}

	submitButton := widget.NewButton("Search", func() {
		defer input.SetText("")
		for i, _ := range countries {
			countryLabel.SetText("country")
			//fmt.Println("başarılı")
			if countries[i].Name == input.Text {
				image = getImage(countries[i].Flags.PNG)
				country.Capital = countries[i].Capital
				country.Name = countries[i].Name
				country.Currencies = countries[i].Currencies

				button := widget.NewButton("Back", func() {
					content := fyne.NewContainerWithLayout(layout.NewMaxLayout(), backgroundContainer, layout.NewSpacer(), container.New(layout.NewVBoxLayout(), container.NewCenter(container.NewHBox(entry, entryButton)), container.New(layout.NewMaxLayout(), tabs)))
					w.SetContent(content)
				})

				cardTitle := canvas.NewText(country.Name, color.White)
				cardTitle.TextSize = 30
				cardTitle.TextStyle.Bold = true
				spacer := layout.NewSpacer()
				spacer.Resize(fyne.NewSize(1, 30))

				content2 := container.NewVBox(
					container.NewCenter(image),
					container.New(layout.NewGridWrapLayout(fyne.NewSize(1, 20)), spacer),

					widget.NewLabelWithStyle(fmt.Sprintf("Name: %s\nCapital: %s\nRegion: %s\nPopulation: %d\nCurrency: %s\nLanguage: %s\n", countries[i].Name, countries[i].Capital, countries[i].Region, countries[i].Population, countries[i].Currencies[0].Name, getLanguages(countries[i].Languages)), fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: true}),

					button,
				)
				cardTitleContainer := fyne.NewContainerWithLayout(layout.NewCenterLayout(), cardTitle)
				cardContentContainer := fyne.NewContainerWithLayout(layout.NewCenterLayout(), content2)

				cardContainer := container.NewVBox(
					container.New(layout.NewGridWrapLayout(fyne.NewSize(1, 20)), spacer),
					cardTitleContainer,
					container.New(layout.NewGridWrapLayout(fyne.NewSize(1, 20)), spacer),
					cardContentContainer,
				)

				card := widget.NewCard("", "", cardContainer)

				w.SetContent(fyne.NewContainerWithLayout(layout.NewMaxLayout(), backgroundContainer, card))

				break
			}
		}
	})

	entryButton = container.NewWithoutLayout(
		container.NewVBox(
			container.NewCenter(&fyne.Container{
				Objects: []fyne.CanvasObject{
					fyne.NewContainerWithLayout(
						&centerLayout{},
						submitButton,
					),
				},
			}),
		),
	)

	tabs = getCountry()

	tabs.Resize(fyne.NewSize(100, 100))

	content := fyne.NewContainerWithLayout(layout.NewMaxLayout(), backgroundContainer, layout.NewSpacer(), container.New(layout.NewVBoxLayout(), container.NewCenter(container.NewHBox(entry, entryButton)), container.New(layout.NewMaxLayout(), tabs)))

	w.SetContent(content)
	w.ShowAndRun()

}

func getLanguages(language []model.Language) string {
	a := ""
	for k, i := range language {
		if k == len(language)-1 {
			a += i.Name
		} else {
			a += i.Name + ", "
		}

	}
	return a
}

func getImage(imageURL string) *canvas.Image {
	response, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("Resim alınamadı:", err)
	}
	defer response.Body.Close()

	imageData, _, err := image.Decode(response.Body)
	if err != nil {
		fmt.Println("Resim decode edilemedi:", err)
	}

	image := canvas.NewImageFromImage(imageData)
	image.FillMode = canvas.ImageFillOriginal
	return image
}

func getCountry() *container.AppTabs {
	var alfabe = [26]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	var delete = [9]string{"Åland Islands", "United States Minor Outlying Islands", "Virgin Islands (British)", "Virgin Islands (U.S.)", "Vatican City", "North Macedonia", "Korea", "Korea (Democratic People's Republic of)", "Korea (Republic of)"}
	tabContents := make([]*container.TabItem, len(alfabe))
	var names []string
	for i, _ := range countries {
		check := "true"
		for _, j := range delete {
			if countries[i].Name == j {
				check = "false"
				break
			}

		}
		if check == "true" {
			names = append(names, countries[i].Name)
		}

	}
	index := 0
	for y, i := range alfabe {
		grid := container.NewGridWithRows(12)
		for true {
			label := widget.NewLabel("")
			if index > len(names)-1 {
				break
			}
			x := names[index]
			if string(x[0]) == i {
				label.SetText(x)
				grid.Add(label)
				index++

			} else {

				break
			}

		}

		tabContents[y] = &container.TabItem{
			Text:    fmt.Sprintf("%s", i),
			Content: grid,
		}
	}

	pages := container.NewAppTabs(tabContents...)
	return pages
}

type centerLayout struct{}

func (l *centerLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) > 0 {
		child := objects[0]
		child.Resize(fyne.NewSize(200, 40))
		child.Move(fyne.NewPos((size.Width-child.Size().Width)/2, (size.Height-child.Size().Height)/2))
	}
}

func (l *centerLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: 200, Height: 40}
}

func loadImageFromURL(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	imageData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return imageData
}

type Resource interface {
	Name() string
	Content() []byte
}
type StaticResource struct {
	StaticName    string
	StaticContent []byte
}

func NewStaticResource(name string, content []byte) *StaticResource {
	return &StaticResource{
		StaticName:    name,
		StaticContent: content,
	}
}
func (r *StaticResource) Name() string {
	return r.StaticName
}
func (r *StaticResource) Content() []byte {
	return r.StaticContent
}

func LoadResourceFromURLString(urlStr string) (Resource, error) {
	res, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	name := filepath.Base(urlStr)
	return NewStaticResource(name, bytes), nil
}
