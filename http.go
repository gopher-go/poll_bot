package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/andrewkav/viber"
)

func serve(v *viber.Viber, ud userDAO) error {
	s, err := newStorage(ud)
	if err != nil {
		return err
	}

	p := generateOurPoll()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleMain(p, v, s, w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		return err
	}
	return nil
}

func isJSON(s []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil
}

func handleMain(p poll, v *viber.Viber, s *storage, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	if !isJSON(bodyBytes) {
		http.Error(w, "Not json response", http.StatusBadRequest)
		return
	}

	c, err := parseCallback(bodyBytes)
	if err != nil {
		log.Printf("Error reading callback: %v for input %v", err, string(bodyBytes))
		http.Error(w, "can't parse body", http.StatusBadRequest)
		return
	}

	// we need it for subscribe
	if c.Event == "webhook" {
		return
	}

	reply, err := generateReplyFor(p, s, c)
	if err != nil {
		log.Printf("Error generating reply: %v for input %v", err, string(bodyBytes))
		http.Error(w, "can't reply", http.StatusBadRequest)
		return
	}
	if reply != nil {
		message := v.NewTextMessage(reply.text)
		if len(reply.options) > 0 {
			message.SetKeyboard(keyboardFromOptions(v, reply.options))
		}
		_, err = v.SendMessage(c.User.ID, message)
		if err != nil {
			log.Printf("Error sending message %v to user id %s", err, c.User.ID)
			http.Error(w, "can't reply", http.StatusBadRequest)
			return
		}
	}
}

func calculateColsAndRows(optCount int) (cols int, rows int) {
	if optCount <= 2 {
		cols = 3
		rows = 2

		return
	}

	cols = 3
	rows = 1

	return
}

func keyboardFromOptions(v *viber.Viber, options []string) *viber.Keyboard {
	ret := v.NewKeyboard("#FFFFFF", true)
	for _, opt := range options {
		cornerRadius := 2
		// columns and rows to occupy by a single button
		cols, rows := calculateColsAndRows(len(options))
		b := &viber.Button{
			Columns:    cols,
			Rows:       rows,
			ActionType: viber.Reply,
			ActionBody: opt,
			Image:      "",
			Text:       fmt.Sprintf(`<font color="#FFFFFF">%s</font>`, opt),
			TextSize:   viber.Medium,
			Frame: &viber.ButtonFrame{
				CornerRadius: &cornerRadius,
			},
			TextVAlign: "",
			TextHAlign: "",
			BgColor:    "#9482F8",
		}
		ret.AddButton(b)
	}
	return ret
}
