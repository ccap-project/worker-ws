package role

import (
  "bytes"
  "crypto/tls"
  "encoding/base64"
  "fmt"
  "net/http"

  g "github.com/xanzy/go-gitlab"
  "github.com/ghodss/yaml"

  "../config"
)

type GitGroup_t struct {
  Name  string
  Path  string
  Group string
}

type Role_t struct {
  ID         int
  Name       string
  Url        string
  Versions   []string
}



/**********************************************
 *
 **********************************************/
func Init(SystemConfig *config.Configuration) {

  transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: SystemConfig.GitlabCfg.TLSInsecureSkipVerify}}
  client    := &http.Client{Transport: transport}

  SystemConfig.Log.Println(SystemConfig.GitlabCfg.Token)
  SystemConfig.GitlabState.Client = g.NewClient(client, SystemConfig.GitlabCfg.Token)

  if err := SystemConfig.GitlabState.Client.SetBaseURL(SystemConfig.GitlabCfg.Url); err != nil {
    SystemConfig.Log.Fatal(err)
  }

  id, err := SearchGroup(SystemConfig, SystemConfig.GitlabCfg.Group)

  if err != nil {
    SystemConfig.Log.Fatal(err)
  }

  SystemConfig.GitlabState.GroupID = id
}

/**********************************************
 *
 **********************************************/
func GetMeta(SystemConfig *config.Configuration, project string, version string) ([]byte, error){

  file, _, err := SystemConfig.GitlabState.Client.RepositoryFiles.GetFile(project, &g.GetFileOptions{FilePath: g.String("meta/main.yml"), Ref: g.String(version)})

  if err != nil {
    SystemConfig.Log.Fatal(err)
    //return 0, err
  }

  data, err := base64.StdEncoding.DecodeString(file.Content)
  if err != nil {
    SystemConfig.Log.Fatal(err)
    //return
  }

  /*
   * Filter galaxy_info node from yaml
   */
  var filtered_data []byte

  filtered_node := string("galaxy_info:")
  filtering_active := 0
  lines := bytes.Split(data, []byte{'\n'})

  for i := range(lines) {
     if filtering_active == 0 && bytes.IndexAny(lines[i], filtered_node) >= 0 {
       filtering_active++
       continue

     } else if filtering_active > 0 {
       x := 0
       for ( x < len(lines[i]) && lines[i][x] == ' ' ) {
         x++
       }
       
       if i == len(lines[i]) {
         continue
       }

       //fmt.Printf("%d - old(%s) new(%s)\n", i, string(lines[i]), string(lines[i][x:]))
       filtered_data = append(filtered_data, lines[i][x:]...)
       filtered_data = append(filtered_data, "\n"...)

     } else {
       filtering_active=0
       filtered_data = append(filtered_data, lines[i]...)
       filtered_data = append(filtered_data, "\n"...)
     }
  }

  //fmt.Println(string(filtered_data))

  m, err := yaml.YAMLToJSON(filtered_data)
  if err != nil {
    SystemConfig.Log.Fatal(err)
    //return
  }

  return m, err
}

/**********************************************
 *
 **********************************************/
func GetParams(SystemConfig *config.Configuration, project string, version string) ([]byte, error){

  file, _, err := SystemConfig.GitlabState.Client.RepositoryFiles.GetFile(project, &g.GetFileOptions{FilePath: g.String("defaults/main.yml"), Ref: g.String(version)})

  if err != nil {
    SystemConfig.Log.Fatal(err)
    //return 0, err
  }

  data, err := base64.StdEncoding.DecodeString(file.Content)
  if err != nil {
    SystemConfig.Log.Fatal(err)
    //return
  }

  m, err := yaml.YAMLToJSON(data)
  if err != nil {
    SystemConfig.Log.Fatal(err)
    //return
  }

  return m, err
}

/**********************************************
 *
 **********************************************/
func SearchGroup(SystemConfig *config.Configuration, group string) (int, error) {

  groups, _, err := SystemConfig.GitlabState.Client.Groups.SearchGroup(group)

  if err != nil {
    return 0, err
  }

  for i := range groups {
    if ( groups[i].Name == group ) {
      return groups[i].ID, nil
    }
  }

  return 0, fmt.Errorf("Group %s not found", group)
}


/**********************************************
 *
 **********************************************/
func List(SystemConfig *config.Configuration, id int) (*[]Role_t, error) {

  roles := []Role_t{}
  projects, _, err := SystemConfig.GitlabState.Client.Groups.ListGroupProjects(id)

  if err != nil {
    SystemConfig.Log.Fatal(err)
    return nil, err
  }

  for i := range projects {
     var role Role_t
     role.ID    = projects[i].ID
     role.Name  = projects[i].Name
     role.Url   = projects[i].HTTPURLToRepo

     tags, _, _ := SystemConfig.GitlabState.Client.Tags.ListTags(projects[i].ID)
    
     for j := range tags {
        role.Versions = append(role.Versions, tags[j].Name)
     }

     //copy(role.Versions, projects[i].TagList)

     //fmt.Println(projects[i].Name)
     //fmt.Println(projects[i].TagList)

     roles = append(roles, role)
  }

  return &roles, err
}
