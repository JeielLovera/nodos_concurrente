package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

/* VARIABLES GLOBALES */
// registro de direcciones IP de la red
var addrs []string

// puertos de servicio
const (
	register_port     = 8000
	notification_port = 8001
)

// direccion de red del nodo
var addr_node string

/* FUNCIONES */
func localAddress() string {
	ifaces, err := net.Interfaces()

	if err != nil {
		fmt.Printf("error")
	}
	fmt.Println(ifaces)
	for _, oiface := range ifaces {
		if strings.Contains(oiface.Name, "local") {
			addrs, err := oiface.Addrs()

			if err != nil {
				fmt.Printf("error")
				continue
			}

			for _, dir := range addrs {
				switch d := dir.(type) {
				case *net.IPNet:
					if strings.HasPrefix(d.IP.String(), "192") {
						return d.IP.String()
					}
				}
			}
		}
	}
	return "127.0.0.1" // ah?
}

/* FUNCIONES COMO SERVIDOR */
func RegisterServer() {
	// ESCUCHAR, en un puerto especifico
	// hostname = ip, puerto_registro
	hostname := fmt.Sprintf("%s:%d", addr_node, register_port)
	listen, _ := net.Listen("tcp", hostname)
	defer listen.Close()

	for {
		// aceptar las conexiones
		conn, _ := listen.Accept()
		go HandleRegister(conn)
	}
}

func HandleRegister(conn net.Conn) {
	// registar en la bitacora al nuevo nodo y notifica a los demás nodos el nuevo miembro
	defer conn.Close()

	// recuperar ip del parametro
	bufferIn := bufio.NewReader(conn)
	ip, _ := bufferIn.ReadString('\n')
	ip = strings.TrimSpace(ip)

	// codificar en formato json
	bytes, _ := json.Marshal(addrs)

	// respuesta al nuevo nodo
	fmt.Fprintf(conn, "%s\n", string(bytes)) // serializar

	// notificar a los nodos
	NotifyAllNodes(ip)

	// actualizar la bitacora local
	addrs = append(addrs, ip)
	fmt.Println(addrs)
}

func NotifyAllNodes(ip string) {
	for _, addr := range addrs {
		Notify(addr, ip)
	}
}

func Notify(addr string, ip string) {
	// comunicacion
	hostremote := fmt.Sprintf("%s:%d", addr, notification_port)
	conn, _ := net.Dial("tcp", hostremote)

	defer conn.Close()

	// envia la ip al host remoto
	fmt.Fprintf(conn, "%s\n", ip)
}

func ListenNotifications() {
	// modo escuchar
	hostname := fmt.Sprintf("%s:%d", addr_node, notification_port)
	listen, _ := net.Listen("tcp", hostname)

	defer listen.Close()

	for {
		conn, _ := listen.Accept()
		go HandleNotification(conn)
	}
}

func HandleNotification(conn net.Conn) {
	defer conn.Close()

	// recuperar lo enviado en la notificacion
	bufferIn := bufio.NewReader(conn)
	ip, _ := bufferIn.ReadString('\n')
	ip = strings.TrimSpace(ip)

	// registrar ip del nuevo nodo en la bitacora local
	addrs = append(addrs, ip)
	fmt.Println(addrs)
}

/* FUNCIONES COMO CLIENTE */
func RegisterClient(hostremote string) {
	// llamada del host remoto
	remote_port := fmt.Sprintf("%s:%d", hostremote, register_port)
	conn, _ := net.Dial("tcp", remote_port)

	defer conn.Close()

	// enviar ip al host remoto
	fmt.Fprintf(conn, "%s\n", addr_node)

	// espera recibir la bitacora del hostremoto
	bufferIn := bufio.NewReader(conn)
	bitacora, _ := bufferIn.ReadString('\n')

	// decodificar
	var arrtemp []string
	json.Unmarshal([]byte(bitacora), &arrtemp)

	// actualizar bitacora local
	addrs = append(arrtemp, hostremote)
	fmt.Println(addrs)

}

/* MAIN */
func main() {
	addr_node = localAddress()
	fmt.Println("IP: ", addr_node)

	// rol de servidor (ESCUCHA)
	go RegisterServer()

	// rol de cliente
	// solicitar unirse a la red
	bufferIn := bufio.NewReader(os.Stdin)
	fmt.Printf("Ingrese ip del nodo a ingresar")

	hostremote, _ := bufferIn.ReadString('\n')
	hostremote = strings.TrimSpace(hostremote)

	// si no es el primer nodo de la red
	if hostremote != "" {
		RegisterClient(hostremote)
	}

	// rol servidor
	ListenNotifications()
}
