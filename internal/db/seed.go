package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/igorzinar/goSocial/internal/store"
	"log"
	"math/rand"
)

var usernames = []string{
	"alice", "bob", "carol", "dave", "eve",
	"frank", "grace", "heidi", "ivan", "judy",
	"mallory", "oscar", "peggy", "trent", "victor",
	"walter", "charlie", "susan", "nancy", "roger",
	"molly", "steve", "wendy", "paul", "lucy",
	"george", "fiona", "harry", "linda", "michael",
	"oliver", "patty", "quinn", "rachel", "sam",
	"tina", "ursula", "vince", "wanda", "xander",
	"yvonne", "zane", "amber", "brian", "celia",
	"dan", "ella", "felix", "gina", "hannah",
	"isabella", "jack", "karen", "leo", "mia",
	"noah", "olivia", "peter", "quentin", "ruby",
	"sophia", "thomas", "ulysses", "victoria", "will",
	"xena", "yasmin", "zoe", "adam", "bella",
	"chris", "diana", "ethan", "faith", "gabriel",
	"harper", "iris", "jake", "kim", "logan",
	"mason", "nina", "owen", "paige", "quincy",
	"rebecca", "shane", "tara", "uma", "violet",
	"wayne", "xavier", "yara", "zeke", "amelia",
}

var titlesList = []string{
	"Top Coding Practices",
	"Understanding APIs",
	"Mastering Data Structures",
	"Intro to Cloud Computing",
	"Building Your First App",
	"JavaScript Tips and Tricks",
	"Beginner’s Guide to Git",
	"The Future of AI",
	"Optimizing Your Code",
	"Getting Started with React",
	"Exploring Machine Learning",
	"UI/UX Design Basics",
	"Guide to Responsive Design",
	"Introduction to Docker",
	"Writing Clean Code",
	"Understanding SEO Basics",
	"Creating a Portfolio Site",
	"Debugging Like a Pro",
	"Learning SQL Fundamentals",
	"Getting Started with Golang",
}

var contentsList = []string{
	"Learn the top coding practices that help developers write cleaner, more efficient code.",
	"A beginner-friendly guide to understanding APIs and how they work.",
	"Master essential data structures to improve coding efficiency and performance.",
	"Discover the basics of cloud computing and its impact on the tech industry.",
	"Step-by-step guide to building your first application from scratch.",
	"JavaScript tips and tricks to level up your web development skills.",
	"A comprehensive guide for beginners on using Git version control.",
	"Exploring the potential and challenges in the future of artificial intelligence.",
	"Learn key techniques to optimize your code and improve application speed.",
	"An introductory guide to starting your journey with React.",
	"Understand the fundamentals of machine learning and its applications.",
	"Basic principles of UI/UX design that enhance user experience.",
	"An overview of responsive design and its importance in web development.",
	"Getting started with Docker and understanding containerized applications.",
	"Best practices for writing clean and maintainable code.",
	"An introductory guide to SEO basics for improving search engine ranking.",
	"Tips for creating a professional portfolio site to showcase your work.",
	"Learn efficient debugging techniques to identify and solve code issues.",
	"A beginner's guide to SQL for managing and querying databases.",
	"An introduction to Go programming language and its unique features.",
}

var tags = []string{
	"coding", "programming", "web development", "API", "data structures",
	"cloud computing", "app development", "JavaScript", "Git", "AI",
	"code optimization", "React", "machine learning", "UI/UX", "responsive design",
	"Docker", "clean code", "SEO", "portfolio", "SQL",
	"debugging", "Golang", "frontend", "backend", "databases",
	"development tips", "tutorial", "best practices", "beginner", "advanced",
}

var commentsList = []string{
	"Great article! Learned a lot from this.",
	"Thanks for the insights, very helpful!",
	"Could you dive deeper into this topic in the next post?",
	"I tried this, but I’m facing some issues. Any advice?",
	"This is exactly what I was looking for. Thanks!",
	"Fantastic breakdown, very easy to understand.",
	"Can you suggest any additional resources?",
	"I disagree with some points, but overall a good read.",
	"I’ve been struggling with this concept—thanks for clarifying!",
	"What tools do you recommend for beginners?",
	"This was a bit advanced for me, but still informative.",
	"Thanks for sharing, looking forward to more posts like this!",
	"Any tips for applying this in a real-world project?",
	"Clear and concise, just what I needed.",
	"This approach saved me a ton of time!",
	"Are there any updates on this topic?",
	"Could you share some example code?",
	"I found a typo in the third section.",
	"Loved the examples, they made it much easier to follow!",
	"This was super helpful, thank you!",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()
	users := generateUsers(500)
	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Fatal("Seed generateUsers ====> ", err)
			return
		}
	}
	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Fatal("Seed generatePosts ====> ", err)
			return
		}
	}

	comments := generateComments(200, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Fatal("Seed generateComments ====> ", err)
			return
		}
	}
	log.Println("Seeding completed")
}

func generateUsers(num int) []*store.User {

	users := make([]*store.User, num)
	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titlesList[rand.Intn(len(titlesList))],
			Content: contentsList[rand.Intn(len(contentsList))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {

	comments := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		comments[i] = &store.Comment{
			ID:      0,
			UserID:  users[rand.Intn(len(users))].ID,
			PostID:  posts[rand.Intn(len(posts))].ID,
			Content: commentsList[rand.Intn(len(commentsList))],
		}
	}

	return comments
}
