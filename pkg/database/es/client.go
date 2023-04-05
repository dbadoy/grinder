package es

import (
	"fmt"

	"github.com/dbadoy/grinder/pkg/database"
	"github.com/elastic/go-elasticsearch"
)

var _ = database.Database(&Client{})

type Client struct {
	conn *elasticsearch.Client
}

func New(urls []string) (*Client, error) {
	fmt.Println(urls)
	c, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: urls,
	})
	return &Client{c}, err
}

func (c *Client) HealthCheck() error {
	_, err := c.conn.Info()
	return err
}

func (c *Client) Insert(key []byte, data database.Data) error {
	panic("need impl")
}

func (c *Client) Put(key []byte, data database.Data) error {
	panic("need impl")
}

func (c *Client) Exist(index string, key []byte) (bool, error) {
	panic("need impl")
}

func (c *Client) Delete(key []byte) error {
	panic("need impl")
}
