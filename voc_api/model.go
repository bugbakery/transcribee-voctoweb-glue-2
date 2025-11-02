package voc_api

type Recording struct {
    Filename     string  `json:"filename"`
    MimeType     string  `json:"mime_type"`
    Language     string  `json:"language"`
    Folder       *string `json:"folder,omitempty"`
    Size         *int    `json:"size,omitempty"`
    Length       *int    `json:"length,omitempty"`
    State        *string `json:"state,omitempty"`
    HighQuality  bool    `json:"high_quality"`
    Width        *int    `json:"width,omitempty"`
    Height       *int    `json:"height,omitempty"`
    UpdatedAt    string  `json:"updated_at"`
    RecordingUrl string  `json:"recording_url"`
}

type Talk struct {
    Guid             string     `json:"guid"`
    Slug             string     `json:"slug"`
    Title            string     `json:"title"`
    Date             string     `json:"date"`
    Subtitle         *string    `json:"subtitle,omitempty"`
    Link             string     `json:"link"`
    Description      string     `json:"description"`
    OriginalLanguage string     `json:"original_language"`
    Persons          []string   `json:"persons"`
    Tags             []string   `json:"tags"`
    ViewCount        int        `json:"view_count"`
    Promoted         bool       `json:"promoted"`
    ReleaseDate      string     `json:"release_date"`
    UpdatedAt        string     `json:"updated_at"`
    Length           int        `json:"length"`
    Duration         int        `json:"duration"`
    ThumbUrl         string     `json:"thumb_url"`
    PosterUrl        string     `json:"poster_url"`
    TimelineUrl      string     `json:"timeline_url"`
    ThumbnailsUrl    string     `json:"thumbnails_url"`
    FrontendLink     string     `json:"frontend_link"`
    Url              string     `json:"url"`
    Related          []string   `json:"related"`
    Recordings       []Recording `json:"recordings"`
}

type Video struct {
    Filename string `json:"filename"`
}

type TalkSummary struct {
    Guid  string `json:"guid"`
    Slug  string `json:"slug"`
    Title string `json:"title"`
    Date  string `json:"date"`
    Video Video  `json:"video"`
}

type Conference struct {
    Id     string         `json:"id"`
    Title  string         `json:"title"`
    Talks []TalkSummary `json:"events"`
}
