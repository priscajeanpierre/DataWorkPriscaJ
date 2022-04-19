package main

import (
	"embed"
	"fmt"
	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xuri/excelize/v2"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"image/png"
	"log"
)

//go:embed graphics/*
var EmbeddedAssets embed.FS

var counter = 0
var demoApp GuiApp
var textWidget *widget.Text

func main() {
	ebiten.SetWindowSize(900, 750)
	ebiten.SetWindowTitle("Minimal Ebiten UI Demo")

	demoApp = GuiApp{AppUI: MakeUIWindow()}

	err := ebiten.RunGame(&demoApp)
	if err != nil {
		log.Fatalln("Error running User Interface Demo", err)
	}
}

func (g GuiApp) Update() error {
	//TODO finish me
	g.AppUI.Update()
	return nil
}

func (g GuiApp) Draw(screen *ebiten.Image) {
	//TODO finish me
	g.AppUI.Draw(screen)
}

func (g GuiApp) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

type GuiApp struct {
	AppUI *ebitenui.UI
}

func MakeUIWindow() (GUIhandler *ebitenui.UI) {
	background := image.NewNineSliceColor(color.Gray16{})
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(background))
	textInfo := widget.TextOptions{}.Text("This is our first Window", basicfont.Face7x13, color.White)

	idle, err := loadImageNineSlice("button-idle.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	hover, err := loadImageNineSlice("button-hover.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	pressed, err := loadImageNineSlice("button-pressed.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	disabled, err := loadImageNineSlice("button-disabled.png", 20, 0)
	if err != nil {
		log.Fatalln(err)
	}
	buttonImage := &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}
	button := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Press Me", basicfont.Face7x13, &widget.ButtonTextColor{
			Idle: color.RGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:  30,
			Right: 30,
		}),
		// ... click handler, etc. ...
		widget.ButtonOpts.ClickedHandler(FunctionNameHere),
	)
	rootContainer.AddChild(button)
	resources, err := newListResources()
	if err != nil {
		log.Println(err)
	}

	allPopulations := loadPopulation()
	dataAsGeneric := make([]interface{}, len(allPopulations))
	for position, population := range allPopulations {
		dataAsGeneric[position] = population
	}

	listWidget := widget.NewList(
		widget.ListOpts.Entries(dataAsGeneric),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			stateName := "%s"
			fmt.Sprintf(stateName, e.(Population).State)
			return e.(Population).State
		}),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(resources.image)),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(resources.track, resources.handle),
			widget.SliderOpts.HandleSize(resources.handleSize),
			widget.SliderOpts.TrackPadding(resources.trackPadding)),
		widget.ListOpts.EntryColor(resources.entry),
		widget.ListOpts.EntryFontFace(resources.face),
		widget.ListOpts.EntryTextPadding(resources.entryPadding),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			//do something when a list item changes
			textWidget.Label = args.Entry.(Population).State
		}))
	rootContainer.AddChild(listWidget)
	textWidget = widget.NewText(textInfo)
	rootContainer.AddChild(textWidget)

	GUIhandler = &ebitenui.UI{Container: rootContainer}
	return GUIhandler
}
func loadImageNineSlice(path string, centerWidth int, centerHeight int) (*image.NineSlice, error) {
	i := loadPNGImageFromEmbedded(path)

	w, h := i.Size()
	return image.NewNineSlice(i,
			[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
			[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight}),
		nil
}

func loadPNGImageFromEmbedded(name string) *ebiten.Image {
	pictNames, err := EmbeddedAssets.ReadDir("graphics")
	if err != nil {
		log.Fatal("failed to read embedded dir ", pictNames, " ", err)
	}
	embeddedFile, err := EmbeddedAssets.Open("graphics/" + name)
	if err != nil {
		log.Fatal("failed to load embedded image ", embeddedFile, err)
	}
	rawImage, err := png.Decode(embeddedFile)
	if err != nil {
		log.Fatal("failed to load embedded image ", name, err)
	}
	gameImage := ebiten.NewImageFromImage(rawImage)
	return gameImage
}

func FunctionNameHere(args *widget.ButtonClickedEventArgs) {
	counter++
	message := fmt.Sprintf("You have pressed the button %d times", counter)
	textWidget.Label = message

}

func loadPopulation() []Population {
	excelFile, err := excelize.OpenFile("countyPopChange2020-2021.xlsx")
	if err != nil {
		log.Fatalln(err)
	}
	sliceOfStates := make([]Population, 50, 55)
	allRows, err := excelFile.GetRows("co-est2021-alldata")
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := excelFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	location := 0

	for number, row := range allRows {

		if number <= 0 {
			continue
		}
		if len(row) <= 1 {
			continue
		}
		if row[0] == "40" {
			fmt.Println(row[5])  //state
			fmt.Println(row[8])  //popEstimate2020
			fmt.Println(row[9])  //popEstimate2021
			fmt.Println(row[10]) //popChange 2020
			fmt.Println(row[11]) //popChange 2021
		}

		currentState := Population{}
		allRows(&currentState.State, &currentState.PopChange2021, &currentState.PopEstimate2021,
			&currentState.PopEstimate2020, &currentState.PopChange2020,
			&currentState.PopPerecentDiff)
		sliceOfStates[location] = currentState
		location++

	}
	return sliceOfStates
}
