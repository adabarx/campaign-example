package data

type BlogPost struct {
	Slug    string
	Title   string
	Content string
	Date    string
}

func GetBlogPosts() []BlogPost {
	return []BlogPost{
		{
			Slug:    "welcome",
			Title:   "Welcome to the Campaign",
			Content: "We're excited to launch this campaign. Our goal is to make a real difference.",
			Date:    "2025-11-02",
		},
		{
			Slug:    "how-to-contribute",
			Title:   "How to Contribute",
			Content: "There are many ways to support our cause. You can donate, volunteer, or share our message.",
			Date:    "2025-11-01",
		},
	}
}
