package psql

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
)

const (
	foreign_key_violation = "23505"
)

type pqErrHandler struct {
	pqErr *pq.Error
}

func (p pqErrHandler) asAlreadyExists() (match bool, err error) {
	match = p.pqErr.Code == foreign_key_violation
	if match {
		err = fmt.Errorf("already exists %s", getFieldFromDetail(p.pqErr))
	}

	return
}

func getFieldFromDetail(pqErr *pq.Error) string {
	return pqErr.Detail[strings.Index(pqErr.Detail, "(")+1 : strings.Index(pqErr.Detail, ")")]
}

func newPQError(pqErr error) pqErrHandler {
	return pqErrHandler{
		pqErr: pqErr.(*pq.Error),
	}
}
