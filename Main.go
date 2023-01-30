package main

import ( "net/http"
        "github.com/gin-gonic/gin"
            "fmt"
            "io/ioutil"
            "strconv"
            "time"
            // "github.com/codingsince1985/checksum"
        )

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

    //examples 
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

////////////////////////////////////////////////////////////
   
    baseUrl := "http://192.168.10.1"
    recordingsPath := "./recordings"

    //setting up router
    router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.10.1"})

    //getClips entry point
	router.GET("/getClips/", func(c *gin.Context) {
         if(c.GetHeader("Authorization")=="TOKEN_STRING_HERE"){
            //Authorized 
          
            // get timeStamp params and checks validation 
            //if fromTimestamp is empty it will be assogned with 0 
            //if untilTimestamp is empty it will be assigned with current time stamp

            fromTimeStamp := c.Query("fromTimestamp")
            untilTimeStamp := c.Query("untilTimestamp")

            if fromTimeStamp != ""{
                _, err1 := strconv.Atoi(fromTimeStamp) 
                if err1 != nil {
                    c.IndentedJSON(http.StatusUnprocessableEntity,"Please insert a valid timeStamp")
                    return
                }
            } else{
                fromTimeStamp = "0"
            }

            if untilTimeStamp != ""{
                _, err2 := strconv.Atoi(untilTimeStamp) 
                if err2 != nil {
                    c.IndentedJSON(http.StatusUnprocessableEntity,"Please insert a valid timeStamp")
                    return 
                
                }
            } else{
                untilTimeStamp = strconv.FormatInt((time.Now().Unix()),10)
            }


            var videos [] video
          
            videosDir, err := ioutil.ReadDir(recordingsPath)
	            if err != nil {
		             fmt.Println(err)
	            }

            // Creating a videos object for every directory

	        for _, vid := range videosDir {
               
               if vid.Name()>fromTimeStamp && vid.Name()<untilTimeStamp {
                    var currentVideo video
                    currentVideo.Id = vid.Name()
                    currentVideo.Timestamp = vid.Name()
           
                    videoPath := recordingsPath+"/"+vid.Name()
                    dirs, err := ioutil.ReadDir(videoPath)
	                if err != nil {
                        c.IndentedJSON(http.StatusInternalServerError,"something went wrong")
                        return
	                }
                    for _, d := range dirs{
                        if !d.IsDir() {
                            currentVideo.Gpx = baseUrl+videoPath+"/"+d.Name()
                        }
                    }     
                    //crating clips for every video

                    var clips [] clip
                    clipPath := videoPath+"/videos"

                    subDirs, err := ioutil.ReadDir(clipPath)
                        if err != nil {
                            c.IndentedJSON(http.StatusInternalServerError,"something went wrong")
                            return
	                     }
            
                    for _, sb := range subDirs{
                        var currentClip clip
                        currentClip.VideoSrc = baseUrl+videoPath+"/videos/"+sb.Name()
                        currentClip.ThumbSrc = baseUrl+videoPath+"/thumbs/"+sb.Name()
                        clips = append(clips,currentClip)
                    }

               
                    currentVideo.Clips = clips
                    videos = append(videos,currentVideo)

	            }
             }
            
            exampleResponse = response{videos}
            c.IndentedJSON(http.StatusOK,exampleResponse)

        } else {
            //UnAuthorized 
            c.IndentedJSON(http.StatusUnauthorized,"Unauthorized")
        }

		
	})
    //Running server
	router.Run()
}
