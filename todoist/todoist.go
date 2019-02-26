package todoist

import (
	"net/http"
	"net/url"
	"os"
	"log"
	"encoding/json"
	"time"
	"errors"
	"io/ioutil"

	"github.com/google/uuid"
)

var syncUrl string = "https://todoist.com/api/v7/sync"

type Project struct {
	Name string
	ID int
	GetItems func() []Item
}

type Item struct {
	ID int
	Content string
	ProjectId int
	DueDate time.Time
}

type Todoist struct {
	token string

	Projects []Project
	Items []Item

	httpClient http.Client
}

// Create a new todoist context.
func New(token string) (*Todoist, error) {
	return &Todoist{
		token: token,
		httpClient: http.Client{},
	}, nil
}

// GET's the complete set of data from the todoist API and loads it into the
// todoist object.
func (t *Todoist) ReadSync() (*Todoist, error) {
	req, err := http.NewRequest("GET", syncUrl, nil)

	if err != nil {
		return t, err
	}

	q := req.URL.Query()
	q.Add("token", t.token)
	q.Add("sync_token", "*")
	q.Add("resource_types", "[\"projects\",\"items\"]")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp,err := t.httpClient.Do(req)

	if err != nil {
		return t, err
	}

	if resp.StatusCode != 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return t, errors.New("response code not 200 "+bodyString)
	}

	var result map[string]interface{}
	// decoder is better when reading from a stream
	json.NewDecoder(resp.Body).Decode(&result)

	t = t.loadProjectsData(result["projects"].([]interface{}))
	t = t.loadItemData(result["items"].([]interface{}))

	return t, nil
}

type Command struct {
	Type string `json:"type"`
	Args interface{} `json:"args,omitempty"`
	UUID string `json:"uuid"`
	TempID string `json:"temp_id,omitempty"`
}

type JSON map[string]interface{}

func (t Todoist) ItemComplete(item Item) {
	// make the data
	args := make(JSON)
	args["id"] = item.ID

	command := Command{
		Type: "item_close",
		UUID: uuid.New().String(),
		Args: args,
	}

	commands := make([]Command, 1)
	commands[0] = command
	json, err := json.Marshal(commands)

	if err != nil {
		log.Fatal(err)
	}

	data := url.Values{}
	data.Set("token", t.token)
	data.Set("commands", string(json[:]))

	//client := http.Client{}

	//log.Println(data.Encode())
	//req, err := http.NewRequest("POST", "https://todoist.com/api/v7/sync", strings.NewReader(data.Encode()))

	//resp,err := client.Do(req)

	resp,err := http.PostForm("https://todoist.com/api/v7/sync", data)

	if err != nil {
		log.Println("fail")
		log.Println(os.Stderr, err)
	}

	if resp == nil {
		log.Println("OH NO")
	}

	// resync or summtin
}

func (t Todoist) ItemAdd(content string, date string, projectID int) error {
	// make the data
	args := make(JSON)
	args["content"] = content
	args["date_string"] = date
	args["project_id"] = projectID

	command := Command{
		Type: "item_add",
		UUID: uuid.New().String(),
		TempID: uuid.New().String(),
		Args: args,
	}

	commands := make([]Command, 1)
	commands[0] = command
	json, err := json.Marshal(commands)

	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("token", t.token)
	data.Set("commands", string(json[:]))

	resp,err := http.PostForm("https://todoist.com/api/v7/sync", data)

	if err != nil {
		return err
	}

	if resp == nil {
		log.Println("OH NO")
	}

	return nil
}

// Parse the JSON data from a sync for projects, load into the todoist
// instance.
func (t *Todoist) loadProjectsData(projects []interface{}) *Todoist {
	t.Projects = nil
	for _,value := range projects {
		project := value.(map[string]interface{})
		t.Projects = append(t.Projects, Project{
			Name: project["name"].(string),
			ID: int(project["id"].(float64)),
		})
	}

	return t
}

// Parse the JSON data from a sync for items.
func (t *Todoist) loadItemData(items []interface{}) *Todoist {
	t.Items = nil
	for _,value := range items {
		item := value.(map[string]interface{})
		// check due_date_utc set...
		var due time.Time

		if item["due_date_utc"] != nil {
			due,_ = time.Parse(
				"Mon 02 Jan 2006 15:04:05 -0700",
				item["due_date_utc"].(string))
		} else {
			due = time.Time{}
		}

		t.Items = append(t.Items, Item{
			ID: int(item["id"].(float64)),
			Content: item["content"].(string),
			ProjectId: int(item["project_id"].(float64)),
			DueDate: due,
		})
	}

	return t
}
