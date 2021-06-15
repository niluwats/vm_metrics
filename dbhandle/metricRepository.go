package dbhandle

import (
	"fmt"
	"log"
	"os"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	resourceGroupName string = os.Getenv("AZURE_GROUP_NAME")
	vmName            string = "vmfirst1"
	db                *mgo.Database
)

type Metrics struct {
	NetworkIn         map[string]interface{} `bson:"network_in_total"`
	NetworkOut        map[string]interface{} `bson:"network_out_total"`
	PercentCpu        map[string]interface{} `bson:"percantage_cpu"`
	AvailableMemBytes map[string]interface{} `bson:"available_memory_bytes"`
}

func getDbCon() *mgo.Database {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB("bethel_dashboard")
	return db
}
func SaveMetrics(m Metrics) {
	fmt.Println(m)
	if getDbCon() == nil {
		db = getDbCon()
	}
	col := db.C("organizations")

	pushQuery := bson.M{"resourcegroup.virtual_machine.$.metrics": m}
	err := col.Update(bson.M{"resourcegroup.resourcegroup_name": resourceGroupName, "resourcegroup.resourcegroup_name.virtual_machine": vmName}, bson.M{"$addToSet": pushQuery})
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
}

func GetVms() {

}
