package g

import(
    "os"
    "io/ioutil" 
    "encoding/json"
    "log"
)


func SaveCacheToFile() {
    event := Events.Clone()
    body, err := json.Marshal(event)
    if err != nil {
        log.Println("SaveCacheToFile Fail:",err)
        return
    }

    if config.SaveFile == "" {
        log.Println("config savefile is empty")
        return    
    }
    err2 := ioutil.WriteFile(config.SaveFile, body, 0666)
    if err2 != nil {
        log.Println("SaveCacheToFile save event list to file Fail!",err2)   
    }
}

func ReadCacheFromFile() {
    fi,err := os.Open( config.SaveFile )    
    if err != nil {
        log.Println("ReadCacheFromFile , cache file open fail:",err)    
        return
    }

    defer fi.Close()

    var events map[string]*EventDto
    fd,err := ioutil.ReadAll(fi)
    jerr := json.Unmarshal(fd,&events)
    if err != nil || jerr != nil {
        log.Println("ReadCacheFromFile, json decode cache file content error: ", err,jerr)    
        return
    }

    if config.Debug {
        log.Println("read event from cache file is:",string(fd))    
    }

    Events.Init(events)
}
