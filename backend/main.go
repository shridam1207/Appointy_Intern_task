package main

import ( //importing the necessary packages
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)



type Article struct { //Construct a basic struct of the response this api would recieve thorugh the POPST request
	Id                 int       `json:"Id"`
	Title              string    `json:"Title"`
	SubTitle           string    `json:"SubTitle"`
	Content            string    `json:"content"`
	Creation_Timestamp time.Time `json:Timestamp`
}

// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
var Articles []Article
var mu sync.Mutex
	var wg sync.WaitGroup
// The Homepage function will display the homepage
//And also since we are not allowed to use mux library, so
//split the url and check if any extra parameter is added to get a specific article
//And if so then call returnSingleArticle function

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	p := strings.Split(r.URL.Path, "/")
	if len(p) == 2 {
		fmt.Fprintf(w, "Welcome to the HomePage!")
		fmt.Println("Endpoint Hit: homePage")
	} else {
		returnSingleArticle(w, r, p[2])
		fmt.Println(p[2])
	}
}

func handleRequests() {

	http.HandleFunc("/", homePage)

	http.HandleFunc("/articles", getorpost)

	http.HandleFunc("/articles/search", searchQuery)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

//Since to accomodate both get and post through a single url the below logic is written

func getorpost(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		wg.Add(1)
		 createNewArticle(w, r ,&wg)
		wg.Wait()
	} else {
	
		 returnAllArticles(w, r)
	
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup ) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	mu.Lock()
	if r.Method == "POST" {
		reqBody, _ := ioutil.ReadAll(r.Body)
		var article Article
		json.Unmarshal(reqBody, &article)
		// update our global Articles array to include
		// our new Article
		article.Id = Articles[len(Articles)-1].Id + 1
		article.Creation_Timestamp = time.Now()
		Articles = append(Articles, article)

		json.NewEncoder(w).Encode(article)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
 	mu.Unlock()
 	wg.Done()
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request, p string) {

	// Loop over all of our Articles
	// if the article.Id equals the key we pass in
	// return the article encoded as JSON
	key, err := strconv.Atoi(p)
	fmt.Println(err)
	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
	fmt.Println("Endpoint Hit: returnSingleArticle")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	//returns all articles encoded as JSON
	
	js, err := json.Marshal(Articles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(js)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	fmt.Println("Endpoint Hit: returnAllArticles")
	
}

//Since our search method has to be case sensitive
//so for each of article struct object, their elements are stored
//in each custom slice and made lower case and appended to a another silce
func convArticleList(a *Article) []string {
	p := make([]string, 0)
	result := make([]string, 0)
	tSplit := strings.Split(a.Title, " ")
	sSplit := strings.Split(a.SubTitle, " ")
	cSplit := strings.Split(a.Content, " ")
	p = append(p, tSplit...)
	p = append(p, sSplit...)
	p = append(p, cSplit...)
	for _, e := range p {
		result = append(result, strings.ToLower(e))
	}
	return result
}

func searchQuery(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["q"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	key := keys[0]
	//converting the parameter value to lowercase
	var key_lower string = strings.ToLower(string(key))

	fmt.Println("Url Param 'key' is: " + string(key_lower))
	//Searching for the parameter in the final slice containing elements of all json objects.
	for _, article := range Articles {
		a := &article
		d := convArticleList(a)
		fmt.Println(d)
		temp := make([]string, 0)
		for i := range d {
			if d[i] == key_lower {
				temp = append(temp, d[i])
			}
		}
		if len(temp) > 0 {
			json.NewEncoder(w).Encode(article)
		}
	}

}

func main() {
	Articles = []Article{
		Article{Id: 1, Title: "Hello", SubTitle: "Article Description", Content: "Article Content", Creation_Timestamp: time.Now()},
		 //Article{Id: 2, Title: "Shridam 2", SubTitle: "Article Description", Content: "Article Content",Creation_Timestamp: time.Now()},
	}
	
	handleRequests()
}

