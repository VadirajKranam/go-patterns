package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/vadiraj/gopher/internal/store"
)

var usernames = []string{
	"ArjunCoder", "DevikaDev", "RohanRocks", "MeeraTechie", "KiranBuilder",
	"SitaSyntax", "VikramGo", "PriyaGuru", "KrishnaKeys", "RadhaRoutines",
	"AmitByte", "AnjaliLoops", "RahulStack", "SnehaStruct", "ManojPointer",
	"NehaGoLang", "VivekCloud", "PreetiHack", "AjayDebug", "RaniRuntime",
	"IshaGoDev", "KartikGopher", "NishaCLI", "ArvindBinary", "MeghaCode",
	"RajGoGuru", "PoojaQuery", "AbhayThread", "GeetaChannel", "AshokCompile",
	"DeepaQueue", "KishoreTask", "LataClosure", "VarunAPI", "JyotiCache",
	"AnandHandler", "ChitraParse", "SureshStream", "KavyaRoutine", "MohitSocket",
	"LeelaScript", "HariNode", "RitikaData", "GaneshModule", "RituFiber",
	"BalaMutex", "KomalPacket", "JayaSemaphore", "TarunHeap", "RekhaLock",
	"SunilGolang",
}

var titles = []string{
	"Go Slices 101",
	"Slices vs Arrays",
	"Append in Go",
	"Slice Capacity",
	"Copying Slices",
	"Slice Basics",
	"Multidimensional Slices",
	"Sorting Slices",
	"Slice Tricks",
	"Slice Debugging",
	"Slices in Go",
	"Slice Headers",
	"Sub-Slicing",
	"Efficient Slices",
	"Slice Gotchas",
	"Go Slices Guide",
	"Go Slice Tips",
	"Slices Made Easy",
	"Slice Examples",
	"Master Slices",
}

var contents = []string{
	"Learn what slices are and how they work in Go, including their dynamic nature.",
	"Understand the key differences between slices and arrays in Go.",
	"Explore how to use the `append` function to add elements to a slice.",
	"Dive into the concept of slice capacity and how to manage it effectively.",
	"Learn how to copy slices efficiently using the `copy` function.",
	"Understand the basics of slice creation, indexing, and iteration.",
	"Discover how to create and work with multidimensional slices in Go.",
	"Learn how to sort slices using the `sort` package in Go.",
	"Explore useful tricks and hacks to manipulate slices in Golang.",
	"Debug common slice issues, including bounds errors and nil slices.",
	"Understand the fundamental role of slices in Go applications.",
	"Dive into the slice header and learn about length, capacity, and pointers.",
	"Learn how to extract subsets of slices using Go's slicing syntax.",
	"Tips for writing efficient code using slices in Go.",
	"Avoid common slice pitfalls, like modifying shared slices.",
	"A quick guide to using slices for beginners in Go programming.",
	"Explore how to handle slices safely in concurrent programs.",
	"Learn how to simplify slice usage with practical tips and examples.",
	"Hands-on examples of how slices are used in real-world Go projects.",
	"Master slice operations and make your Go programs more robust.",
}

var tags = []string{
	"golang", "slices", "arrays", "golang-tips", "go-programming",
	"golang-basics", "golang-advanced", "programming", "backend", "data-structures",
	"code-efficiency", "go-language", "slice-tutorials", "go-debugging", "golang-best-practices",
	"golang-concurrency", "code-optimization", "go-syntax", "golang-guide", "coding-tricks",
}

var rand_comments = []string{
	"Great article! Learned a lot about Go slices.",
	"I had no idea about slice capacity. Thanks for the explanation!",
	"Awesome tips on slicing and appending data in Go.",
	"This really cleared up my confusion about slices vs arrays.",
	"Loved the practical examples! Keep it up.",
	"Never thought about slices this way. Very insightful.",
	"Can’t wait to try these slice tricks in my projects.",
	"The example on copying slices was super helpful.",
	"I’ve been struggling with slice issues, this guide helped a lot.",
	"This helped me understand slice initialization in Go.",
	"Really useful content for someone new to Go.",
	"I had no clue about multidimensional slices. Great explanation.",
	"I always made mistakes with slices, your post helped me avoid them!",
	"This was a fantastic overview of slices in Go.",
	"Your posts are always so clear and easy to follow!",
	"It’s great to see slice tips for performance optimization.",
	"This blog makes Go slices look way simpler than I thought.",
	"Looking forward to more articles on Go data structures!",
	"I didn't know slices could be this efficient. Thanks for the post.",
	"Excellent breakdown of the slice header concept.",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()
	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			log.Println("Error creating user:", err)
			return
		}

	}
	tx.Commit()
	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating post:", err)
			return
		}
	}
	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}
	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			UserName: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@seeding.com",
			Role: store.Role{
				Name: "user",
			},
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
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
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
			UserID:  users[rand.Intn(len(users))].ID,
			PostID:  posts[rand.Intn(len(posts))].ID,
			Content: rand_comments[rand.Intn(len(rand_comments))],
		}
	}
	return comments
}
