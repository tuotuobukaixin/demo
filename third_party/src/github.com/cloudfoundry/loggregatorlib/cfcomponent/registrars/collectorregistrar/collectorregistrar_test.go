package collectorregistrar_test

import (
	"github.com/cloudfoundry/loggregatorlib/cfcomponent/registrars/collectorregistrar"

	"errors"
	"sync/atomic"
	"time"

	"github.com/apcera/nats"
	"github.com/cloudfoundry/gosteno"
	"github.com/cloudfoundry/loggregatorlib/cfcomponent"
	"github.com/cloudfoundry/loggregatorlib/loggertesthelper"
	"github.com/cloudfoundry/yagnats"
	"github.com/cloudfoundry/yagnats/fakeyagnats"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Collectorregistrar", func() {
	Describe("Run", func() {
		var (
			fakeClient                  *fakeClient
			component                   cfcomponent.Component
			registrar                   *collectorregistrar.CollectorRegistrar
			doneChan                    chan struct{}
			errorProvider               func() error
			fakeClientProviderCallCount int32
		)

		BeforeEach(func() {
			fakeClient = newFakeClient()
			component, _ = cfcomponent.NewComponent(loggertesthelper.Logger(), "compType", 3, nil, 9999, []string{"username", "password"}, nil)
			component.UUID = "OurUUID"
			errorProvider = func() error {
				return nil
			}
			fakeClientProviderCallCount = 0
			fakeClientProvider := func(*gosteno.Logger, *cfcomponent.Config) (yagnats.NATSConn, error) {
				atomic.AddInt32(&fakeClientProviderCallCount, 1)
				return fakeClient, errorProvider()
			}
			registrar = collectorregistrar.NewCollectorRegistrar(fakeClientProvider, component, 10*time.Millisecond, nil)
			doneChan = make(chan struct{})

			go func() {
				defer close(doneChan)
				registrar.Run()
			}()
		})

		AfterEach(func() {
			registrar.Stop()
			Eventually(doneChan).Should(BeClosed())
		})

		Context("with no errors", func() {
			It("periodically announces itself via NATS", func() {
				var messages []*nats.Msg
				Eventually(func() int {
					messages = fakeClient.PublishedMessages(collectorregistrar.AnnounceComponentMessageSubject)
					return len(messages)
				}).Should(BeNumerically(">", 1))

				for _, message := range messages {
					Expect(message.Data).To(MatchRegexp(`^\{"type":"compType","index":3,"host":"[^:]*:9999","uuid":"3-OurUUID","credentials":\["username","password"\]\}$`))
				}
			})

			It("reuses the client connection", func() {
				Eventually(func() int32 { return atomic.LoadInt32(&fakeClientProviderCallCount) }).Should(BeEquivalentTo(1))
				Consistently(func() int32 { return atomic.LoadInt32(&fakeClientProviderCallCount) }).Should(BeEquivalentTo(1))
			})
		})

		Context("with errors", func() {
			Context("from the client provider", func() {
				BeforeEach(func() {
					fakeError := errors.New("fake error")
					errorProvider = func() error {
						returnedError := fakeError
						fakeError = nil
						return returnedError
					}
				})

				It("recovers when the client provider recovers", func() {
					Eventually(fakeClient.PublishedMessageCount).Should(BeNumerically(">", 0))
				})

				It("disconnects the client", func() {
					Eventually(fakeClient.Closed).Should(BeTrue())
				})
			})
		})
	})
})

type fakeClient struct {
	*fakeyagnats.FakeNATSConn
	closed bool
}

func newFakeClient() *fakeClient {
	return &fakeClient{
		FakeNATSConn: fakeyagnats.Connect(),
	}
}

func (f *fakeClient) Close() {
	f.Lock()
	defer f.Unlock()
	f.closed = true
}

func (f *fakeClient) Closed() bool {
	f.Lock()
	defer f.Unlock()
	return f.closed
}
