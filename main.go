package main

func main() {

	Addr := &ApiServer{ListenAddr: "0.0.0.0:80", InstanceOfReading: &Reading{}}
	Runserver(Addr)

}