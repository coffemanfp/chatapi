package auth

import (
	"fmt"
	"net/http"

	"github.com/coffemanfp/chat/server/handlers"
	"github.com/coffemanfp/chat/users"
)

const systemHandlerName handlerName = "system"

type systemUserReader struct {
	reader handlers.RequestReader
	writer handlers.ResponseWriter
}

func (s systemUserReader) read(w http.ResponseWriter, r *http.Request) (user users.User, err error) {
	err = s.reader.JSON(r, &user)
	if err != nil {
		fmt.Printf("failed for %s\n", err)
		s.writer.JSON(w, http.StatusBadRequest, handlers.Hash{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
	}
	return
}
