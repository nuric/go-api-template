package handlers

import (
	"fmt"
	"net/http"

	"github.com/nuric/go-api-template/utils"
)

/* Key things to note:
- Request and response types are nearby the handler for easy debugging, you know
what is coming and going.
- We validate the request explicitly so it knows what is expected.
*/

type GreetingRequest struct {
	Name string `json:"name"`
}

type GreetingResponse struct {
	Greeting string `json:"greeting"`
}

func (r GreetingRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

func GreetingHandler(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValid[GreetingRequest](r)
	/* You'll find this error checking pattern repeats in all handlers. It's not
	 * that much if it really annoys you then you can do it inside the
	 * DecodeValid function. Either way is fine. I prefer explicit error checking
	 * in case we can do something else. */
	if err != nil {
		utils.Encode(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	// Construct the response and encode as JSON.
	response := GreetingResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Name),
	}
	utils.Encode(w, http.StatusOK, response)
}
