package data

//go:generate ../../../makelens data blog.go blog_generated.go
type Comment struct {
	Title   string
	Content string
}

type Rating struct {
	Author string
	Stars  int
}

type BlogPost struct {
	Content  string
	Comments []Comment
	Ratings  []Rating
}

func NewBlogPost(content string, comments []Comment) BlogPost {
	return BlogPost{
		Content:  content,
		Comments: comments,
	}
}

func NewComment(title string, content string) Comment {
	return Comment{
		Title:   title,
		Content: content,
	}
}
