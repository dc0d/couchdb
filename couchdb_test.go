package couchdb

import yaml "gopkg.in/yaml.v2"

func init() {
	data := `
couchdb:
  user: admin
  password: abc2016
  url: http://172.17.0.2:5984/
  dbname: dummy
`
	err := yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		panic(err)
	}

	client, err = NewAuthClient(conf.CouchDB.User, conf.CouchDB.Password, conf.CouchDB.URL)
	if err != nil {
		panic(err)
	}

	c, err = NewAuthClient(conf.CouchDB.User, conf.CouchDB.Password, conf.CouchDB.URL)
	if err != nil {
		panic(err)
	}
	db = c.Use("dummy")

	cView, err = NewAuthClient(conf.CouchDB.User, conf.CouchDB.Password, conf.CouchDB.URL)
	if err != nil {
		panic(err)
	}
	dbView = cView.Use("gotest")
}

var client *Client
var c *Client
var db Database

var cView *Client
var dbView Database

var conf struct {
	CouchDB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		URL      string `yaml:"url"`
		Name     string `yaml:"dbname"`
	}
}
