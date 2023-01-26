package main // identifies that main.go is a standalone file and not part of a package

// imports required libraries
import (
    "net/http"
    "github.com/gin-gonic/gin"
    "math"
)

// creates a struct to represent a sensor's location
// each sensor has an x and y coordinate
type coordinate struct {
    X       float64   `json:"x"`
    Y       float64   `json:"y"`
}

// distance calculates the euclidean distance between two coordinates
func distance(c1 coordinate, c2 coordinate) float64 {
    // x dist
    dx := c2.X - c1.X
    // y dist
    dy := c2.Y - c1.Y
    return math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
}

// sensor represents data about a sensor
// each sensor has a name(string), tag(list of strings), and location(float64)
type sensor struct {
    Name     string  `json:"name"`
	Tag		 []string `json:"tag"`
	Location coordinate  `json:"location"`
}

// sensors slice to store initial sensor data
var sensors = []sensor {
    {Name: "Sensor_1", Tag: []string{"tag1"}, Location: coordinate{X: 60.00, Y: 90.00}},
    {Name: "Sensor_2", Tag: []string{"tag_2"}, Location: coordinate{X: 0, Y: 0}},
    {Name: "Sensor_3", Tag: []string{"tag1", "tag2"}, Location: coordinate{X: 159.12, Y: 7.13}},
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

    // binds JSON payload to a new sensor called sensor
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
func getSensorByLocation(location coordinate) (sensor, error) {

    // initializes closestSensor to first sensor in slice
    var closestSensor sensor
    closestSensor = sensors[0]
    
    var minDist float64

    // loops through slice to see if any sensors have the exact same location as the provided location
    for _, sensor := range sensors {
        if sensor.Location.X == location.X && sensor.Location.Y == location.Y {
            closestSensor = sensor
            return closestSensor, nil
        }
    }

    // sets minDist to the distance between the first sensor and location
    minDist = distance(sensors[0].Location, location)

    // loops through slice to find the closest sensor the location if none match location exactly
    for _, sensor := range sensors {
        distance := distance(sensor.Location, location)
        if minDist == 0 || distance < minDist {
            closestSensor = sensor
            minDist = distance
        } 
    }

    return closestSensor, nil
}

// sensorHandler is used with getSensorByLocation to handle a GET request
// specific to handling GET request, validating parameters, and calling getSensorByLocation
func sensorHandler(c *gin.Context) {

    var coord coordinate


	if err := c.ShouldBindJSON(&coord); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // gets closest sensor in sensor slice to the provided location in JSON payload
    closestSensor, err := getSensorByLocation(coord)

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

// deleteSensor takes in a JSON payload with a sensor name and deletes that sensor
func deleteSensor(c *gin.Context) {
    // finds and stores the sensor name from the JSON payload
    name := c.Param("name")
 
    // flag variable to check if a sensor has been found
    found := false
    
    // loops through sensor slice to find sensor with corresponding name
    for i, _ := range sensors {
        if name == sensors[i].Name {
            // removes the sensor from the slice and adjusts flag variable
            sensors = append(sensors[:i], sensors[i+1:]...)
            found = true
        }
    }
 
    // returns JSON and message depending on whether the sensor was deleted or not
    if found {
        c.IndentedJSON(http.StatusOK, gin.H{"success": "Sensor deleted"})
    } else {
        c.IndentedJSON(http.StatusOK, gin.H{"error": "Sensor not found"})
    }
   
 }

func main() {
    router := gin.Default()
    router.GET("/sensors", getSensors) // GET list of all sensors
	router.GET("/sensors/name/:name", getSensorByName) // GET specific sensor by name
    router.GET("/sensors/location", sensorHandler) // GET sensor by closest location
	router.POST("/sensors", postSensors) // POST a new sensor
	router.PATCH("/sensors/:name", updateSensor) // PATCH an existing sensors location
    router.DELETE("/sensors/:name", deleteSensor) // DELETE an existing sensor
    router.Run("localhost:8080")
}
