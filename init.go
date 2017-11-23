package chameleon

import (
    "net/http"
)

func init() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/sms/incoming", receiveSMSHandler)
}
