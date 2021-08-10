package main

type config struct {
	HTTP struct {
		Addr string `default:":80"`
	}
	DB struct {
		Addr     string `required:"true"`
		User     string `required:"true"`
		Password string `required:"true"`
		Name     string `required:"true"`
	}
}
