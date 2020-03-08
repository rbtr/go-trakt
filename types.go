package trakt

type IDs struct {
	Trakt int    `json:"trakt,omitempty"`
	Slug  string `json:"slug,omitempty"`
	TVDB  int    `json:"tvdb,omitempty"`
	IMDB  string `json:"imdb,omitempty"`
	TMDb  int    `json:"tmdb,omitempty"`
}

type Movie struct {
	Title string `json:"title,omitempty"`
	Year  int    `json:"year,omitempty"`
	IDs   IDs    `json:"ids,omitempty"`
}

type Show struct {
	Title string `json:"title,omitempty"`
	Year  int    `json:"year,omitempty"`
	IDs   IDs    `json:"ids,omitempty"`
}

type Season struct {
	Number int `json:"number,omitempty"`
	IDs    IDs `json:"ids,omitempty"`
}

type Episode struct {
	Season int    `json:"season,omitempty"`
	Number int    `json:"number,omitempty"`
	Title  string `json:"title,omitempty"`
	IDs    IDs    `json:"ids,omitempty"`
}

type Person struct {
	Name string `json:"name,omitempty"`
	IDs  IDs    `json:"ids,omitempty"`
}
