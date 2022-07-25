package api

import (
	"encoding/json"
	"errors"
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
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
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

		lastId = groupsChunk[len(groupsChunk)-1].Id
	}

	return groups, nil
}

func (ac Client) doGetGroupList(limit int, lastId int) ([]models.Group, error) {
	var req *http.Request
	var response *http.Response
	var err error

	if req, err = http.NewRequest(http.MethodGet, "https://my.prom.ua/api/v1/groups/list", nil); err != nil {
		return nil, fmt.Errorf("client: could not create request: %s", err)
	}
	req.Header.Set("Authorization", "Bearer "+ac.apiKey)

	// Process query parameters
	q := url.Values{}
	if limit > 0 {
		q.Add("limit", strconv.Itoa(limit))
	}
	if lastId > 0 {
		q.Add("last_id", strconv.Itoa(lastId))
	}
	req.URL.RawQuery = q.Encode()

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	if response, err = client.Do(req); err != nil {
		return nil, fmt.Errorf("client: error making http request: %s", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error on closing body: %s", err)
			os.Exit(1)
		}
	}(response.Body)

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

		lastId = productsChunk[len(productsChunk)-1].Id
	}

	return products, nil
}

func (ac Client) doGetProductList(limit int, lastId int, groupId int) ([]models.Product, error) {
	var req *http.Request
	var response *http.Response
	var err error

	if req, err = http.NewRequest(http.MethodGet, "https://my.prom.ua/api/v1/products/list", nil); err != nil {
		return nil, errors.New(fmt.Sprintf("client: could not create request: %s", err))
	}

	// Add authorization
	req.Header.Set("Authorization", "Bearer "+ac.apiKey)

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

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	if response, err = client.Do(req); err != nil {
		return nil, errors.New(fmt.Sprintf("client: error making http request: %s", err))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error on closing body: %s", err)
			os.Exit(1)
		}
	}(response.Body)

	if response.StatusCode == http.StatusOK {
		bodyBytes, errRead := ioutil.ReadAll(response.Body)
		if errRead != nil {
			return nil, errors.New(fmt.Sprintf("error in reading response body: %s", err))
		}

		data := models.Products{}

		if err = json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, errors.New(fmt.Sprintf("error on unmarshaling json: %s", err))
		}

		return data.Products, nil
	}

	return nil, errors.New(fmt.Sprintf("result code is not 200, it is %d", response.StatusCode))
}
