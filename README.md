# Problema de maximización de una función por medio del Algoritmo de Búsqueda Armónica

Programa que busca el punto máximo Z de una función por medio de Búsqueda Armónica.

![Busqueda Armónica](/busquedaArmonica.png)


### Compilación
	go build busquedaArmonica.go

### Ejecución
	./busquedaArmonica < config.txt

#### **Importante**
Es necesario pasarle un archivo con la configuración del problema.  El archivo debe contener lo siguiente y ese orden:

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
