package module

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/AvicennaJr/Nuru/object"
)

var NetFunctions = map[string]object.ModuleFunction{}

func init() {
	NetFunctions["peruzi"] = getRequest
	NetFunctions["tuma"] = postRequest
}

func getRequest(args []object.Object) object.Object {

	if len(args) > 3 {
		return &object.Error{Message: "Hatuhitaji hoja zaidi ya 3."}
	}

	if args[0].Type() != object.STRING_OBJ {
		return &object.Error{Message: "Link iwe ndani ya \"\". Mfano: \"https://google.com\""}
	}

	url := args[0].Inspect()

	req, err := http.NewRequest("GET", url, nil)

	// var responseBody *bytes.Buffer

	if len(args) == 2 {

		switch v := args[1].(type) {

		case *object.Byte:

			// responseBody = bytes.NewBuffer(v.Value)

		case *object.Dict:
			for _, val := range v.Pairs {
				req.Header.Set(val.Key.Inspect(), val.Value.Inspect())
			}
			// input := args[0].Inspect()

			// jsonBody, err := json.Marshal(input)

			// if err != nil {
			// 	return &object.Error{Message: "Huku format query yako vizuri."}
			// }

			// responseBody = bytes.NewBuffer(jsonBody)

		default:
			return &object.Error{Message: "Data unayoruhusiwa kutuma ni Kamusi (Dict) au Bytes."}
		}

	}

	if err != nil {
		return &object.Error{Message: "Tumeshindwa kutuma request."}
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return &object.Error{Message: "Tumeshindwa kutuma request."}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &object.Error{Message: "Tumeshindwa kusoma majibu."}
	}

	return &object.String{Value: string(body)}
}

func postRequest(args []object.Object) object.Object {

	if len(args) != 2 {
		return &object.Error{Message: "Tunahitaji hoja mbili."}
	}
	var url string

	switch link := args[0].(type) {
	case *object.String:
		url = link.Value
	default:
		return &object.Error{Message: "Hii sio link sahihi. Link iwe ndani ya \"\". Mfano: \"https://google.com\""}
	}

	var responseBody *bytes.Buffer

	switch v := args[1].(type) {

	case *object.Byte:

		responseBody = bytes.NewBuffer(v.Value)

	case *object.Dict:
		input := args[0].Inspect()

		jsonBody, err := json.Marshal(input)

		if err != nil {
			return &object.Error{Message: "Huku format query yako vizuri."}
		}

		responseBody = bytes.NewBuffer(jsonBody)

	default:
		return &object.Error{Message: "Data unayoruhusiwa kutuma ni Kamusi (Dict) au Bytes."}
	}

	resp, err := http.Post(url, "application/json", responseBody)

	if err != nil {
		return &object.Error{Message: "Tumeshindwa kupost data. Huenda huna bando?"}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &object.Error{Message: "Tumeshindwa kusoma majibu yaliyo rudishwa."}
	}

	return &object.String{Value: string(body)}
}
