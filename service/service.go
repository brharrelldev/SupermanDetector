package service

import (
	"encoding/json"
	"fmt"
	"github.com/brharrelldev/SupermanDetector/databases"
	"github.com/gorilla/mux"
	"github.com/umahmood/haversine"
	"io/ioutil"
	"log"
	"net/http"
)

type Server struct {
	Port        string
	routes      *mux.Router
	Request     *SuperRequest
	Response    *SuperResponse
	SupermanDBs *databases.SupermanDatabases
}
type SuperRequest struct {
	Username      string `json:"username"`
	IpAddress     string `json:"ip_address,omitempty"`
	EventUUID     string `json:"event_uuid,omitempty"`
	UnixTimeStamp int    `json:"unix_timestamp,omitempty"`
}

type SuperResponse struct {
	CurrentGeo struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
		Radius    int     `json:"radius"`
	} `json:"currentGeo"`
	TravelToCurrentGeoSuspicious   bool `json:"travelToCurrentGeoSuspicious"`
	TravelFromCurrentGeoSuspicious bool `json:"travelFromCurrentGeoSuspicious"`
	PrecedingIpAccess              struct {
		IpAddress string  `json:"ip"`
		Speed     float64    `json:"speed"`
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
		Radius    int     `json:"radius"`
		Timestamp int32   `json:"timestamp"`
	} `json:"precedingIpAccess"`

	SubsequentIpAccess struct {
		IpAddress string  `json:"ip"`
		Speed     float64    `json:"speed"`
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
		Radius    int     `json:"radius"`
		TimeStamp int32   `json:"timestamp"`
	} `json:"subsequenIpAccess"`
}

func (srv *Server) LookupResultsHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string

	incoming, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(incoming, &requestBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username := requestBody["username"]

	loginDB, err := databases.NewDBClient(&databases.SupermanDatabases{
		LoginDBClient: &databases.LoginDBClient{
			DBFile: srv.SupermanDBs.LoginDBClient.DBFile,
		},
		MMDB: databases.MMDB{
			MMFile: srv.SupermanDBs.MMDB.MMFile,
		},
	})

	if err != nil {
		log.Fatalf("could unknown %v", err)
	}

	locations, err := loginDB.RetrieveEventsByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	resp := &SuperResponse{}
	for idx, location  :=  range locations{
		results, err  := loginDB.LookupIp(location["ipaddress"].(string))
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch idx {
		case 0:
			resp.PrecedingIpAccess.Latitude = results.Location.Latitude
			resp.PrecedingIpAccess.Longitude = results.Location.Longitude
			resp.PrecedingIpAccess.Timestamp = locations[0]["timestamp"].(int32)
		case 1:
			resp.CurrentGeo.Latitude = results.Location.Longitude
			resp.CurrentGeo.Longitude = results.Location.Latitude
			resp.CurrentGeo.Radius = 20

		case 2:
			resp.SubsequentIpAccess.Latitude = results.Location.Latitude
			resp.SubsequentIpAccess.Longitude = results.Location.Longitude
			resp.SubsequentIpAccess.TimeStamp = locations[2]["timestamp"].(int32)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	current := haversine.Coord{Lat: resp.CurrentGeo.Latitude, Lon: resp.CurrentGeo.Longitude}
	preceding := haversine.Coord{Lat: resp.PrecedingIpAccess.Latitude, Lon: resp.PrecedingIpAccess.Longitude}
	subsequent := haversine.Coord{Lat: resp.SubsequentIpAccess.Latitude, Lon: resp.SubsequentIpAccess.Longitude}



	travelToDistance, _ := haversine.Distance(preceding, current)
	travelFromDistance, _ := haversine.Distance(current, subsequent)

	resp.PrecedingIpAccess.Speed = travelFromDistance
	resp.SubsequentIpAccess.Speed = travelToDistance

	if travelFromDistance > 500{
		resp.TravelFromCurrentGeoSuspicious = true
	}

	if travelToDistance > 500{
		resp.TravelToCurrentGeoSuspicious = true
	}

	respBody, err := json.Marshal(resp)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respBody)




}

func (srv *Server) loadRoutes() *mux.Router {

	r := mux.NewRouter()
	r.HandleFunc("/v1", srv.LookupResultsHandler)

	return r
}

func (srv *Server) StartServer() error {

	srv.routes = srv.loadRoutes()

	log.Printf("Listening on port %s", srv.Port)
	return http.ListenAndServe(fmt.Sprintf(":%s", srv.Port), srv.routes)
}
