/**
*	Autor: Martín Alejandro Pérez Güendulain
*	Correo: mperezguendulain@gmail.com
*	Descripción del Programa: Programa que busca el punto maximo Z dentro de un límite, por medio de Busqueda Armónica
*	Ejecutar el programa de la siguiente forma:
*		go run busquedaArmonica.go < config.txt
*	Se puede generar el ejecutable de la siguiente forma:
*		go build busquedaArmonica.go
*	Y se ejecutaría así:
*		./busquedaArmonica < config.txt
*	Importante: Es necesario pasarle un arhivo con la configuración del problema
*	El archivo debe contener lo siguiente y ese orden:

	Número de Iteraciónes: 100
	Tamaño de Población Nueva: 4000
	Número de Bits para representar Coordenadas X, Y: 32
	HMS: 20
	PAR: 0.1
	HMCR: 0.7
	Espacio de Busqueda:
	Xmin: -10
	Xmax: 10
	Ymin: -10
	Ymax: 10
**/

package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"time"
)

type Config struct {
	HMS                   int
	PAR                   float64
	HMCR                  float64
	NumDeBitsParaRepCoord int
	NumIteraciones        int
	TamPoblacionNueva     int
	Xmin                  int
	Xmax                  int
	Ymin                  int
	Ymax                  int
}

type Solucion struct {
	XBin  []bool
	YBin  []bool
	Costo float64
}

type Poblacion []Solucion

func (p Poblacion) Len() int           { return len(p) }
func (p Poblacion) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Poblacion) Less(i, j int) bool { return p[i].Costo > p[j].Costo }

func main() {
	rand.Seed(time.Now().Unix())
	fmt.Println("Algoritmo de Busqueda Armónica")
	var HM, HMNuevo, mejoresSoluciones []Solucion
	var mejorSol Solucion
	config := readConfig()
	mejoresSoluciones = make([]Solucion, 0, config.NumIteraciones)
	tamXTrozoX := getTamTrozo(config.NumDeBitsParaRepCoord, float64(config.Xmin), float64(config.Xmax))
	tamXTrozoY := getTamTrozo(config.NumDeBitsParaRepCoord, float64(config.Ymin), float64(config.Ymax))
	printConfig(config)
	HM = getSol(config.HMS, config.NumDeBitsParaRepCoord, tamXTrozoX, tamXTrozoY, config.Xmin, config.Ymin)
	mejorSol = getMejorSol(HM)
	for i := 0; i < config.NumIteraciones; i++ {
		HMNuevo = getSol(config.TamPoblacionNueva, config.NumDeBitsParaRepCoord, tamXTrozoX, tamXTrozoY, config.Xmin, config.Ymin)
		for j := 0; j < config.TamPoblacionNueva; j++ {
			// Para X
			if rand.Float64() < config.HMCR {
				HMNuevo[j].XBin = mejorSol.XBin
			}
			if rand.Float64() < config.PAR {
				muta(HMNuevo[j].XBin)
			}
			// Para Y
			if rand.Float64() < config.HMCR {
				HMNuevo[j].YBin = mejorSol.YBin
			}
			if rand.Float64() < config.PAR {
				muta(HMNuevo[j].YBin)
			}
		}
		HMNuevo = append(HM, HMNuevo...)
		sort.Sort(Poblacion(HMNuevo))
		mejorSol = HMNuevo[0]
		mejoresSoluciones = append(mejoresSoluciones, mejorSol)
		HM = HMNuevo[:config.HMS]
	}
	printPoblacion(mejoresSoluciones, tamXTrozoX, tamXTrozoY, config.Xmin, config.Ymin)
	generaArchivosSol(mejoresSoluciones)
}

// Función que genera los archivos HTML y JS con la solución del problema del viajero y abre el navegador google-chrome para mostrar los resultados
func generaArchivosSol(mejoresSoluciones []Solucion) {
	archivoJS := "Solucion/js/init.js"
	os.Remove(archivoJS)
	generaJSSol(mejoresSoluciones, archivoJS)
	exec.Command("google-chrome", "Solucion/index.html").Start()
}

// Función que crea el archivo Javascript con la solución del problema del agente viajero
func generaJSSol(mejoresSoluciones []Solucion, archivoJs string) {
	fdJS, errJS := os.OpenFile(archivoJs, os.O_CREATE|os.O_WRONLY, 0666)
	if errJS != nil {
		fmt.Println("Error:", errJS)
		return
	}
	defer fdJS.Close()

	fmt.Fprintf(fdJS, "mtz_avance_opt = [\n")
	for i := 0; i < len(mejoresSoluciones); i++ {
		fmt.Fprintf(fdJS, "	[%d, %f],\n", i+1, mejoresSoluciones[i].Costo)
	}
	fmt.Fprintf(fdJS, "];\n")
}

// Funcón que hace un toggle a un bit del arreglo dado
func muta(ADN []bool) {
	pos := rand.Intn(len(ADN))
	pos2 := rand.Intn(len(ADN))
	ADN[pos] = !ADN[pos]
	ADN[pos2] = !ADN[pos2]
}

// Función que retorna la mejor solución de una población dada
func getMejorSol(poblacion []Solucion) Solucion {
	mejorSol := poblacion[0]
	for i := 1; i < len(poblacion); i++ {
		if poblacion[i].Costo > mejorSol.Costo {
			mejorSol = poblacion[i]
		}
	}
	return mejorSol
}

// Devuelve el tamaño de cada trozo de los 2^NumBits-1, esto para que despues al convertir un numero binario de 16 bits solo tenga que multiplicarlo por este tamaño para obtener la coordenada que representa
func getTamTrozo(numBitsParaXY int, min, max float64) float64 {
	trozos := math.Pow(2, float64(numBitsParaXY)) - 1
	distancia := math.Abs(min - max)
	return distancia / trozos
}

// Función que imprime la población que se le pase como parametro
func printPoblacion(poblacion []Solucion, tamXTrozoX, tamXTrozoY float64, limiteXInf, limiteYInf int) {
	fmt.Println("limiteXInf:", limiteXInf)
	fmt.Println("limiteYInf:", limiteYInf)
	for i := 0; i < len(poblacion); i++ {
		x, y := getCoordReales(poblacion[i].XBin, poblacion[i].YBin, tamXTrozoX, tamXTrozoY, limiteXInf, limiteYInf)
		for j := 0; j < len(poblacion[i].XBin); j++ {
			if poblacion[i].XBin[j] {
				fmt.Print("1")
			} else {
				fmt.Print("0")
			}
		}
		fmt.Println(" =>", x)
		for j := 0; j < len(poblacion[i].YBin); j++ {
			if poblacion[i].YBin[j] {
				fmt.Print("1")
			} else {
				fmt.Print("0")
			}
		}
		fmt.Println(" =>", y)
		fmt.Println("Costo:", poblacion[i].Costo)
		fmt.Println()
	}
}

// Función que devuelve una población generada aleatoriamente
func getSol(numElem, tamNumBin int, tamXTrozoX, tamXTrozoY float64, limiteXInf, limiteYInf int) []Solucion {
	poblacion := make([]Solucion, numElem)
	for i := 0; i < numElem; i++ {
		poblacion[i].XBin = getNumBin(tamNumBin)
		poblacion[i].YBin = getNumBin(tamNumBin)
		poblacion[i].Costo = evalua(poblacion[i].XBin, poblacion[i].YBin, tamXTrozoX, tamXTrozoY, limiteXInf, limiteYInf)
	}
	return poblacion
}

// Función de evaluación, devuelve el costo de evaluar sus argumentos en la función
func evalua(numXBin, numYBin []bool, tamXTrozoX, tamXTrozoY float64, limiteXInf, limiteYInf int) float64 {
	x, y := getCoordReales(numXBin, numYBin, tamXTrozoX, tamXTrozoY, limiteXInf, limiteYInf)
	return ((math.Pow(x, 2.0) - math.Pow(y, 2.0)) * (math.Sin(x + y))) + x + y
}

// Función que regresa las coordenadas que representan X y Y en representación binaria
func getCoordReales(numXBin, numYBin []bool, tamXTrozoX, tamXTrozoY float64, limiteXInf, limiteYInf int) (float64, float64) {
	xInt := binToDec(numXBin)
	yInt := binToDec(numYBin)

	var x, y float64
	if limiteXInf < 0 {
		x = xInt*tamXTrozoX + float64(limiteXInf)
	} else {
		x = xInt*tamXTrozoX - float64(limiteXInf)
	}
	if limiteYInf < 0 {
		y = yInt*tamXTrozoY + float64(limiteYInf)
	} else {
		y = yInt*tamXTrozoY - float64(limiteYInf)
	}
	return x, y
}

// Función que convierte un número binario a decimal
func binToDec(numBin []bool) float64 {
	var num float64
	for i := 0; i < len(numBin); i++ {
		if numBin[i] {
			num += math.Pow(2, float64(i))
		}
	}
	return num
}

// Función para obtener un numero binario contruido aleatoriamente de tamaño numBits
func getNumBin(numBits int) []bool {
	numBin := make([]bool, numBits)
	for i := 0; i < numBits; i++ {
		if rand.Float64() < 0.5 {
			numBin[i] = true
		}
	}
	return numBin
}

// Función para imprimir la configuración inicial
func printConfig(config Config) {
	fmt.Printf("Número de Iteraciónes: %d\n", config.NumIteraciones)
	fmt.Printf("Tamaño de Población Nueva: %d\n", config.TamPoblacionNueva)
	fmt.Printf("Número de Bits para representar Coordenadas X, Y: %d\n", config.NumDeBitsParaRepCoord)
	fmt.Printf("HMS: %d\n", config.HMS)
	fmt.Printf("PAR: %f\n", config.PAR)
	fmt.Printf("HMCR: %f\n", config.HMCR)
	fmt.Printf("Espacio de Busqueda:\n")
	fmt.Printf("Xmin: %d\n", config.Xmin)
	fmt.Printf("Xmax: %d\n", config.Xmax)
	fmt.Printf("Ymin: %d\n", config.Ymin)
	fmt.Printf("Ymax: %d\n", config.Ymax)
}

// Función para leer la configuración inicial
func readConfig() Config {
	var config = Config{}
	fmt.Scanf("Número de Iteraciónes: %d\n", &config.NumIteraciones)
	fmt.Scanf("Tamaño de Población Nueva: %d\n", &config.TamPoblacionNueva)
	fmt.Scanf("Número de Bits para representar Coordenadas X, Y: %d\n", &config.NumDeBitsParaRepCoord)
	fmt.Scanf("HMS: %d\n", &config.HMS)
	fmt.Scanf("PAR: %f\n", &config.PAR)
	fmt.Scanf("HMCR: %f\n", &config.HMCR)
	fmt.Scanf("Espacio de Busqueda:\n")
	fmt.Scanf("Xmin: %d\n", &config.Xmin)
	fmt.Scanf("Xmax: %d\n", &config.Xmax)
	fmt.Scanf("Ymin: %d\n", &config.Ymin)
	fmt.Scanf("Ymax: %d\n", &config.Ymax)

	return config
}
