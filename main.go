package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
	"strconv"
	"fmt"
)

// album represents data about a record album.
type sensor struct {
    Name     string  `json:"name"`
	Tag		 string `json:"tag"`
	Location float64  `json:"location"`
}

// sensors slice to seed record album data.
var sensors = []sensor{
    {Name: "Sensor_1", Tag: "This is a tag", Location: 30.00},
    {Name: "Sensor_2", Tag: "This is a tag_2", Location: 60.00},
    {Name: "Sensor_3", Tag: "This is a tag_3", Location: 90.00},
}

// getSensor responds with the list of all sensors as JSON.
func getSensors(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, sensors)
}

// postSensors adds an album from JSON received in the request body.
func postSensors(c *gin.Context) {
    var newSensor sensor

    // Call BindJSON to bind the received JSON to
    // newSensor.
    if err := c.BindJSON(&newSensor); err != nil {
        return
    }

    // Add the new album to the slice.
    sensors = append(sensors, newSensor)
    c.IndentedJSON(http.StatusCreated, newSensor)
}

func updateSensor(c *gin.Context) {
	name := c.Param("name")
	updatedTag := c.Param("tag")
    updatedLocation := c.Param("location")

	f1, _ := strconv.ParseFloat(updatedLocation, 64)

	fmt.Println(updatedLocation)
	fmt.Println(updatedTag)

	for i := range sensors {
		if sensors[i].Name == name {
			sensors[i].Location = f1
			sensors[i].Tag = updatedTag
			c.JSON(http.StatusOK, gin.H{"success": "sensor updated"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "sensor not found"})
	}
}


// getSensorByName locates the album whose Name value matches the name
// parameter sent by the client, then returns that album as a response.
func getSensorByName(c *gin.Context) {
    name := c.Param("name")

    // Loop over the list of sensors, looking for
    // an sensor whose ID value matches the parameter.
    for _, a := range sensors {
        if a.Name == name {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sensor not found"})
}

func main() {
    router := gin.Default()
    router.GET("/sensors", getSensors)
	router.GET("/sensors/:name", getSensorByName)
	router.POST("/sensors", postSensors)
	router.PATCH("/sensors/:name", updateSensor)
    router.Run("localhost:8080")
}