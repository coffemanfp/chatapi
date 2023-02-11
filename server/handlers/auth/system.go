package auth

import (
	"fmt"
	"net/http"

	"github.com/coffemanfp/chat/account"
	"github.com/coffemanfp/chat/server/handlers"
)

const systemHandlerName handlerName = "system"

type systemAccountReader struct {
	reader handlers.RequestReader
	writer handlers.ResponseWriter
}

func (s systemAccountReader) read(w http.ResponseWriter, r *http.Request) (account account.Account, err error) {
	err = s.reader.JSON(r, &account)
	if err != nil {
		fmt.Printf("failed for %s\n", err)
		s.writer.JSON(w, http.StatusBadRequest, handlers.Hash{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
	}
	return
}
