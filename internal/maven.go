/*
* Copyright 2024 Google LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* 	https://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func searchMaven(sha1 string) (*Dependency, error) {
	searchURL := fmt.Sprintf("https://search.maven.org/solrsearch/select?q=1:%s&rows=20&wt=json", sha1)
	resp, err := http.Get(searchURL)
	if err != nil {
			return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
			return nil, err
	}

	var parsedReply map[string]interface{}
	err = json.Unmarshal(body, &parsedReply)
	if err != nil {
			return nil, err
	}

	response := parsedReply["response"].(map[string]interface{})
	numFound := int(response["numFound"].(float64))
	if numFound == 1 {
			docs := response["docs"].([]interface{})
			jarInfo := docs[0].(map[string]interface{})
			dependency := &Dependency{
					GroupId:    jarInfo["g"].(string),
					ArtifactId: jarInfo["a"].(string),
					Version:    jarInfo["v"].(string),
			}
			return dependency, nil
	}

	return nil, nil
}