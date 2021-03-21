package yaml

import (
    "fmt"
    "sedb/entity"
    apperror "sedb/entity/error"
    "gopkg.in/yaml.v2"
    "os"
)
// Location of configs
const config = "config/config.yaml"

// ConfigYaml has relative structure as in file
type ConfigYaml struct{
    Server struct {
        Port string `yaml:"port"`
        Host string `yaml:"host"`
    } `yaml:"server"`
    Database struct {
        Host string `yaml:"host"`
        Port string `yaml:"port"`
        User string `yaml:"user"`
        Pass string `yaml:"pass"`
        Dbname string `yaml:"dbname"`
    } `yaml:"db"`
}

// Storage used as a configuration values
type Storage struct {
    Storage entity.Values
}

// SetConfig saves configuration in Storage
func (s *Storage) SetConfig(v entity.Values) {
    s.Storage = v
}

// GetConfig retrieves configurations
func (s Storage) GetConfig() entity.Values {
    if nil == s.Storage {
        s.initStorage()
    }

    return s.Storage
}

// GetVal return single config value
func (s Storage) GetVal(prefix string, key string) (string, bool) {
   conf := s.GetConfig()
   path := fmt.Sprintf("%s_%s", prefix, key)
   val, ok := conf[path]
   if !ok {
       return "", false
   }
   return val, true
}

// Initialize configuration storage
func (s *Storage) initStorage() {
   v := make(entity.Values)
   file := fmt.Sprintf("%s/%s", mustGetRoot(), config)
   f, err := os.Open(file)
   if err != nil {
       apperror.ProcessServErr(err)
   }
   defer f.Close()

   var fc ConfigYaml
   decoder := yaml.NewDecoder(f)
   err = decoder.Decode(&fc)
   if err != nil {
       apperror.ProcessServErr(err)
   }
   v["server_host"] = fc.Server.Host
   v["server_port"] = fc.Server.Port
   v["db_host"] = fc.Database.Host
   v["db_port"] = fc.Database.Port
   v["db_user"] = fc.Database.User
   v["db_pass"] = fc.Database.Pass
   v["db_dbname"] = fc.Database.Dbname
   s.Storage = v
}

// Return cwd location
func mustGetRoot() string {
    d, err := os.Getwd()
    if err != nil {
        panic("Can't get root directory")
    }

    return d
}
