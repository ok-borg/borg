package conf

type Mysql struct {
	Addr string `json:"addr"`
	Ids  string `json:"ids"`
}

type Github struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type Conf struct {
	EsAddr    string `json:"esaddr"`
	Github    Github `json:"github"`
	Sitemap   string `json:"sitemap"`
	Analytics string `json:"analytics"`
	Mysql     Mysql  `json:"mysql"`
}
