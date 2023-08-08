package main

import (
	"JurassicPark/core"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

func main() {
	router := gin.Default()

	db := core.Connect()
	defer db.Close()

	router.GET("/dinosaurs", func(c *gin.Context) {
		species := c.Query("species")
		dinosaurs, err := core.GetDinosaurs(species, db)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"dinosaurs": dinosaurs})
		}
	})

	router.GET("/cages", func(c *gin.Context) {
		powerStatus := c.Query("powerStatus")
		cages, err := core.GetCages(powerStatus, db)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"cages": cages})
		}
	})

	router.GET("/cages/:number", func(c *gin.Context) {
		number, err := strconv.Atoi(c.Param("number"))
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid cage number")
			return
		}
		cage, err := core.GetCage(number, db)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, cage)
		}
	})

	router.PUT("/cages/:number/:powerStatus", func(c *gin.Context) {
		powerStatus := c.Param("powerStatus")
		number, err := strconv.Atoi(c.Param("number"))
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid cage number")
			return
		}

		dinosaurs, err := core.GetDinosaursInCage(number, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(dinosaurs) > 0 {
			c.String(http.StatusBadRequest, "Cannot change the state of a cage with dinosaurs in it")
			return
		}

		if err := core.UpdateCage(number, powerStatus, db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.String(http.StatusOK, "Cage %d successfully set to %s", number, powerStatus)
		}
	})

	router.POST("/dinosaurs", func(c *gin.Context) {
		var dinosaur core.Dinosaur
		if err := c.BindJSON(&dinosaur); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cage, err := core.GetCage(dinosaur.Cage, db)
		if err != nil {
			if err == sql.ErrNoRows {
				// Given more time, I would implement a custom error type (e.x. NoCageFoundError) and have
				// core.GetCage raise that to avoid having sql implementation details here.
				c.JSON(http.StatusNotFound, gin.H{"error": "Cage not found"})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		if err = core.DestinationCageIsValid(dinosaur, cage, db); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error assigning dinosaur to cage": err.Error()})
			return
		}

		if err := core.CreateDinosaur(dinosaur, db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, dinosaur)
		}
	})

	router.GET("/dinosaurs/:name", func(c *gin.Context) {
		name := c.Param("name")
		dinosaur, err := core.GetDinosaur(name, db)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Dinosaur not found"})
				return
			}
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, dinosaur)
		}
	})

	router.POST("/cages", func(c *gin.Context) {
		var cage core.Cage
		if err := c.BindJSON(&cage); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := core.CreateCage(cage, db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, cage)
		}
	})

	router.PUT("/dinosaurs/:name/:cage", func(c *gin.Context) {
		name := c.Param("name")
		cageNumber, err := strconv.Atoi(c.Param("cage"))
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid cage number")
			return
		}

		dinosaur, err := core.GetDinosaur(name, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cage, err := core.GetCage(cageNumber, db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err = core.DestinationCageIsValid(dinosaur, cage, db); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error assigning dinosaur to cage": err})
			return
		}

		if err := core.MoveDinosaur(name, cageNumber, db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.String(http.StatusOK, "%s successfully moved to cage %d", name, cageNumber)
		}
	})

	router.GET("/cages/:number/dinosaurs", func(c *gin.Context) {
		number, err := strconv.Atoi(c.Param("number"))
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid cage number")
			return
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dinosaurs, err := core.GetDinosaursInCage(number, db)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"dinosaurs": dinosaurs})
		}
	})

	router.Run() // Port 8080 by default
}
