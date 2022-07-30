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

	"github.com/fatih/color"

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

// GetGroupList gets groups with limit and last ID
func (c Client) GetGroupList(limit int, lastId int) (groups []models.Group, err error) {
	for {
		groupsChunk, err := c.doGetGroupList(limit, lastId)
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

func (c Client) doGetGroupList(limit int, lastId int) ([]models.Group, error) {
	var response *http.Response
	var err error

	req := createRequest(http.MethodGet, "https://my.prom.ua/api/v1/groups/list", c.apiKey, nil)

	// Process query parameters
	q := url.Values{}
	if limit > 0 {
		q.Add("limit", strconv.Itoa(limit))
	}
	if lastId > 0 {
		q.Add("last_id", strconv.Itoa(lastId))
	}
	req.URL.RawQuery = q.Encode()

	response, err = c.httpClient.Do(req)
	defer closeBody(response.Body)

	if err != nil {
		return nil, fmt.Errorf("http client: помилка при відправці запиту: %s", err)
	}

	if response.StatusCode == http.StatusOK {
		bodyBytes, errRead := ioutil.ReadAll(response.Body)
		if errRead != nil {
			return nil, fmt.Errorf("помилка на читанні відповіді: %s", err)
		}

		data := models.Groups{}
		if err = json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("помилка на декодуванні json: %s", err)
		}

		return data.Groups, nil
	}

	return nil, formatStatusCodeError(response.StatusCode)
}

// GetProductList gets products with limit, last ID and group ID
func (c Client) GetProductList(limit int, lastId int, groupId int) (products []models.Product, err error) {
	for {
		productsChunk, err := c.doGetProductList(limit, lastId, groupId)
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

func (c Client) doGetProductList(limit int, lastId int, groupId int) ([]models.Product, error) {
	var response *http.Response
	var err error

	req := createRequest(http.MethodGet, "https://my.prom.ua/api/v1/products/list", c.apiKey, nil)

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

	response, err = c.httpClient.Do(req)
	defer closeBody(response.Body)

	if err != nil {
		return nil, fmt.Errorf("http client: помилка при відправці запиту: %s", err)
	}

	if response.StatusCode == http.StatusOK {
		bodyBytes, errRead := ioutil.ReadAll(response.Body)
		if errRead != nil {
			return nil, fmt.Errorf("помилка на читанні відповіді: %s", err)
		}

		data := models.Products{}
		if err = json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("помилка на декодуванні json: %s", err)
		}

		return data.Products, nil
	}

	return nil, formatStatusCodeError(response.StatusCode)
}

// UpdateProduct updates products
func (c Client) UpdateProduct(products []models.ProductUpdate) (*models.ProductUpdateResponse, error) {
	var response *http.Response
	var err error

	jsonStr, _ := json.Marshal(products)
	req := createRequest(http.MethodPost, "https://my.prom.ua/api/v1/products/update", c.apiKey, bytes.NewBuffer(jsonStr))

	response, err = c.httpClient.Do(req)
	defer closeBody(response.Body)

	if err != nil {
		return nil, fmt.Errorf("http client: помилка при відправці запиту: %s", err)
	}

	if response.StatusCode == http.StatusOK {
		bodyBytes, errRead := ioutil.ReadAll(response.Body)
		if errRead != nil {
			return nil, fmt.Errorf("помилка на читанні відповіді: %s", err)
		}

		data := models.ProductUpdateResponse{}
		if err = json.Unmarshal(bodyBytes, &data); err != nil {
			return nil, fmt.Errorf("помилка на декодуванні json: %s", err)
		}

		return &data, nil
	}

	return nil, formatStatusCodeError(response.StatusCode)
}

func createRequest(method string, url string, apiKey string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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

func formatStatusCodeError(statusCode int) error {
	return fmt.Errorf("неуспішний запит, отриманий статус код %s", color.New(color.FgRed).Sprintf("%d", statusCode))
}
