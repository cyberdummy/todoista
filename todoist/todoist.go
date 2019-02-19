package todoist

// API Stuff
import (
	"net/http"
	"os"
	"log"
	"encoding/json"
)

type Project struct {
	Name string
	Id int
}

type Item struct {
	Content string
	ProjectId int
}

type Todoist struct {
	HttpClient http.Client
	Token string

	Projects []Project
	Items []Item
}

func (t Todoist) Sync() Todoist {
	client := http.Client{}

	req, err := http.NewRequest("GET", "https://todoist.com/api/v7/sync", nil)

	q := req.URL.Query()
	q.Add("token", t.Token)
	q.Add("sync_token", "*")
	q.Add("resource_types", "[\"projects\",\"items\"]")

	req.URL.RawQuery = q.Encode()

	resp,err := client.Do(req)

	if err != nil {
		log.Println(os.Stderr, "Broken")
	}

	var result map[string]interface{}
	// decoder is better when reading from a stream
	json.NewDecoder(resp.Body).Decode(&result)
	t = t.LoadProjectsData(result["projects"].([]interface{}))

	//log.Println(result["items"])
	t = t.LoadItemData(result["items"].([]interface{}))

	//log.Println(t.Projects)
	//log.Println(t.Items)

	return t
}

func (t Todoist) LoadProjectsData(projects []interface{}) Todoist {
	for _,value := range projects {
		project := value.(map[string]interface{})
		t.Projects = append(t.Projects, Project{
			Name: project["name"].(string),
			Id: int(project["id"].(float64)),
		})
	}

	return t
}

func (t Todoist) LoadItemData(items []interface{}) Todoist {
	for _,value := range items {
		item := value.(map[string]interface{})
		t.Items = append(t.Items, Item{
			Content: item["content"].(string),
			ProjectId: int(item["project_id"].(float64)),
		})
	}

	return t
}
