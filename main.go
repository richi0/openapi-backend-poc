package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	generated "openapi/generated"
	"strings"

	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type DB struct {
	pets    []generated.Pet
	counter int64
}

func (d *DB) AddPet(newPet *generated.NewPet) generated.Pet {
	pet := generated.Pet{Id: d.counter, Name: newPet.Name, Tag: newPet.Tag}
	d.pets = append(d.pets, pet)
	d.counter++
	return pet
}

func (d *DB) DeletePet(id int64) {
	filteredPets := make([]generated.Pet, 0)
	for _, pet := range d.pets {
		if pet.Id != id {
			filteredPets = append(filteredPets, pet)
		}
	}
	d.pets = filteredPets
}

func NewDB() *DB {
	return &DB{pets: make([]generated.Pet, 0), counter: 0}
}

type server struct {
	db *DB
}

func (s server) AddPet(c echo.Context) error {
	var pet generated.NewPet
	err := json.NewDecoder(c.Request().Body).Decode(&pet)
	if err != nil {
		log.Error("cannot read body")
		return nil
	}
	return c.JSON(http.StatusOK, s.db.AddPet(&pet))
}

func (s server) DeletePet(c echo.Context, id int64) error {
	s.db.DeletePet(id)
	return c.String(http.StatusOK, fmt.Sprintf("Deleted %d", id))
}

func (s server) FindPetById(c echo.Context, id int64) error {
	return nil
}

func (s server) FindPets(c echo.Context, params generated.FindPetsParams) error {
	return c.JSON(http.StatusOK, s.db.pets)
}

var swaggerSpec, _ = generated.GetSwagger()

var basicAuth = middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
		return true, nil
	}
	return false, nil
})

var validation = oapimiddleware.OapiRequestValidatorWithOptions(swaggerSpec, &oapimiddleware.Options{Skipper: func(c echo.Context) bool {
	if strings.Contains(c.Request().URL.Path, "/documentation") {
		return true
	}
	return false
}})

var cors = middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOrigins: []string{"*"},
})

func main() {
	e := echo.New()
	e.Use(cors)
	e.Use(validation)
	e.File("/documentation", "documentation/index.html", basicAuth)
	generated.RegisterHandlers(e, server{NewDB()})
	e.Logger.Fatal(e.Start(":8080"))
}
