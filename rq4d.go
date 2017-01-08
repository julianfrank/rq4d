package main

import (
	"flag"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"time"
)

var (
	rdbServiceName *string // This should be the service Name used by the rqLite Cluster
	rdbHTTPPort    *string // Port for use with the http parameter
	rdbAdvHTTPPort *string // Port for use with the advhttp parameter
	rdbRAFTPort    *string // Port for use with the raft parameter
	rdbAdvRAFTPort *string // Port for use with the advraft parameter
	rdbExecutable  *string // Executable as applicable in your Container environment
	rdbDBdir       *string // DB Directory as applicable in your Container environment
)

func main() {
	rdbServiceName = flag.String("sername", "rdb", "This should be the service Name used by the rqLite Cluster")
	rdbHTTPPort = flag.String("http", ":4001", "Port for use with the http parameter")
	//rdbAdvHTTPPort = flag.String("httpadv", ":4001", "Port for use with the advhttp parameter")
	rdbRAFTPort = flag.String("raft", ":4002", "Port for use with the raft parameter")
	//rdbAdvRAFTPort = flag.String("raftadv", ":4002", "Port for use with the advraft parameter")
	rdbExecutable = flag.String("exec", "/go/bin/rqlite/rqlited", "Executable as applicable in your Container environment")
	rdbDBdir = flag.String("db", "/db", "DB Directory as applicable in your Container environment")
	flag.Parse()

	// Lookup and Build localIPv4 Interface Table
	localInterfaces := GetLocalInterfaceList()
	//Lookup and Build DNS Table
	dnsTable := GetDNSTable(*rdbServiceName)
	//Choose the IP in the Interface table that matches with the DNS Table and use that as the Local IP for advertising
	localIP := GetLocalIP(localInterfaces, dnsTable)
	//Get the List of other hosts in the network
	otherHosts := GetOtherHostList(localIP, dnsTable)
	log.Print(localInterfaces, dnsTable, localIP, otherHosts)

	cmd := new(exec.Cmd)
	switch len(dnsTable) {
	case 0:
		{
			log.Panic("Something is wrong...Host Count of ", *rdbServiceName, " is Zero!!")
		}
	case 1:
		{
			cmd = exec.Command(
				*rdbExecutable,
				"-http", localIP+*rdbHTTPPort,
				"-raft", localIP+*rdbRAFTPort,
				"-ondisk", *rdbDBdir+"/"+os.Getenv("HOSTNAME"))
			log.Println("This is the Seed Server for ", *rdbServiceName, "with IP:", localIP)
			log.Println("Seeding Cluster with Parameters : ", cmd.Args)
		}
	default:
		{
			seedIP := GetActiveMaster(otherHosts, *rdbHTTPPort)
			cmd = exec.Command(
				*rdbExecutable,
				"-http", localIP+*rdbHTTPPort,
				"-raft", localIP+*rdbRAFTPort,
				"-join", "http://"+seedIP+*rdbHTTPPort,
				"-ondisk", *rdbDBdir+"/"+os.Getenv("HOSTNAME"))
			log.Println("There are", len(otherHosts), " Other Servers in Cluster ", *rdbServiceName, ", My Server has IP:", localIP)
			log.Println("Joining Cluster with Parameters : ", cmd.Args)
		}
	}

	stdout, _ := cmd.StdoutPipe()
	go io.Copy(os.Stdout, stdout)
	stderr, _ := cmd.StderrPipe()
	go io.Copy(os.Stderr, stderr)

	err := cmd.Start()
	if err != nil {
		log.Println("Error in Executing RQLite ", err)
	}

	defer log.Println(cmd.Wait())
}

// GetActiveMaster returns the remote host that responds to tcp dial on its http port
func GetActiveMaster(otherHostList []string, httpPort string) string {
	for i := 0; i < 7; i++ {
		for _, otherHost := range otherHostList {
			log.Print("Attempting Remote Host Active Check #", i, " to ", otherHost+httpPort)
			_, err := net.Dial("tcp", otherHost+httpPort)
			if err == nil {
				return otherHost
			}
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			rest := time.Duration(r1.Intn(7)) * time.Second
			log.Println(" Remote Host Activity Failed ... Will Retry after ", rest.String())
			time.Sleep(rest)
		}
	}
	return "0.0.0.0"
}

// GetOtherHostList returns the list of other hosts excluding the host
func GetOtherHostList(localIP string, dnsTable []string) []string {
	var otherHostList []string
	for _, dnsIP := range dnsTable {
		if dnsIP != localIP {
			otherHostList = append(otherHostList, dnsIP)
		}
	}
	return otherHostList
}

// GetLocalIP returns the ip that is both in the interfacelist and the dns List
func GetLocalIP(infList []string, dnsList []string) string {
	for _, infIP := range infList {
		for _, dnsIP := range dnsList {
			if infIP == dnsIP {
				return dnsIP
			}
		}
	}
	return ""
}

// GetDNSTable returns the List of IPs Registed for this Hostname
func GetDNSTable(host string) []string {
	dnsList, err := net.LookupHost(host)
	if err != nil {
		log.Fatal("Error during DNS Host lookup -> ", err)
	}
	//log.Print(" DNS Array for ", host, " -> ", dnsList)
	return dnsList
}

// GetLocalInterfaceList returns the non loopback local IPv4 List of the host
// Code Taken from http://stackoverflow.com/a/31551220/5362821 and modified
func GetLocalInterfaceList() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(" Address Retreival Failed")
	}
	//log.Print(" This Server has the following Interface Addresses ", addrs)
	var infArray []string
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				infArray = append(infArray, ipnet.IP.String())
			}
		}
	}
	//log.Print(" Interface Array Built ", infArray)
	return infArray
}
