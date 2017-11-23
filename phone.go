package chameleon

import (
    "html/template"
    "fmt"
    "net/http"
    "os"
    "log"
    "regexp"
    "strconv"

    "golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"

    "bitbucket.org/ckvist/twilio/twiml"
    "bitbucket.org/ckvist/twilio/twirest"
)

var (
	twilioClient = twirest.NewClient(mustGetenv("TWILIO_ACCOUNT_SID"),
		                             mustGetenv("TWILIO_AUTH_TOKEN"))
	twilioNumber = mustGetenv("TWILIO_NUMBER")

    // This is the regular expressions to match incoming messages to
	locationMessage = regexp.MustCompile("([0-9]{3})[\\s,]*([0-9]{3})")
	checkpointMessage = regexp.MustCompile("\\b([a-zA-Z]{6})\\b")
	bonusMessage = regexp.MustCompile("\\b([0-9]{2})\\b")
)

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

func handleLocationMessage(c context.Context, team *datastore.Key, northing string, easting string) string {
    // These are verified with regular expressions, so won't fail
    north, _ = strconv.Atoi(northing)
    east, _ = strconv.Atoi(easting)

    // Location must be between 720 810 and 780 730
    swapped = false
    if east > 780 && east < 720 {
        // They might be the wrong way round
        swapped = true
        north, east = east, north
    }

    // Now test the location
    if east < 720 || east > 780 || north < 730 || north > 810 {
        response = "You are not in the play area - move between 720 810 and 780 730"
        // Swap back to log since they aren't in the play area
        if swapped {
            north, east = east, north
        }
    } else {
        response = message.respond("Thank you, your location has been recorded.")
    }
    
    location := Location {
        northing: north
        easting: east,
        time: time.Now()
    }
    key := datastore.NewIncompleteKey(c, "Location", team)
    _, err = datastore.Put(c, key, &location);
    if err != nil {
        response = "Sorry, there was an error with the server, please try again"
    }
    return response
}

func handleCheckpointMessage(c context.Context, team *datastore.Key, checkpoint string) string {
    return "checkpoint"
}

func handleBonusMessage(c context.Context, team *datastore.Key, bonus string) string {
    return "bonus"
}

func receiveSMSHandler(w http.ResponseWriter, r *http.Request) {
    ctx := appengine.NewContext(r)

	sender := r.FormValue("From")
	body := r.FormValue("Body")
	var response string
	
	team := mobileToTeam(ctx, sender)
	if team == nil {
    	response = "Your number is not registered with a team. If this is incorrect, please call 0121 280 0589."
    }

	if response == "" {
        location := locationMessage.FindStringSubmatch(body)
        if location != nil {
            response = handleLocationMessage(ctx, team, location[1], location[2])
        }
    }

    if response == "" {
        checkpoint := checkpointMessage.FindStringSubmatch(body)
        if checkpoint != nil {
            response = handleCheckpointMessage(ctx, team, checkpoint[1])
        }
    }

    if response == "" {
        bonus := bonusMessage.FindStringSubmatch(body)
        if bonus != nil {
            response = handleBonusMessage(ctx, team, bonus[1])
        }
    }

	if response == "" {
	    response = "THIS IS AN AUTOMATED SYSTEM. Please message only your six figure grid reference or six letter waypoint code. For help please call 0121 280 0589."
	}

	resp := twiml.NewResponse()
	resp.Action(twiml.Message{
		Body: response,
		From: twilioNumber,
		To:   sender,
	})
	resp.Send(w)
}

var guestbookTemplate = template.Must(template.New("book").Parse(`
<html>
  <head>
    <title>Go Guestbook</title>
  </head>
  <body>
    {{range .}}
      {{with .Author}}
        <p><b>{{.}}</b> wrote:</p>
      {{else}}
        <p>An anonymous person wrote:</p>
      {{end}}
      <pre>{{.Content}}</pre>
    {{end}}
    <form action="/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
  </body>
</html>
`))

func guestbookKey(c context.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}

func sendSMSHandler(w http.ResponseWriter, r *http.Request) {
        to := r.FormValue("to")
        if to == "" {
                http.Error(w, "Missing 'to' parameter.", http.StatusBadRequest)
                return
        }

        msg := twirest.SendMessage{
                Text: "Hello from App Engine!",
                From: twilioNumber,
                To:   to,
        }

        resp, err := twilioClient.Request(msg)
        if err != nil {
                http.Error(w, fmt.Sprintf("Could not send SMS: %v", err), 500)
                return
        }

        fmt.Fprintf(w, "SMS sent successfully. Response:\n%#v", resp.Message)
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Hello, %v!", u)
}

