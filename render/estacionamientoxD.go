package render

import (
	"Parking_simulator/models"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"math/rand"
	"sync"
	"time"
)

func IniciarLlegada(numVehiculos int) {
	appMain := app.New()
	rand.Seed(time.Now().UnixNano())

	park := models.NuevoEstacionamiento(20)
	ventana, contenedor := configurarVentana(appMain)
	ejecutarEstacionamiento(park, contenedor, numVehiculos)

	ventana.ShowAndRun()
}

func configurarVentana(appMain fyne.App) (fyne.Window, *fyne.Container) {
	vent := appMain.NewWindow("Estacionamiento Up xD")
	vent.Resize(fyne.NewSize(800, 500))
	vent.SetFixedSize(true)

	imgRecurso, _ := fyne.LoadResourceFromPath("assets/bg.png")
	img := canvas.NewImageFromResource(imgRecurso)
	img.Resize(fyne.NewSize(800, 500))

	cont := container.NewWithoutLayout(img)
	vent.SetContent(cont)
	return vent, cont
}

func ejecutarEstacionamiento(park *models.Estacionamiento, contenedor *fyne.Container, numVehiculos int) {
	var coord sync.WaitGroup
	go func() {
		time.Sleep(2 * time.Second)
		for i := 1; i <= numVehiculos; i++ {
			coord.Add(1)
			go park.VehiculoEntra(i, &coord, contenedor)
			time.Sleep(time.Duration(rand.ExpFloat64() * float64(time.Second)))
		}
		coord.Wait()
	}()
}
