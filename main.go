package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	generated "openapi/generated"
	"strings"

	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type server struct {
}

func (s server) AddPet(c echo.Context) error {
	var pet generated.NewPet
	err := json.NewDecoder(c.Request().Body).Decode(&pet)
	if err != nil {
		log.Error("cannot read body")
		return nil
	}
	fmt.Println(pet.Name)
	return nil
}

func (s server) DeletePet(c echo.Context, id int64) error {
	return nil
}

func (s server) FindPetById(c echo.Context, id int64) error {
	return nil
}

func (s server) FindPets(c echo.Context, params generated.FindPetsParams) error {
	return nil
}

var basicAuth = middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
		return true, nil
	}
	return false, nil
})

func main() {
	e := echo.New()
	swaggerSpec, err := generated.GetSwagger()
	if err != nil {
		panic("Cannot load swaggerSpec")
	}
	validation := oapimiddleware.OapiRequestValidatorWithOptions(swaggerSpec, &oapimiddleware.Options{Skipper: func(c echo.Context) bool {
		if strings.Contains(c.Request().URL.Path, "/documentation") {
			return true
		}
		return false
	}})
	e.Use(validation)
	e.File("/documentation", "documentation/index.html", basicAuth)
	generated.RegisterHandlers(e, server{})
	e.Logger.Fatal(e.Start(":8080"))
}
