//go:build integration
// +build integration

package test

func StressTestSimpleWebGamesService(t *testing.T) {
	postBody, _ := json.Marshal(BodyCreate{size: 3})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:8080/create", "application/json", responseBody)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb, resp.Header.Get("Set-Cookie"))
}
