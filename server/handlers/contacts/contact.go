package contacts

import (
	"log"
	"net/http"
	"strconv"

	"github.com/coffemanfp/chat/config"
	"github.com/coffemanfp/chat/database"
	sErrors "github.com/coffemanfp/chat/errors"
	"github.com/coffemanfp/chat/server/handlers"
)

type ContactHandler struct {
	config     config.ConfigInfo
	repository database.ContactRepository
	writer     handlers.ResponseWriter
	reader     handlers.RequestReader
}

func NewContactHandler(repo database.ContactRepository, r handlers.RequestReader, w handlers.ResponseWriter, conf config.ConfigInfo) (c ContactHandler) {
	return ContactHandler{
		reader:     r,
		writer:     w,
		repository: repo,
		config:     conf,
	}
}

func (c ContactHandler) GetByRange(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting contacts by range...")
	limitP, offsetP := r.URL.Query().Get("limit"), r.URL.Query().Get("offset")
	limit, err := strconv.Atoi(limitP)
	if err != nil {
		c.handleError(w, err)
		return
	}
	offset, err := strconv.Atoi(offsetP)
	if err != nil {
		c.handleError(w, err)
		return
	}

	id := r.Context().Value("id").(int)

	contacts, err := c.repository.GetByRange(id, limit, offset)
	if err != nil {
		c.handleError(w, err)
		return
	}

	c.writer.JSON(w, http.StatusOK, contacts)
}

func (c ContactHandler) handleError(w http.ResponseWriter, err error) {
	hErr, ok := err.(sErrors.ClientError)
	if !ok {
		log.Println(err)
		c.writer.JSON(w, http.StatusInternalServerError, handlers.Hash{
			"message": sErrors.SERVER_ERROR_MESSAGE,
		})
		return
	}
	c.writer.JSON(w, hErr.HTTPCode(), handlers.Hash{
		"message": hErr.Error(),
	})
}
