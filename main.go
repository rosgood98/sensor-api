package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
   // "strconv"
    //"fmt"
)

// album represents data about a record album.
type sensor struct {
    Name     string  `json:"name"`
	Tag		 []string `json:"tag"`
	Location float64  `json:"location"`
}

// sensors slice to seed record album data.
var sensors = []sensor{
    {Name: "Sensor_1", Tag: []string{"tag1"}, Location: 30.00},
    {Name: "Sensor_2", Tag: []string{"tag_2"}, Location: 60.00},
    {Name: "Sensor_3", Tag: []string{"tag1", "tag2"}, Location: 90.00},
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

    sensors = append(sensors, newSensor)
    c.IndentedJSON(http.StatusCreated, newSensor)
}

// updateSensor takes a JSON and updates an already stored sensor's information with the provided
// information
func updateSensor(c *gin.Context) {
    var sensor sensor
    if err := c.ShouldBindJSON(&sensor); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request body"})
        return
    }

    // check if sensor with given name exists
    found := false
    for i := range sensors {
        if sensors[i].Name == sensor.Name {
            sensors[i].Location = sensor.Location
            sensors[i].Tag = sensor.Tag
            found = true
            break
        }
    }

    if found {
        c.JSON(200, gin.H{"success": "Sensor updated"})
    } else {
        c.JSON(404, gin.H{"error": "Sensor not found"})
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