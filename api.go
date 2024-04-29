package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	ListenAddr string
	InstanceOfReading *Reading
	// loadRequest time.Time

}
type ApiError struct{
	Error string
}
type ApiFunc func (http.ResponseWriter,*http.Request) error
type ApiMethodFunc func (http.ResponseWriter,*http.Request) (ApiFunc,error)

type Reading struct{
	DistanceInCm float64 `json:"distanceInCm"`
}

func Runserver(s *ApiServer) {
	fmt.Println("Sever on ", s.ListenAddr)
	fmt.Println("network ip", GetDhcpIp())
	printGoApi()
	router := mux.NewRouter()
	router.HandleFunc("/",makeHttpHandlefunc(POST(s.test)))
	router.HandleFunc("/get",makeHttpHandlefunc(GET(s.test)))
	router.HandleFunc("/UltrasonicGET",makeHttpHandlefunc(GET(s.UltrasonicGET)))
	router.HandleFunc("/UltrasonicPUT",makeHttpHandlefunc(POST(s.UltrasonicPUT)))

	router.HandleFunc("/static/{id}", func(w http.ResponseWriter, r *http.Request) {
		val := mux.Vars(r)["id"]
		http.ServeFile(w, r, fmt.Sprintf("./templ/css/%v", val))
	})

	http.ListenAndServe(s.ListenAddr, router)
}

func (s *ApiServer) test(w http.ResponseWriter, r *http.Request) error{
	WriteJson(w,200,"hello")
	return nil
}
func (s *ApiServer) UltrasonicPUT(w http.ResponseWriter, r *http.Request) error{
	UltrasonicStructReading := &Reading{}
	err:= ReadJson(r,UltrasonicStructReading)
	if err!=nil{
		return err
	}
	WriteJson(w,200,UltrasonicStructReading.DistanceInCm)
	s.InstanceOfReading = UltrasonicStructReading
	return nil
}
func (s *ApiServer) UltrasonicGET(w http.ResponseWriter, r *http.Request) error{
	WriteJson(w,200,s.InstanceOfReading)
	return nil
}


func printGoApi(){
	fmt.Println(`   
 ________  ________          ________  ________  ___     
|\   ____\|\   __  \        |\   __  \|\   __  \|\  \    
\ \  \___|\ \  \|\  \       \ \  \|\  \ \  \|\  \ \  \   
 \ \  \  __\ \  \\\  \       \ \   __  \ \   ____\ \  \  
  \ \  \|\  \ \  \\\  \       \ \  \ \  \ \  \___|\ \  \ 
   \ \_______\ \_______\       \ \__\ \__\ \__\    \ \__\
    \|_______|\|_______|        \|__|\|__|\|__|     \|__|
																`)	
                                              
}

func logline(r* http.Request){
	host:=r.RemoteAddr
	log.Println(r.URL,host[:strings.Index(host,":")],r.UserAgent())
}



func WriteJson(w http.ResponseWriter,status int, args any) error{
	w.Header().Add("Content-Type","application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(args)

}
func makeHttpHandlefunc(e ApiFunc) http.HandlerFunc {
	
	return func (w http.ResponseWriter,r *http.Request){
		logline(r)
		if err:=e(w,r);err!=nil{
								//
			WriteJson(w,http.StatusBadRequest,ApiError{Error: err.Error()})					
		}

	}

}

func POST(e ApiFunc) (ApiFunc){

	return func(w http.ResponseWriter, r *http.Request) error {

		if r.Method=="POST"{
			err:= e(w,r)
			return err
		}else{
			return fmt.Errorf("invalid method POST expected")
		}

		
	}

	}
func GET(e ApiFunc) (ApiFunc){

	return func(w http.ResponseWriter, r *http.Request) error {

		if r.Method=="GET"{
			err:= e(w,r)
			return err
		}else{
			return fmt.Errorf("invalid method GET expected")
		}

		
	}

	}

func  ReadJson[T Reading](r* http.Request,e *T) (error) {
		if err := json.NewDecoder(r.Body).Decode(e); err!=nil{
		return err
	}
	defer r.Body.Close()
	return nil
}

func NewApiServerAddr(ListenAddr string) *ApiServer { //constuctor functions
	instance := &ApiServer{
		ListenAddr: ListenAddr,
	
		}
	return instance
}

func GetDhcpIp() string{
	netInterfaceAddresses, _ := net.InterfaceAddrs()
	for _,val:= range netInterfaceAddresses{
		ip := val.String()
		
		if(strings.HasPrefix(ip,"192") || strings.HasPrefix(ip,"172")|| strings.HasPrefix(ip,"10")){
			return ip[:strings.Index(ip,"/")]
		}
	}
	return ""




}



