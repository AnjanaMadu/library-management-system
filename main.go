package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-adodb"
	"github.com/olekukonko/tablewriter"
)

func returnPg(typee, msg, ref string) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Info</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/tailwindcss/dist/tailwind.min.css">
</head>
<body class="bg-gray-800">
  <main class="py-8">
    <div class="container mx-auto">
      <h1 class="text-3xl font-bold text-white text-center mb-8">%s</h1>
      <p class="text-xl text-gray-100 text-center mb-8">%s</p>
      <a href="%s" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-full">Go Back</a>
    </div>
  </main>
</body>
</html>
	`, typee, msg, ref)
}

func checkLogin(inp string) bool {
	data, err := os.ReadFile("login.txt")
	var login string
	if err != nil {
		login = "admin:admin"
	} else {
		login = string(data)
	}
	return inp == login
}

func main() {
	// Create router
	r := gin.Default()

	// Connect to the database
	db, err := sql.Open("adodb", "Provider=Microsoft.ACE.OLEDB.12.0;Data Source=database.accdb")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Serve static files
	r.GET("/", func(ctx *gin.Context) {
		c, err := ctx.Cookie("login")
		if err != nil || !checkLogin(c) {
			ctx.Redirect(302, "/login")
			return
		}
		ctx.File("./public/index.html")
	})
	r.GET("/books", func(ctx *gin.Context) {
		c, err := ctx.Cookie("login")
		if err != nil || !checkLogin(c) {
			ctx.Redirect(302, "/login")
			return
		}
		ctx.File("./public/books.html")
	})
	r.GET("/members", func(ctx *gin.Context) {
		c, err := ctx.Cookie("login")
		if err != nil || !checkLogin(c) {
			ctx.Redirect(302, "/login")
			return
		}
		ctx.File("./public/members.html")
	})
	r.GET("/borrowed", func(ctx *gin.Context) {
		c, err := ctx.Cookie("login")
		if err != nil || !checkLogin(c) {
			ctx.Redirect(302, "/login")
			return
		}
		ctx.File("./public/borrowed.html")
	})
	r.GET("/login", func(ctx *gin.Context) {
		ctx.File("./public/login.html")
	})

	r.POST("/login", func(ctx *gin.Context) {
		if checkLogin(ctx.PostForm("username") + ":" + ctx.PostForm("password")) {
			ctx.SetCookie("login", ctx.PostForm("username")+":"+ctx.PostForm("password"), 3600, "/", "", false, true)
			ctx.Redirect(302, "/")
		} else {
			ctx.Redirect(302, "/login?error=1")
		}
	})

	r.GET("/didntReturn", func(ctx *gin.Context) {
		// Execute a SELECT statement
		rows, err := db.Query("SELECT bookId, memId, dateBorrowed FROM borrowed WHERE dateReturned IS NULL")
		if err != nil {
			println(err.Error())
		}
		defer rows.Close()

		table := tablewriter.NewWriter(ctx.Writer)
		table.SetHeader([]string{"Book ID", "Member ID", "Date Taken"})
	
		// Iterate over the rows and print the values of the columns
		for rows.Next() {
			var bookId string
			var memId string
			var dateBorrowed sql.NullTime
			err := rows.Scan(&bookId, &memId, &dateBorrowed)
			if err != nil {
				println(err.Error())
			}
			println(bookId, memId, dateBorrowed.Time.Format("2006-01-02 06:15:04 PM"))
			table.Append([]string{bookId, memId, dateBorrowed.Time.Format("2006-01-02 06:15:04 PM")})
		}

		if table.NumLines() == 0 {
			ctx.String(200, "No books have to be returned")
		} else {
			table.Render()
		}
	})

	// API functions
	r.GET("/api/:fnc", func(ctx *gin.Context) {
		fnc := ctx.Param("fnc")
		ref := ctx.GetHeader("Referer")
		if ref == "" {
			ref = "/"
		}

		ctx.Header("Content-Type", "text/html")

		switch fnc {
		case "addBook":
			bookId := getNextPk(db, "books")
			_, err := db.Exec("INSERT INTO books (bookId, bookName, author) VALUES (?, ?, ?)", bookId, ctx.Query("name"), ctx.Query("author"))
			if err != nil {
				ctx.String(500, returnPg("Error", err.Error(), ref))
				return
			}
			ctx.String(200, returnPg("Info", "Book added successfully. Book ID: "+bookId, ref))

		case "removeBook":
			_, err := db.Exec("DELETE FROM books WHERE bookId = ?", ctx.Query("id"))
			if err != nil {
				ctx.String(500, returnPg("Error", err.Error(), ref))
				return
			}
			ctx.String(200, returnPg("Info", "Book removed successfully", ref))

		case "addMember":
			memId := getNextPk(db, "members")
			_, err := db.Exec("INSERT INTO members (memId, memName, grade) VALUES (?, ?, ?)", memId, ctx.Query("name"), ctx.Query("grade"))
			if err != nil {
				ctx.String(500, returnPg("Error", err.Error(), ref))
				return
			}
			ctx.String(200, returnPg("Info", "Member added successfully. Member ID: "+memId, ref))

		case "removeMember":
			_, err := db.Exec("DELETE FROM members WHERE memId = ?", ctx.Query("id"))
			if err != nil {
				ctx.String(500, returnPg("Error", err.Error(), ref))
				return
			}
			ctx.String(200, returnPg("Info", "Member removed successfully", ref))

		case "markTake":
			cTime := time.Now()
			_, err := db.Exec("INSERT INTO borrowed (bookId, memId, dateBorrowed, dateReturned) VALUES (?, ?, ?, ?)", ctx.Query("bookId"), ctx.Query("memId"), cTime, nil)
			if err != nil {
				ctx.String(500, returnPg("Error", err.Error(), ref))
				return
			}
			ctx.String(200, returnPg("Info", "Book take marked. Time: "+cTime.Format("2006-01-02 15:04:05"), ref))

		case "markReturn":
			cTime := time.Now()
			_, err := db.Exec("UPDATE borrowed SET dateReturned = ? WHERE bookId = ? AND memId = ?", cTime, ctx.Query("bookId"), ctx.Query("memId"))
			if err != nil {
				ctx.String(500, returnPg("Error", err.Error(), ref))
				return
			}
			ctx.String(200, returnPg("Info", "Book return marked. Time: "+cTime.Format("2006-01-02 15:04:05"), ref))

		}
	})

	// Start server
	r.Run(":8080")
}

func getNextPk(db *sql.DB, table string) string {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count)
	return strconv.Itoa(count + 1)
}
