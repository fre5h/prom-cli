package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/fre5h/prom-cli/internal/models"
)

type Client struct {
	apiKey     string
	httpClient http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (ac Client) GetGroupList(limit int, lastId int) ([]models.Group, error) {
	var groups []models.Group

	for {
		groupsChunk, err := ac.doGetGroupList(limit, lastId)
		if err != nil {
			return nil, err
		}
		if len(groupsChunk) == 0 {
			break
		}

		groups = append(groups, groupsChunk...)

		if len(groups) >= limit {
			break
		}

		lastId = groupsChunk[len(groupsChunk)-1].ID
	}

	return groups, nil
}

func (ac Client) doGetGroupList(limit int, lastId int) ([]models.Group, error) {
	var response *http.Response
	var err error

	req := createRequest(http.MethodGet, "https://my.prom.ua/api/v1/groups/list", ac.apiKey, nil)

	// Process query parameters
	q := url.Values{}
	if limit > 0 {
		q.Add("limit", strconv.Itoa(limit))
	}
	if lastId > 0 {
		q.Add("last_id", strconv.Itoa(lastId))
	}
	req.URL.RawQuery = q.Encode()

	response, err = ac.httpClient.Do(req)
	defer closeBody(response.Body)

	if err != nil {
		return nil, fmt.Errorf("http client: помилка на створенні запиту: %s", err)
	}

	if response.StatusCode == http.StatusOK {
		bodyBytes, errRead := ioutil.ReadAll(response.Body)
		if errRead != nil {
			return nil, fmt.Errorf("error in reading response body: %s", err)
		}

		data := models.Groups{}

		if err = json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("error on unmarshaling json: %s", err)
		}

		return data.Groups, nil
	}

	return nil, fmt.Errorf("result code is not 200, it is %d", response.StatusCode)
}

func (ac Client) GetProductList(limit int, lastId int, groupId int) ([]models.Product, error) {
	var products []models.Product

	for {
		productsChunk, err := ac.doGetProductList(limit, lastId, groupId)
		if err != nil {
			return nil, err
		}
		if len(productsChunk) == 0 {
			break
		}

		products = append(products, productsChunk...)

		if len(products) >= limit {
			break
		}

		lastId = productsChunk[len(productsChunk)-1].ID
	}

	return products, nil
}

func (ac Client) doGetProductList(limit int, lastId int, groupId int) ([]models.Product, error) {
	var response *http.Response
	var err error

	req := createRequest(http.MethodGet, "https://my.prom.ua/api/v1/products/list", ac.apiKey, nil)

	// Process query parameters
	q := url.Values{}
	if limit > 0 {
		q.Add("limit", strconv.Itoa(limit))
	}
	if lastId > 0 {
		q.Add("last_id", strconv.Itoa(lastId))
	}
	if groupId > 0 {
		q.Add("group_id", strconv.Itoa(groupId))
	}
	req.URL.RawQuery = q.Encode()

	response, err = ac.httpClient.Do(req)
	defer closeBody(response.Body)

	if err != nil {
		return nil, fmt.Errorf("http client: помилка на створенні запиту: %s", err)
	}

	if response.StatusCode == http.StatusOK {
		bodyBytes, errRead := ioutil.ReadAll(response.Body)
		if errRead != nil {
			return nil, fmt.Errorf("error in reading response body: %s", err)
		}

		data := models.Products{}

		if err = json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("error on unmarshaling json: %s", err)
		}

		return data.Products, nil
	}

	return nil, fmt.Errorf("result code is not 200, it is %d", response.StatusCode)
}

func (ac Client) UpdateProduct(products []models.ProductUpdate) error {
	var response *http.Response
	var err error

	jsonStr, _ := json.Marshal(products)

	body := bytes.NewBuffer(jsonStr)

	req := createRequest(http.MethodPost, "https://my.prom.ua/api/v1/products/update", ac.apiKey, body)

	if response, err = ac.httpClient.Do(req); err != nil {
		return fmt.Errorf("client: error making http request: %s", err)
	}

	defer closeBody(response.Body)

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("result code is not 200, it is %d", response.StatusCode)
	}

	return nil
}

func createRequest(method string, url string, apiKey string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Add authorization
	req.Header.Set("Authorization", "Bearer "+apiKey)

	return req
}

func closeBody(body io.ReadCloser) {
	err := body.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
