package main

import (
    "database/sql"
    "fmt"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "log"
    "flag"
    "os"
    "github.com/phungocnguyen/platformOps-EC/tree/phu-branch-1"
    _"github.com/lib/pq"
)

type Config struct {
    Dbname   string `json:"dbname"`
    Username string `json:"username"`
    Password string `json:"password"`
    Sslmode  string `json:"sslmode"`
    Location string `json:"location"`
}

func (c Config) GetDbname() string{
    return c.Dbname
}

func (c Config) GetUsername() string{
    return c.Username
}

func (c Config) GetPassword() string{
    return c.Password
}

func (c Config) GetSslmode() string{
    return c.Sslmode
}

func (c Config) GetLocation() string{
    return c.Location
}

func main() {
    var excelFileName, configFile string

        flag.StringVar(&excelFileName, "i", "", "Input excel baseline file. If missing, program will exit.")
        flag.StringVar(&configFile, "c", "", "Configuration file. If missing, program will exit.")
        flag.Parse()

        if excelFileName == "" {
                 fmt.Println("Missing input excel baseline. Program will exit.")
                 os.Exit(1)
        }

        if configFile == "" {
                 fmt.Println("Missing configuration file. Program will exit.")
                 os.Exit(1)
        }


    fmt.Println("Loading Excel file ", excelFileName)

    baseline, controls := services.LoadFromExcel(excelFileName)

    //fmt.Println(baseline)
    //fmt.Println(controls[0])


    fmt.Println("Loading config file")

    config := getConfig(configFile)

    fmt.Println("Connecting to database ")

    connStr := getConnStr(config)
    db, err := sql.Open("postgres", connStr)
      if err != nil {
       log.Fatal(err)
      }

    fmt.Println("Inserting Baseline")

    baseline_id := services.InsertBaseline(db, baseline)

    services.ReadBaselineAll(db)

    fmt.Println("Inserting controls")
    for i:=0; i<len(controls); i++ {
        controls[i].SetBaseline_id(baseline_id)
        services.InsertControl(db, controls[i])

    }

    //services.ReadControlByBaselineId(db, baseline_id)
    fmt.Println("Done inserting Baseline and Controls.  Check DB")

}

func getConfig(configFile string) Config {
    raw, err := ioutil.ReadFile(configFile)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    var c []Config
    errj := json.Unmarshal(raw, &c)
    if errj != nil {
        		fmt.Println("error parsing json input", err)
        }
    return c[0]
}

func getConnStr (config Config) string {
    var buffer bytes.Buffer
    buffer.WriteString("postgres://")
    buffer.WriteString(config.GetUsername())
    buffer.WriteString(":")
    buffer.WriteString(config.GetPassword())
    buffer.WriteString("@")
    buffer.WriteString(config.GetLocation())
    buffer.WriteString("/")
    buffer.WriteString(config.GetDbname())
    buffer.WriteString("?sslmode=")
    buffer.WriteString(config.GetSslmode())

    fmt.Println(buffer.String())

    return buffer.String()
}