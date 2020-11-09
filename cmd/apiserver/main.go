package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/VIPowERuS/nsu_postman/internal/app/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {

	flag.Parse()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	s := apiserver.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	/*
		h := sha1.New()
		//h.Write([]byte("test"))
		//h.Reset()
		h.Write([]byte("qweasd22"))
		bs := h.Sum([]byte{})
		fmt.Printf("%x\n", string(bs))
		f6bb6d326a3826e18df674d05e6fa2bdd8518284
	*/
}
