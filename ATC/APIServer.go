package main

import (
	"encoding/json"
	"fmt"
	"github.com/ARMmaster17/Captain/ATC/DB"
	"github.com/gorilla/handlers"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Maintains the state of an APIServer instance. Provides caching of DB connections to speed up queries. The APIServer
// is responsible for responding to external HTTP requests, and updating the state of the Captain stack using the
// state database.
type APIServer struct {
	// HTTP router that responds to all incoming requests. Mostly stateless.
	router *mux.Router
	// Cached database connection instance to speed up queries and ease the load on the database. According to the
	// Gorm documentation, sharing an instance also provides thread-safe locking when working with local databases
	// such as Sqlite3.
	db     *gorm.DB
}

// Initializes the APIServer instance. Connects to the database and caches the connection for use by HTTP route
// handlers. CRUD routes for all REST... objects are mapped to the HTTP router.
func (a *APIServer) Start() error {
	a.router = mux.NewRouter()
	var err error
	if a.db == nil {
		a.db, err = DB.ConnectToDB()
		if err != nil {
			return fmt.Errorf("unable to connect to database: %w", err)
		}
	}
	a.registerHandlers()
	return nil
}

// Starts the HTTP APIServer on a new thread that listens on the specified port.
func (a *APIServer) Serve(port int) {
	corsAO := handlers.AllowedOrigins([]string{"*"})
	corsAM := handlers.AllowedMethods([]string{"GET","POST","PUT","DELETE","OPTIONS"})
	corsAH := handlers.AllowedHeaders([]string{"Content-Type"})
	go http.ListenAndServe(fmt.Sprintf(":%d", port), handlers.CORS(corsAO, corsAM)(a.router))
}

// Registers the HTTP REST routes for each of the 4 models stored in the state database. Each model is reponsible for
// mapping its own CRUD endpoints.
func (a *APIServer) registerHandlers() {
	a.registerAirspaceHandlers()
	a.registerFlightHandlers()
	a.registerFormationHandlers()
	//a.registerPlaneHandlers()
}

// Responds to an HTTP REST request with a generic 500 error. Good for hiding internal error messages that might pose
// a security threat, or may not be necessary to pass on to the end user. If this response is used, it is expected that
// the endpoint logic handles logging the error stack to stdout or stderr.
func (a *APIServer) respondWithError(w http.ResponseWriter) {
	a.respondWithErrorMessage(w, "Internal server error")
}

// Responds to an HTTP REST request with a 500 error code and the specified message. Good for when things break, but
// telling the end user would be both safe and helpful (e.g. validation errors). If this response is used, it is
// expected that the endpoint logic handles logging the error stack to stdout or stderr.
func (a *APIServer) respondWithErrorMessage(w http.ResponseWriter, message string) {
	a.respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": message})
}

// Responds to an HTTP REST request with a success code and the specified payload, which will be converted into JSON.
func (a *APIServer) respondOKWithJSON(w http.ResponseWriter, payload interface{}) {
	a.respondWithJSON(w, http.StatusOK, payload)
}

// Responds to an HTTP REST request with the specified code, and the specified payload, which will be converted into a
// JSON object.
func (a *APIServer) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Airspace REST handlers, functions and helpers                                                                     //
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type RESTAirspace struct {
	ID        uint
	HumanName string
	NetName   string
}

func (a *APIServer) registerAirspaceHandlers() {
	a.router.HandleFunc("/airspaces", a.getAirspaces).Methods("GET")
	a.router.HandleFunc("/airspace", a.createAirspace).Methods("POST")
	a.router.HandleFunc("/airspace/{id:[0-9]+}", a.getAirspace).Methods("GET")
	a.router.HandleFunc("/airspace/{id:[0-9]+}", a.updateAirspace).Methods("PUT")
	a.router.HandleFunc("/airspace/{id:[0-9]+}", a.deleteAirspace).Methods("DELETE")
}

// swagger:operation GET /airspaces airspace GetAirspaces
// Get all airspaces managed by this ATC instance.
// Gets a list of all airspaces stored in the state database. Does not auto-populate the Flight field.
// ---
// produces:
// - application/json
// responses:
//   '200':
//     description: Request processed
//     schema:
//       type: array
//       items:
//         properties:
//           ID:
//             type: integer
//             description: Unique airspace ID in state database.
//           HumanName:
//             type: string
//             description: Human-readable name for airspace.
//           NetName:
//             type: string
//             description: Name used for DNS name building, and internal queries against the state database.
//   '500':
//     description: Internal server error, possibly a database error.
func (a *APIServer) getAirspaces(w http.ResponseWriter, r *http.Request) {
	var airspaces []Airspace
	result := a.db.Find(&airspaces)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	var restAirspaces []RESTAirspace
	for i := range airspaces {
		restAirspaces = append(restAirspaces, RESTAirspace{
			ID:        airspaces[i].ID,
			HumanName: airspaces[i].HumanName,
			NetName:   airspaces[i].NetName,
		})
	}
	a.respondOKWithJSON(w, restAirspaces)
}

// swagger:operation POST /airspace airspace CreateAirspace
// Creates an airspace
// Creates an isolated environment for flights and formations to be provisioned.
// ---
// produces:
// - application/json
// parameters:
// - name: Airspace
//   in: body
//   description: Human-readable name for this airspace
//   schema:
//     required:
//       - HumanName
//       - NetName
//     type: object
//     properties:
//       HumanName:
//         type: string
//         description: Human-readable name for airspace.
//       NetName:
//         type: string
//         description: Name used for DNS name building, and internal queries against the state database.
// responses:
//   '201':
//     description: Request processed
//     schema:
//       type: object
//       properties:
//         ID:
//           type: integer
//           description: Unique airspace ID in state database.
//         HumanName:
//           type: string
//           description: Human-readable name for airspace.
//         NetName:
//           type: string
//           description: Name used for DNS name building, and internal queries against the state database.
//   '500':
//     description: Internal server error, possibly a database or validation error.
func (a *APIServer) createAirspace(w http.ResponseWriter, r *http.Request) {
	var as RESTAirspace
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&as); err != nil {
		a.respondWithErrorMessage(w, "Invalid Airspace object")
		return
	}
	defer r.Body.Close()
	airspace := Airspace{
		HumanName: as.HumanName,
		NetName:   as.NetName,
	}
	result := a.db.Create(&airspace)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	as.ID = airspace.ID
	a.respondWithJSON(w, http.StatusCreated, as)
}

// swagger:operation GET /airspace/{id} airspace GetAirspace
// Get an airspace managed by this ATC instance.
// Gets an airspace stored in the state database. Does not auto-populate the Flight field.
// ---
// produces:
// - application/json
// paramters:
// - name: id
//   in: path
//   schema:
//     type: integer
//   required: true
//   description: Unique ID of the airspace to get.
// responses:
//   '200':
//     description: Request processed
//     schema:
//       type: object
//       properties:
//         ID:
//           type: integer
//           description: Unique airspace ID in state database.
//         HumanName:
//           type: string
//           description: Human-readable name for airspace.
//         NetName:
//           type: string
//           description: Name used for DNS name building, and internal queries against the state database.
//   '500':
//     description: Internal server error, possibly a database error.
func (a *APIServer) getAirspace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error().Stack().Err(err).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	var airspace Airspace
	result := a.db.First(&airspace, id)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	a.respondOKWithJSON(w, RESTAirspace{
		ID:        airspace.ID,
		HumanName: airspace.HumanName,
		NetName:   airspace.NetName,
	})
}

// swagger:operation PUT /airspace/{id} airspace UpdateAirspace
// Updates an airspace.
// Updates the properties of an airspace. Note that only the HumanName can be changed after creation.
// ---
// produces:
// - application/json
// parameters:
// - name: id
//   in: path
//   schema:
//     type: integer
//   required: true
//   description: Unique ID of the airspace to get.
// - name: Airspace
//   in: body
//   description: Human-readable name for this airspace
//   schema:
//     required:
//       - HumanName
//     type: object
//     properties:
//       HumanName:
//         type: string
//         description: Human-readable name for airspace.
// responses:
//   '200':
//     description: Request processed
//     schema:
//       type: object
//       properties:
//         ID:
//           type: integer
//           description: Unique airspace ID in state database.
//         HumanName:
//           type: string
//           description: Human-readable name for airspace.
//         NetName:
//           type: string
//           description: Name used for DNS name building, and internal queries against the state database.
//   '500':
//     description: Internal server error, possibly a database or validation error.
func (a *APIServer) updateAirspace(w http.ResponseWriter, r *http.Request) {
	var as RESTAirspace
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&as); err != nil {
		a.respondWithErrorMessage(w, "Invalid Airspace object")
		return
	}
	defer r.Body.Close()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error().Stack().Err(err).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	var airspace Airspace
	result := a.db.First(&airspace, id)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	airspace.HumanName = as.HumanName
	airspace.NetName = as.NetName
	result = a.db.Save(&airspace)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	a.respondWithJSON(w, http.StatusCreated, as)
}

// swagger:operation DELETE /airspace/{id} airspace DeleteAirspace
// Get an airspace managed by this ATC instance.
// Gets an airspace stored in the state database. Does not auto-populate the Flight field.
// ---
// produces:
// - application/json
// parameters:
// - name: id
//   in: path
//   schema:
//     type: integer
//   required: true
//   description: Unique ID of the airspace to get.
// responses:
//   '200':
//     description: Request processed
//   '500':
//     description: Internal server error, possibly a database error.
func (a *APIServer) deleteAirspace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error().Stack().Err(err).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	result := a.db.Delete(&Airspace{}, id)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Flight REST handlers, functions and helpers                                                                       //
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type RESTFlight struct {
	AirspaceID uint
	ID         uint
	Name       string
}

func (a *APIServer) registerFlightHandlers() {
	a.router.HandleFunc("/flights", a.getFlights).Methods("GET")
	a.router.HandleFunc("/airspace/{aid:[0-9+]}/flights", a.getFlightsInAirspace).Methods("GET")
	a.router.HandleFunc("/flight", a.createFlight).Methods("POST")
	a.router.HandleFunc("/flight/{id:[0-9]+}", a.getFlight).Methods("GET")
	a.router.HandleFunc("/flight/{id:[0-9]+}", a.updateFlight).Methods("PUT")
	a.router.HandleFunc("/flight/{id:[0-9]+}", a.deleteFlight).Methods("DELETE")
}

func (a *APIServer) getFlights(w http.ResponseWriter, r *http.Request) {
	var flights []Flight
	result := a.db.Find(&flights)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	var restFlights []RESTFlight
	for i := range flights {
		restFlights = append(restFlights, RESTFlight{
			ID:   flights[i].ID,
			Name: flights[i].Name,
		})
	}
	a.respondOKWithJSON(w, restFlights)
}

func (a *APIServer) getFlightsInAirspace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	aid, err := strconv.Atoi(vars["aid"])
	if err != nil {
		log.Error().Stack().Err(err).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}

	var flights []Flight
	result := a.db.Where("airspace_id = ?", aid).Find(&flights)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	var restFlights []RESTFlight
	for i := range flights {
		restFlights = append(restFlights, RESTFlight{
			ID:   flights[i].ID,
			Name: flights[i].Name,
		})
	}
	a.respondOKWithJSON(w, restFlights)
}

func (a *APIServer) createFlight(w http.ResponseWriter, r *http.Request) {
	var as RESTFlight
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&as); err != nil {
		a.respondWithErrorMessage(w, "Invalid Flight object")
		return
	}
	defer r.Body.Close()
	// TODO: Verify that AirspaceID exists.
	flight := Flight{
		AirspaceID: int(as.AirspaceID),
		Name:       as.Name,
	}
	result := a.db.Create(&flight)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	as.ID = flight.ID
	a.respondWithJSON(w, http.StatusCreated, as)
}

func (a *APIServer) getFlight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error().Stack().Err(err).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	var flight Flight
	result := a.db.First(&flight, id)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	a.respondOKWithJSON(w, RESTFlight{
		AirspaceID: uint(flight.AirspaceID),
		ID:         flight.ID,
		Name:       flight.Name,
	})
}

// swagger:operation PUT /flight/{id} flight UpdateFlight
// Updates a flight.
// Updates the properties of a flight. Note that only the Name property can be changed after creation.
// ---
// produces:
// - application/json
// parameters:
// - name: id
//   in: path
//   schema:
//     type: integer
//   required: true
//   description: Unique ID of the airspace to get.
// - name: Flight
//   in: body
//   description: Human-readable name for this flight
//   schema:
//     required:
//       - Name
//     type: object
//     properties:
//       Name:
//         type: string
//         description: Human-readable name for flight.
// responses:
//   '200':
//     description: Request processed
//     schema:
//       type: object
//       properties:
//         NetName:
//           type: string
//           description: Name used for identification in API calls and other Captain tooling.
//   '500':
//     description: Internal server error, possibly a database or validation error.
func (a *APIServer) updateFlight(w http.ResponseWriter, r *http.Request) {
	var as RESTFlight
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&as); err != nil {
		a.respondWithErrorMessage(w, "Invalid Flight object")
		return
	}
	defer r.Body.Close()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error().Stack().Err(err).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	var flight Flight
	result := a.db.First(&flight, id)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	flight.Name = as.Name
	result = a.db.Save(&flight)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	a.respondWithJSON(w, http.StatusCreated, as)
}

// swagger:operation DELETE /flight/{id} flight DeleteFlight
// Deletes a flight from the state database.
// Deletes a flight, and any dependent formations and planes.
// ---
// produces:
// - application/json
// parameters:
// - name: id
//   in: path
//   schema:
//     type: integer
//   required: true
//   description: Unique ID of the flight to delete.
// responses:
//   '200':
//     description: Request processed
//   '500':
//     description: Internal server error, possibly a database error.
func (a *APIServer) deleteFlight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error().Stack().Err(err).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	result := a.db.Delete(&Flight{}, id)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Formation REST handlers, functions and helpers                                                                    //
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type RESTFormation struct {
	FlightID    int
	ID          uint
	Name        string
	CPU         int
	RAM         int
	Disk        int
	BaseName    string
	Domain      string
	TargetCount int
}

func convertToRESTFormation(as Formation) RESTFormation {
	return RESTFormation{
		FlightID:    as.FlightID,
		Name:        as.Name,
		CPU:         as.CPU,
		RAM:         as.RAM,
		Disk:        as.Disk,
		BaseName:    as.BaseName,
		Domain:      as.Domain,
		TargetCount: as.TargetCount,
	}
}

func (a *APIServer) registerFormationHandlers() {
	a.router.HandleFunc("/formations", a.getFormations).Methods("GET")
	a.router.HandleFunc("/flight/{fid:[0-9+]}/formations", a.getFormationsInFlight).Methods("GET")
	a.router.HandleFunc("/formation", a.createFormation).Methods("POST")
	a.router.HandleFunc("/formation/{id:[0-9]+}", a.getFormation).Methods("GET")
	a.router.HandleFunc("/formation/{id:[0-9]+}", a.updateFormation).Methods("PUT")
	a.router.HandleFunc("/formation/{id:[0-9]+}", a.deleteFormation).Methods("DELETE")
}

func (a *APIServer) getFormations(w http.ResponseWriter, r *http.Request) {
	var formations []Formation
	result := a.db.Find(&formations)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	var restFormations []RESTFormation
	for i := range formations {
		restFormations = append(restFormations, convertToRESTFormation(formations[i]))
	}
	a.respondOKWithJSON(w, restFormations)
}

func (a *APIServer) getFormationsInFlight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	aid, err := strconv.Atoi(vars["fid"])
	if err != nil {
		log.Error().Stack().Err(err).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}

	var formations []Formation
	result := a.db.Where("flight_id = ?", aid).Find(&formations)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	var restFormations []RESTFormation
	for i := range formations {
		restFormations = append(restFormations, convertToRESTFormation(formations[i]))
	}
	a.respondOKWithJSON(w, restFormations)
}

func (a *APIServer) createFormation(w http.ResponseWriter, r *http.Request) {
	var as RESTFormation
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&as); err != nil {
		a.respondWithErrorMessage(w, "Invalid Formation object")
		return
	}
	defer r.Body.Close()
	// TODO: Verify that FlightID exists.
	formation := Formation{
		FlightID:    as.FlightID,
		Name:        as.Name,
		CPU:         as.CPU,
		RAM:         as.RAM,
		Disk:        as.Disk,
		BaseName:    as.BaseName,
		Domain:      as.Domain,
		TargetCount: as.TargetCount,
	}
	result := a.db.Create(&formation)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	as.ID = formation.ID
	a.respondWithJSON(w, http.StatusCreated, as)
}

func (a *APIServer) getFormation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error().Stack().Err(err).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	var as Formation
	result := a.db.First(&as, id)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	a.respondOKWithJSON(w, convertToRESTFormation(as))
}

func (a *APIServer) updateFormation(w http.ResponseWriter, r *http.Request) {
	var as RESTFormation
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&as); err != nil {
		a.respondWithErrorMessage(w, "Invalid Formation object")
		return
	}
	defer r.Body.Close()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error().Stack().Err(err).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	var formation Formation
	result := a.db.First(&formation, id)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	// Only the TargetCount can be changed. When rolling releases are implemented in the future some values like CPU and RAM will also be able to be changed.
	formation.TargetCount = as.TargetCount
	result = a.db.Save(&formation)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	a.respondOKWithJSON(w, convertToRESTFormation(formation))
}

func (a *APIServer) deleteFormation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error().Stack().Err(err).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
	result := a.db.Delete(&Formation{}, id)
	if result.Error != nil {
		log.Error().Stack().Err(result.Error).Msgf("unable to process request %s", r.RequestURI)
		a.respondWithError(w)
		return
	}
}
