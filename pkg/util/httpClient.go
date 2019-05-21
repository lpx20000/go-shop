package util

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

func HttpGet(urls string, data map[string]string) (result []byte, err error) {
	params := url.Values{}
	requestUrl, err := url.Parse(urls)
	if err != nil {
		return
	}

	for index, value := range data {
		params.Set(index, value)
	}

	requestUrl.RawQuery = params.Encode()
	urlPath := requestUrl.String()

	resp, err := http.Get(urlPath)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)
	return
}

//func httpDo() {
//	client := &http.Client{}
//
//	req, err := http.NewRequest("POST", "http://www.01happy.com/demo/accept.php", strings.NewReader("name=cjb"))
//	if err != nil {
//		// handle error
//	}
//
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//	req.Header.Set("Cookie", "name=anny")
//
//	resp, err := client.Do(req)
//
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		// handle error
//	}
//
//	fmt.Println(string(body))
//}

//func httpPostForm() {
//	resp, err := http.PostForm("http://www.01happy.com/demo/accept.php",
//		url.Values{"key": {"Value"}, "id": {"123"}})
//
//	if err != nil {
//		// handle error
//	}
//
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		// handle error
//	}
//
//	fmt.Println(string(body))
//
//}

//func httpPost() {
//	resp, err := http.Post("http://www.01happy.com/demo/accept.php",
//		"application/x-www-form-urlencoded",
//		strings.NewReader("name=cjb"))
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		// handle error
//	}
//
//	fmt.Println(string(body))
//}
