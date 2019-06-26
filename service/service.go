package service

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	port     string
	routes   *mux.Router
	Request  *SuperRequest
	Response *SuperResponse
}
type SuperRequest struct {
	Username      string `json:"username"`
	IpAddress     string `json:"ip_address"`
	EventUUID     string `json:"event_uuid"`
	UnixTimeStamp int    `json:"unix_timestamp"`
}

type SuperResponse struct {
	CurrentGeo struct {
		Latitude  float32 `json:"lat"`
		Longitude float32 `json:"lon"`
		Radius    int     `json:"radius"`
	} `json:"currentGeo"`
	TravelToCurrentGeoSuspicious   bool `json:"travelToCurrentGeoSuspicious"`
	TravelFromCurrentGeoSuspicious bool `json:"travelFromCurrentGeoSuspicious"`
	PrecedingIpAccess              struct {
		IpAddress string  `json:"ip"`
		Speed     int     `json:"speed"`
		Latitude  float32 `json:"lat"`
		Longitude float32 `json:"lon"`
		Radius    int     `json:"radius"`
		Timestamp int32   `json:"timestamp"`
	} `json:"precedingIpAccess"`

	SubsequentIpAccess struct {
		IpAddress string  `json:"ip"`
		Speed     int     `json:"speed"`
		Latitude  float32 `json:"lat"`
		Longitude float32 `json:"lon"`
		Radius    int     `json:"radius"`
		TimeStamp int32   `json:"timestamp"`
	} `json:"subsequenIpAccess"`
}

func (srv *Server) LookupResultsHandler(w http.ResponseWriter, r *http.Request) {


}

func (srv *Server) loadRoutes() *mux.Router {

	r := mux.NewRouter()
	r.HandleFunc("/v1", srv.LookupResultsHandler)

	return r
}

func (srv *Server) StartServer(port string) error {

	srv.routes = srv.loadRoutes()

	log.Printf("Listening on port %s", port)
	return http.ListenAndServe(fmt.Sprintf("%v", port), srv.routes)
}
