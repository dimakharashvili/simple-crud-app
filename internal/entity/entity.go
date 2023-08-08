package entity

type RedditPost struct {
	UUID     string
	Title    string
	Likes    uint32
	Comments []*Comment
}

type Comment struct {
	UUID  string
	Body  string
	Likes uint32
}
