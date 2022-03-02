//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

//Package //(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

//Package caputilities ...
package caputilities

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/PluginCiscoACI/capresponse"
	"github.com/ODIM-Project/PluginCiscoACI/config"
	"github.com/fsnotify/fsnotify"
)

// Status holds the Status of plugin it will be intizaied during startup time
var Status capresponse.Status

// PluginStartTime hold the time from which plugin started
var PluginStartTime time.Time

// TrackConfigFileChanges monitors the config changes using fsnotfiy
func TrackConfigFileChanges(configFilePath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = watcher.Add(configFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	go func() {
		for {
			select {
			case fileEvent, ok := <-watcher.Events:
				if !ok {
					continue
				}
				if fileEvent.Op&fsnotify.Write == fsnotify.Write || fileEvent.Op&fsnotify.Remove == fsnotify.Remove {
					log.Debug("Modified file: " + fileEvent.Name)
					// update the plugin config
					if err := config.SetConfiguration(); err != nil {
						log.Error("While trying to set configuration, got: " + err.Error())
					}
				}
				//Reading file to continue the watch
				watcher.Add(configFilePath)
			case err, _ := <-watcher.Errors:
				if err != nil {
					log.Error(err.Error())
					defer watcher.Close()
				}
			}
		}
	}()
}
