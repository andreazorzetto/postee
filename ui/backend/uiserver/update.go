package uiserver

import (
	"encoding/json"
	"github.com/aquasecurity/postee/alertmgr"
	"github.com/ghodss/yaml"
	"io"
	"log"
	"net/http"
	"os"
)

func (srv *uiServer) updateConfig(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	inputJson, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can't read JSON string", http.StatusBadRequest)
		return
	}
	plugins := []alertmgr.PluginSettings{}
	if err := json.Unmarshal(inputJson, &plugins); err != nil {
		http.Error(w, "Can't read JSON string", http.StatusBadRequest)
		return
	}
	b, err := yaml.Marshal(&plugins)
	if err != nil {
		http.Error(w, "Can't convert PluginSettings into yaml string", http.StatusBadRequest)
		return
	}
	if err := os.Rename(srv.cfgPath, srv.cfgPath+".copy"); err != nil {
		log.Printf("rename file error %v", err)
		http.Error(w, "Can't remove data from the config file for overwrite", http.StatusBadRequest)
		return
	}
	f, err := os.Create(srv.cfgPath)
	if err != nil {
		log.Printf("create file error %v", err)
		http.Error(w, "Can't open the config file for overwrite", http.StatusBadRequest)
		return
	}
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		log.Printf("write file error %v", err)
		http.Error(w, "Can't write to the config file for overwrite", http.StatusBadRequest)
		return
	}
	os.RemoveAll(srv.cfgPath + ".copy")
}
