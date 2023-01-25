package main

import ( "net/http"
        "github.com/gin-gonic/gin" )

type event struct{
    Name                string  `json:"name"`
    Date                int     `json:"date"`
}

type clip struct{
    Name                string  `json:"name"`
    VideoSrc            string  `json:"videoSrc"`
    ThumbSrc            string  `json:"thumbSrc"`
    Checksum            string  `json:"checksum"`
}

type video struct{
    Id                  string  `json:"id"`
    Timestamp           string  `json:"timestamp"`
    Clips               []clip  `json:"clips"`
    Gpx                 string  `json:"gpx"`
    HaveEmergencyEvent  bool    `json:"haveEmergncyEvent"` 
    Event               []event `json:"event"`
    OnSubDir            bool    `json:"onSubDir"`

}
type response struct{
    Videos              []video  `json:"videos"`  
}

func main() {
    exampleEvents := [] event {
                     { "buttonPress",
                        1629730257 }}
    
    exampleClips := [] clip { 
                    { "front",
                      "http://192.168.10.1/recordings/10510/videos/0.mp4",
                      "http://192.168.10.1/recordings/10510/thumbs/0.jpg",
                      "5f039b4ef0058a1d652f13d612375a5b" }}

    exampleVideos := [] video {
                            { "10510", 
                            "10510",
                            exampleClips,
                            "http://192.168.10.1/recordings/10510/route.geojson",
                            true,
                            exampleEvents,
                            false }}

    exampleResponse := response{exampleVideos}
    router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.10.1"})

    //getClips entry point
	router.GET("/getClips/", func(c *gin.Context) {
        if(c.GetHeader("Authorization")=="TOKEN_STRING_HERE"){
            //Authorized 
            c.IndentedJSON(http.StatusOK,exampleResponse)

        } else{
                //UnAuthorized 
                c.IndentedJSON(http.StatusUnauthorized,"Unauthorized")
        }

		
	})
    //Running server
	router.Run()
}
