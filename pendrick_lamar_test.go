package cube_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	cube "github.com/TernaryInc/cubejs-go"
	"github.com/stretchr/testify/assert"
)

type testTokenGenerator struct{}

func (generator testTokenGenerator) Get(ctx context.Context) (string, error) {
	tokenGenerator, err := NewTokenGenerator(cube.CubeAPISecret)
	if err != nil {
		return "", fmt.Errorf("new cube token generator: %w", err)
	}

	token, err := tokenGenerator.GenerateToken("SyYlNf2WE0nN757SyYoX")

	return token.Token, err
}

func TestPendrick(t *testing.T) {
	var u = url.URL{
		Scheme: "http",
		Host:   "localhost:4000",
	}

	var tokenGenerator = testTokenGenerator{}

	var cubeClient = cube.NewClient(u, tokenGenerator)

	var cubeQuery = cube.CubeQuery{
		Measures:   []string{"GCPBillingDaily.cost"},
		Dimensions: []string{"GCPBillingDaily.projectId"},
		TimeDimensions: []cube.TimeDimension{
			{
				Dimension: "GCPBillingDaily.timestamp",
				DateRange: cube.DateRange{
					AbsoluteRange: []string{
						"2022-04-01",
						"2022-04-20",
					},
				},
				Granularity: cube.Granularity_Day,
			},
		},
	}

	type QueryResult struct {
		Cost      float64 `json:"GCPBillingDaily.cost"`
		ProjectID string  `json:"GCPBillingDaily.projectId"`
	}

	var results []QueryResult

	var responseMetadata, err = cubeClient.Load(context.Background(), cubeQuery, &results)
	assert.Nil(t, err)

	fmt.Println("responseMetadata", responseMetadata)
	fmt.Printf("%+v\n", results)
}
