package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mileusna/viber"
)

func serve(v *viber.Viber) error {
	s, err := newStorage()
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
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return err
	}
	return nil
}

func isJSON(s []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil
}

func handleMain(p poll, v *viber.Viber, s *Storage, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	if !isJSON(bytes) {
		http.Error(w, "Not json response", http.StatusBadRequest)
		return
	}

	c, err := parseCallback(bytes)
	if err != nil {
		log.Printf("Error reading callback: %v for input %v", err, string(bytes))
		http.Error(w, "can't parse body", http.StatusBadRequest)
		return
	}

	reply, err := generateReplyFor(p, s, c)
	if err != nil {
		log.Printf("Error generating reply: %v for input %v", err, string(bytes))
		http.Error(w, "can't reply", http.StatusBadRequest)
		return
	}
	if reply != nil {
		log.Printf("Got reply %v", *reply)
		message := v.NewTextMessage(reply.text)
		if len(reply.options) > 0 {
			message.SetKeyboard(keyboardFromOptions(v, reply.options))
		}
		_, err = v.SendMessage(c.User.Id, message)
		if err != nil {
			log.Printf("Error sending message %v", err)
			http.Error(w, "can't reply", http.StatusBadRequest)
			return
		}
	}
}

func keyboardFromOptions(v *viber.Viber, options []string) *viber.Keyboard {
	ret := v.NewKeyboard("#FFFFFF", true)
	colSize := len(options)
	for _, opt := range options {
		b := v.NewTextButton(colSize, 1, viber.Reply, opt, opt)
		ret.AddButton(b)
	}
	return ret
}
