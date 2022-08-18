package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var client = ZtClient{
	Client:  &http.Client{Timeout: config.Zerotier.Timeout},
	baseURL: config.Zerotier.URL}

type ZtClient struct {
	Client  *http.Client
	baseURL string
}

func (c *ZtClient) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	full_url := strings.TrimSuffix(config.Zerotier.URL, "/")
	if strings.HasPrefix(path, "/") {
		full_url += path
	} else {
		full_url += "/" + path
	}

	u, err := url.Parse(full_url)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}

	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Add("Authorization", "token "+config.Zerotier.Token)

	return req, nil
}

func (c *ZtClient) Request(method, path string, body interface{}) (*http.Response, error) {
	req, err := c.NewRequest(method, path, body)
	if err != nil {
		return nil, nil
	}

	return c.Client.Do(req)
}

func parse_response(resp *http.Response) []byte {
	if resp.StatusCode != 200 {
		fmt.Println("status:", resp.Status)
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return data
}

func zt_req(method, path string, obj, body interface{}) {

	resp, err := client.Request(method, path, body)
	if err != nil {
		fmt.Println("** error:", err)
		return
	}

	if method == http.MethodPost {
		defer resp.Body.Close()
	}

	data := parse_response(resp)

	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("** error:", err)
	}
}

func zt_get_uid() string {
	// get uid by create a tmp network
	body := make(map[string]interface{})
	obj := make(map[string]interface{})
	zt_req(http.MethodPost, "/network", &obj, &body)

	// once get resp, remove the network
	if nid, ok := obj["id"]; ok {
		if v, ok := nid.(string); ok {
			client.Request(http.MethodDelete, "/network/"+v, nil)
		}
	}

	// parse and return uid
	if uid, ok := obj["ownerId"]; ok {
		if v, ok := uid.(string); ok {
			return v
		}
	}

	return ""
}

func zt_user_record() {
	// get uid if not specified
	uid := config.Zerotier.UID
	if uid == "" {
		uid = zt_get_uid()
	}

	// get user record
	obj := make(map[string]interface{})
	zt_req(http.MethodGet, "/user/"+uid, &obj, nil)

	fmt.Println(dumps(obj, args.Format))
}

func zt_network_list() {
	obj := make([]ZtNetInfo, 0)
	zt_req(http.MethodGet, "/network", &obj, nil)

	display_networks(obj)
}

func zt_network_create() {
	if args.Noop {
		fmt.Printf(">> create args:\n%s\n", dumps(config.Zerotier.Net, args.Format))
		return
	}

	body := make(map[string]interface{})
	obj := make(map[string]interface{})
	zt_req(http.MethodPost, "/network", &obj, &body)

	if nid, ok := obj["id"]; ok {
		if value, ok := nid.(string); ok {
			obj := ZtNetInfo{}
			body := config.Zerotier.Net
			zt_req(http.MethodPost, "/network/"+value, &obj, &body)

			fmt.Println(dumps(obj, args.Format))
		}
	}
}

func zt_network_del(nid string) {
	resp, err := client.Request(http.MethodDelete, "/network/"+nid, nil)
	if err != nil {
		fmt.Println("** error:", err)
		return
	}

	fmt.Println("resp:", resp.Status)
}

func zt_network_member_list(nid string) {
	obj := make([]ZtNetMemberInfo, 0)
	zt_req(http.MethodGet, "/network/"+nid+"/member", &obj, nil)

	display_network_members(obj)
}

func zt_network_member_set(nid string, mid string) {
	if args.Noop {
		fmt.Printf(">> set args:\n%s\n", dumps(config.Zerotier.Netm, args.Format))
		return
	}

	path := "/network/" + nid + "/member/" + mid
	obj := ZtNetMemberInfo{}
	body := config.Zerotier.Netm
	zt_req(http.MethodPost, path, &obj, &body)

	fmt.Println(dumps(obj, args.Format))
}

func zt_network_member_del(nid string, mid string) {
	path := "/network/" + nid + "/member/" + mid
	resp, err := client.Request(http.MethodDelete, path, nil)
	if err != nil {
		fmt.Println("** error:", err)
		return
	}

	fmt.Println("resp:", resp.Status)
}
