package main // identifies that main.go is a standalone file and not part of a package

// imports required library
import (
    "net/http"
    "github.com/gin-gonic/gin"
    "math"
    "strconv"
)

// sensor represents data about a sensor
// each sensor has a name(string), tag(list of strings), and location(float64)
type sensor struct {
    Name     string  `json:"name"`
	Tag		 []string `json:"tag"`
	Location float64  `json:"location"`
}

// sensors slice to store initial sensor data
var sensors = []sensor{
    {Name: "Sensor_1", Tag: []string{"tag1"}, Location: 30.00},
    {Name: "Sensor_2", Tag: []string{"tag_2"}, Location: 60.00},
    {Name: "Sensor_3", Tag: []string{"tag1", "tag2"}, Location: 90.00},
}

// getSensor responds with the list of all sensors as JSON
// handles GET request
// *gin.Context is a object containing information about the current HTTP request
func getSensors(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, sensors)
}

// postSensors adds an sensor from JSON received in the request body
// handles POST request
func postSensors(c *gin.Context) {
    // Creates a newSensor of type sensor
    var newSensor sensor

    // Call BindJSON to bind the received JSON to newSensor
    if err := c.BindJSON(&newSensor); err != nil {
        // if the bind failed, send an Indented JSON with an error message
        // an Indented JSON is just a JSON but made more readable to humans
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // adds the new sensor to the slice of sensors
    sensors = append(sensors, newSensor)
    // sends Indented JSON with successfull message and the new sensor
    c.IndentedJSON(http.StatusCreated, newSensor)
}

// updateSensor takes a JSON and updates an already stored sensor's information with the provided information
// handles a PATCH request
func updateSensor(c *gin.Context) {
    var sensor sensor

    if err := c.ShouldBindJSON(&sensor); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Invalid request body"})
        return
    }

    // flag variable for later use
    found := false

    // loops through slice of sensors to see if the sensor we want exists
    for i := range sensors {
        if sensors[i].Name == sensor.Name {
            // if the sensor exists, update its info and break from the loop
            sensors[i].Location = sensor.Location
            sensors[i].Tag = sensor.Tag
            found = true
            break
        }
    }

    // sends a JSON and message depending on whether the sensor was found
    if found {
        c.IndentedJSON(http.StatusOK, gin.H{"success": "Sensor updated"})
    } else {
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Sensor not found"})
    }
}

// getSensorByLocation takes in a location and returns the closest sensor as well as an error
// used with sensorHandler to handle a GET request
func getSensorByLocation(location float64) (sensor, error) {
    var closestSensor sensor
    var minDist float64

    for _, sensor := range sensors {
        // calculates the distance between each sensor in the slice and the location
        distance := math.Abs(sensor.Location - location)

        if minDist == 0 || distance < minDist {
            // sets closestSensor to the sensor in the slice closest to the location
            closestSensor = sensor
            minDist = distance
        } 
    }

    return closestSensor, nil
}

// sensorHandler is used with getSensorByLocation to handle a GET request
// specific to handling GET request, validating parameters, and calling getSensorByLocation
func sensorHandler(c *gin.Context) {
    // stores location from JSON payload as a string
    location, err := strconv.ParseFloat(c.Param("location"), 64)

    if err != nil {
        // sends JSON and error message if location could not be determined
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid location"})
        return
    }

    // gets closest sensor in sensor slice to the provided location in JSON payload
    closestSensor, err := getSensorByLocation(location)

    // sends JSON and error message if closest sensor could not be found
    if err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    // sends final JSON and message with the closest sensor
    c.IndentedJSON(http.StatusOK, closestSensor)
}

// getSensorByName locates the sensor in the slice with the name we want
// parameter sent by the client, then returns that album as a response.
func getSensorByName(c *gin.Context) {
    // gets the name from the JSON payload
    name := c.Param("name")

    // Loop over the list of sensors, looking for
    // an sensor whose name value matches the name in the JSON payload
    for _, sensor := range sensors {
        if sensor.Name == name {
            // sends a JSON with a code and sensor if the correct one was found
            c.IndentedJSON(http.StatusOK, sensor)
            return
        }
    }
    // sends a JSON and error message if the sensor was not found
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sensor not found"})
}

func main() {
    router := gin.Default()
    router.GET("/sensors", getSensors)
	router.GET("/sensors/name/:name", getSensorByName)
    router.GET("/sensors/location/:location", sensorHandler)
	router.POST("/sensors", postSensors)
	router.PATCH("/sensors/:name", updateSensor)
    router.Run("localhost:8080")
}