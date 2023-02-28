package streamer

import (
	"encoding/json"
	"fmt"
)

func (s *streamer) handlerWebSocket(data receiveData) error {
	var req payload
	err := json.Unmarshal(data.payload, &req)
	if err != nil {
		return err
	}

	fmt.Printf("payload: %v\n", data.payload)
	fmt.Printf("method: %s\n", req.Method)
	fmt.Printf("args: %v\n", req.Args)

	switch req.Method {
	case "message":
		fmt.Printf("message received\n")
		sendMessage, ok := req.Args["message"].(string)
		if !ok {
			return fmt.Errorf("invalid type for message\n")
		}
		s.send([]byte(sendMessage), func(_ *client) bool { return true })
	case "join":
		fmt.Printf("join received\n")
		processJoin(s, data.id, req.Args)
	case "input":
		fmt.Printf("input received\n")
		processInput(s, data.id, req.Args)
	default:
		fmt.Printf("invalid method")
	}

	return nil
}
