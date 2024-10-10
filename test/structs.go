package main

import "time"

type response struct {
	Data struct {
		Presentation struct {
			Typename              string `json:"__typename"`
			StayProductDetailPage struct {
				Typename string `json:"__typename"`
				Reviews  struct {
					Typename string `json:"__typename"`
					Reviews  []struct {
						Typename      string    `json:"__typename"`
						CollectionTag any       `json:"collectionTag"`
						Comments      string    `json:"comments"`
						ID            string    `json:"id"`
						Language      string    `json:"language"`
						CreatedAt     time.Time `json:"createdAt"`
						Reviewee      struct {
							Typename           string `json:"__typename"`
							Deleted            bool   `json:"deleted"`
							FirstName          string `json:"firstName"`
							HostName           string `json:"hostName"`
							ID                 string `json:"id"`
							PictureURL         string `json:"pictureUrl"`
							ProfilePath        string `json:"profilePath"`
							IsSuperhost        bool   `json:"isSuperhost"`
							UserProfilePicture struct {
								Typename      string `json:"__typename"`
								BaseURL       string `json:"baseUrl"`
								OnPressAction struct {
									Typename string `json:"__typename"`
									URL      string `json:"url"`
								} `json:"onPressAction"`
							} `json:"userProfilePicture"`
						} `json:"reviewee"`
						Reviewer struct {
							Typename           string `json:"__typename"`
							Deleted            bool   `json:"deleted"`
							FirstName          string `json:"firstName"`
							HostName           string `json:"hostName"`
							ID                 string `json:"id"`
							PictureURL         string `json:"pictureUrl"`
							ProfilePath        string `json:"profilePath"`
							IsSuperhost        bool   `json:"isSuperhost"`
							UserProfilePicture struct {
								Typename      string `json:"__typename"`
								BaseURL       string `json:"baseUrl"`
								OnPressAction struct {
									Typename string `json:"__typename"`
									URL      string `json:"url"`
								} `json:"onPressAction"`
							} `json:"userProfilePicture"`
						} `json:"reviewer"`
						ReviewHighlight           string `json:"reviewHighlight"`
						HighlightType             string `json:"highlightType"`
						LocalizedDate             string `json:"localizedDate"`
						LocalizedRespondedDate    string `json:"localizedRespondedDate"`
						LocalizedReviewerLocation string `json:"localizedReviewerLocation"`
						LocalizedReview           any    `json:"localizedReview"`
						Rating                    int    `json:"rating"`
						RatingAccessibilityLabel  string `json:"ratingAccessibilityLabel"`
						RecommendedNumberOfLines  any    `json:"recommendedNumberOfLines"`
						Response                  string `json:"response"`
						RoomTypeListingTitle      any    `json:"roomTypeListingTitle"`
						HighlightedReviewSentence []any  `json:"highlightedReviewSentence"`
						HighlightReviewMentioned  any    `json:"highlightReviewMentioned"`
						ShowMoreButton            struct {
							Typename         string `json:"__typename"`
							Title            string `json:"title"`
							LoggingEventData struct {
								Typename            string `json:"__typename"`
								LoggingID           string `json:"loggingId"`
								Experiments         []any  `json:"experiments"`
								EventData           any    `json:"eventData"`
								EventDataSchemaName any    `json:"eventDataSchemaName"`
								Section             any    `json:"section"`
								Component           any    `json:"component"`
							} `json:"loggingEventData"`
						} `json:"showMoreButton"`
						SubtitleItems           []any `json:"subtitleItems"`
						Channel                 any   `json:"channel"`
						ReviewMediaItems        []any `json:"reviewMediaItems"`
						IsHostHighlightedReview any   `json:"isHostHighlightedReview"`
						ReviewPhotoUrls         []any `json:"reviewPhotoUrls"`
					} `json:"reviews"`
					Metadata struct {
						Typename                     string `json:"__typename"`
						ReviewsCount                 int    `json:"reviewsCount"`
						ModalSubtitle                any    `json:"modalSubtitle"`
						IsReviewsSearchResults       any    `json:"isReviewsSearchResults"`
						IsReviewsHighlightTagResults any    `json:"isReviewsHighlightTagResults"`
						IsAutoTranslateOn            bool   `json:"isAutoTranslateOn"`
						EndCursor                    any    `json:"endCursor"`
						Experiments                  []any  `json:"experiments"`
						ReviewTags                   []any  `json:"reviewTags"`
					} `json:"metadata"`
				} `json:"reviews"`
			} `json:"stayProductDetailPage"`
		} `json:"presentation"`
	} `json:"data"`
	Extensions struct {
		TraceID string `json:"traceId"`
	} `json:"extensions"`
}
