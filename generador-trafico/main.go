package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var wg sync.WaitGroup
var url string

type Caso struct {
	Name         string `json:"name"`
	Location     string `json:"location"`
	Age          int    `json:"age"`
	Infectedtype string `json:"infectedtype"`
	State        string `json:"state"`
}

type CasosContenedor struct {
	Casos []Caso `json:"Casos"`
}

func (t Caso) toString() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(bytes)
}

func getCasos(path string) CasosContenedor {
	var casoContenedor CasosContenedor
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	json.Unmarshal(raw, &casoContenedor)
	return casoContenedor
}

func main() {
	var gorutinas int
	var solicitudes int
	var path string

	fmt.Println("Ingrese la url:")
	fmt.Scanf("%s", &url)
	fmt.Println("Cantidad de gorutinas a utilizar:")
	fmt.Scanf("%d", &gorutinas)
	fmt.Println("Cantidad de solicitudes:")
	fmt.Scanf("%d", &solicitudes)
	fmt.Println("Ruta del archivo:")
	fmt.Scanf("%s", &path)

	if gorutinas > solicitudes {
		fmt.Println("Las gorutinas no pueden ser mayor a las solicitudes")
		return
	}

	casos := getCasos(path)

	if solicitudes > len(casos.Casos) {
		fmt.Println("Las solicitudes son mayores a las contenidas en el archivo")
		fmt.Println("Se enviaran el total de solicitudes que contiene el archivo.")
		solicitudes = len(casos.Casos)
	}

	indice := 0
	rango := solicitudes / gorutinas
	faltante := (solicitudes % gorutinas)

	wg.Add(gorutinas)

	for gorutinas > 0 {
		indiceFinal := 0
		if gorutinas == 1 {
			indiceFinal = indice + rango + faltante
		} else {
			indiceFinal = indice + rango
		}
		go enviarCasos(casos.Casos, indice, indiceFinal)

		if indiceFinal == solicitudes {
			break
		}
		indice += rango
		gorutinas--
	}
	wg.Wait()
	fmt.Println("===================Terminando programa===================")
}

func enviarCasos(casos []Caso, indiceInicial int, indiceFinal int) {
	defer wg.Done()
	for indiceInicial < indiceFinal {
		push(casos[indiceInicial])
		indiceInicial += 1
	}
}

func push(caso Caso) {
	jsonReq, err := json.Marshal(caso)
	if url == "" {
		url = "http://casos.covid19so1.tk/caso"
	}
	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
}
