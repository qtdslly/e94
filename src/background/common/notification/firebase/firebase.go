package notification

import (
	"bytes"
	"common/logger"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	firebaseCMUrl = "https://fcm.googleapis.com/fcm/send"
	firebaseKey   = "AAAAJQIih4M:APA91bGlmiK-eJw4lc_uozSvdbvlN4ZoFpAHRpiIohQamLw-jtS8O-pr3gF4TM5wCgslEo1yvc7SvymVUlzEHcAmhHi7CrpSlltjG1FmgjOwX6rFR7hvytG-c_yPoA2wfhG3KjvRq_sMTmxZVN6lEgRApxuoMNtCHw"
)

/*
	{
	    "to": "/topics/all",
	    "notification": {
		"body": "topics/all test"
	    },
	    "priority": 10
	}
*/
type pushRequest struct {
	To           string `json:"to"`
	Notification struct {
		Body string `json:"body"`
	} `json:"notification"`
	Priority uint32 `json:"priority"`
}

func Push(body string) error {
	request := &pushRequest{}
	request.To = "/topics/all"
	request.Notification.Body = body
	request.Priority = 10

	b, err := json.Marshal(&request)
	if err != nil {
		logger.Error(err)
		return err
	}

	payload := bytes.NewReader(b)
	req, err := http.NewRequest("POST", firebaseCMUrl, payload)
	if err != nil {
		logger.Error(err)
		return err
	}

	req.Header.Add("authorization", "key="+firebaseKey)
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
		return err
	}

	defer res.Body.Close()
	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Debug("fcm response: ", string(resp))
	return nil
}
