package client

import "strconv"

const limit = 100

type Pagination struct {
	offset int
	params map[string]string
}

type Paginated interface {
	Environment
	getEndpoint() string
}

func (p *Pagination) getParams() map[string]string {
	return map[string]string{
		"limit":  strconv.Itoa(limit),
		"offset": strconv.Itoa(p.offset),
	}
}

// returns true if there is more data.
func (p *Pagination) next(currentPageSize int) bool {
	p.offset += currentPageSize

	return currentPageSize == limit
}

// params - additional params. may be nil.
func getAll[P Paginated](client *ApiClient, params map[string]string) ([]P, error) {
	p := Pagination{
		offset: 0,
		params: params,
	}

	var allResults []P

	for {
		pageParams := p.getParams()
		for k, v := range params {
			pageParams[k] = v
		}

		var pageResults []P

		err := client.http.Get(P{}.getEndpoint(), pageParams, &pageResults)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, pageResults...)

		if more := p.next(len(pageResults)); !more {
			break
		}
	}

	return allResults, nil
}
