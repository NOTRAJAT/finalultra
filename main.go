package main

func main() {

	Addr := &ApiServer{ListenAddr: "127.0.0.1:80", InstanceOfReading: &Reading{}}
	Runserver(Addr)

}