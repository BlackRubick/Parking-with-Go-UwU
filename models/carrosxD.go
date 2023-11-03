package models

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"image/color"
	"math/rand"
	"sync"
	"time"
)

type Estacionamiento struct {
	Espacios       []bool
	MutexCapacidad *sync.Mutex
	entrada        int
	salida         int
}

func NuevoEstacionamiento(numEspacios int) *Estacionamiento {
	return &Estacionamiento{
		Espacios:       make([]bool, numEspacios),
		MutexCapacidad: &sync.Mutex{},
	}
}

func (e *Estacionamiento) estaLleno() bool {
	for _, ocupado := range e.Espacios {
		if !ocupado {
			return false
		}
	}
	return true
}

func (e *Estacionamiento) estacionar(id int, carro *canvas.Circle) (estacionado bool) {
	e.MutexCapacidad.Lock()
	defer e.MutexCapacidad.Unlock()

	if e.estaLleno() {
		fmt.Printf("Carro %d bloqueado; estacionamiento lleno.\n", id)
		return false
	}

	for i, ocupado := range e.Espacios {
		if !ocupado {
			e.Espacios[i] = true
			carro.Move(fyne.NewPos(float32(160+20*e.entrada), 200))
			e.entrada--
			fmt.Printf("Carro %d estacionado en espacio %d.\n", id, i)
			moverVehiculo(carro, i)
			return true
		}
	}
	return false
}

func (e *Estacionamiento) VehiculoEntra(id int, wg *sync.WaitGroup, w *fyne.Container) {
	defer wg.Done()

	fmt.Printf("Carro %d llega.\n", id)
	e.entrada++

	carro := crearVehiculo(w)

	for !e.estacionar(id, carro) {
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(dormirRandom())

	e.salirVehiculo(id, carro)
}

func crearVehiculo(w *fyne.Container) *canvas.Circle {
	carro := canvas.NewCircle(color.Black)
	carro.Resize(fyne.NewSize(20, 20))
	w.Add(carro)
	carro.Move(fyne.NewPos(30, 260))
	w.Refresh()
	return carro
}

func moverVehiculo(carro *canvas.Circle, espacio int) {
	x := 200
	y := 200
	if espacio != 0 {
		x = 50*espacio + x
	}
	if espacio > 10 {
		x = 65*(espacio-10) + x
		y = 400
	}
	carro.Move(fyne.NewPos(float32(x), float32(y)))
}

func dormirRandom() time.Duration {
	min := 1
	max := 5
	return time.Duration(min+rand.Intn(max-min+1)) * time.Second
}

func (e *Estacionamiento) salirVehiculo(id int, carro *canvas.Circle) {
	e.MutexCapacidad.Lock()
	defer e.MutexCapacidad.Unlock()

	for i, ocupado := range e.Espacios {
		if ocupado {
			e.Espacios[i] = false
			break
		}
	}

	fmt.Printf("Carro %d se dirige a la salida\n", id)
	e.salida++
	carro.Move(fyne.NewPos(float32(20*e.salida), 200))

	time.Sleep(time.Second * time.Duration(rand.Intn(2)))
	fmt.Printf("Carro %d sale del estacionamiento\n", id)
	carro.Hidden = true
	e.salida--
	carro.Move(fyne.NewPos(120, 200))
	time.Sleep(200 * time.Millisecond)
	carro.Hidden = true
}
