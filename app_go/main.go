package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	//"github.com/shirou/gopsutil/process"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/passwd/user"
)

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "index.html")
}

func obtenerEstado(caracter string) (estado string) {
	if caracter == "0" {
		cantidadRunning += 1
		return "Running"
	} else if caracter == "1" || caracter == "2" {
		cantidadSleeping += 1
		return "Sleeping"
	} else if caracter == "8" {
		cantidadStoped += 1
		return "Stop"
	} else if caracter == "32" {
		cantidadZombie += 1
		return "Zombie"
	} else if caracter == "W" {
		return "Wait"
	} else if caracter == "L" {
		return "Lock"
	}
	return "Estado indefinido"
}

type PROCESO struct {
	PID        string
	Usuario    string
	Estado     string
	Memoria    string
	Nombre     string
	PID_Parent string
}
type struct_datos struct {
	TotalProcesos    int
	TotalEjecucion   int
	TotalSuspendidos int
	TotalDetenidos   int
	TotalZombie      int
	Procesos         []PROCESO
}

var cantidadRunning, cantidadSleeping, cantidadStoped, cantidadZombie int
var arreglo_procesos []PROCESO

func datosProcesosHandler(response http.ResponseWriter, request *http.Request) {
	data_memoria, _ := ioutil.ReadFile("/proc/procs_grupo24")
	archivo := string(data_memoria)
	lineas := strings.Split(archivo, "\n")
	cantidadRunning = 0
	cantidadSleeping = 0
	cantidadStoped = 0
	cantidadZombie = 0
	arreglo_procesos = nil
	lista_procesos := lineas
	for i := 1; i < len(lista_procesos)-1; i++ {
		proceso_unico := strings.Split(lista_procesos[i], ",")
		usuario := proceso_unico[3]
		u, _ := user.LookupId(proceso_unico[3])
		usuario = u.Username
		estado := proceso_unico[2]
		memoria := proceso_unico[5]
		nombre := proceso_unico[1]
		arreglo_procesos = append(arreglo_procesos, PROCESO{PID: proceso_unico[0], Usuario: usuario, Estado: obtenerEstado(estado), Memoria: memoria, Nombre: nombre, PID_Parent: proceso_unico[4]})
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	datos := struct_datos{
		TotalProcesos:    len(arreglo_procesos),
		TotalEjecucion:   cantidadRunning,
		TotalSuspendidos: cantidadSleeping,
		TotalDetenidos:   cantidadStoped,
		TotalZombie:      cantidadZombie,
		Procesos:         arreglo_procesos,
	}
	datos_json, _ := json.Marshal(datos)
	response.Write(datos_json)
}
func killProcesoHandler(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "index.html") //Cargo nuevamente la página principal
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func memoriaHandler(response http.ResponseWriter, request *http.Request) {

	http.ServeFile(response, request, "memoria.html")
}
func datosmemoriaHandler(response http.ResponseWriter, request *http.Request) {
	data_memoria, _ := ioutil.ReadFile("/proc/mem_grupo24")
	archivo := string(data_memoria)
	lineas := strings.Split(archivo, "\n")
	linea_memtotal := strings.Split(lineas[0], ",")
	mem_total, _ := strconv.Atoi(strings.TrimSpace(linea_memtotal[1]))
	linea_memusada := strings.Split(lineas[1], ",")
	mem_usada, _ := strconv.Atoi(strings.TrimSpace(linea_memusada[1]))
	var total, consumida, porcentaje_consumo, megabytes float64
	megabytes = 1024
	total = (float64)(mem_total) / megabytes     //Tamaño en MB
	consumida = (float64)(mem_usada) / megabytes //Tamaño en MB
	porcentaje_consumo = ((consumida * 100) / total)
	//fmt.Println(total)
        //fmt.Println(consumida)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	type MEMORIA struct {
		Total      float64
		Consumida  float64
		Porcentaje float64
	}
	datos := MEMORIA{Total: total, Consumida: consumida, Porcentaje: porcentaje_consumo}
	datos_json, _ := json.Marshal(datos)
	response.Write(datos_json)
}

var validPath = regexp.MustCompile("^/(kill|save|view)/([a-zA-Z0-9]+)$")

func killHandler(w http.ResponseWriter, r *http.Request, pid string) {
	arg0 := pid
	val0, _ := strconv.ParseInt(arg0, 10, 32)
	pidfind := int(val0)
	//signal := 9
	syscall.Kill(pidfind, 9)
	http.Redirect(w, r, "/", http.StatusFound)
}
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

var router = mux.NewRouter()

func main() {
	router.HandleFunc("/", indexPageHandler)                  //Página principal de la aplicación
	router.HandleFunc("/datosprocesos", datosProcesosHandler) //Página principal de la aplicación
	router.HandleFunc("/memoria", memoriaHandler)
	router.HandleFunc("/datosmemoria", datosmemoriaHandler)
	http.HandleFunc("/kill/", makeHandler(killHandler))
	http.Handle("/", router)
	fmt.Println("Servidor corriendo en http://localhost:8081/")
	http.ListenAndServe(":8081", nil)
}
