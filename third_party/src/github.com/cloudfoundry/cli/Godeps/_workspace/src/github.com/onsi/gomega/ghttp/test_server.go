/*
Package ghttp supports testing HTTP clients by providing a test server (simply a thin wrapper around httptest's server) that supports
registering multiple handlers.  Incoming requests are not routed between the different handlers
- rather it is merely the order of the handlers that matters.  The first request is handled by the first
registered handler, the second request by the second handler, etc.

The intent here is to have each handler *verify* that the incoming request is valid.  To accomplish, ghttp
also provides a collection of bite-size handlers that each perform one aspect of request verification.  These can
be composed together and registered with a ghttp server.  The result is an expressive language for describing
the requests generated by the client under test.

Here's a simple example, note that the server handler is only defined in one BeforeEach and then modified, as required, by the nested BeforeEaches.
A more comprehensive example is available at https://onsi.github.io/gomega/#_testing_http_clients

	var _ = Describe("A Sprockets Client", func() {
		var server *ghttp.Server
		var client *SprocketClient
		BeforeEach(func() {
			server = ghttp.NewServer()
			client = NewSprocketClient(server.URL(), "skywalker", "tk427")
		})

		AfterEach(func() {
			server.Close()
		})

		Describe("fetching sprockets", func() {
			var statusCode int
			var sprockets []Sprocket
			BeforeEach(func() {
				statusCode = http.StatusOK
				sprockets = []Sprocket{}
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/sprockets"),
					ghttp.VerifyBasicAuth("skywalker", "tk427"),
					ghttp.RespondWithJSONEncodedPtr(&statusCode, &sprockets),
				))
			})

			Context("when requesting all sprockets", func() {
				Context("when the response is succesful", func() {
					BeforeEach(func() {
						sprockets = []Sprocket{
							NewSprocket("Alfalfa"),
							NewSprocket("Banana"),
						}
					})

					It("should return the returned sprockets", func() {
						Ω(client.Sprockets()).Should(Equal(sprockets))
					})
				})

				Context("when the response is missing", func() {
					BeforeEach(func() {
						statusCode = http.StatusNotFound
					})

					It("should return an empty list of sprockets", func() {
						Ω(client.Sprockets()).Should(BeEmpty())
					})
				})

				Context("when the response fails to authenticate", func() {
					BeforeEach(func() {
						statusCode = http.StatusUnauthorized
					})

					It("should return an AuthenticationError error", func() {
						sprockets, err := client.Sprockets()
						Ω(sprockets).Should(BeEmpty())
						Ω(err).Should(MatchError(AuthenticationError))
					})
				})

				Context("when the response is a server failure", func() {
					BeforeEach(func() {
						statusCode = http.StatusInternalServerError
					})

					It("should return an InternalError error", func() {
						sprockets, err := client.Sprockets()
						Ω(sprockets).Should(BeEmpty())
						Ω(err).Should(MatchError(InternalError))
					})
				})
			})

			Context("when requesting some sprockets", func() {
				BeforeEach(func() {
					sprockets = []Sprocket{
						NewSprocket("Alfalfa"),
						NewSprocket("Banana"),
					}

					server.WrapHandler(0, ghttp.VerifyRequest("GET", "/sprockets", "filter=FOOD"))
				})

				It("should make the request with a filter", func() {
					Ω(client.Sprockets("food")).Should(Equal(sprockets))
				})
			})
		})
	})
*/
package ghttp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	. "github.com/onsi/gomega"
)

func new() *Server {
	return &Server{
		AllowUnhandledRequests:     false,
		UnhandledRequestStatusCode: http.StatusInternalServerError,
		writeLock:                  &sync.Mutex{},
	}
}

// NewServer returns a new `*ghttp.Server` that wraps an `httptest` server.  The server is started automatically.
func NewServer() *Server {
	s := new()
	s.HTTPTestServer = httptest.NewServer(s)
	return s
}

// NewTLSServer returns a new `*ghttp.Server` that wraps an `httptest` TLS server.  The server is started automatically.
func NewTLSServer() *Server {
	s := new()
	s.HTTPTestServer = httptest.NewTLSServer(s)
	return s
}

type Server struct {
	//The underlying httptest server
	HTTPTestServer *httptest.Server

	//Defaults to false.  If set to true, the Server will allow more requests than there are registered handlers.
	AllowUnhandledRequests bool

	//The status code returned when receiving an unhandled request.
	//Defaults to http.StatusInternalServerError.
	//Only applies if AllowUnhandledRequests is true
	UnhandledRequestStatusCode int

	receivedRequests []*http.Request
	requestHandlers  []http.HandlerFunc

	writeLock *sync.Mutex
	calls     int
}

//URL() returns a url that will hit the server
func (s *Server) URL() string {
	return s.HTTPTestServer.URL
}

//Close() should be called at the end of each test.  It spins down and cleans up the test server.
func (s *Server) Close() {
	server := s.HTTPTestServer
	s.HTTPTestServer = nil
	server.Close()
}

//ServeHTTP() makes Server an http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()
	defer func() {
		recover()
	}()

	if s.calls < len(s.requestHandlers) {
		s.requestHandlers[s.calls](w, req)
	} else {
		if s.AllowUnhandledRequests {
			ioutil.ReadAll(req.Body)
			req.Body.Close()
			w.WriteHeader(s.UnhandledRequestStatusCode)
		} else {
			Ω(req).Should(BeNil(), "Received Unhandled Request")
		}
	}
	s.receivedRequests = append(s.receivedRequests, req)
	s.calls++
}

//ReceivedRequests is an array containing all requests received by the server (both handled and unhandled requests)
func (s *Server) ReceivedRequests() []*http.Request {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	return s.receivedRequests
}

//AppendHandlers will appends http.HandlerFuncs to the server's list of registered handlers.  The first incoming request is handled by the first handler, the second by the second, etc...
func (s *Server) AppendHandlers(handlers ...http.HandlerFunc) {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	s.requestHandlers = append(s.requestHandlers, handlers...)
}

//SetHandler overrides the registered handler at the passed in index with the passed in handler
//This is useful, for example, when a server has been set up in a shared context, but must be tweaked
//for a particular test.
func (s *Server) SetHandler(index int, handler http.HandlerFunc) {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	s.requestHandlers[index] = handler
}

//GetHandler returns the handler registered at the passed in index.
func (s *Server) GetHandler(index int) http.HandlerFunc {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	return s.requestHandlers[index]
}

//WrapHandler combines the passed in handler with the handler registered at the passed in index.
//This is useful, for example, when a server has been set up in a shared context but must be tweaked
//for a particular test.
//
//If the currently registered handler is A, and the new passed in handler is B then
//WrapHandler will generate a new handler that first calls A, then calls B, and assign it to index
func (s *Server) WrapHandler(index int, handler http.HandlerFunc) {
	existingHandler := s.GetHandler(index)
	s.SetHandler(index, CombineHandlers(existingHandler, handler))
}
