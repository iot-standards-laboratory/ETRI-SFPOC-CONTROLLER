package consulapi

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/hashicorp/consul/api"
)

const maxConnectAttempts = 100
const acl_token = "62KY7J2YIOycoeuzgOOEsWwg7ZuE7J+dIOOEuOOFj+uKlOuCoAoK7Jyg7KaI66as7ZWYIOy5tOugjCBI7Lu1IOuvvOqwkOuwlOuUlCDshLHqsJDqsJzrsJwKCjMwME1BQU4tNzgwIOyXkHJvIOyduO2UjOujqOyWuOyE"

var client *api.Client = nil

func Connect(consulAgent string) error {
	consulConfig := api.DefaultConfig()
	consulAgentUrl, err := url.Parse(consulAgent)

	if err != nil {
		glog.Infof("Error parsing consul url")
		return err
	}

	if consulAgentUrl.Host != "" {
		consulConfig.Address = consulAgentUrl.Host
	}

	if consulAgentUrl.Scheme != "" {
		consulConfig.Scheme = consulAgentUrl.Scheme
	}
	consulConfig.Token = acl_token

	client, err = api.NewClient(consulConfig)
	if err != nil {
		glog.Info("Error creating consul client")
		return err
	}

	for attempt := 1; attempt < maxConnectAttempts; attempt++ {
		if _, err := client.Agent().Self(); err == nil {
			break
		}

		glog.Infof("[Attempt: %d] Attempting access to consul after 5 seconds sleep", attempt)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to consul agent: %v, error: %v", consulAgent, err)
	}

	glog.Infof("Consul agent found: %v", consulAgent)

	return nil
}

func GetStates() (map[string]bool, error) {
	states := make(map[string]bool)

	catalogs, meta, err := client.Catalog().Services(&api.QueryOptions{})
	if err != nil {
		return nil, err
	}
	glog.Info(meta)

	for k, _ := range catalogs {
		status, _, _ := client.Health().Checks(k, nil)
		states[k] = strings.Compare(status.AggregatedStatus(), "passing") == 0
	}

	return states, nil
}

func GetStatus(name string) (string, error) {
	s, _, err := client.Agent().AgentHealthServiceByID(name)
	return s, err
}

func Monitor(h func(string), ctx context.Context) error {
	stopCh := make(chan struct{})
	notificationCh, err := client.Agent().Monitor("", stopCh, &api.QueryOptions{})
	if err != nil {
		return err
	}

	for {
		select {
		case param := <-notificationCh:
			if h != nil {
				h(param)
			}
		case <-ctx.Done():
			stopCh <- struct{}{}
			return nil
		}
	}
}

func registerEntity(name, endpoint string) (err error) {
	agent := client.Agent()
	uri, err := url.Parse(endpoint)
	if err != nil {
		return
	}
	host, port, err := net.SplitHostPort(uri.Host)
	if err != nil {
		return
	}
	port_i, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		return
	}

	// server check the health
	// service := &consulapi.AgentServiceRegistration{
	// 	Name:    name,
	// 	Port:    int(port_i),
	// 	Address: host,
	// 	Check: &consulapi.AgentServiceCheck{
	// 		// Args:     []string{"curl", "localhost:8080"},
	// 		HTTP:     endpoint,
	// 		Timeout:  "10s",
	// 		Interval: "10s",
	// 	},
	// }

	// client report the health
	entity := &api.AgentServiceRegistration{
		Name:    name,
		Port:    int(port_i),
		Address: host,
		Check: &api.AgentServiceCheck{
			TTL: ttl.String(),
		},
	}

	if err = agent.ServiceRegister(entity); err != nil {
		return
	}

	// client.Agent().CheckRegister(&api.AgentCheckRegistration{
	// 	ServiceID: name,
	// })

	// s, _, _ := client.Agent().AgentHealthServiceByID(name)
	// log.Fatalln(s)

	return nil
}
