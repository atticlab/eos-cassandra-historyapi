package routes

import (
	"encoding/json"
	"eos-cassandra-historyapi/storage"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)


func (r *Router) handleGetControlledAccounts() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		bytes, err := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		if err != nil {
			writeErrorResponse(writer, http.StatusInternalServerError, "Internal service error")
			log.Println("Failed to read request body. Error: " + err.Error())
			return
		}

		var args storage.GetControlledAccountsArgs
		if len(bytes) > 0 {
			if err = json.Unmarshal(bytes, &args); err != nil {
				if _, ok := err.(*json.SyntaxError); ok {
					writeErrorResponse(writer, http.StatusBadRequest, "Invalid json in request body")
					log.Println("Invalid request body. Error: " + err.Error())
					return
				}
			}
		}

		response, errorResult := r.historyStorage.GetControlledAccounts(args)
		if errorResult != nil {
			writeErrorResponse(writer, errorResult.Code, errorResult.Error())
			log.Println("Got error from IHistoryStorage.GetControlledAccounts(). Error: " + errorResult.Error())
			return
		}
		b, err := json.Marshal(response)
		if err != nil {
			writeErrorResponse(writer, http.StatusInternalServerError, "Internal service error")
			log.Println("Failed to marshal response. Error: " + err.Error())
			return
		}
		fmt.Fprintf(writer, string(b))
	}
}