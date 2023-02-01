package main

import ( "net/http"
        "github.com/gin-gonic/gin"
            "io/ioutil"
            "strconv"
            "time"
            "path/filepath"
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
   
    recordingsPath := "./recordings"
    AuthorizationToken := "TOKEN_STRING_HERE"

    //setting up router
    router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.10.1"})

    //getClips entry point
	router.GET("/getClips/", func(c *gin.Context) {

        // checking for authorization token 
         if(c.GetHeader("Authorization")!=AuthorizationToken){
                c.IndentedJSON(http.StatusUnauthorized,"Unauthorized")
                return
         }

         // get timeStamp params and checks validation 
         //if fromTimestamp is empty it will be assogned with 0 
         //if untilTimestamp is empty it will be assigned with current time stamp

            fromTimeStamp := c.Query("fromTimestamp")
            untilTimeStamp := c.Query("untilTimestamp")

            if fromTimeStamp != ""{
                _, err1 := strconv.Atoi(fromTimeStamp) 
                if err1 != nil {
                    c.IndentedJSON(http.StatusUnprocessableEntity,"Please insert a valid timeStamp : fromTimestamp")
                    return
                }
            } else {
                fromTimeStamp = "0"
            }

            if untilTimeStamp != ""{
                _, err2 := strconv.Atoi(untilTimeStamp) 
                if err2 != nil {
                    c.IndentedJSON(http.StatusUnprocessableEntity,"Please insert a valid timeStamp : untilTimestamp")
                    return 
                
                }
            } else {
                untilTimeStamp = strconv.FormatInt((time.Now().Unix()),10)
            }


            var videos [] video
          
            videosDir, err := ioutil.ReadDir(recordingsPath)
	            if err != nil {
                    c.IndentedJSON(http.StatusInternalServerError,"something went wrong")
                    return
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
                            abs,_ := filepath.Abs(videoPath+"/"+d.Name())
                            currentVideo.Gpx = abs
                        }
                    }     
                    //crating clips for every video

                    var clips [] clip
                    clipPath := videoPath+"/videos"

                    files, err := ioutil.ReadDir(clipPath)
                        if err != nil {
                            c.IndentedJSON(http.StatusInternalServerError,"something went wrong")
                            return
	                     }
            
                    for _, file := range files{
                        var currentClip clip
                        absVid,_ := filepath.Abs(videoPath+"/videos/"+file.Name())
                        absThmb,_ := filepath.Abs(videoPath+"/thumbs/"+file.Name())
                        currentClip.VideoSrc = absVid
                        currentClip.ThumbSrc = absThmb
                        clips = append(clips,currentClip)
                    }

               
                    currentVideo.Clips = clips
                    videos = append(videos,currentVideo)

	            }
             }
            
            exampleResponse := response{videos}
            c.IndentedJSON(http.StatusOK,exampleResponse)
		
	})
    //Running server
	router.Run()
}
