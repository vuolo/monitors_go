package proxies

import (
  "math/rand"
)

// ########################################### UTILITY VARIABLES
var datacenterProxies = []string{
  "http://207.229.93.66:1070",
  "http://207.229.93.66:1071",
  "http://207.229.93.66:1068",
  "http://207.229.93.66:1073",
  "http://207.229.93.66:1069",
  "http://207.229.93.66:1074",
  "http://207.229.93.66:1072",
  "http://207.229.93.66:1066",
  "http://207.229.93.66:1059",
  "http://207.229.93.66:1065",
  "http://207.229.93.66:1060",
  "http://207.229.93.66:1064",
  "http://207.229.93.66:1067",
  "http://207.229.93.66:1056",
  "http://207.229.93.66:1063",
  "http://207.229.93.66:1057",
  "http://207.229.93.66:1058",
  "http://207.229.93.66:1061",
  "http://207.229.93.66:1062",
  "http://207.229.93.66:1055",
  "http://207.229.93.66:1053",
  "http://207.229.93.66:1054",
  "http://207.229.93.66:1052",
  "http://207.229.93.66:1046",
  "http://207.229.93.66:1047",
  "http://207.229.93.66:1048",
  "http://207.229.93.66:1051",
  "http://207.229.93.66:1050",
  "http://207.229.93.66:1044",
  "http://207.229.93.66:1049",
  "http://207.229.93.66:1045",
  "http://207.229.93.66:1040",
  "http://207.229.93.66:1032",
  "http://207.229.93.66:1033",
  "http://207.229.93.66:1037",
  "http://207.229.93.66:1038",
  "http://207.229.93.66:1039",
  "http://207.229.93.66:1042",
  "http://207.229.93.66:1035",
  "http://207.229.93.66:1041",
  "http://207.229.93.66:1034",
  "http://207.229.93.66:1043",
  "http://207.229.93.66:1036",
  "http://207.229.93.66:1026",
  "http://207.229.93.66:1029",
  "http://207.229.93.66:1031",
  "http://207.229.93.66:1025",
  "http://207.229.93.66:1030",
  "http://207.229.93.66:1028",
  "http://207.229.93.66:1027",
}

var residentialProxies = []string{
	"",
}

// ########################################### UTILITY FUNCTIONS
func GrabProxy() string {
	n := rand.Int() % len(datacenterProxies)
	return datacenterProxies[n]
}

func GrabResidentialProxy() string {
	n := rand.Int() % len(residentialProxies)
	return residentialProxies[n]
}