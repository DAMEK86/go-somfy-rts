package shutters

import (
	"net/http"
	"strconv"
)

type Controller interface {
	CreateEndpoints(router *http.ServeMux)
}

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{
		service: service,
	}
}

func (c *controller) CreateEndpoints(router *http.ServeMux) {
	router.HandleFunc("/up", c.up)
	router.HandleFunc("/down", c.down)
	router.HandleFunc("/my", c.my)
	router.HandleFunc("/program", c.program)
}

func (c *controller) up(w http.ResponseWriter, r *http.Request) {
	addr, err := c.getAddressFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.service.MoveUp(addr, 1)
	w.WriteHeader(http.StatusOK)
}

func (c *controller) down(w http.ResponseWriter, r *http.Request) {
	addr, err := c.getAddressFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.service.MoveDown(addr, 1)
	w.WriteHeader(http.StatusOK)
}

func (c *controller) my(w http.ResponseWriter, r *http.Request) {
	addr, err := c.getAddressFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.service.MoveMy(addr, 1)
	w.WriteHeader(http.StatusOK)
}

func (c *controller) program(w http.ResponseWriter, r *http.Request) {
	addr, err := c.getAddressFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.service.Program(addr, 1)
	w.WriteHeader(http.StatusOK)
}

func (c *controller) getAddressFromRequest(r *http.Request) (uint32, error) {
	addrString := r.FormValue("addr")
	addr, err := strconv.ParseUint(addrString, 10, 24)
	if err != nil {
		return 0, err
	}
	return uint32(addr), nil
}
