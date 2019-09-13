package routes

import (
	"encoding/json"
	"eos-cassandra-historyapi/storage"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const MaxTxHashesPerRequest = 500

func (r *Router) handleGetTransactions() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		bytes, err := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		if err != nil {
			writeErrorResponse(writer, http.StatusInternalServerError, "Internal service error")
			log.Println("Failed to read request body. Error: " + err.Error())
			return
		}

		var args storage.GetTransactionsArgs
		if len(bytes) > 0 {
			if err = json.Unmarshal(bytes, &args); err != nil {
				if _, ok := err.(*json.SyntaxError); ok {
					writeErrorResponse(writer, http.StatusBadRequest, "Invalid json in request body")
					log.Println("Invalid request body. Error: " + err.Error())
					return
				}
			}
		}

		if len(args.IDs) > MaxTxHashesPerRequest {
			writeErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("max limit of txs per request: %d", MaxTxHashesPerRequest))
			return
		}

		response, errorResult := r.historyStorage.GetTransactions(args)
		if errorResult != nil {
			writeErrorResponse(writer, errorResult.Code, errorResult.Error())
			log.Println("Got error from IHistoryStorage.GetTransaction(). Error: " + errorResult.Error())
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