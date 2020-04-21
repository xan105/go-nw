//go:generate goversioninfo -platform-specific=true -icon=icon.ico -manifest=manifest.xml
package main

import(
    "log"
    "os"
    "os/exec"
    "syscall"
    "path/filepath"
    "crypto/sha256"
    "encoding/json"
    "io"
    "encoding/hex"
    "io/ioutil"
    "flag"
)

func fileExist(path string) bool {
    _, err := os.Stat(path)
    if err == nil { return true }
    if os.IsNotExist(err) { return false }
    return true
}

func checkSum(filePath string) (result string, err error) {
    file, err := os.Open(filePath)
    if err != nil { return }
    defer file.Close()

    hash := sha256.New()
    _, err = io.Copy(hash, file)
    if err != nil { return }

    result = hex.EncodeToString(hash.Sum(nil))
    return
}

func main(){ 
    
    configFilePtr := flag.String("config", "nw.json", "configuration file")
    flag.Parse()
    
    self, err := os.Executable()
    if err != nil { log.Fatal(err) }
    selfPath := filepath.Dir(self)

    var configFilePath string

    if(filepath.IsAbs(*configFilePtr)) {
      configFilePath = *configFilePtr
    } else {
      configFilePath = filepath.Join(selfPath,*configFilePtr)
    }

    jsonFile, err := os.Open(configFilePath)
    if err != nil { log.Fatal(err) }
    defer jsonFile.Close()
    
    type File struct {
      Filepath  string `json:"file"`
      Sum   string `json:"sum256"`
      Size  int64 `json:"size"`
    }
    
    type Config struct {
      Binary string `json:"bin"`
      Args   string `json:"args"`
      Cwd    string `json:"cwd"`
      FileCheck []File `json:"fileCheck"`
    }
    
    b, _ := ioutil.ReadAll(jsonFile)
    
    var config Config
    
    err = json.Unmarshal(b, &config)
    if err != nil { log.Fatal(err) }

    binPath := filepath.Join(selfPath,config.Binary)
    if fileExist(binPath) {
      
      if (config.FileCheck != nil && len(config.FileCheck) > 0) {
      
        for i := 0; i < len(config.FileCheck); i++ {

          var file string
          
          if (len(config.FileCheck[i].Filepath) > 0) {
            file = filepath.Join(selfPath,config.FileCheck[i].Filepath)
          } else {
            file = binPath
          }
          stats, err := os.Stat(file)
          if err != nil { log.Fatal(err) }
          if os.IsNotExist(err) { log.Fatal(err) }
          
          if (stats.Size() != config.FileCheck[i].Size ){ log.Fatal(err) }
          
          sum, err := checkSum(file)
          if err != nil { log.Fatal(err) }
          if ( sum != config.FileCheck[i].Sum ){ log.Fatal(err) }
          
        }
      }
      
      cmd := exec.Command(binPath,config.Args)
      
      if (len(config.Cwd) > 0) {
        cmd.Dir = filepath.Join(selfPath,config.Cwd)
      } else {
        cmd.Dir = filepath.Dir(binPath)
      }

      cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
      err := cmd.Start()
      if err != nil { log.Fatal(err) }

    }
}