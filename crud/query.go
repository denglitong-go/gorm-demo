package crud

import "fmt"

func ShowQuery() {
	initDB()
	defer closeDB()
	pingDB()

	// retrieveSingleObject()
	// retrievingObjectsWithPrimaryKey()
	// retrieveWithConditions()

	scan()
}

// db.Table() will keep the state previously statement queries,
// so we need to run db.Table() for each statement separately.
func retrieveSingleObject() {
	album := Album{ID: 2}
	// Get the first record ordered by primary key
	// SELECT * FROM album ORDER BY id LIMIT 1;
	db.Table(albumTable).First(&album)
	fmt.Printf("first album by ID: %v\n", album)

	// Get one record, no specified order
	// SELECT * FROM album LIMIT 1
	album2 := Album{}
	db.Table(albumTable).Take(&album2)
	fmt.Printf("first album in table: %v\n", album2)

	album3 := Album{}
	db.Table(albumTable).Last(&album3)
	fmt.Printf("last album in table: %v\n", album3)
}

func retrievingObjectsWithPrimaryKey() {
	album := Album{}
	// select * from album where id = 3
	db.Table(albumTable).First(&album, 3)
	fmt.Println("album id=3", album)

	album = Album{}
	db.Table(albumTable).First(&album, "24")
	fmt.Println("album id=24", album)

	var albums []Album
	db.Table(albumTable).Find(&albums, []int{1, 2, 3})
	fmt.Println("albums id=[1,2,3]", albums)

	album = Album{}
	db.Table(albumTable).First(&album, "id = ?", "1")
	fmt.Println("album id=? ?==1", album)

	albums = []Album{}
	db.Table(albumTable).Find(&albums)
	fmt.Println("all albums", albums)
}

func retrieveWithConditions() {
	album := Album{}
	db.Table(albumTable).Where("title = ?", "Sarah Vaughan").First(&album)
	fmt.Println("first title = ?", album)

	albums := []Album{}
	db.Table(albumTable).Where("title = ?", "Giant Steps").Find(&albums)
	fmt.Println("all title = ?", albums)

	// IN
	db.Table(albumTable).Where("title IN ?", []string{"a", "b"}).Find(&albums)
	// LIKE
	db.Table(albumTable).Where("title LIKE ?", "%some title%").Find(&albums)
	// AND
	db.Table(albumTable).Where("title = ? AND price <= ?", "a", "b").Find(&albums)
	// BETWEEN
	db.Table(albumTable).Where("price BETWEEN ? and ?", "a", "b").Find(&albums)

	album = Album{ID: 10}
	// select * from album where id = 10 and id = 20 order by id desc limit 1;
	db.Table(albumTable).Where("id = ?", 20).Find(&album)
}

func specifyStructSearchFields() {
	album := Album{
		Title: "a",
		Price: 0,
	}
	// GORM automatically ignore the zero value fields if we don't specify it in query
	// select * from album where title = "a"
	db.Table(albumTable).Where(&album).Find(&album)
	// select * from album where title = "a" and price = 0
	db.Table(albumTable).Where(&album, "Title", "Price").Find(&album)
	// select * from album where title = "a" and price = 0
	db.Table(albumTable).Where(&map[string]interface{}{
		"Title": "a",
		"Price": 0,
	}).First(&album)
}

func inlineCondition() {
	var album Album
	// select * from album where id = 1
	db.Table(albumTable).First(&album, "id = ?", "1")
	// select * from album where title = 'a'
	db.Table(albumTable).First(&album, "title = ?", "a")
	// select * from album where title = 'a'
	db.Table(albumTable).First(&album, Album{Title: "a"})
	// select * from album where title = 'a'
	db.Table(albumTable).First(&album, map[string]interface{}{"Title": "a"})
}

// Not Conditions
func notConditions() {
	var album Album
	var albums []Album
	// select * from album where not title = "a" order by id limit 1
	db.Table(albumTable).Not("title = ?", "a").First(&album)
	// select * from album where not title in ("a", "b")
	db.Table(albumTable).Not(map[string]interface{}{
		"title": []string{"a", "b"},
	}).Find(&albums)
	// select * from album where not title = "a" order by id limit 1
	db.Table(albumTable).Not(Album{Title: "a"}).First(&album)
	// select * from album where not id in (1,2,3)
	db.Not([]int64{1, 2, 3}).Find(&albums)
}

// Or Conditions

// Selecting Specific Fields

// Order

// Limit & Offset

// Group By & Having

// Distinct

// Joins

// Joins Preloading

// Joins a Derived Table

// Scan
func scan() {
	type Result struct {
		Number int64
		Name   string
		Author string
	}
	var result Result
	db.Table(albumTable).Select("id as number", "title as name", "artist as author").Scan(&result)
	// scan {1 Blue Train John Coltrane}
	fmt.Println("scan", result)

	query := "select id as number, title as name, artist as author from album where id = ?"
	db.Raw(query, 2).Scan(&result)
	// scan {2 Giant Steps John Coltrane}
	fmt.Println("scan", result)
}
