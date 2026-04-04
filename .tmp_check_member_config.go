package main

import (
  "fmt"
  members "balance/app/modules/members"
  appConf "balance/config"
  "balance/internal/config"
)

func main() {
  confMod := config.New(&appConf.App)
  mconf := config.Conf[members.Config](confMod.Svc)
  c := mconf.Val
  fmt.Println("CID empty:", c.GoogleClientID == "")
  fmt.Println("Secret empty:", c.GoogleClientSecret == "")
  fmt.Println("Redirect empty:", c.GoogleRedirectURL == "")
  fmt.Println("CID:", c.GoogleClientID)
  fmt.Println("Redirect:", c.GoogleRedirectURL)
}
