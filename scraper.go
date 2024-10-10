package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"sync"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/data-harvesters/goapify"
)

var (
	apiKeyRe = regexp.MustCompile(`(?m)"api_config":{"key":"(.*?)"`)
)

type scraper struct {
	actor *goapify.Actor
	input *input

	apiKey string

	client tls_client.HttpClient
}

func newScraper(input *input, actor *goapify.Actor) (*scraper, error) {
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_124),
		tls_client.WithNotFollowRedirects(),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	return &scraper{
		actor:  actor,
		input:  input,
		client: client,
	}, nil
}

func (s *scraper) Run() {
	fmt.Println("beginning scrapping...")

	apiKey, err := s.getApiKey(s.input.RoomIds[0])
	if err != nil {
		fmt.Printf("Failed to get api key: %v\n", err)
		return
	}
	s.apiKey = apiKey

	var wg sync.WaitGroup
	for _, roomId := range s.input.RoomIds {
		r, err := s.scrapeReviews(roomId)
		if err != nil {
			fmt.Printf("%s: failed to get total reviews: %v\n", roomId, err)
			continue
		}
		totalReviews := r.Data.Presentation.StayProductDetailPage.Reviews.Metadata.ReviewsCount

		wg.Add(1)
		go func() {
			defer wg.Done()

			s.startScrape(roomId, totalReviews)
		}()
	}
	fmt.Println("succesfully scraped all reviews")
}

func (s *scraper) startScrape(roomId string, totalReviews int) {
	q := NewQueue()
	finished := make(chan bool)

	totalReviewsScrapped := 0

	go func() {
		for {
			i := q.Pop()
			if i == nil {
				break
			}
			fmt.Printf("%s: scrapping offset: %d\n", roomId, i.(int))

			resp, err := s.scrapeReviews(i.(string))
			if err != nil {
				fmt.Printf("%s: Failed to scrape agents: %v\n", roomId, err)
				continue
			}
			fmt.Printf("%s: scrapping offset: %d for %d reviews\n", roomId, i.(int), len(resp.Data.
				Presentation.StayProductDetailPage.Reviews.Reviews))

			totalReviewsScrapped += len(resp.Data.
				Presentation.StayProductDetailPage.Reviews.Reviews)

			err = s.actor.Output(resp.Data.
				Presentation.StayProductDetailPage.Reviews.Reviews)
			if err != nil {
				fmt.Printf("%s: Failed to send output: %v\n", roomId, err)
				continue
			}
			fmt.Println("%s: succesfully sent output!", roomId)

			time.Sleep(500 * time.Millisecond)
		}

		finished <- true
	}()

	for i := s.input.Offset; i < totalReviews; i += s.input.Limit - 1 {
		go func(i int) {
			q.Push(i)
		}(i)
	}
	<-finished
	fmt.Printf("%s: succesfully scraped all reviews! Total: %d\n", roomId, totalReviewsScrapped)

}

func (s *scraper) scrapeReviews(roomId string) (*response, error) {
	u, err := url.Parse(`https://www.airbnb.com/api/v3/StaysPdpReviewsQuery/dec1c8061483e78373602047450322fd474e79ba9afa8d3dbbc27f504030f91d?operationName=StaysPdpReviewsQuery&locale=en&currency=USD&variables={"id":"U3RheUxpc3Rpbmc6MTQxMjY2NTc=","pdpReviewsRequest":{"fieldSelector":"for_p3_translation_only","forPreview":false,"limit":24,"offset":"0","showingTranslationButton":false,"first":24,"sortingPreference":"MOST_RECENT","checkinDate":"2024-12-06","checkoutDate":"2024-12-08","numberOfAdults":"1","numberOfChildren":"0","numberOfInfants":"0","numberOfPets":"0"}}&extensions={"persistedQuery":{"version":1,"sha256Hash":"dec1c8061483e78373602047450322fd474e79ba9afa8d3dbbc27f504030f91d"}}`)
	if err != nil {
		return nil, err
	}

	id := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("StayListing:%s", roomId)))

	query := u.Query()
	query.Set("variables", fmt.Sprintf(`{"id":"%s","pdpReviewsRequest":{"fieldSelector":"for_p3_translation_only","forPreview":false,"limit":%d,"offset":"%v","showingTranslationButton":false,"first":20,"sortingPreference":"MOST_RECENT","checkinDate":"2024-12-06","checkoutDate":"2024-12-08","numberOfAdults":"1","numberOfChildren":"0","numberOfInfants":"0","numberOfPets":"0"}}`,
		id, s.input.Limit, s.input.Offset))
	u.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Host", "www.airbnb.com")
	req.Header.Add("Sec-Ch-Ua", "\"Chromium\";v=\"127\", \"Not)A;Brand\";v=\"99\"")
	req.Header.Add("X-Airbnb-Supports-Airlock-V2", "true")
	req.Header.Add("X-Csrf-Token", "")
	req.Header.Add("X-Airbnb-Api-Key", s.apiKey)
	req.Header.Add("Accept-Language", "en-US")
	req.Header.Add("Sec-Ch-Ua-Platform-Version", "\"\"")
	req.Header.Add("X-Niobe-Short-Circuited", "true")
	req.Header.Add("Dpr", "1.25")
	req.Header.Add("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Add("Device-Memory", "8")
	req.Header.Add("X-Airbnb-Graphql-Platform-Client", "minimalist-niobe")
	req.Header.Add("X-Client-Version", "411703a6982e9d5a88e44a77a1f98b22c9d3a6dc")
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Add("X-Csrf-Without-Token", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.6533.100 Safari/537.36")
	req.Header.Add("Viewport-Width", "1536")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Ect", "3g")
	req.Header.Add("X-Airbnb-Graphql-Platform", "web")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	// req.Header.Add("Referer", "https://www.airbnb.com/rooms/14126657?adults=1&children=0&enable_m3_private_room=true&infants=0&pets=0&search_mode=regular_search&check_in=2024-12-06&check_out=2024-12-08&source_impression_id=p3_1728567041_P3v0p6U0kS9rLco5&previous_page_section_name=1000&federated_search_id=d27fcca1-ae48-4344-8dbe-87f8e99f3f84")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Priority", "u=1, i")

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to scrape reviews: %d %s", res.StatusCode, string(b))
	}

	var response response
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *scraper) getApiKey(roomId string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.airbnb.com/rooms/%s/reviews", roomId), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Host", "www.airbnb.com")
	req.Header.Add("Sec-Ch-Ua", "\"Chromium\";v=\"127\", \"Not)A;Brand\";v=\"99\"")
	req.Header.Add("Accept-Language", "en-US")
	req.Header.Add("Sec-Ch-Ua-Platform-Version", "\"\"")
	req.Header.Add("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Add("Device-Memory", "8")
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.6533.100 Safari/537.36")
	req.Header.Add("Viewport-Width", "1536")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Ect", "3g")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Priority", "u=1, i")

	res, err := s.client.Do(req)
	if err != nil {
		return "", err

	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get api key: %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	matches := apiKeyRe.FindAllStringSubmatch(string(b), -1)

	if len(matches) == 0 {
		return "", errors.New("failed to find api key")
	}
	apiKey := matches[0][1]

	return apiKey, nil
}
