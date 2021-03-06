package router

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/go-errors/errors"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	"github.com/octavore/naga/service"
	"github.com/octavore/nagax/logger"
	"github.com/octavore/nagax/router"

	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/server/router/middleware/gzip"
	logger2 "github.com/ketchuphq/ketchup/server/router/middleware/logger"
)

var ErrNotFound = fmt.Errorf("not found")

type Module struct {
	*router.Module
	Logger *logger.Module
}

func (m *Module) Init(c *service.Config) {
	c.Setup = func() error {
		m.Module.Middleware.Set(
			logger2.New(m.Logger.Infof),
			gzip.Default,
		)
		return nil
	}
}

func (m *Module) Subrouter(path string) *httprouter.Router {
	r := httprouter.New()
	m.Handle(path, r)
	return r
}

func (m *Module) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	m.Middleware.ServeHTTP(rw, req)
}

func Proto(rw http.ResponseWriter, pb proto.Message) error {
	rw.Header().Set("Content-Type", "application/json")
	err := JSON.Marshal(rw, pb)
	if err != nil {
		return fmt.Errorf("router: %v", err)
	}
	return nil
}

func (m *Module) InternalError(rw http.ResponseWriter, err error) {
	if err == ErrNotFound {
		NotFound(rw)
		return
	}

	switch e := err.(type) {
	case *errors.Error:
		s := e.StackFrames()[0]
		f := path.Base(s.File)
		m.Logger.Errorf("[%s/%s:%d] %v", s.Package, f, s.LineNumber, e.Error())
	default:
		m.Logger.Error(e)
	}
	rw.WriteHeader(http.StatusInternalServerError)
	err = Proto(rw, &api.Error{
		Code:   api.ErrorCode_INTERNAL_SERVER_ERROR.Enum(),
		Detail: proto.String("Internal server error."),
	})
	if err != nil {
		log.Println("router: ", err)
	}
}

func NotFound(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNotFound)
	err := Proto(rw, &api.Error{
		Code:   api.ErrorCode_NOT_FOUND.Enum(),
		Detail: proto.String("Not found."),
	})
	if err != nil {
		log.Println("router: ", err)
	}
}

var JSON = &jsonpb.Marshaler{
	EnumsAsInts: false,
	Indent:      "  ",
}
